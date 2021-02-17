package capabililty

import (
	"context"
	"encoding/json"
	"fmt"
	"lualsp/protocol"
	"path"
	"sync"

	"github.com/go-errors/errors"
)

type serverState int

const (
	serverCreated      = serverState(iota)
	serverInitializing // set once the server has received "initialize" request
	serverInitialized  // set once the server has received "initialized" request
	serverShutDown
)

func (s serverState) String() string {
	switch s {
	case serverCreated:
		return "created"
	case serverInitializing:
		return "initializing"
	case serverInitialized:
		return "initialized"
	case serverShutDown:
		return "shutDown"
	}
	return fmt.Sprintf("(unknown state: %d)", int(s))
}

func NewServer() *Server {
	return &Server{
		state: serverCreated,
	}
}

type Server struct {
	stateMu sync.Mutex
	state   serverState
	// notifications generated before serverInitialized
	notifications []*protocol.ShowMessageParams
}

func (s *Server) initialize(ctx context.Context, params *protocol.ParamInitialize) (*protocol.InitializeResult, error) {
	s.stateMu.Lock()
	if s.state >= serverInitializing {
		defer s.stateMu.Unlock()
		return nil, errors.Errorf("%w: initialize called while server in %v state", jsonrpc2.ErrInvalidRequest, s.state)
	}
	s.state = serverInitializing
	s.stateMu.Unlock()

	s.clientPID = int(params.ProcessID)
	s.progress.supportsWorkDoneProgress = params.Capabilities.Window.WorkDoneProgress

	options := s.session.Options()
	defer func() { s.session.SetOptions(options) }()

	if err := s.handleOptionResults(ctx, source.SetOptions(options, params.InitializationOptions)); err != nil {
		return nil, err
	}
	options.ForClientCapabilities(params.Capabilities)

	folders := params.WorkspaceFolders
	if len(folders) == 0 {
		if params.RootURI != "" {
			folders = []protocol.WorkspaceFolder{{
				URI:  string(params.RootURI),
				Name: path.Base(params.RootURI.SpanURI().Filename()),
			}}
		}
	}
	for _, folder := range folders {
		uri := span.URIFromURI(folder.URI)
		if !uri.IsFile() {
			continue
		}
		s.pendingFolders = append(s.pendingFolders, folder)
	}
	// gopls only supports URIs with a file:// scheme, so if we have no
	// workspace folders with a supported scheme, fail to initialize.
	if len(folders) > 0 && len(s.pendingFolders) == 0 {
		return nil, fmt.Errorf("unsupported URI schemes: %v (gopls only supports file URIs)", folders)
	}

	var codeActionProvider interface{} = true
	if ca := params.Capabilities.TextDocument.CodeAction; len(ca.CodeActionLiteralSupport.CodeActionKind.ValueSet) > 0 {
		// If the client has specified CodeActionLiteralSupport,
		// send the code actions we support.
		//
		// Using CodeActionOptions is only valid if codeActionLiteralSupport is set.
		codeActionProvider = &protocol.CodeActionOptions{
			CodeActionKinds: s.getSupportedCodeActions(),
		}
	}
	var renameOpts interface{} = true
	if r := params.Capabilities.TextDocument.Rename; r.PrepareSupport {
		renameOpts = protocol.RenameOptions{
			PrepareProvider: r.PrepareSupport,
		}
	}

	goplsVersion, err := json.Marshal(debug.VersionInfo())
	if err != nil {
		return nil, err
	}

	return &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			CallHierarchyProvider: true,
			CodeActionProvider:    codeActionProvider,
			CompletionProvider: protocol.CompletionOptions{
				TriggerCharacters: []string{"."},
			},
			DefinitionProvider:         true,
			TypeDefinitionProvider:     true,
			ImplementationProvider:     true,
			DocumentFormattingProvider: true,
			DocumentSymbolProvider:     true,
			WorkspaceSymbolProvider:    true,
			ExecuteCommandProvider: protocol.ExecuteCommandOptions{
				Commands: options.SupportedCommands,
			},
			FoldingRangeProvider:      true,
			HoverProvider:             true,
			DocumentHighlightProvider: true,
			DocumentLinkProvider:      protocol.DocumentLinkOptions{},
			ReferencesProvider:        true,
			RenameProvider:            renameOpts,
			SignatureHelpProvider: protocol.SignatureHelpOptions{
				TriggerCharacters: []string{"(", ","},
			},
			TextDocumentSync: &protocol.TextDocumentSyncOptions{
				Change:    protocol.Incremental,
				OpenClose: true,
				Save: protocol.SaveOptions{
					IncludeText: false,
				},
			},
			Workspace: protocol.WorkspaceGn{
				WorkspaceFolders: protocol.WorkspaceFoldersGn{
					Supported:           true,
					ChangeNotifications: "workspace/didChangeWorkspaceFolders",
				},
			},
		},
		ServerInfo: struct {
			Name    string `json:"name"`
			Version string `json:"version,omitempty"`
		}{
			Name:    "gopls",
			Version: string(goplsVersion),
		},
	}, nil
}

func (s *Server) initialized(ctx context.Context, params *protocol.InitializedParams) error {
	s.stateMu.Lock()
	if s.state >= serverInitialized {
		defer s.stateMu.Unlock()
		return errors.Errorf("%w: initialized called while server in %v state", jsonrpc2.ErrInvalidRequest, s.state)
	}
	s.state = serverInitialized
	s.stateMu.Unlock()

	for _, not := range s.notifications {
		s.client.ShowMessage(ctx, not)
	}
	s.notifications = nil

	options := s.session.Options()
	defer func() { s.session.SetOptions(options) }()

	if err := s.addFolders(ctx, s.pendingFolders); err != nil {
		return err
	}
	s.pendingFolders = nil

	if options.ConfigurationSupported && options.DynamicConfigurationSupported {
		registrations := []protocol.Registration{
			{
				ID:     "workspace/didChangeConfiguration",
				Method: "workspace/didChangeConfiguration",
			},
			{
				ID:     "workspace/didChangeWorkspaceFolders",
				Method: "workspace/didChangeWorkspaceFolders",
			},
		}
		if options.SemanticTokens {
			registrations = append(registrations, semanticTokenRegistration())
		}
		if err := s.client.RegisterCapability(ctx, &protocol.RegistrationParams{
			Registrations: registrations,
		}); err != nil {
			return err
		}
	}
	return nil
}

package capabililty

import (
	"context"
	"fmt"
	"lualsp/auxiliary"
	"lualsp/protocol"
	"path"
	"sync"

	"github.com/sourcegraph/jsonrpc2"
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
		state:    serverCreated,
		progress: progressManager{},
	}
}

type Server struct {
	//服务器状态机
	stateMu sync.Mutex
	state   serverState
	//客户端信息
	clientPID int
	init      protocol.InnerInitializeParams

	//进度条管理
	progress progressManager

	//工作区文件夹
	folders []protocol.WorkspaceFolder
	// notifications generated before serverInitialized
	notifications []*protocol.ShowMessageParams
}

func (s *Server) Initialize(ctx context.Context, params *protocol.ParamInitialize) (*protocol.InitializeResult, error) {
	s.stateMu.Lock()
	if s.state >= serverInitializing {
		defer s.stateMu.Unlock()
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidRequest, Message: s.state.String()}
	}
	s.state = serverInitializing
	s.stateMu.Unlock()

	//初始化参数
	s.init = params.InnerInitializeParams
	s.clientPID = int(params.ProcessID)
	s.progress.supportpro = params.Capabilities.Window.WorkDoneProgress

	folders := params.WorkspaceFolders
	if len(folders) == 0 {
		if params.RootURI != "" {
			folders = []protocol.WorkspaceFolder{{
				URI:  string(params.RootURI),
				Name: path.Base(string(auxiliary.URIFromURI(string(params.RootURI)))),
			}}
		}
	}
	for _, folder := range folders {
		uri := auxiliary.URIFromURI(folder.URI)
		if !uri.IsFile() {
			continue
		}
		s.folders = append(s.folders, folder)
	}
	if len(folders) > 0 && len(s.folders) == 0 {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidRequest}
	}
	/*
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
	*/
	return &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			TextDocumentSync: &protocol.TextDocumentSyncOptions{
				Change:    protocol.Incremental,
				OpenClose: true,
				WillSave:  true,
				Save: protocol.SaveOptions{
					IncludeText: true,
				},
			},
			CompletionProvider: protocol.CompletionOptions{
				//TriggerCharacters: []string{"."},
				ResolveProvider: true,
			},
			/*
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
			*/
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
			Name:    "lua-language-server-go",
			Version: "0.1",
		},
	}, nil
}

func (s *Server) Initialized(ctx context.Context, params *protocol.InitializedParams) error {
	s.stateMu.Lock()
	if s.state >= serverInitialized {
		defer s.stateMu.Unlock()
		return &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidRequest, Message: s.state.String()}
		//return errors.Errorf("%w: initialized called while server in %v state", jsonrpc2.ErrInvalidRequest, s.state)
	}
	s.state = serverInitialized
	s.stateMu.Unlock()
	/*
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
		}*/
	return nil
}

func (s *Server) Exit(context.Context) error {
	return nil
}
func (s *Server) Shutdown(context.Context) error {
	return nil
}

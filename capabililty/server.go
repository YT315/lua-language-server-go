package capabililty

import (
	"context"
	"fmt"
	"lualsp/analysis"
	"lualsp/auxiliary"
	"lualsp/protocol"
	"path"
	"strings"
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

	// notifications generated before serverInitialized
	notifications []*protocol.ShowMessageParams

	//工程
	project *analysis.Project
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

	//初始化工作区
	s.project = &analysis.Project{
		State: analysis.ProjectCreated,
	}
	//工作区文件夹
	folders := params.WorkspaceFolders
	if len(folders) == 0 {
		if params.RootURI != "" {
			folders = []protocol.WorkspaceFolder{{
				URI:  string(params.RootURI),
				Name: path.Base(string(params.RootURI)),
			}}
		}
	}
	for _, folder := range folders {
		if !strings.HasPrefix(folder.URI, "file://") {
			continue
		}
		path := auxiliary.UriToPath(folder.URI)
		s.project.Workspaces = append(s.project.Workspaces, &analysis.Workspace{
			RootPath: path,
			Files:    make(map[string]*analysis.File),
			Project:  s.project,
		})
	}
	if len(folders) > 0 && len(s.project.Workspaces) == 0 {
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidRequest}
	}
	//开始扫描文件
	go s.project.Scan(ctx)
	//返回配置
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
			CompletionProvider: &protocol.CompletionOptions{
				TriggerCharacters: []string{
					".",  //属性获取
					"=",  //赋值类型提示
					"==", //比较类型提示
					"!=",
				},
				ResolveProvider: true,
			},

			CodeActionProvider: true, // codeActionProvider,
			/*	CallHierarchyProvider: true,
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
			Workspace: &protocol.WorkspaceGn{
				WorkspaceFolders: &protocol.WorkspaceFoldersGn{
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

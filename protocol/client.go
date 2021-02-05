package protocol

import (
	"context"
)

type Client interface {
	ShowMessage(context.Context, *ShowMessageParams) error
	LogMessage(context.Context, *LogMessageParams) error
	Event(context.Context, *interface{}) error
	PublishDiagnostics(context.Context, *PublishDiagnosticsParams) error
	Progress(context.Context, *ProgressParams) error
	WorkspaceFolders(context.Context) ([]WorkspaceFolder /*WorkspaceFolder[] | null*/, error)
	Configuration(context.Context, *ParamConfiguration) ([]interface{}, error)
	WorkDoneProgressCreate(context.Context, *WorkDoneProgressCreateParams) error
	RegisterCapability(context.Context, *RegistrationParams) error
	UnregisterCapability(context.Context, *UnregistrationParams) error
	ShowMessageRequest(context.Context, *ShowMessageRequestParams) (*MessageActionItem /*MessageActionItem | null*/, error)
	ApplyEdit(context.Context, *ApplyWorkspaceEditParams) (*ApplyWorkspaceEditResponse, error)
}

func (s *clientDispatcher) ShowMessage(ctx context.Context, params *ShowMessageParams) error {
	return s.Conn.Notify(ctx, "window/showMessage", params)
}

func (s *clientDispatcher) LogMessage(ctx context.Context, params *LogMessageParams) error {
	return s.Conn.Notify(ctx, "window/logMessage", params)
}

func (s *clientDispatcher) Event(ctx context.Context, params *interface{}) error {
	return s.Conn.Notify(ctx, "telemetry/event", params)
}

func (s *clientDispatcher) PublishDiagnostics(ctx context.Context, params *PublishDiagnosticsParams) error {
	return s.Conn.Notify(ctx, "textDocument/publishDiagnostics", params)
}

func (s *clientDispatcher) Progress(ctx context.Context, params *ProgressParams) error {
	return s.Conn.Notify(ctx, "$/progress", params)
}
func (s *clientDispatcher) WorkspaceFolders(ctx context.Context) ([]WorkspaceFolder /*WorkspaceFolder[] | null*/, error) {
	var result []WorkspaceFolder /*WorkspaceFolder[] | null*/
	if err := Call(ctx, s.Conn, "workspace/workspaceFolders", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *clientDispatcher) Configuration(ctx context.Context, params *ParamConfiguration) ([]interface{}, error) {
	var result []interface{}
	if err := Call(ctx, s.Conn, "workspace/configuration", params, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *clientDispatcher) WorkDoneProgressCreate(ctx context.Context, params *WorkDoneProgressCreateParams) error {
	return Call(ctx, s.Conn, "window/workDoneProgress/create", params, nil) // Call, not Notify
}

func (s *clientDispatcher) RegisterCapability(ctx context.Context, params *RegistrationParams) error {
	return Call(ctx, s.Conn, "client/registerCapability", params, nil) // Call, not Notify
}

func (s *clientDispatcher) UnregisterCapability(ctx context.Context, params *UnregistrationParams) error {
	return Call(ctx, s.Conn, "client/unregisterCapability", params, nil) // Call, not Notify
}

func (s *clientDispatcher) ShowMessageRequest(ctx context.Context, params *ShowMessageRequestParams) (*MessageActionItem /*MessageActionItem | null*/, error) {
	var result *MessageActionItem /*MessageActionItem | null*/
	if err := Call(ctx, s.Conn, "window/showMessageRequest", params, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *clientDispatcher) ApplyEdit(ctx context.Context, params *ApplyWorkspaceEditParams) (*ApplyWorkspaceEditResponse, error) {
	var result *ApplyWorkspaceEditResponse
	if err := Call(ctx, s.Conn, "workspace/applyEdit", params, &result); err != nil {
		return nil, err
	}
	return result, nil
}

package protocol

import (
	"context"

	"github.com/sourcegraph/jsonrpc2"
)

type (
	Client interface {
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
	clientDispatcher struct {
		Conn jsonrpc2.JSONRPC2
	}
)

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
	if err := s.Conn.Call(ctx, "workspace/workspaceFolders", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *clientDispatcher) Configuration(ctx context.Context, params *ParamConfiguration) ([]interface{}, error) {
	var result []interface{}
	if err := s.Conn.Call(ctx, "workspace/configuration", params, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *clientDispatcher) WorkDoneProgressCreate(ctx context.Context, params *WorkDoneProgressCreateParams) error {
	return s.Conn.Call(ctx, "window/workDoneProgress/create", params, nil) // Call, not Notify
}

func (s *clientDispatcher) RegisterCapability(ctx context.Context, params *RegistrationParams) error {
	return s.Conn.Call(ctx, "client/registerCapability", params, nil) // Call, not Notify
}

func (s *clientDispatcher) UnregisterCapability(ctx context.Context, params *UnregistrationParams) error {
	return s.Conn.Call(ctx, "client/unregisterCapability", params, nil) // Call, not Notify
}

func (s *clientDispatcher) ShowMessageRequest(ctx context.Context, params *ShowMessageRequestParams) (*MessageActionItem /*MessageActionItem | null*/, error) {
	var result *MessageActionItem /*MessageActionItem | null*/
	if err := s.Conn.Call(ctx, "window/showMessageRequest", params, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *clientDispatcher) ApplyEdit(ctx context.Context, params *ApplyWorkspaceEditParams) (*ApplyWorkspaceEditResponse, error) {
	var result *ApplyWorkspaceEditResponse
	if err := s.Conn.Call(ctx, "workspace/applyEdit", params, &result); err != nil {
		return nil, err
	}
	return result, nil
}

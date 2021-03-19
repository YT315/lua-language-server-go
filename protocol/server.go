package protocol

// Package protocol contains data types and code for LSP jsonrpcs
// generated automatically from vscode-languageserver-node
// commit: 901fd40345060d159f07d234bbc967966a929a34
// last fetched Mon Oct 26 2020 09:10:42 GMT-0400 (Eastern Daylight Time)

// Code generated (see typescript/README.md) DO NOT EDIT.

import (
	"context"
	"encoding/json"

	"github.com/sourcegraph/jsonrpc2"
	errors "golang.org/x/xerrors"
)

type Server interface {
	DidChangeWorkspaceFolders(context.Context, *DidChangeWorkspaceFoldersParams) error
	WorkDoneProgressCancel(context.Context, *WorkDoneProgressCancelParams) error
	Initialized(context.Context, *InitializedParams) error
	Exit(context.Context) error
	DidChangeConfiguration(context.Context, *DidChangeConfigurationParams) error
	DidOpen(context.Context, *DidOpenTextDocumentParams) error
	DidChange(context.Context, *DidChangeTextDocumentParams) error
	DidClose(context.Context, *DidCloseTextDocumentParams) error
	DidSave(context.Context, *DidSaveTextDocumentParams) error
	WillSave(context.Context, *WillSaveTextDocumentParams) error
	DidChangeWatchedFiles(context.Context, *DidChangeWatchedFilesParams) error
	SetTrace(context.Context, *SetTraceParams) error
	LogTrace(context.Context, *LogTraceParams) error
	Implementation(context.Context, *ImplementationParams) (Definition /*Definition | DefinitionLink[] | null*/, error)
	TypeDefinition(context.Context, *TypeDefinitionParams) (Definition /*Definition | DefinitionLink[] | null*/, error)
	DocumentColor(context.Context, *DocumentColorParams) ([]ColorInformation, error)
	ColorPresentation(context.Context, *ColorPresentationParams) ([]ColorPresentation, error)
	FoldingRange(context.Context, *FoldingRangeParams) ([]FoldingRange /*FoldingRange[] | null*/, error)
	Declaration(context.Context, *DeclarationParams) (Declaration /*Declaration | DeclarationLink[] | null*/, error)
	SelectionRange(context.Context, *SelectionRangeParams) ([]SelectionRange /*SelectionRange[] | null*/, error)
	PrepareCallHierarchy(context.Context, *CallHierarchyPrepareParams) ([]CallHierarchyItem /*CallHierarchyItem[] | null*/, error)
	IncomingCalls(context.Context, *CallHierarchyIncomingCallsParams) ([]CallHierarchyIncomingCall /*CallHierarchyIncomingCall[] | null*/, error)
	OutgoingCalls(context.Context, *CallHierarchyOutgoingCallsParams) ([]CallHierarchyOutgoingCall /*CallHierarchyOutgoingCall[] | null*/, error)
	SemanticTokensFull(context.Context, *SemanticTokensParams) (*SemanticTokens /*SemanticTokens | null*/, error)
	SemanticTokensFullDelta(context.Context, *SemanticTokensDeltaParams) (interface{} /* SemanticTokens | SemanticTokensDelta | nil*/, error)
	SemanticTokensRange(context.Context, *SemanticTokensRangeParams) (*SemanticTokens /*SemanticTokens | null*/, error)
	SemanticTokensRefresh(context.Context) error
	Initialize(context.Context, *ParamInitialize) (*InitializeResult, error)
	Shutdown(context.Context) error
	WillSaveWaitUntil(context.Context, *WillSaveTextDocumentParams) ([]TextEdit /*TextEdit[] | null*/, error)
	Completion(context.Context, *CompletionParams) (*CompletionList /*CompletionItem[] | CompletionList | null*/, error)
	Resolve(context.Context, *CompletionItem) (*CompletionItem, error)
	Hover(context.Context, *HoverParams) (*Hover /*Hover | null*/, error)
	SignatureHelp(context.Context, *SignatureHelpParams) (*SignatureHelp /*SignatureHelp | null*/, error)
	Definition(context.Context, *DefinitionParams) (Definition /*Definition | DefinitionLink[] | null*/, error)
	References(context.Context, *ReferenceParams) ([]Location /*Location[] | null*/, error)
	DocumentHighlight(context.Context, *DocumentHighlightParams) ([]DocumentHighlight /*DocumentHighlight[] | null*/, error)
	DocumentSymbol(context.Context, *DocumentSymbolParams) ([]interface{} /*SymbolInformation[] | DocumentSymbol[] | null*/, error)
	CodeAction(context.Context, *CodeActionParams) ([]CodeAction /*(Command | CodeAction)[] | null*/, error)
	ResolveCodeAction(context.Context, *CodeAction) (*CodeAction, error)
	Symbol(context.Context, *WorkspaceSymbolParams) ([]SymbolInformation /*SymbolInformation[] | null*/, error)
	CodeLens(context.Context, *CodeLensParams) ([]CodeLens /*CodeLens[] | null*/, error)
	ResolveCodeLens(context.Context, *CodeLens) (*CodeLens, error)
	DocumentLink(context.Context, *DocumentLinkParams) ([]DocumentLink /*DocumentLink[] | null*/, error)
	ResolveDocumentLink(context.Context, *DocumentLink) (*DocumentLink, error)
	Formatting(context.Context, *DocumentFormattingParams) ([]TextEdit /*TextEdit[] | null*/, error)
	RangeFormatting(context.Context, *DocumentRangeFormattingParams) ([]TextEdit /*TextEdit[] | null*/, error)
	OnTypeFormatting(context.Context, *DocumentOnTypeFormattingParams) ([]TextEdit /*TextEdit[] | null*/, error)
	Rename(context.Context, *RenameParams) (*WorkspaceEdit /*WorkspaceEdit | null*/, error)
	PrepareRename(context.Context, *PrepareRenameParams) (*Range /*Range | { range: Range, placeholder: string } | { defaultBehavior: boolean } | null*/, error)
	ExecuteCommand(context.Context, *ExecuteCommandParams) (interface{} /*any | null*/, error)
	Moniker(context.Context, *MonikerParams) ([]Moniker /*Moniker[] | null*/, error)
	NonstandardRequest(ctx context.Context, method string, params interface{}) (interface{}, error)
}

func serverDispatch(ctx context.Context, server Server, reply func(context.Context, interface{}, error) error, r jsonrpc2.Request) (bool, error) {
	switch r.Method {
	case "workspace/didChangeWorkspaceFolders": // notif
		var params DidChangeWorkspaceFoldersParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		err := server.DidChangeWorkspaceFolders(ctx, &params)
		return true, reply(ctx, nil, err)
	case "window/workDoneProgress/cancel": // notif
		var params WorkDoneProgressCancelParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		err := server.WorkDoneProgressCancel(ctx, &params)
		return true, reply(ctx, nil, err)
	case "initialized": // notif
		var params InitializedParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		err := server.Initialized(ctx, &params)
		return true, reply(ctx, nil, err)
	case "exit": // notif
		err := server.Exit(ctx)
		return true, reply(ctx, nil, err)
	case "workspace/didChangeConfiguration": // notif
		var params DidChangeConfigurationParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		err := server.DidChangeConfiguration(ctx, &params)
		return true, reply(ctx, nil, err)
	case "textDocument/didOpen": // notif
		var params DidOpenTextDocumentParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		err := server.DidOpen(ctx, &params)
		return true, reply(ctx, nil, err)
	case "textDocument/didChange": // notif
		var params DidChangeTextDocumentParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		err := server.DidChange(ctx, &params)
		return true, reply(ctx, nil, err)
	case "textDocument/didClose": // notif
		var params DidCloseTextDocumentParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		err := server.DidClose(ctx, &params)
		return true, reply(ctx, nil, err)
	case "textDocument/didSave": // notif
		var params DidSaveTextDocumentParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		err := server.DidSave(ctx, &params)
		return true, reply(ctx, nil, err)
	case "textDocument/willSave": // notif
		var params WillSaveTextDocumentParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		err := server.WillSave(ctx, &params)
		return true, reply(ctx, nil, err)
	case "workspace/didChangeWatchedFiles": // notif
		var params DidChangeWatchedFilesParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		err := server.DidChangeWatchedFiles(ctx, &params)
		return true, reply(ctx, nil, err)
	case "$/setTrace": // notif
		var params SetTraceParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		err := server.SetTrace(ctx, &params)
		return true, reply(ctx, nil, err)
	case "$/logTrace": // notif
		var params LogTraceParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		err := server.LogTrace(ctx, &params)
		return true, reply(ctx, nil, err)
	case "textDocument/implementation": // req
		var params ImplementationParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.Implementation(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/typeDefinition": // req
		var params TypeDefinitionParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.TypeDefinition(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/documentColor": // req
		var params DocumentColorParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.DocumentColor(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/colorPresentation": // req
		var params ColorPresentationParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.ColorPresentation(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/foldingRange": // req
		var params FoldingRangeParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.FoldingRange(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/declaration": // req
		var params DeclarationParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.Declaration(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/selectionRange": // req
		var params SelectionRangeParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.SelectionRange(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/prepareCallHierarchy": // req
		var params CallHierarchyPrepareParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.PrepareCallHierarchy(ctx, &params)
		return true, reply(ctx, resp, err)
	case "callHierarchy/incomingCalls": // req
		var params CallHierarchyIncomingCallsParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.IncomingCalls(ctx, &params)
		return true, reply(ctx, resp, err)
	case "callHierarchy/outgoingCalls": // req
		var params CallHierarchyOutgoingCallsParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.OutgoingCalls(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/semanticTokens/full": // req
		var params SemanticTokensParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.SemanticTokensFull(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/semanticTokens/full/delta": // req
		var params SemanticTokensDeltaParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.SemanticTokensFullDelta(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/semanticTokens/range": // req
		var params SemanticTokensRangeParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.SemanticTokensRange(ctx, &params)
		return true, reply(ctx, resp, err)
	case "workspace/semanticTokens/refresh": // req
		if len(*r.Params) > 0 {
			return true, reply(ctx, nil, errors.Errorf("expected no params"))
		}
		err := server.SemanticTokensRefresh(ctx)
		return true, reply(ctx, nil, err)
	case "initialize": // req
		var params ParamInitialize
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.Initialize(ctx, &params)
		return true, reply(ctx, resp, err)
	case "shutdown": // req
		if len(*r.Params) > 0 {
			return true, reply(ctx, nil, errors.Errorf("expected no params"))
		}
		err := server.Shutdown(ctx)
		return true, reply(ctx, nil, err)
	case "textDocument/willSaveWaitUntil": // req
		var params WillSaveTextDocumentParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.WillSaveWaitUntil(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/completion": // req
		var params CompletionParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.Completion(ctx, &params)
		return true, reply(ctx, resp, err)
	case "completionItem/resolve": // req
		var params CompletionItem
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.Resolve(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/hover": // req
		var params HoverParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.Hover(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/signatureHelp": // req
		var params SignatureHelpParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.SignatureHelp(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/definition": // req
		var params DefinitionParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.Definition(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/references": // req
		var params ReferenceParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.References(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/documentHighlight": // req
		var params DocumentHighlightParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.DocumentHighlight(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/documentSymbol": // req
		var params DocumentSymbolParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.DocumentSymbol(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/codeAction": // req
		var params CodeActionParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.CodeAction(ctx, &params)
		return true, reply(ctx, resp, err)
	case "codeAction/resolve": // req
		var params CodeAction
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.ResolveCodeAction(ctx, &params)
		return true, reply(ctx, resp, err)
	case "workspace/symbol": // req
		var params WorkspaceSymbolParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.Symbol(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/codeLens": // req
		var params CodeLensParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.CodeLens(ctx, &params)
		return true, reply(ctx, resp, err)
	case "codeLens/resolve": // req
		var params CodeLens
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.ResolveCodeLens(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/documentLink": // req
		var params DocumentLinkParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.DocumentLink(ctx, &params)
		return true, reply(ctx, resp, err)
	case "documentLink/resolve": // req
		var params DocumentLink
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.ResolveDocumentLink(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/formatting": // req
		var params DocumentFormattingParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.Formatting(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/rangeFormatting": // req
		var params DocumentRangeFormattingParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.RangeFormatting(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/onTypeFormatting": // req
		var params DocumentOnTypeFormattingParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.OnTypeFormatting(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/rename": // req
		var params RenameParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.Rename(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/prepareRename": // req
		var params PrepareRenameParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.PrepareRename(ctx, &params)
		return true, reply(ctx, resp, err)
	case "workspace/executeCommand": // req
		var params ExecuteCommandParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.ExecuteCommand(ctx, &params)
		return true, reply(ctx, resp, err)
	case "textDocument/moniker": // req
		var params MonikerParams
		if err := json.Unmarshal(*r.Params, &params); err != nil {
			return true, reply(ctx, nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: err.Error()})
		}
		resp, err := server.Moniker(ctx, &params)
		return true, reply(ctx, resp, err)

	default:
		if r.Method[0] == '$' {
			return true, nil
		}
		return false, nil
	}
}

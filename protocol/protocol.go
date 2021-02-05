package protocol

import (
	"context"

	errors "golang.org/x/xerrors"
)

func sendParseError(ctx context.Context, reply jsonrpc2.Replier, err error) error {
	return reply(ctx, nil, errors.Errorf("%w: %s", jsonrpc2.ErrParse, err))
}

package capabililty

import (
	"context"
	"lualsp/protocol"
)

func newProgressManager() {

}

type (
	//进度管理
	progressManager struct {
		supportpro  bool
		allProgress map[protocol.ProgressToken]*progress
	}
	//进度
	progress struct {
		cancel     bool
		Percentage float64
	}
)

func (s *Server) WorkDoneProgressCancel(context.Context, *protocol.WorkDoneProgressCancelParams) error {
	return nil
}

package capabililty

import (
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

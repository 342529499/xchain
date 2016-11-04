package server

import (
	pb "github.com/1851616111/xchain/pkg/protos"
	cm "github.com/1851616111/xchain/pkg/server/message_cons_mng"
	"sync"
)

var (
	xcodeOnce *sync.Once = new(sync.Once)
	xcode     *xcode
)

func GetXCode() *xcode {
	if xcode == nil {
		xcode = &xcode{
			m: map[string]cm.Connection{},
		}
	}

	return &xcode
}

type xcode struct {
	m map[string]cm.Connection
}

func (c *xcode) Transport(s pb.CodeService_ServiceServer) error {

	return nil
}

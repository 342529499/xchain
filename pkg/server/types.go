package server

import (
	pb "github.com/1851616111/xchain/pkg/protos"
	cm "github.com/1851616111/xchain/pkg/server/message_cons_mng"
)

type recvConnetionMetadata struct {
	client pb.EndPoint
	con    cm.Connection
	errCh  chan error
	doneCh chan struct{}
}

type lounchConnectionMetadata struct {
	targetAddress string
	con           cm.Connection
	errCh         chan error
	doneCh        chan struct{}
}

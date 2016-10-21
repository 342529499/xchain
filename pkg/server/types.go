package server

import (
	pb "github.com/1851616111/xchain/pkg/protos"
	cm "github.com/1851616111/xchain/pkg/server/connection_manager"
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

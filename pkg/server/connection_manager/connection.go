package connection_manager

import (
	pb "github.com/1851616111/xchain/pkg/protos"
)

type Connection interface {
	Send(*pb.Message) error
	Recv() (*pb.Message, error)
}

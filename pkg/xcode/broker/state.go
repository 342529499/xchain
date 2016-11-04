package broker

import (
	pb "github.com/1851616111/xchain/pkg/protos"
)

type Stating struct {
	command    *pb.XCodeStateSpec
	responseCh chan interface{}
}

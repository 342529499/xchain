package instruction_cons_mng

import (
	pb "github.com/1851616111/xchain/pkg/protos"
)

type Connection interface {
	Send(*pb.Instruction) error
	Recv() (*pb.Instruction, error)
}

package broker

import (
	pb "github.com/1851616111/xchain/pkg/protos"
	"github.com/golang/protobuf/proto"
)

func NewStartInstruction(name, address, language string) *pb.Instruction {
	start := &pb.Instruction_Start{
		Name:    name,
		Address: address,
		Lang:    language,
	}

	b, _ := proto.Marshal(start)

	return &pb.Instruction{
		Action:     pb.Action_Request,
		Type:       pb.Instruction_START,
		Identifier: getIdentifier(),
		Payload:    b,
	}
}

func parseStartInstruction(start *pb.Instruction_Start) func(i *pb.Instruction) error {
	return func(i *pb.Instruction) error {
		if start == nil {
			return ERRPointerObjectNil
		}

		return proto.Unmarshal(i.Payload, start)
	}
}

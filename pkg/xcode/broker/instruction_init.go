package broker

import (
	pb "github.com/1851616111/xchain/pkg/protos"
	"github.com/golang/protobuf/proto"
)

func parseInitInstruction(spec *pb.Instruction_Init) func(i *pb.Instruction) error {
	return func(i *pb.Instruction) error {
		if spec == nil {
			return ERRPointerObjectNil
		}

		return proto.Unmarshal(i.Payload, spec)
	}
}

func NewInitInstruction(function string, args []string) *pb.Instruction {
	init := &pb.Instruction_Init{
		Init: &pb.InstructionUnit{
			TransactionID: "123456",
			Function:      function,
			Args:          args,
		}}

	b, _ := proto.Marshal(init)

	return &pb.Instruction{
		Action:     pb.Action_Request,
		Type:       pb.Instruction_INIT,
		Identifier: getIdentifier(),
		Payload:    b,
	}
}

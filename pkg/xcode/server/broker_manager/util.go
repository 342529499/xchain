package server

import (
	pb "github.com/1851616111/xchain/pkg/protos"
	errors "github.com/1851616111/xchain/pkg/xcode/broker"
	"github.com/golang/protobuf/proto"
)

func parseStateInstruction(from *pb.Instruction) (*pb.Instruction_State, error) {
	if from == nil {
		return nil, errors.ERRPointerObjectNil
	}

	state := &pb.Instruction_State{}
	if err := proto.Unmarshal(from.Payload, state); err != nil {
		return nil, err
	}

	return state, nil
}

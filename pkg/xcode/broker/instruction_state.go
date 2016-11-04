package broker

import (
	pb "github.com/1851616111/xchain/pkg/protos"
	"github.com/golang/protobuf/proto"
)

func genGetStateInstruction(key string) ([]byte, error) {
	state := &pb.Instruction_State{
		Type: pb.Instruction_State_GET,
		Key:  key,
	}

	b, _ := proto.Marshal(state)

	i := &pb.Instruction{
		Action:     pb.Action_Request,
		Type:       pb.Instruction_STATE,
		Identifier: getIdentifier(),
		Payload:    b,
	}

	returnInstruction, err := execute(i)
	if err != nil {
		return nil, err
	}

	if err := parseStateInstruction(returnInstruction)(state); err != nil {
		return nil, err
	}

	return state.Value, nil
}

func genPutStateInstruction(key string, value []byte) error {
	state := &pb.Instruction_State{
		Type:  pb.Instruction_State_PUT,
		Key:   key,
		Value: value,
	}

	b, _ := proto.Marshal(state)

	i := &pb.Instruction{
		Action:     pb.Action_Request,
		Type:       pb.Instruction_STATE,
		Identifier: getIdentifier(),
		Payload:    b,
	}

	returnInstruction, err := execute(i)
	if err != nil {
		return err
	}

	if !IsOKInstruction(returnInstruction) {
		return ERRUnExpectedResponsePayload
	}

	return nil
}

func genDelStateInstruction(key string) error {
	state := &pb.Instruction_State{
		Type: pb.Instruction_State_DELETE,
		Key:  key,
	}

	b, _ := proto.Marshal(state)

	i := &pb.Instruction{
		Action:     pb.Action_Request,
		Type:       pb.Instruction_STATE,
		Identifier: getIdentifier(),
		Payload:    b,
	}

	returnInstruction, err := execute(i)
	if err != nil {
		return err
	}

	if !IsOKInstruction(returnInstruction) {
		return ERRUnExpectedResponsePayload
	}

	return nil
}

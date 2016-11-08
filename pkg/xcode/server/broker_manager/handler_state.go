package server

import (
	pb "github.com/1851616111/xchain/pkg/protos"
)
func (m *manager) handleState(id string, state *pb.Instruction_State, response ResponseWriter) {
	switch state.Type {
		case pb.Instruction_State_PUT:

		case pb.Instruction_State_GET:

		case pb.Instruction_State_DELETE:
	}

	return
}

package broker

import (
	pb "github.com/1851616111/xchain/pkg/protos"
	"github.com/golang/protobuf/proto"
	"os"
	"syscall"
	"time"
)

const (
	KEY_LISTENER_ADDRESS     = "XCHAIN_XCODE_LISTENER_ADDRESS"
	KEY_STAND_ALONE_TTL      = "XCHAIN_XCODE_STAND_ALONE_TTL"
	DEFAULT_LISTENER_ADDRESS = "0.0.0.0:10692"
)

func getListenerAddress() string {
	if address, found := getEnvAddress(KEY_LISTENER_ADDRESS); found {
		return address
	} else {
		return DEFAULT_LISTENER_ADDRESS
	}
}

func getEnvAddress(key string) (string, bool) {
	return syscall.Getenv(key)
}

func (s *coderService) watchStart() {
	ttl, _ := getEnvAddress(KEY_STAND_ALONE_TTL)
	if ttl == "" {
		return
	}

	duration, err := time.ParseDuration(ttl)
	if err != nil {
		logger.Printf("parse env key:%s, value:%s err", KEY_STAND_ALONE_TTL, ttl)
		return
	}

	for {
		select {
		case <-s.startCh:
			close(s.startCh)
			return
		case <-time.Tick(duration):
			logger.Printf("no found xcode server input request for %d", duration)
			os.Exit(1)
		}
	}
}

func parseInvokeInstruction(spec *pb.Instruction_Invoke) func(i *pb.Instruction) error {
	return func(i *pb.Instruction) error {
		if spec == nil {
			return ERRPointerObjectNil
		}

		return proto.Unmarshal(i.Payload, spec)
	}
}

func parseQueryInstruction(spec *pb.Instruction_Query) func(i *pb.Instruction) error {
	return func(i *pb.Instruction) error {
		if spec == nil {
			return ERRPointerObjectNil
		}

		return proto.Unmarshal(i.Payload, spec)
	}
}

func parseStateInstruction(from *pb.Instruction) func(target *pb.Instruction_State) error {

	return func(target *pb.Instruction_State) error {
		if from == nil {
			return ERRPointerObjectNil
		}

		if err := proto.Unmarshal(from.Payload, target); err != nil {
			return err
		}

		return nil
	}
}

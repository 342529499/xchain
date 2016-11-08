package broker

import (
	"fmt"
	pb "github.com/1851616111/xchain/pkg/protos"
	"google.golang.org/grpc/peer"
	"io"
	"strings"
)

type InstructionHandler func(i *pb.Instruction) error

func (s *coderService) handle(codeStream pb.CodeService_ExecuteServer) error {

	rsp := make(chan *pb.Instruction, 30)
	stop := make(chan struct{}, 1)
	go func() {
		for {
			select {
			case msg := <-rsp:
				codeStream.Send(msg)
			case rspMsg := <-s.sendCh:
				codeStream.Send(rspMsg)
			case <-stop:
				return
			}
		}
	}()

	for {
		instruction, err := codeStream.Recv()
		if err == io.EOF {
			fmt.Println("read eof")
		}

		if err != nil {
			return err
		}

		if !s.started {
			if isInstructionStart(instruction) {
				p, _ := peer.FromContext(codeStream.Context())
				s.handleStart(instruction, p.Addr.String(), rsp)
			}

			continue
		}

		if !s.isAuthorized(instruction, codeStream) {
			continue
		}

		if found := s.handleResponse(instruction); found {
			continue
		}

		switch instruction.Type {
		case pb.Instruction_KEEPALIVE:
			go s.handleKeepAlive(instruction, rsp)

		case pb.Instruction_INIT:
			go s.handleInit(instruction, rsp)

		case pb.Instruction_INVOKE:
		case pb.Instruction_QUERY:
		case pb.Instruction_STATE:
			fmt.Printf("xcode broker recv state operate %#v\n", *instruction)

		default:
			logger.Printf("recv unknown type msg: %#v\n", *instruction)
		}

	}
	return nil
}

func isInstructionStart(i *pb.Instruction) bool {
	return i.Type == pb.Instruction_START
}

func (s *coderService) start(i *pb.Instruction, codeStream pb.CodeService_ExecuteServer) {
	p, _ := peer.FromContext(codeStream.Context())
	s.nodeAddress = p.Addr.String()
}

func (s *coderService) isAuthorized(i *pb.Instruction, codeStream pb.CodeService_ExecuteServer) bool {
	p, _ := peer.FromContext(codeStream.Context())
	return strings.Contains(p.Addr.String(), s.nodeAddress)
}
func (s *coderService) handleResponse(i *pb.Instruction) bool {
	if i.Action == pb.Action_Response {
		c, found := s.IDCh[i.Identifier]
		if found {
			c <- i
		}

		return found
	}

	return false
}

func (s *coderService) handleStart(i *pb.Instruction, address string, response chan *pb.Instruction) {
	s.Lock()
	defer s.Unlock()

	//已经启动，退出
	if s.started {
		return
	}

	//没有启动的处理
	s.startCh <- struct{}{}
	startInstruction := new(pb.Instruction_Start)
	if err := parseStartInstruction(startInstruction)(i); err != nil {
		res := newReturnInstruction(i)
		res.Type = pb.Instruction_ERROR
		res.Payload = []byte(err.Error())

		response <- res
		return
	}

	if !strings.Contains(address, startInstruction.Address) {
		return
	}

	s.started = true
	s.nodeAddress = startInstruction.Address
	s.nodeName = startInstruction.Name
	s.language = startInstruction.Lang

	logger.Println("xcode has been started")
	response <- ReturnOKInstruction(i)
	return
}

func (s *coderService) handleKeepAlive(i *pb.Instruction, response chan *pb.Instruction) {
	return
}

//init
//invoke
//query
func (s *coderService) handleInit(i *pb.Instruction, response chan *pb.Instruction) {
	s.Lock()
	defer s.Unlock()

	spec := new(pb.Instruction_Init)
	//这里有一些小问题，类型被改成error后服务端收到后需要存储identifier同时还要存储消息的类型
	if err := parseInitInstruction(spec)(i); err != nil {
		response <- ReturnErrorInstruction(i, err)
		return
	}

	var responseInstruction *pb.Instruction
	impl := newInstructionImpl(spec)
	_, err := s.c.Init(impl, impl.function, impl.args)
	if err != nil {
		responseInstruction = ReturnErrorInstruction(i, err)
	} else {
		responseInstruction = ReturnOKInstruction(i)
	}

	response <- responseInstruction
	return
}
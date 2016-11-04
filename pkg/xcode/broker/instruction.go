package broker

import (
	pb "github.com/1851616111/xchain/pkg/protos"
	"time"
)

type instructionImpl struct {
	txID     string
	function string
	args     []string
}

func newInstructionImpl(init *pb.Instruction_Init) *instructionImpl {
	return &instructionImpl{
		txID:     init.Init.TransactionID,
		function: init.Init.Function,
		args:     init.Init.Args,
	}
}

func (i *instructionImpl) GetStringArgs() []string {
	return i.args
}

func (i *instructionImpl) GetTxID() string {
	return i.txID
}

func (i *instructionImpl) GetState(key string) ([]byte, error) {
	return genGetStateInstruction(key)
}

func (i *instructionImpl) PutState(key string, value []byte) error {
	return genPutStateInstruction(key, value)
}

func (i *instructionImpl) DelState(key string) error {
	return genDelStateInstruction(key)
}

func execute(i *pb.Instruction) (*pb.Instruction, error) {
	//无须加锁的原因是Identifier是递增的。
	//chan length 1 is engough, in case of block
	cs.IDCh[i.Identifier] = make(chan *pb.Instruction, 5)
	cs.sendCh <- i

	select {
	case <-time.Tick(time.Second * 30):
		delete(cs.IDCh, i.Identifier)
		return nil, ERRWaitReturnInstructionTimeOut

	case returnInstruction := <-cs.IDCh[i.Identifier]:
		return returnInstruction, nil
	}

}

func newReturnInstruction(i *pb.Instruction) *pb.Instruction {
	if i == nil {
		return nil
	}

	return &pb.Instruction{
		Action:     pb.Action_Response,
		Type:       i.Type,
		Identifier: i.Identifier,
	}
}

func ReturnErrorInstruction(i *pb.Instruction, err error) *pb.Instruction {
	res := newReturnInstruction(i)
	res.Type = pb.Instruction_ERROR
	res.Payload = []byte(err.Error())
	return res
}

func ReturnOKInstruction(i *pb.Instruction) *pb.Instruction {
	res := newReturnInstruction(i)
	res.Payload = []byte("ok")
	return res
}

func IsOKInstruction(i *pb.Instruction) bool {
	return string(i.Payload) == "ok"
}

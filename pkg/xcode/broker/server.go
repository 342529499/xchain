package broker

import (
	pb "github.com/1851616111/xchain/pkg/protos"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"sync"
)

var (
	logger = log.New(os.Stdout, "[xcode]", log.LstdFlags)
	cs     *coderService
)

//心跳来维持broker与xcode server的连接。
//当没有连接的时长秒数超过KEY_STAND_ALONE_TTL时，认为服务已断开，可以停止服务。
//当没有连接的时长秒数KEY_STAND_ALONE_TTL==""(空值)时，认为即使服务断开通信也不推出服务。
type coderService struct {
	sync.RWMutex // 防止被网络轰炸导致阻塞或panic
	started      bool
	startCh      chan struct{}
	nodeName     string
	nodeAddress  string
	language     string
	c            XCoder

	sendCh chan *pb.Instruction

	//由于使用了identifier机制，所以可以并发的get/put state。
	//首先区分txid， 其次区分idetifier
	IDCh map[int64]chan *pb.Instruction

	//增加ID TIMEOUT来防止没有返回的chan没有手动删除而导致oom
}

func (s *coderService) Execute(codeStream pb.CodeService_ExecuteServer) error {

	return s.handle(codeStream)
}

func newCoderServiceServer(coder XCoder) *coderService {
	if coder == nil {
		log.Fatal("coder interface not exist.")
	}

	cs = &coderService{
		startCh: make(chan struct{}, 1),
		c:       coder,
		IDCh:    map[int64]chan *pb.Instruction{},
		sendCh:  make(chan *pb.Instruction, 1000),
	}

	return cs
}

func StartCoderService(coder XCoder) error {
	lis, err := net.Listen("tcp", getListenerAddress())
	if err != nil {
		logger.Printf("new grpc server listen address err %v\n", err)
		return err
	}

	server := grpc.NewServer()

	service := newCoderServiceServer(coder)
	pb.RegisterCodeServiceServer(server, service)

	go service.watchStart()

	if err := server.Serve(lis); err != nil {
		return err
	}

	return nil
}

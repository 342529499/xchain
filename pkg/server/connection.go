package server

import (
	"log"
	"os"

	pb "github.com/1851616111/xchain/pkg/protos"
	"google.golang.org/grpc"

	"code.google.com/p/go.net/context"
	"sync"
)

var (
	logger                         = log.New(os.Stderr, "connection manager", log.LstdFlags)
	consManager *connectionsManager = newConnectionsManager()
)

type pair struct {
	littler pb.EndPoint
	bigger  pb.EndPoint
}

type connectionsManager struct {
	locker        sync.RWMutex
	localEndPoint *pb.EndPoint
	m             map[pair]Connection
}

func newConnectionsManager() *connectionsManager {
	manager := new(connectionsManager)
	manager.m = map[pair]Connection{}
	return manager
}


type Connection interface {
	Send(*pb.Message) error
	Recv() (*pb.Message, error)
}

func SetLocalEndPoint(localEndPoint *pb.EndPoint) {
	consManager.localEndPoint = localEndPoint
}

func GetLocalEndPoint() *pb.EndPoint {
	return consManager.localEndPoint
}

func GetConnectionsManager() *connectionsManager {
	return consManager
}

func ServerAddConnection(pair pair, con Connection) {
	consManager.locker.Lock()
	defer consManager.locker.Unlock()

	consManager.m[pair] = con
}

func PrintConnectionManager() {
	logger.Printf("%#v\n", *consManager)
}

func (m *connectionsManager) Join(targetNet string) error {
	m.locker.Lock()
	defer m.locker.Unlock()
	conn, err := grpc.Dial(targetNet, grpc.WithInsecure())
	if err != nil {
		logger.Printf("join net(%s) dial err: %v.", targetNet, err)
		return err
	}

	cli := pb.NewNetClient(conn)

	stream, err := cli.Connect(context.Background())
	if err != nil {
		logger.Printf("join net(%s) connect err: %v.", targetNet, err)
	}

	agreement := newHandshakeAgreement(*m.localEndPoint, targetNet)
	if err := agreement.handlerJoin(stream); err != nil {
		return err
	}

	return nil
}

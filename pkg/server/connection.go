package server

import (
	"log"
	"os"

	pb "github.com/1851616111/xchain/pkg/protos"
	"google.golang.org/grpc"

	"code.google.com/p/go.net/context"
	"sync"
	"time"
)

const (
	default_Keepalive_TimeDuration = time.Second * 5
	default_Keepalive_FailTime = 3
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
	simpleM       map[string]Connection

	keepaliveFailTime uint8
	keepaliveDuration time.Duration
}

func newConnectionsManager() *connectionsManager {
	manager := new(connectionsManager)
	manager.keepaliveDuration = default_Keepalive_TimeDuration
	manager.m = map[pair]Connection{}
	manager.simpleM = map[string]Connection{}
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
	logger.Println("total connections num  is %d\n", len(consManager.simpleM))
	return consManager
}

func ServerAddConnection(pair pair, con Connection) {
	consManager.m[pair] = con
	consManager.simpleM[pair.bigger.Address] = con
}

func GetConnection(pair pair) Connection {
	return consManager.m[pair]
}
func PrintConnectionManager() {
	logger.Printf("%#v\n", *consManager)
}

//客户端主动加入
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

	m.KeepaliveNewConnection(stream)

	return nil
}

func (m *connectionsManager) KeepaliveNewConnection(con Connection) (err error) {
	go func(){
		var counter uint8 = 0
		ping := makeKeepaliveMsg()
		for ;counter < m.keepaliveFailTime; {
			if err = con.Send(ping); err != nil {
				counter ++
			}
		}

		return
	}()

	return
}
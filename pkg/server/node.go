package server

import (
	pb "github.com/1851616111/xchain/pkg/protos"
	cm "github.com/1851616111/xchain/pkg/server/message_cons_mng"
	"sync"
	"time"
)

var (
	singleton sync.Once
	node      *Node

	//开发环境时为20秒
	develop_Ping_Duration time.Duration = time.Second * 60
	done                  struct{}      = struct{}{}
)

func newNode(local pb.EndPoint) *Node {
	singleton.Do(func() {
		node = &Node{
			localEndPoint:  local,
			epManager:      newEndPointManager(),
			netManager:     newNetManager(),
			keyExistMarker: map[string]interface{}{},

			recvConnectCh: make(chan *recvConnetionMetadata, 30),

			lounchClientCh: make(chan *lounchConnectionMetadata, 30),

			pingDuration: develop_Ping_Duration,

			pingMaxFailureTimes: map[string]uint{},

			runningClientDoneCH: map[string]chan struct{}{},


		}
	})

	return node
}

func getNode() *Node {
	return node
}

//先假设netManager的server端是稳定的。
type Node struct {
	sync.RWMutex

	localEndPoint pb.EndPoint

	epManager *EndPointManager

	keyExistMarker map[string]interface{}

	netManager *NetManager

	//server成功与client握手后的通知，用于定时任务，及信息记录
	recvConnectCh chan *recvConnetionMetadata

	//发起新的客户端接入请求
	lounchClientCh chan *lounchConnectionMetadata

	//ping的间隔，定期取回其他节点的node List 信息
	pingDuration time.Duration

	pingMaxFailureTimes map[string]uint

	runningClientDoneCH map[string]chan struct{}

}

type EndPointManager struct {

	//验证节点ID List
	ValidatorList []string
	//非验证节点ID List
	NonValidateList []string

	IDToAddress map[string]string
	AddressToID map[string]string
}

//网络节点接受某台node加入的接口
func (n *Node) Accept(ep pb.EndPoint, con cm.Connection) error {
	n.Lock()
	defer n.Unlock()

	key := ep.Id

	if err := n.netManager.serverAdd(key, con); err != nil {
		return err
	}
	n.epManager.addEndPoint(ep)

	return nil
}

//加入网络（某个node）的接口
func (n *Node) Connect(ep pb.EndPoint, con cm.Connection) error {
	n.Lock()
	defer n.Unlock()

	key := ep.Id

	if err := n.netManager.clientAdd(key, con); err != nil {
		return err
	}
	n.epManager.addEndPoint(ep)

	return nil
}

//取消节点在本节点的加入
func (n *Node) CancelAccept(ep pb.EndPoint) {
	n.Lock()
	defer n.Unlock()

	key := ep.Id

	n.epManager.delEndPoint(ep.Id)
	n.netManager.delete(key)
	return
}

//断开与某台机器的连接
func (n *Node) DisConnect(ep pb.EndPoint, con cm.Connection) {
	n.Lock()
	defer n.Unlock()

	n.epManager.delEndPoint(ep.Id)
	n.netManager.delete(ep.Id)
}

func (n *Node) Exist(address string) bool {
	n.RLock()
	defer n.RUnlock()

	_, exist := n.epManager.AddressToID[address]
	return exist
}

func (n *Node) GetLocalEndPoint() *pb.EndPoint {
	return &n.localEndPoint
}

func (n *Node) RemoteClient(id string) {
	n.Lock()
	defer n.Unlock()

	if _, exists := n.pingMaxFailureTimes[id]; !exists {
		n.pingMaxFailureTimes[id] = 1
	}

	n.pingMaxFailureTimes[id]++

	if n.pingMaxFailureTimes[id] == 3 {
		n.netManager.delete(id)
		n.epManager.delEndPoint(id)
		n.runningClientDoneCH[id] <- done
	}
}

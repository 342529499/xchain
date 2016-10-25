package server

import (
	pb "github.com/1851616111/xchain/pkg/protos"
	cm "github.com/1851616111/xchain/pkg/server/connection_manager"
	"sync"
	"time"
)

var (
	singleton sync.Once
	node      *Node

	//开发环境时为20秒
	develop_Ping_Duration time.Duration = time.Second * 20
)

func newNode(local pb.EndPoint) *Node {
	singleton.Do(func() {
		node = &Node{
			localEndPoint:  local,
			epManager:      newEndPointManager(),
			netManager:     newNetManager(),
			keyExistMarker: map[string]interface{}{},

			recvConnectCh: make(chan *recvConnetionMetadata, 30),

			lounchConnectCh: make(chan *lounchConnectionMetadata, 30),

			pingDuration: develop_Ping_Duration,
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
	lounchConnectCh chan *lounchConnectionMetadata

	pingDuration time.Duration
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

	n.epManager.delEndPoint(ep)
	n.netManager.delete(key)
	return
}

//断开与某台机器的连接
func (n *Node) DisConnect(ep pb.EndPoint, con cm.Connection) {
	n.Lock()
	defer n.Unlock()

	key := ep.Id

	n.epManager.delEndPoint(ep)
	n.netManager.delete(key)
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

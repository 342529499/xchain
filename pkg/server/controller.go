package server

import (
	pb "github.com/1851616111/xchain/pkg/protos"
	cm "github.com/1851616111/xchain/pkg/server/connection_manager"
	"log"
	"os"
)

var logger = log.New(os.Stderr, "controller:", log.LstdFlags)

func (n *Node) RunController() {
	for {
		select {
		case rc := <-n.recvConnectCh:
			logger.Printf("node controller: recv connection %v\n", rc)

			if err := n.Accept(rc.client, rc.con); err != nil {
				rc.errCh <- err
			} else {
				rc.doneCh <- struct{}{}
			}

		case task := <-n.lounchConnectCh:

			var successFn = func(target pb.EndPoint, con cm.Connection) error {
				return n.Connect(target, con)
			}
			//接收到一个作为客户端发起连接的tash时
			//先调用实际的握手handle流程，当握手成功后
			//通过successfn回调n.Connect() 将信息加入到node节点上
			//并将task的结果返回给调用者
			if err := n.startAndJoin(task.targetAddress, successFn); err != nil {
				log.Printf("node controller: lounch connection err:%v\n", err)
				task.errCh <- err
			}
			task.doneCh <- struct{}{}
			continue

			logger.Printf("node controller: success launch connection for %s\n", task.targetAddress)
		}
	}
}

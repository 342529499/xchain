package server

import (
	pb "github.com/1851616111/xchain/pkg/protos"
	cm "github.com/1851616111/xchain/pkg/server/message_cons_mng"
	"log"
	"os"
	"time"
)

var (
	logger = log.New(os.Stderr, "[controller]", log.LstdFlags)
)

func (n *Node) RunController() {

	successFunc := func(target pb.EndPoint, con cm.Connection) error {
		if err := n.Connect(target, con); err != nil {
			return err
		}

		done := make(chan struct{}, 1)
		n.runningClientDoneCH[target.Id] = done
		go clientConnectionHandler(target, con, done)

		return nil
	}

	pingFunc := func(id string, con cm.Connection) error {
		if err := con.Send(makePingReqMsg()); err != nil {
			logger.Printf("broadcast ping id:%s err %v\n", id, err)
			n.RemoteClient(id)
			return err
		}

		return nil
	}

	for {
		select {
		case rc := <-n.recvConnectCh:
			logger.Printf("node controller: recv connection %v\n", rc)

			if err := n.Accept(rc.client, rc.con); err != nil {
				rc.errCh <- err
			} else {
				rc.doneCh <- done
			}

		case task := <-n.lounchClientCh:

			logger.Printf("connect to entry point %s\n", task.targetAddress)
			//接收到一个作为客户端发起连接的tash时
			//先调用实际的握手handle流程，当握手成功后
			//通过successfn回调n.Connect() 将信息加入到node节点上
			//并将task的结果返回给调用者
			if err := n.startAndJoin(task.targetAddress, successFunc); err != nil {
				log.Printf("node controller: lounch connection err:%v\n", err)
				task.errCh <- err
			}

			task.doneCh <- struct{}{}
			continue
			logger.Printf("node controller: success launch connection for %s\n", task.targetAddress)

		case <-time.Tick(n.pingDuration):
			if err := n.netManager.BroadcastFunc(true, pingFunc); err != nil {
				logger.Printf("broadcast ping err %v\n", err)
			}
		}
	}
}

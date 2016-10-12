package service

import (
	pingpb "github.com/1851616111/xchain/pkg/protos"
	ttlcache "github.com/dadgar/onecache/ttlstore"
	//"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"

	"fmt"
	"log"
	"os"
	"time"
)

const (
	ExpireSec_Service_Ping_Nodes = time.Second * 120
)

var (
	logger        = log.New(os.Stderr, "ping", log.LstdFlags)
	cache, _      = ttlcache.New(50, logger)
	Ping     ping = ping{
		cache: cache,
		encoder: func(obj *Node) ([]byte, error) {
			if len(obj.Address) == 0 {
				return nil, fmt.Errorf("cache node(%v) msg, param address is nil.", obj)
			}
			return []byte(obj.Address), nil
		},
	}
)

type ping struct {
	cache *ttlcache.DataStore

	encoder func(obj *Node) ([]byte, error)
}

func (s *Server) Ping(c context.Context, msg *pingpb.Message) (*pingpb.Message, error) {
	//
	//node := &Node{}
	//if err := proto.Unmarshal(msg.Payload, node); err != nil {
	//	return nil, err
	//}
	//
	////本peer服务第一个启动，有新的peer ping本peer
	////
	//if len(ping.cache.List()) == 0 {
	//
	//}
	//
	//ping.cache.Set(node.Id, ping.encoder(node), ExpireSec_Service_Ping_Nodes, 0)

	return nil, nil
}

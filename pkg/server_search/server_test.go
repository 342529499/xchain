package server
//
//import (
//	"net"
//	"testing"
//
//	"fmt"
//	pb "github.com/1851616111/xchain/pkg/protos"
//	"time"
//)
//
////func TestGetConnectionsManager(t *testing.T) {
////	SetLocalEndPoint(&pb.EndPoint{
////		Id:"michael",
////		Address:GetLocalIP(),
////	})
////
////	go func() {
////		time.Sleep(time.Second* 5)
////		TestGetConnection(t)
////	}()
////
////	if err := NewAndStartGrpcServer(&ServerOptions{
////		Address: "0.0.0.0:10690",
////	}); err != nil {
////		t.Fatalf("new and start default grpc server err %v", err)
////	}
////
////}
//
//func TestGetConnectionsManager(t *testing.T) {
//
//	SetLocalEndPoint(&pb.EndPoint{
//		Id:      "jessie",
//		Address: GetLocalIP(),
//	})
//
//	manager := GetConnectionsManager()
//
//	err := manager.Join("192.168.0.110:10690")
//	if err != nil {
//		println(err.Error())
//	}
//
//	time.Sleep(time.Second * 3)
//
//	for _, v := range manager.m {
//
//		for {
//
//			fmt.Println("-----------------")
//			time.Sleep(time.Second)
//			v.Send(&pb.Message{Payload: []byte{1, 1, 1}})
//			msg, err := v.Recv()
//			if err != nil {
//				fmt.Println(err)
//			}
//
//			fmt.Printf("---->%s\n", msg)
//		}
//
//	}
//
//}
//
//func tGetConnection(t *testing.T) {
//	manager := GetConnectionsManager()
//	for k, v := range manager.simpleM {
//		err := v.Send(&pb.Message{Payload: []byte(k)})
//		if err != nil {
//			fmt.Println(err)
//		}
//		fmt.Printf("send %v\n", k)
//	}
//}
//
////// GetLocalIP returns the non loopback local IP of the host
//func GetLocalIP() string {
//	addrs, err := net.InterfaceAddrs()
//	if err != nil {
//		return ""
//	}
//	for _, address := range addrs {
//		// check the address type and if it is not a loopback then display it
//		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
//			if ipnet.IP.To4() != nil {
//				return ipnet.IP.String()
//			}
//		}
//	}
//	return ""
//}

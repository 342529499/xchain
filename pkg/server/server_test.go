package server

import (
	"net"
	"testing"

	pb "github.com/1851616111/xchain/pkg/protos"
)

func TestGetConnectionsManager(t *testing.T) {
	SetLocalEndPoint(&pb.EndPoint{
		Id:"michael",
		Address:GetLocalIP(),
	})


		if err := NewAndStartGrpcServer(&ServerOptions{
			Address: "0.0.0.0:10690",
		}); err != nil {
			t.Fatalf("new and start default grpc server err %v", err)
		}
	//}()

	//



}
//
//func TestunLimitSendMsg(t *testing.T) {
//
//	SetLocalEndPoint(&pb.EndPoint{
//		Id:"michael",
//		Address:GetLocalIP(),
//	})
//
//
//	manager := GetConnectionsManager()
//
//	err := manager.Join("192.168.174.132:10690")
//	if err != nil {
//		t.Error(err)
//	}
//
//	for _, v := range manager.m {
//		v.Send(&pb.Message{Payload:[]byte{"1","2","3"}})
//	}
//
//}
//// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback then display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
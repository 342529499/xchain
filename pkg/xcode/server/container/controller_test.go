package container

import (
	"fmt"
	pb "github.com/1851616111/xchain/pkg/protos"
	"testing"
	"time"
)

func TestController_Start(t *testing.T) {
	ctl := GetController()
	ctl.ping = time.Second * 5
	ctl.Run()

	//testController_Dispatch_Deploy_Localhost(ctl, t)
	//testController_Dispatch_Deploy_Http(ctl, t)
	testController_DeployValidate(ctl, t)
	testController_Deploy(ctl, t)
	select {}
}

func testController_Deploy(ctl *Controller, t *testing.T) {
	spec := &pb.XCodeSpec{
		Type: pb.XCodeSpec_GOLANG,
		XcodeID: &pb.XCodeID{
			Path: "github.com/1851616111/xchain/example/example01",
		},

		XcodeMsg: &pb.XCodeInput{
			Args: ToXCodeArgs("f"),
		},
	}
	fmt.Println("------01")
	err := ctl.Deploy(spec)
	if err != nil {
		t.Fatalf("deploy err %v\n", err)
	}
	fmt.Println("------02")
}

func testController_DeployValidate(ctl *Controller, t *testing.T) {
	spec := &pb.XCodeSpec{
		Type: pb.XCodeSpec_GOLANG,
		XcodeID: &pb.XCodeID{
			Path: "github.com/1851616111/xchain/example/example01",
		},

		XcodeMsg: &pb.XCodeInput{
			Args: ToXCodeArgs("f"),
		},
	}
	err := ctl.DeployValidate(spec)
	if err != nil && err != ErrDeployWorkDuplicated {
		t.Fatalf("deploy validate err %v\n", err)
	}
}

//func testController_Dispatch_Deploy_Localhost(ctl *controller, t *testing.T) {
//
//	spec := &pb.XCodeSpec{
//		Type: pb.XCodeSpec_GOLANG,
//		XcodeID: &pb.XCodeID{
//			Path: "github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example01",
//		},
//
//		XcodeMsg: &pb.XCodeInput{
//			Args: ToXCodeArgs("f"),
//		},
//	}
//	work := &Worker{
//		act:  Job_Action_BuildImage,
//		id:   spec.XcodeID.Path,
//		lang: Job_Language_Go,
//		metadata: &WorkSpec{
//			XCodeSpec: spec,
//			PeerID:    "test",
//		},
//
//		opts: &docker.BuildImageOptions{
//			Name: "123",
//		},
//		resultCh: make(chan interface{}),
//	}
//
//	ctl.Dispatch(work)
//}
//
//func testController_Dispatch_Deploy_Http(ctl *controller, t *testing.T) {
//
//	spec := &pb.XCodeSpec{
//		Type: pb.XCodeSpec_GOLANG,
//		XcodeID: &pb.XCodeID{
//			Path: "https://github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example02",
//		},
//
//		XcodeMsg: &pb.XCodeInput{
//			Args: ToXCodeArgs("f"),
//		},
//	}
//	work := &Worker{
//		act:  Job_Action_BuildImage,
//		id:   spec.XcodeID.Path,
//		lang: Job_Language_Go,
//		metadata: &WorkSpec{
//			XCodeSpec: spec,
//			PeerID:    "test",
//		},
//
//		opts: &docker.BuildImageOptions{
//			Name: "789",
//		},
//		resultCh: make(chan interface{}),
//	}
//
//	ctl.Dispatch(work)
//}

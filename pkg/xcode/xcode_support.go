package xcode

import "fmt"



import (
	"fmt"
	pb "github.com/1851616111/xchain/protos"
	"io"
	"time"
)

type XChainName string

const (
	XChain_DefaultName           XChainName    = "default"
	XChain_DefaultStartupTimeout time.Duration = time.Millisecond * 5000
	XChain_DefaultInstallPath    string        = "/opt/gopath/bin/"
	XChain_DafaultPeerAddress    string        = "0.0.0.0:7051"
)

//type XCodeSupportConfig struct {
//	xcodeStartUpTimeout time.Duration
//	//xcode               (j)
//}

//func NewXCode_SupportDefaultConfig()

type XCodeWorker struct {
}

func (worker *XCodeWorker) Load(stream pb.XCode_LoadServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		fmt.Printf("recive msg %s\n", in)
	}
}

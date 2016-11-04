package container

import (
	pb "github.com/1851616111/xchain/pkg/protos"
)

type Interface interface {
	Deploy(spec *pb.XCodeSpec) error
	//Start(ctxt context.Context, ccid ccintf.CCID, args []string, env []string, attachstdin bool, attachstdout bool, reader io.Reader) error
	//Stop(ctxt context.Context, ccid ccintf.CCID, timeout uint, dontkill bool, dontremove bool) error
	//Destroy(ctxt context.Context, ccid ccintf.CCID, force bool, noprune bool) error
}

type WorkSpec struct {
	XCodeSpec *pb.XCodeSpec
	PeerID    string
}

package container

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	pb "github.com/1851616111/xchain/pkg/protos"
	"github.com/1851616111/xchain/pkg/xcode/golang"
	"io"
)

func GetXCodePackageBytes(spec *pb.XCodeSpec) (io.Reader, error) {
	if spec == nil || spec.XcodeID == nil {
		return nil, fmt.Errorf("invalid chaincode spec")
	}

	inputbuf := bytes.NewBuffer(nil)
	gw := gzip.NewWriter(inputbuf)
	tw := tar.NewWriter(gw)

	goLang := &golang.Platform{}

	err := goLang.WritePackage(spec, tw)
	if err != nil {
		return nil, err
	}

	tw.Close()
	gw.Close()

	if err != nil {
		return nil, err
	}

	return inputbuf, nil
}

func ToXCodeArgs(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}

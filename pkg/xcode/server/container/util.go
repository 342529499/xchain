package container

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	pb "github.com/1851616111/xchain/pkg/protos"
	"github.com/1851616111/xchain/pkg/xcode/server/golang"
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

func genContainerName(spec *pb.XCodeSpec) string {
	var paramStr string = "no"
	if len(spec.XcodeMsg.Args) > 0 {
		m := md5.New()

		for _, b := range spec.XcodeMsg.Args {
			m.Write(b)
		}
		s := m.Sum(nil)
		paramStr = hex.EncodeToString(s)
	}

	return fmt.Sprintf("XCODE-%s-%s-%s", spec.Type.String(), spec.XcodeID.Path, paramStr)
}

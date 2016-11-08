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
	//"github.com/1851616111/go-dockerclient"
	"github.com/1851616111/go-dockerclient"
	"strings"
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

func genCodeID(spec *pb.XCodeSpec) string {
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

func genDockerID(CodeID string) string {
	s := strings.Replace(CodeID, ":", "_", -1)
	return strings.Replace(s, "/", "_", -1)

}

func convertLabelToFilter(targetFilter *map[string][]string, label map[string]string) {
	if len(label) == 0 {
		return
	}

	if *targetFilter == nil {
		*targetFilter = map[string][]string{}
	}

	labelFilterSlice := []string{}
	for k, v := range label {
		labelFilterSlice = append(labelFilterSlice, fmt.Sprintf("%s=%s", k, v))
	}


	(*targetFilter)["label"] = labelFilterSlice
}

func getDockerHostConfig() *docker.HostConfig {
	return &docker.HostConfig{
		Memory: 2147483648,
	}
}

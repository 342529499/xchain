/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package golang

import (
	"archive/tar"
	"fmt"
	pb "github.com/1851616111/xchain/pkg/protos"
	"strings"
	"time"
)

//tw is expected to have the chaincode in it from GenerateHashcode. This method
//will just package rest of the bytes
func writeChaincodePackage(spec *pb.XCodeSpec, tw *tar.Writer) error {

	var urlLocation string
	if strings.HasPrefix(spec.XcodeID.Path, "http://") {
		urlLocation = spec.XcodeID.Path[7:]
	} else if strings.HasPrefix(spec.XcodeID.Path, "https://") {
		urlLocation = spec.XcodeID.Path[8:]
	} else {
		urlLocation = spec.XcodeID.Path
	}

	if urlLocation == "" {
		return fmt.Errorf("empty url location")
	}

	if strings.LastIndex(urlLocation, "/") == len(urlLocation)-1 {
		urlLocation = urlLocation[:len(urlLocation)-1]
	}
	toks := strings.Split(urlLocation, "/")
	if toks == nil || len(toks) == 0 {
		return fmt.Errorf("cannot get path components from %s", urlLocation)
	}

	chaincodeGoName := toks[len(toks)-1]
	if chaincodeGoName == "" {
		return fmt.Errorf("could not get chaincode name from path %s", urlLocation)
	}

	//let the executable's name be chaincode ID's name
	//newRunLine := fmt.Sprintf("RUN go install %s && mv $GOPATH/bin/%s $GOPATH/bin/%s", urlLocation, chaincodeGoName, spec.XcodeID.Name)
	dockerFile := fmt.Sprintf(`from hyperledger/fabric-baseimage:x86_64-0.1.0
	#from utxo:0.1.0
	COPY src $GOPATH/src
	WORKDIR $GOPATH/src/%s
	RUN go build
	ENTRYPOINT ["./%s"]`, urlLocation, chaincodeGoName)

	dockerFileSize := int64(len([]byte(dockerFile)))

	//Make headers identical by using zero time
	var zeroTime time.Time
	tw.WriteHeader(&tar.Header{Name: "Dockerfile", Size: dockerFileSize, ModTime: zeroTime, AccessTime: zeroTime, ChangeTime: zeroTime})
	tw.Write([]byte(dockerFile))
	err := WriteGopathSrc(tw, urlLocation)
	if err != nil {
		return fmt.Errorf("Error writing Chaincode package contents: %s", err)
	}
	return nil
}

package main

import (
	pb "github.com/1851616111/xchain/pkg/protos"
	"github.com/golang/protobuf/proto"
	"fmt"
)



func main() {

	var ep *pb.EndPoint = &pb.EndPoint{
		Id:"miasdfasdfasdf",
		Address:"456sdfasdfasdfasdf",
	}

	b, err := proto.Marshal(ep)
	if err != nil {
		fmt.Println(err)
	}

	newEP := &pb.EndPoint{
	}

	proto.Unmarshal(b, newEP)

	fmt.Printf("%v\n", newEP)


	var hs *pb.HandShake = &pb.HandShake{
		EndPoint: ep,
	}



	b, err = proto.Marshal(hs)
	if err != nil {
		fmt.Println(err)
	}

	newhs := &pb.HandShake{}
	proto.Unmarshal(b, newhs)
	fmt.Printf("%v\n", newhs)

}
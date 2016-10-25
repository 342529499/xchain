package server

import (
	pb "github.com/1851616111/xchain/pkg/protos"
	sliceutil "github.com/1851616111/xchain/pkg/util/slice"
	"log"
	"os"
)

var (
	endPointLog = log.New(os.Stderr, "[endpoint]", log.LstdFlags)
)

func newEndPointManager() *EndPointManager {
	m := new(EndPointManager)
	m.IDToAddress = map[string]string{}
	m.AddressToID = map[string]string{}
	return m
}

func (m *EndPointManager) findNewEndPointHandler(epList []*pb.EndPoint, handler func(*pb.EndPoint)) {
	for _, newEP := range epList {
		if !m.ifExistEndPoint(newEP) {
			handler(newEP)
		}
	}
}

func (m *EndPointManager) ifExistEndPoint(epList *pb.EndPoint) bool {
	_, exist := m.IDToAddress[epList.Id]
	return exist
}

func (m *EndPointManager) addEndPoint(ep pb.EndPoint) {
	var key string = ep.Id

	switch ep.Type {
	case pb.EndPoint_VALIDATOR:
		m.ValidatorList = append(m.ValidatorList, ep.Id)
	case pb.EndPoint_NON_VALIDATOR:
		m.NonValidateList = append(m.NonValidateList, ep.Id)
	}

	m.IDToAddress[key] = ep.Address
	m.AddressToID[ep.Address] = key
}

func (m *EndPointManager) delEndPoint(ep pb.EndPoint) {
	var key string = ep.Id

	switch ep.Type {
	case pb.EndPoint_VALIDATOR:
		m.ValidatorList = sliceutil.RemoveSliceElement(m.ValidatorList, key)
	case pb.EndPoint_NON_VALIDATOR:
		m.ValidatorList = sliceutil.RemoveSliceElement(m.NonValidateList, key)
	}

	delete(m.IDToAddress, key)
	delete(m.AddressToID, ep.Address)
}

func (m *EndPointManager) list() []*pb.EndPoint {

	if Is_Develop_Mod {
		endPointLog.Printf("validate:%v\n", m.ValidatorList)
		endPointLog.Printf("non-validate:%v\n", m.NonValidateList)
	}

	validateEPs, nonValidateEPs := []*pb.EndPoint{}, []*pb.EndPoint{}
	rangeValidateFunc := func(idx int, id string) error {
		validateEPs = append(validateEPs, &pb.EndPoint{
			Id:      id,
			Address: m.IDToAddress[id],
			Type:    pb.EndPoint_VALIDATOR,
		})
		return nil
	}

	exec := true
	sliceutil.RangeSlice(m.ValidatorList, &exec, rangeValidateFunc)

	rangeNonValidateFunc := func(idx int, id string) error {
		nonValidateEPs = append(nonValidateEPs, &pb.EndPoint{
			Id:      id,
			Address: m.IDToAddress[id],
			Type:    pb.EndPoint_NON_VALIDATOR,
		})
		return nil
	}

	sliceutil.RangeSlice(m.NonValidateList, &exec, rangeNonValidateFunc)

	return append(validateEPs, nonValidateEPs...)
}

func ListWithLocalEP(l []*pb.EndPoint, local *pb.EndPoint) []*pb.EndPoint {
	return append(l, local)
}

func ListWithOutLocalEP(l []*pb.EndPoint, local *pb.EndPoint) []*pb.EndPoint {
	if len(l) == 0 {
		return l
	}

	for idx, ep := range l {
		if ep == local {
			return append(l[:idx], l[idx+1])
		}
	}

	return l
}

func printEPList(l []*pb.EndPoint) {

	if Is_Develop_Mod {
		endPointLog.Println("endpoints list:")
		for _, v := range l {
			endPointLog.Printf("{id:\"%s\",address:\"%s\",type\":\"%d\"}\n", v.Id, v.Address, v.Type)
		}
	}

}

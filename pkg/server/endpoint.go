package server

import (
	pb "github.com/1851616111/xchain/pkg/protos"
	sliceutil "github.com/1851616111/xchain/pkg/util/slice"
	"log"
	"os"
)

var (
	endPointLog = log.New(os.Stderr, "[Event]", log.LstdFlags)
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

func (m *EndPointManager) delEndPoint(delID string) {
	address, exist := m.IDToAddress[delID]
	if !exist {
		return
	}

	for idx, id := range m.ValidatorList {
		if id == delID {
			m.ValidatorList = append(m.ValidatorList[:idx], m.ValidatorList[idx+1:]...)
		}
	}

	for idx, id := range m.NonValidateList {
		if id == delID {
			m.NonValidateList = append(m.NonValidateList[:idx], m.NonValidateList[idx+1:]...)
		}
	}

	delete(m.IDToAddress, delID)
	delete(m.AddressToID, address)
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

func ListWithOutLocalEP(l []*pb.EndPoint) []*pb.EndPoint {
	if len(l) == 0 {
		return l
	}

	local := node.GetLocalEndPoint()

	for idx, ep := range l {
		if ep.Id == local.Id {
			return append(l[:idx], l[idx+1:]...)
		}
	}

	return l
}

// Code generated by protoc-gen-go.
// source: message.proto
// DO NOT EDIT!

/*
Package protos is a generated protocol buffer package.

It is generated from these files:
	message.proto
	type.proto

It has these top-level messages:
	Message
	EndPoint
	HandShake
	HandShakeResponse
*/
package protos

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/timestamp"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Action int32

const (
	Action_Request  Action = 0
	Action_Response Action = 1
)

var Action_name = map[int32]string{
	0: "Request",
	1: "Response",
}
var Action_value = map[string]int32{
	"Request":  0,
	"Response": 1,
}

func (x Action) String() string {
	return proto.EnumName(Action_name, int32(x))
}
func (Action) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Message_Type int32

const (
	Message_UNDEFINED           Message_Type = 0
	Message_Net_HANDSHAKE       Message_Type = 1
	Message_Net_PING            Message_Type = 2
	Message_Contract_Create     Message_Type = 3
	Message_Contract_Run        Message_Type = 4
	Message_Contract_Query      Message_Type = 5
	Message_Code_State_Put      Message_Type = 6
	Message_Code_State_Get      Message_Type = 7
	Message_Ledger_Block_Create Message_Type = 8
	Message_Ledger_Block_Sync   Message_Type = 9
	Message_Identity_Share      Message_Type = 10
)

var Message_Type_name = map[int32]string{
	0:  "UNDEFINED",
	1:  "Net_HANDSHAKE",
	2:  "Net_PING",
	3:  "Contract_Create",
	4:  "Contract_Run",
	5:  "Contract_Query",
	6:  "Code_State_Put",
	7:  "Code_State_Get",
	8:  "Ledger_Block_Create",
	9:  "Ledger_Block_Sync",
	10: "Identity_Share",
}
var Message_Type_value = map[string]int32{
	"UNDEFINED":           0,
	"Net_HANDSHAKE":       1,
	"Net_PING":            2,
	"Contract_Create":     3,
	"Contract_Run":        4,
	"Contract_Query":      5,
	"Code_State_Put":      6,
	"Code_State_Get":      7,
	"Ledger_Block_Create": 8,
	"Ledger_Block_Sync":   9,
	"Identity_Share":      10,
}

func (x Message_Type) String() string {
	return proto.EnumName(Message_Type_name, int32(x))
}
func (Message_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type Message struct {
	Action    Action                     `protobuf:"varint,1,opt,name=action,enum=protos.Action" json:"action,omitempty"`
	Type      Message_Type               `protobuf:"varint,2,opt,name=type,enum=protos.Message_Type" json:"type,omitempty"`
	Payload   []byte                     `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
	Timestamp *google_protobuf.Timestamp `protobuf:"bytes,4,opt,name=timestamp" json:"timestamp,omitempty"`
}

func (m *Message) Reset()                    { *m = Message{} }
func (m *Message) String() string            { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()               {}
func (*Message) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Message) GetTimestamp() *google_protobuf.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func init() {
	proto.RegisterType((*Message)(nil), "protos.Message")
	proto.RegisterEnum("protos.Action", Action_name, Action_value)
	proto.RegisterEnum("protos.Message_Type", Message_Type_name, Message_Type_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for Net service

type NetClient interface {
	Connect(ctx context.Context, opts ...grpc.CallOption) (Net_ConnectClient, error)
}

type netClient struct {
	cc *grpc.ClientConn
}

func NewNetClient(cc *grpc.ClientConn) NetClient {
	return &netClient{cc}
}

func (c *netClient) Connect(ctx context.Context, opts ...grpc.CallOption) (Net_ConnectClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Net_serviceDesc.Streams[0], c.cc, "/protos.net/connect", opts...)
	if err != nil {
		return nil, err
	}
	x := &netConnectClient{stream}
	return x, nil
}

type Net_ConnectClient interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ClientStream
}

type netConnectClient struct {
	grpc.ClientStream
}

func (x *netConnectClient) Send(m *Message) error {
	return x.ClientStream.SendMsg(m)
}

func (x *netConnectClient) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Net service

type NetServer interface {
	Connect(Net_ConnectServer) error
}

func RegisterNetServer(s *grpc.Server, srv NetServer) {
	s.RegisterService(&_Net_serviceDesc, srv)
}

func _Net_Connect_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(NetServer).Connect(&netConnectServer{stream})
}

type Net_ConnectServer interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ServerStream
}

type netConnectServer struct {
	grpc.ServerStream
}

func (x *netConnectServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

func (x *netConnectServer) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Net_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protos.net",
	HandlerType: (*NetServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "connect",
			Handler:       _Net_Connect_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: fileDescriptor0,
}

func init() { proto.RegisterFile("message.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 382 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x5c, 0x92, 0xdf, 0xb2, 0x92, 0x50,
	0x14, 0xc6, 0x0f, 0x07, 0x82, 0xe3, 0x3a, 0xea, 0xd9, 0x67, 0x59, 0x13, 0xe3, 0x4d, 0x8d, 0xcd,
	0x34, 0x4e, 0x17, 0x58, 0x76, 0xe3, 0xad, 0xa9, 0xa9, 0x53, 0x31, 0x06, 0x76, 0xcd, 0x20, 0xac,
	0xc8, 0x49, 0xf7, 0x26, 0xd8, 0x5c, 0xf0, 0x14, 0x3d, 0x61, 0xef, 0xd2, 0x16, 0xc4, 0x26, 0xaf,
	0x98, 0xf5, 0x7d, 0xbf, 0xf5, 0xe7, 0x03, 0xa0, 0x73, 0xa4, 0x3c, 0x0f, 0x13, 0x72, 0xd2, 0x4c,
	0x48, 0x81, 0x66, 0xf5, 0xc8, 0xfb, 0x2f, 0x12, 0x21, 0x92, 0x03, 0x8d, 0xaa, 0x72, 0x57, 0x7c,
	0x1f, 0xc9, 0xbd, 0x02, 0x65, 0x78, 0x4c, 0x6b, 0x70, 0xf0, 0x5b, 0x07, 0xeb, 0x4b, 0xdd, 0x8a,
	0xaf, 0xc1, 0x0c, 0x23, 0xb9, 0x17, 0xdc, 0xd6, 0x5e, 0x6a, 0xc3, 0xee, 0xb8, 0x5b, 0x33, 0xb9,
	0x33, 0xad, 0x54, 0xef, 0xec, 0xe2, 0x10, 0x0c, 0x59, 0xa6, 0x64, 0xdf, 0x56, 0xd4, 0xd3, 0x86,
	0x3a, 0x8f, 0x71, 0xb6, 0xca, 0xf3, 0x2a, 0x02, 0x6d, 0xb0, 0xd2, 0xb0, 0x3c, 0x88, 0x30, 0xb6,
	0x75, 0x05, 0xb7, 0xbd, 0xa6, 0xc4, 0x09, 0xb4, 0x2e, 0xa7, 0xd8, 0x86, 0xf2, 0xee, 0xc7, 0x7d,
	0xa7, 0x3e, 0xd6, 0x69, 0x8e, 0x75, 0xb6, 0x0d, 0xe1, 0xfd, 0x83, 0x07, 0x7f, 0x34, 0x30, 0x4e,
	0x2b, 0xb0, 0x03, 0xad, 0x6f, 0xee, 0x7c, 0xf1, 0x71, 0xed, 0x2e, 0xe6, 0xec, 0x06, 0x1f, 0xa1,
	0xe3, 0x92, 0x0c, 0x56, 0x53, 0x77, 0xee, 0xaf, 0xa6, 0x9f, 0x16, 0x4c, 0xc3, 0x36, 0xdc, 0x9d,
	0xa4, 0xcd, 0xda, 0x5d, 0xb2, 0x5b, 0xec, 0xc1, 0xc3, 0x4c, 0x70, 0x99, 0xa9, 0x14, 0xc1, 0x2c,
	0xa3, 0x50, 0x12, 0xd3, 0x91, 0x41, 0xfb, 0x22, 0x7a, 0x05, 0x67, 0x06, 0x22, 0x74, 0x2f, 0xca,
	0xd7, 0x82, 0xb2, 0x92, 0x3d, 0xa9, 0xb5, 0x98, 0x02, 0x5f, 0xaa, 0xae, 0x60, 0x53, 0x48, 0x66,
	0x5e, 0x69, 0x4b, 0x92, 0xcc, 0xc2, 0xe7, 0xd0, 0xfb, 0x4c, 0x71, 0x42, 0x59, 0xf0, 0xe1, 0x20,
	0xa2, 0x9f, 0xcd, 0x9a, 0x3b, 0x7c, 0x06, 0x8f, 0xff, 0x19, 0x7e, 0xc9, 0x23, 0xd6, 0x3a, 0xcd,
	0x58, 0xc7, 0xc4, 0xe5, 0x5e, 0x96, 0x81, 0xff, 0x23, 0xcc, 0x88, 0xc1, 0x9b, 0x57, 0x60, 0xd6,
	0xef, 0x1b, 0xef, 0xc1, 0xf2, 0xe8, 0x57, 0xa1, 0x72, 0xab, 0x78, 0x2a, 0x8b, 0x47, 0x79, 0x2a,
	0x78, 0x4e, 0x4c, 0x1b, 0x4f, 0x40, 0xe7, 0x24, 0xf1, 0x1d, 0x58, 0x91, 0xe0, 0x9c, 0x22, 0x89,
	0x0f, 0x57, 0x9f, 0xa1, 0x7f, 0x2d, 0x0c, 0x6e, 0x86, 0xda, 0x5b, 0x6d, 0x57, 0xff, 0x19, 0xef,
	0xff, 0x06, 0x00, 0x00, 0xff, 0xff, 0xa2, 0x6a, 0x49, 0x27, 0x31, 0x02, 0x00, 0x00,
}

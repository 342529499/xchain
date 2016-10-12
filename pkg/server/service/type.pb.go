// Code generated by protoc-gen-go.
// source: type.proto
// DO NOT EDIT!

/*
Package service is a generated protocol buffer package.

It is generated from these files:
	type.proto

It has these top-level messages:
	Node
*/
package service

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Type int32

const (
	Type_UNDEFINED     Type = 0
	Type_VALIDATOR     Type = 1
	Type_NON_VALIDATOR Type = 2
)

var Type_name = map[int32]string{
	0: "UNDEFINED",
	1: "VALIDATOR",
	2: "NON_VALIDATOR",
}
var Type_value = map[string]int32{
	"UNDEFINED":     0,
	"VALIDATOR":     1,
	"NON_VALIDATOR": 2,
}

func (x Type) String() string {
	return proto.EnumName(Type_name, int32(x))
}
func (Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Node struct {
	Id      string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Address string `protobuf:"bytes,2,opt,name=address" json:"address,omitempty"`
	Type    Type   `protobuf:"varint,3,opt,name=type,enum=service.Type" json:"type,omitempty"`
}

func (m *Node) Reset()                    { *m = Node{} }
func (m *Node) String() string            { return proto.CompactTextString(m) }
func (*Node) ProtoMessage()               {}
func (*Node) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func init() {
	proto.RegisterType((*Node)(nil), "service.Node")
	proto.RegisterEnum("service.Type", Type_name, Type_value)
}

func init() { proto.RegisterFile("type.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 157 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x2a, 0xa9, 0x2c, 0x48,
	0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2f, 0x4e, 0x2d, 0x2a, 0xcb, 0x4c, 0x4e, 0x55,
	0x0a, 0xe6, 0x62, 0xf1, 0xcb, 0x4f, 0x49, 0x15, 0xe2, 0xe3, 0x62, 0xca, 0x4c, 0x91, 0x60, 0x54,
	0x60, 0xd4, 0xe0, 0x0c, 0x02, 0xb2, 0x84, 0x24, 0xb8, 0xd8, 0x13, 0x53, 0x52, 0x8a, 0x52, 0x8b,
	0x8b, 0x25, 0x98, 0xc0, 0x82, 0x30, 0xae, 0x90, 0x22, 0x17, 0x0b, 0xc8, 0x20, 0x09, 0x66, 0xa0,
	0x30, 0x9f, 0x11, 0xaf, 0x1e, 0xd4, 0x24, 0xbd, 0x10, 0xa0, 0x60, 0x10, 0x58, 0x4a, 0xcb, 0x9c,
	0x8b, 0x05, 0xc4, 0x13, 0xe2, 0xe5, 0xe2, 0x0c, 0xf5, 0x73, 0x71, 0x75, 0xf3, 0xf4, 0x73, 0x75,
	0x11, 0x60, 0x00, 0x71, 0xc3, 0x1c, 0x7d, 0x3c, 0x5d, 0x1c, 0x43, 0xfc, 0x83, 0x04, 0x18, 0x85,
	0x04, 0xb9, 0x78, 0xfd, 0xfc, 0xfd, 0xe2, 0x11, 0x42, 0x4c, 0x49, 0x6c, 0x60, 0xd7, 0x19, 0x03,
	0x02, 0x00, 0x00, 0xff, 0xff, 0xef, 0xd2, 0x4d, 0xaa, 0xab, 0x00, 0x00, 0x00,
}
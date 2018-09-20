// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protos/tendermint/tendermint.proto

package tendermint

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

type Tx struct {
	Method               string   `protobuf:"bytes,1,opt,name=method,proto3" json:"method,omitempty"`
	Param                string   `protobuf:"bytes,2,opt,name=param,proto3" json:"param,omitempty"`
	Nonce                string   `protobuf:"bytes,3,opt,name=nonce,proto3" json:"nonce,omitempty"`
	Signature            []byte   `protobuf:"bytes,4,opt,name=signature,proto3" json:"signature,omitempty"`
	NodeID               string   `protobuf:"bytes,5,opt,name=nodeID,proto3" json:"nodeID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Tx) Reset()         { *m = Tx{} }
func (m *Tx) String() string { return proto.CompactTextString(m) }
func (*Tx) ProtoMessage()    {}
func (*Tx) Descriptor() ([]byte, []int) {
	return fileDescriptor_tendermint_7b6b82f5ac49f186, []int{0}
}
func (m *Tx) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Tx.Unmarshal(m, b)
}
func (m *Tx) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Tx.Marshal(b, m, deterministic)
}
func (dst *Tx) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Tx.Merge(dst, src)
}
func (m *Tx) XXX_Size() int {
	return xxx_messageInfo_Tx.Size(m)
}
func (m *Tx) XXX_DiscardUnknown() {
	xxx_messageInfo_Tx.DiscardUnknown(m)
}

var xxx_messageInfo_Tx proto.InternalMessageInfo

func (m *Tx) GetMethod() string {
	if m != nil {
		return m.Method
	}
	return ""
}

func (m *Tx) GetParam() string {
	if m != nil {
		return m.Param
	}
	return ""
}

func (m *Tx) GetNonce() string {
	if m != nil {
		return m.Nonce
	}
	return ""
}

func (m *Tx) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func (m *Tx) GetNodeID() string {
	if m != nil {
		return m.NodeID
	}
	return ""
}

type Query struct {
	FnName               string   `protobuf:"bytes,1,opt,name=fnName,proto3" json:"fnName,omitempty"`
	Param                string   `protobuf:"bytes,2,opt,name=param,proto3" json:"param,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Query) Reset()         { *m = Query{} }
func (m *Query) String() string { return proto.CompactTextString(m) }
func (*Query) ProtoMessage()    {}
func (*Query) Descriptor() ([]byte, []int) {
	return fileDescriptor_tendermint_7b6b82f5ac49f186, []int{1}
}
func (m *Query) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Query.Unmarshal(m, b)
}
func (m *Query) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Query.Marshal(b, m, deterministic)
}
func (dst *Query) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Query.Merge(dst, src)
}
func (m *Query) XXX_Size() int {
	return xxx_messageInfo_Query.Size(m)
}
func (m *Query) XXX_DiscardUnknown() {
	xxx_messageInfo_Query.DiscardUnknown(m)
}

var xxx_messageInfo_Query proto.InternalMessageInfo

func (m *Query) GetFnName() string {
	if m != nil {
		return m.FnName
	}
	return ""
}

func (m *Query) GetParam() string {
	if m != nil {
		return m.Param
	}
	return ""
}

func init() {
	proto.RegisterType((*Tx)(nil), "Tx")
	proto.RegisterType((*Query)(nil), "Query")
}

func init() {
	proto.RegisterFile("protos/tendermint/tendermint.proto", fileDescriptor_tendermint_7b6b82f5ac49f186)
}

var fileDescriptor_tendermint_7b6b82f5ac49f186 = []byte{
	// 164 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x2a, 0x28, 0xca, 0x2f,
	0xc9, 0x2f, 0xd6, 0x2f, 0x49, 0xcd, 0x4b, 0x49, 0x2d, 0xca, 0xcd, 0xcc, 0x2b, 0x41, 0x62, 0xea,
	0x81, 0x25, 0x95, 0xea, 0xb8, 0x98, 0x42, 0x2a, 0x84, 0xc4, 0xb8, 0xd8, 0x72, 0x53, 0x4b, 0x32,
	0xf2, 0x53, 0x24, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0xa0, 0x3c, 0x21, 0x11, 0x2e, 0xd6, 0x82,
	0xc4, 0xa2, 0xc4, 0x5c, 0x09, 0x26, 0xb0, 0x30, 0x84, 0x03, 0x12, 0xcd, 0xcb, 0xcf, 0x4b, 0x4e,
	0x95, 0x60, 0x86, 0x88, 0x82, 0x39, 0x42, 0x32, 0x5c, 0x9c, 0xc5, 0x99, 0xe9, 0x79, 0x89, 0x25,
	0xa5, 0x45, 0xa9, 0x12, 0x2c, 0x0a, 0x8c, 0x1a, 0x3c, 0x41, 0x08, 0x01, 0x90, 0x0d, 0x79, 0xf9,
	0x29, 0xa9, 0x9e, 0x2e, 0x12, 0xac, 0x10, 0x1b, 0x20, 0x3c, 0x25, 0x53, 0x2e, 0xd6, 0xc0, 0xd2,
	0xd4, 0xa2, 0x4a, 0x90, 0x82, 0xb4, 0x3c, 0xbf, 0xc4, 0xdc, 0x54, 0x98, 0x13, 0x20, 0x3c, 0xec,
	0x4e, 0x48, 0x62, 0x03, 0xbb, 0xde, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x0c, 0x6a, 0x27, 0x7b,
	0xe3, 0x00, 0x00, 0x00,
}

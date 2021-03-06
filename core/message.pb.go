// Code generated by protoc-gen-go. DO NOT EDIT.
// source: message.proto

/*
Package core is a generated protocol buffer package.

It is generated from these files:
	message.proto

It has these top-level messages:
	SerializedMessageData
	SerializedMessage
*/
package core

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

type SerializedMessageData struct {
	Data             []byte            `protobuf:"bytes,1,req,name=Data" json:"Data,omitempty"`
	Metadata         map[string][]byte `protobuf:"bytes,2,rep,name=Metadata" json:"Metadata,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	XXX_unrecognized []byte            `json:"-"`
}

func (m *SerializedMessageData) Reset()                    { *m = SerializedMessageData{} }
func (m *SerializedMessageData) String() string            { return proto.CompactTextString(m) }
func (*SerializedMessageData) ProtoMessage()               {}
func (*SerializedMessageData) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *SerializedMessageData) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *SerializedMessageData) GetMetadata() map[string][]byte {
	if m != nil {
		return m.Metadata
	}
	return nil
}

type SerializedMessage struct {
	StreamID         *uint64                `protobuf:"varint,1,req,name=StreamID" json:"StreamID,omitempty"`
	Data             *SerializedMessageData `protobuf:"bytes,2,req,name=Data" json:"Data,omitempty"`
	PrevStreamID     *uint64                `protobuf:"varint,3,opt,name=PrevStreamID" json:"PrevStreamID,omitempty"`
	OrigStreamID     *uint64                `protobuf:"varint,4,opt,name=OrigStreamID" json:"OrigStreamID,omitempty"`
	Timestamp        *int64                 `protobuf:"varint,5,opt,name=Timestamp" json:"Timestamp,omitempty"`
	Original         *SerializedMessageData `protobuf:"bytes,6,opt,name=Original" json:"Original,omitempty"`
	XXX_unrecognized []byte                 `json:"-"`
}

func (m *SerializedMessage) Reset()                    { *m = SerializedMessage{} }
func (m *SerializedMessage) String() string            { return proto.CompactTextString(m) }
func (*SerializedMessage) ProtoMessage()               {}
func (*SerializedMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *SerializedMessage) GetStreamID() uint64 {
	if m != nil && m.StreamID != nil {
		return *m.StreamID
	}
	return 0
}

func (m *SerializedMessage) GetData() *SerializedMessageData {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *SerializedMessage) GetPrevStreamID() uint64 {
	if m != nil && m.PrevStreamID != nil {
		return *m.PrevStreamID
	}
	return 0
}

func (m *SerializedMessage) GetOrigStreamID() uint64 {
	if m != nil && m.OrigStreamID != nil {
		return *m.OrigStreamID
	}
	return 0
}

func (m *SerializedMessage) GetTimestamp() int64 {
	if m != nil && m.Timestamp != nil {
		return *m.Timestamp
	}
	return 0
}

func (m *SerializedMessage) GetOriginal() *SerializedMessageData {
	if m != nil {
		return m.Original
	}
	return nil
}

func init() {
	proto.RegisterType((*SerializedMessageData)(nil), "serializedMessageData")
	proto.RegisterType((*SerializedMessage)(nil), "serializedMessage")
}

func init() { proto.RegisterFile("message.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 269 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x8f, 0x41, 0x4b, 0xc3, 0x40,
	0x10, 0x85, 0x99, 0x4d, 0x5a, 0xd2, 0x69, 0x0a, 0xba, 0xa8, 0x84, 0xe2, 0x61, 0x09, 0x1e, 0x16,
	0x0f, 0x39, 0xe4, 0x24, 0x7a, 0x11, 0xa9, 0x07, 0x0f, 0x45, 0x59, 0x3d, 0x79, 0x5b, 0xec, 0x50,
	0x16, 0x93, 0xa6, 0x6c, 0xd6, 0x42, 0xfd, 0x49, 0xfe, 0x3f, 0xef, 0xb2, 0x4d, 0x5d, 0x2d, 0x16,
	0x4f, 0x3b, 0xf3, 0xe6, 0x9b, 0x9d, 0xf7, 0x70, 0x54, 0x53, 0xdb, 0xea, 0x39, 0x15, 0x4b, 0xdb,
	0xb8, 0x26, 0xff, 0x00, 0x3c, 0x6e, 0xc9, 0x1a, 0x5d, 0x99, 0x77, 0x9a, 0x4d, 0xbb, 0xd9, 0x44,
	0x3b, 0xcd, 0x39, 0xc6, 0xfe, 0xcd, 0x40, 0x30, 0x99, 0xaa, 0x4d, 0xcd, 0xaf, 0x31, 0x99, 0x92,
	0xd3, 0x33, 0xaf, 0x33, 0x11, 0xc9, 0x61, 0x79, 0x56, 0xec, 0xdd, 0x2e, 0xbe, 0xb1, 0xdb, 0x85,
	0xb3, 0x6b, 0x15, 0xb6, 0xc6, 0x57, 0x38, 0xda, 0x19, 0xf1, 0x03, 0x8c, 0x5e, 0x69, 0x9d, 0x81,
	0x00, 0x39, 0x50, 0xbe, 0xe4, 0x47, 0xd8, 0x5b, 0xe9, 0xea, 0x8d, 0x32, 0x26, 0x40, 0xa6, 0xaa,
	0x6b, 0x2e, 0xd9, 0x05, 0xe4, 0x9f, 0x80, 0x87, 0x7f, 0xce, 0xf1, 0x31, 0x26, 0x8f, 0xce, 0x92,
	0xae, 0xef, 0x26, 0x1b, 0xb3, 0xb1, 0x0a, 0x3d, 0x3f, 0xdf, 0x86, 0x60, 0x82, 0xc9, 0x61, 0x79,
	0xb2, 0xdf, 0xec, 0x36, 0x5c, 0x8e, 0xe9, 0x83, 0xa5, 0x55, 0xf8, 0x2b, 0x12, 0x20, 0x63, 0xb5,
	0xa3, 0x79, 0xe6, 0xde, 0x9a, 0x79, 0x60, 0xe2, 0x8e, 0xf9, 0xad, 0xf1, 0x53, 0x1c, 0x3c, 0x99,
	0x9a, 0x5a, 0xa7, 0xeb, 0x65, 0xd6, 0x13, 0x20, 0x23, 0xf5, 0x23, 0xf0, 0x12, 0x13, 0x4f, 0x9b,
	0x85, 0xae, 0xb2, 0xbe, 0x80, 0x7f, 0x5c, 0x05, 0xee, 0xa6, 0xff, 0x1c, 0xbf, 0x34, 0x96, 0xbe,
	0x02, 0x00, 0x00, 0xff, 0xff, 0xc4, 0xc6, 0x18, 0xde, 0xbc, 0x01, 0x00, 0x00,
}

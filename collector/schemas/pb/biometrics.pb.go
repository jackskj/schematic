// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pb/biometrics.proto

package pb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

//  Sample message to be sent from an apple watch over grpc to the collector
type Biometrics struct {
	WatchId              int32    `protobuf:"varint,1,opt,name=watch_id,json=watchId,proto3" json:"watch_id,omitempty"`
	UserName             string   `protobuf:"bytes,2,opt,name=user_name,json=userName,proto3" json:"user_name,omitempty"`
	Heartbeat            int32    `protobuf:"varint,3,opt,name=heartbeat,proto3" json:"heartbeat,omitempty"`
	Temperature          int32    `protobuf:"varint,4,opt,name=temperature,proto3" json:"temperature,omitempty"`
	Latitude             float32  `protobuf:"fixed32,5,opt,name=latitude,proto3" json:"latitude,omitempty"`
	Longitude            float32  `protobuf:"fixed32,6,opt,name=longitude,proto3" json:"longitude,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Biometrics) Reset()         { *m = Biometrics{} }
func (m *Biometrics) String() string { return proto.CompactTextString(m) }
func (*Biometrics) ProtoMessage()    {}
func (*Biometrics) Descriptor() ([]byte, []int) {
	return fileDescriptor_a02968e76442d0ef, []int{0}
}

func (m *Biometrics) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Biometrics.Unmarshal(m, b)
}
func (m *Biometrics) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Biometrics.Marshal(b, m, deterministic)
}
func (m *Biometrics) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Biometrics.Merge(m, src)
}
func (m *Biometrics) XXX_Size() int {
	return xxx_messageInfo_Biometrics.Size(m)
}
func (m *Biometrics) XXX_DiscardUnknown() {
	xxx_messageInfo_Biometrics.DiscardUnknown(m)
}

var xxx_messageInfo_Biometrics proto.InternalMessageInfo

func (m *Biometrics) GetWatchId() int32 {
	if m != nil {
		return m.WatchId
	}
	return 0
}

func (m *Biometrics) GetUserName() string {
	if m != nil {
		return m.UserName
	}
	return ""
}

func (m *Biometrics) GetHeartbeat() int32 {
	if m != nil {
		return m.Heartbeat
	}
	return 0
}

func (m *Biometrics) GetTemperature() int32 {
	if m != nil {
		return m.Temperature
	}
	return 0
}

func (m *Biometrics) GetLatitude() float32 {
	if m != nil {
		return m.Latitude
	}
	return 0
}

func (m *Biometrics) GetLongitude() float32 {
	if m != nil {
		return m.Longitude
	}
	return 0
}

func init() {
	proto.RegisterType((*Biometrics)(nil), "pb.Biometrics")
}

func init() { proto.RegisterFile("pb/biometrics.proto", fileDescriptor_a02968e76442d0ef) }

var fileDescriptor_a02968e76442d0ef = []byte{
	// 184 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x8f, 0x3d, 0x8e, 0x83, 0x30,
	0x10, 0x46, 0x65, 0x76, 0x61, 0x61, 0xb6, 0x73, 0x1a, 0xe7, 0xa7, 0xb0, 0x52, 0x51, 0x25, 0x45,
	0x6e, 0x90, 0x2e, 0x4d, 0x0a, 0x2e, 0x80, 0x6c, 0x18, 0x05, 0x4b, 0x18, 0x5b, 0x66, 0x50, 0x6e,
	0x97, 0xb3, 0x45, 0x38, 0x12, 0xa4, 0x7c, 0xef, 0x69, 0x46, 0xfa, 0x60, 0xe3, 0xf5, 0x59, 0x1b,
	0x67, 0x91, 0x82, 0x69, 0xc6, 0x93, 0x0f, 0x8e, 0x1c, 0x4f, 0xbc, 0x3e, 0xbe, 0x18, 0xc0, 0x75,
	0x09, 0x7c, 0x0b, 0xf9, 0x53, 0x51, 0xd3, 0xd5, 0xa6, 0x15, 0x4c, 0xb2, 0x32, 0xad, 0xfe, 0x22,
	0xdf, 0x5a, 0xbe, 0x87, 0x62, 0x1a, 0x31, 0xd4, 0x83, 0xb2, 0x28, 0x12, 0xc9, 0xca, 0xa2, 0xca,
	0x67, 0x71, 0x57, 0x16, 0xf9, 0x01, 0x8a, 0x0e, 0x55, 0x20, 0x8d, 0x8a, 0xc4, 0x4f, 0x3c, 0x5c,
	0x05, 0x97, 0xf0, 0x4f, 0x68, 0x3d, 0x06, 0x45, 0x53, 0x40, 0xf1, 0x1b, 0xfb, 0xb7, 0xe2, 0x3b,
	0xc8, 0x7b, 0x45, 0x86, 0xa6, 0x16, 0x45, 0x2a, 0x59, 0x99, 0x54, 0x0b, 0xcf, 0xbf, 0x7b, 0x37,
	0x3c, 0x3e, 0x31, 0x8b, 0x71, 0x15, 0x3a, 0x8b, 0x5b, 0x2e, 0xef, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x70, 0x37, 0x81, 0xe5, 0xe2, 0x00, 0x00, 0x00,
}
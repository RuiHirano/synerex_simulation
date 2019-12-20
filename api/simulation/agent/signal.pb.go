// Code generated by protoc-gen-go. DO NOT EDIT.
// source: simulation/agent/signal.proto

package agent

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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Signal struct {
	Status               *SignalStatus `protobuf:"bytes,4,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *Signal) Reset()         { *m = Signal{} }
func (m *Signal) String() string { return proto.CompactTextString(m) }
func (*Signal) ProtoMessage()    {}
func (*Signal) Descriptor() ([]byte, []int) {
	return fileDescriptor_65489d1b2b9fdc09, []int{0}
}

func (m *Signal) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Signal.Unmarshal(m, b)
}
func (m *Signal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Signal.Marshal(b, m, deterministic)
}
func (m *Signal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Signal.Merge(m, src)
}
func (m *Signal) XXX_Size() int {
	return xxx_messageInfo_Signal.Size(m)
}
func (m *Signal) XXX_DiscardUnknown() {
	xxx_messageInfo_Signal.DiscardUnknown(m)
}

var xxx_messageInfo_Signal proto.InternalMessageInfo

func (m *Signal) GetStatus() *SignalStatus {
	if m != nil {
		return m.Status
	}
	return nil
}

type SignalStatus struct {
	Age                  string   `protobuf:"bytes,1,opt,name=age,proto3" json:"age,omitempty"`
	Sex                  string   `protobuf:"bytes,2,opt,name=sex,proto3" json:"sex,omitempty"`
	Name                 string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SignalStatus) Reset()         { *m = SignalStatus{} }
func (m *SignalStatus) String() string { return proto.CompactTextString(m) }
func (*SignalStatus) ProtoMessage()    {}
func (*SignalStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_65489d1b2b9fdc09, []int{1}
}

func (m *SignalStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SignalStatus.Unmarshal(m, b)
}
func (m *SignalStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SignalStatus.Marshal(b, m, deterministic)
}
func (m *SignalStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SignalStatus.Merge(m, src)
}
func (m *SignalStatus) XXX_Size() int {
	return xxx_messageInfo_SignalStatus.Size(m)
}
func (m *SignalStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_SignalStatus.DiscardUnknown(m)
}

var xxx_messageInfo_SignalStatus proto.InternalMessageInfo

func (m *SignalStatus) GetAge() string {
	if m != nil {
		return m.Age
	}
	return ""
}

func (m *SignalStatus) GetSex() string {
	if m != nil {
		return m.Sex
	}
	return ""
}

func (m *SignalStatus) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func init() {
	proto.RegisterType((*Signal)(nil), "api.signal.Signal")
	proto.RegisterType((*SignalStatus)(nil), "api.signal.SignalStatus")
}

func init() { proto.RegisterFile("simulation/agent/signal.proto", fileDescriptor_65489d1b2b9fdc09) }

var fileDescriptor_65489d1b2b9fdc09 = []byte{
	// 186 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x8f, 0x31, 0x0b, 0xc2, 0x30,
	0x10, 0x46, 0xa9, 0x2d, 0x05, 0xa3, 0x83, 0x64, 0xca, 0x22, 0x94, 0x4e, 0x9d, 0x12, 0x51, 0x44,
	0x70, 0x74, 0xf0, 0x07, 0xb4, 0x9b, 0x8b, 0x5c, 0x25, 0xa4, 0x81, 0x36, 0x09, 0xbd, 0x14, 0xea,
	0xbf, 0x97, 0xa6, 0x8a, 0xe2, 0x74, 0x8f, 0x77, 0x6f, 0xf9, 0xc8, 0x16, 0x75, 0x37, 0xb4, 0xe0,
	0xb5, 0x35, 0x02, 0x94, 0x34, 0x5e, 0xa0, 0x56, 0x06, 0x5a, 0xee, 0x7a, 0xeb, 0x2d, 0x25, 0xe0,
	0x34, 0x9f, 0x4d, 0x7e, 0x26, 0x69, 0x15, 0x88, 0xee, 0x48, 0x8a, 0x1e, 0xfc, 0x80, 0x2c, 0xc9,
	0xa2, 0x62, 0xb5, 0x67, 0xfc, 0x9b, 0xf1, 0xb9, 0xa9, 0xc2, 0xbf, 0x7c, 0x77, 0xf9, 0x95, 0xac,
	0x7f, 0x3d, 0xdd, 0x90, 0x18, 0x94, 0x64, 0x51, 0x16, 0x15, 0xcb, 0x72, 0xc2, 0xc9, 0xa0, 0x1c,
	0xd9, 0x62, 0x36, 0x28, 0x47, 0x4a, 0x49, 0x62, 0xa0, 0x93, 0x2c, 0x0e, 0x2a, 0xf0, 0xe5, 0x74,
	0x3b, 0x2a, 0xed, 0x9b, 0xa1, 0xe6, 0x0f, 0xdb, 0x09, 0x7c, 0x1a, 0xd9, 0xcb, 0xf1, 0x73, 0xef,
	0xd0, 0xba, 0x06, 0x04, 0x38, 0x2d, 0xfe, 0x57, 0xd5, 0x69, 0xd8, 0x73, 0x78, 0x05, 0x00, 0x00,
	0xff, 0xff, 0x39, 0x49, 0x27, 0xb8, 0xf0, 0x00, 0x00, 0x00,
}
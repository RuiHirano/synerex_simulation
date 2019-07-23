// Code generated by protoc-gen-go. DO NOT EDIT.
// source: area/area.proto

package area

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

type AreaService struct {
	NumAdult             uint32   `protobuf:"varint,5,opt,name=num_adult,json=numAdult,proto3" json:"num_adult,omitempty"`
	NumChild             uint32   `protobuf:"varint,6,opt,name=num_child,json=numChild,proto3" json:"num_child,omitempty"`
	AmountPrice          uint32   `protobuf:"varint,8,opt,name=amount_price,json=amountPrice,proto3" json:"amount_price,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AreaService) Reset()         { *m = AreaService{} }
func (m *AreaService) String() string { return proto.CompactTextString(m) }
func (*AreaService) ProtoMessage()    {}
func (*AreaService) Descriptor() ([]byte, []int) {
	return fileDescriptor_6bf477ccdfe50b7a, []int{0}
}

func (m *AreaService) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AreaService.Unmarshal(m, b)
}
func (m *AreaService) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AreaService.Marshal(b, m, deterministic)
}
func (m *AreaService) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AreaService.Merge(m, src)
}
func (m *AreaService) XXX_Size() int {
	return xxx_messageInfo_AreaService.Size(m)
}
func (m *AreaService) XXX_DiscardUnknown() {
	xxx_messageInfo_AreaService.DiscardUnknown(m)
}

var xxx_messageInfo_AreaService proto.InternalMessageInfo

func (m *AreaService) GetNumAdult() uint32 {
	if m != nil {
		return m.NumAdult
	}
	return 0
}

func (m *AreaService) GetNumChild() uint32 {
	if m != nil {
		return m.NumChild
	}
	return 0
}

func (m *AreaService) GetAmountPrice() uint32 {
	if m != nil {
		return m.AmountPrice
	}
	return 0
}

func init() {
	proto.RegisterType((*AreaService)(nil), "api.area.AreaService")
}

func init() { proto.RegisterFile("area/area.proto", fileDescriptor_6bf477ccdfe50b7a) }

var fileDescriptor_6bf477ccdfe50b7a = []byte{
	// 168 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4f, 0x2c, 0x4a, 0x4d,
	0xd4, 0x07, 0x11, 0x7a, 0x05, 0x45, 0xf9, 0x25, 0xf9, 0x42, 0x1c, 0x89, 0x05, 0x99, 0x7a, 0x20,
	0xbe, 0x52, 0x16, 0x17, 0xb7, 0x63, 0x51, 0x6a, 0x62, 0x70, 0x6a, 0x51, 0x59, 0x66, 0x72, 0xaa,
	0x90, 0x34, 0x17, 0x67, 0x5e, 0x69, 0x6e, 0x7c, 0x62, 0x4a, 0x69, 0x4e, 0x89, 0x04, 0xab, 0x02,
	0xa3, 0x06, 0x6f, 0x10, 0x47, 0x5e, 0x69, 0xae, 0x23, 0x88, 0x0f, 0x93, 0x4c, 0xce, 0xc8, 0xcc,
	0x49, 0x91, 0x60, 0x83, 0x4b, 0x3a, 0x83, 0xf8, 0x42, 0x8a, 0x5c, 0x3c, 0x89, 0xb9, 0xf9, 0xa5,
	0x79, 0x25, 0xf1, 0x05, 0x45, 0x99, 0xc9, 0xa9, 0x12, 0x1c, 0x60, 0x79, 0x6e, 0x88, 0x58, 0x00,
	0x48, 0xc8, 0x49, 0x3b, 0x4a, 0x33, 0x3d, 0xb3, 0x24, 0xa3, 0x34, 0x49, 0x2f, 0x39, 0x3f, 0x57,
	0xbf, 0xb8, 0x32, 0x2f, 0xb5, 0x28, 0xb5, 0x02, 0x46, 0xc7, 0x27, 0xe6, 0x14, 0x64, 0x24, 0xea,
	0x27, 0x16, 0x64, 0x82, 0x1d, 0x9a, 0xc4, 0x06, 0x76, 0xa9, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff,
	0x40, 0x80, 0x5b, 0xea, 0xbc, 0x00, 0x00, 0x00,
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// source: simulation/area/area.proto

package area

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	common "github.com/synerex/synerex_alpha/api/simulation/common"
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

type Area struct {
	Id                   uint64          `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                 string          `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	DuplicateArea        []*common.Coord `protobuf:"bytes,3,rep,name=duplicate_area,json=duplicateArea,proto3" json:"duplicate_area,omitempty"`
	ControlArea          []*common.Coord `protobuf:"bytes,4,rep,name=control_area,json=controlArea,proto3" json:"control_area,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *Area) Reset()         { *m = Area{} }
func (m *Area) String() string { return proto.CompactTextString(m) }
func (*Area) ProtoMessage()    {}
func (*Area) Descriptor() ([]byte, []int) {
	return fileDescriptor_c8212774c97f163d, []int{0}
}

func (m *Area) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Area.Unmarshal(m, b)
}
func (m *Area) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Area.Marshal(b, m, deterministic)
}
func (m *Area) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Area.Merge(m, src)
}
func (m *Area) XXX_Size() int {
	return xxx_messageInfo_Area.Size(m)
}
func (m *Area) XXX_DiscardUnknown() {
	xxx_messageInfo_Area.DiscardUnknown(m)
}

var xxx_messageInfo_Area proto.InternalMessageInfo

func (m *Area) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Area) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Area) GetDuplicateArea() []*common.Coord {
	if m != nil {
		return m.DuplicateArea
	}
	return nil
}

func (m *Area) GetControlArea() []*common.Coord {
	if m != nil {
		return m.ControlArea
	}
	return nil
}

func init() {
	proto.RegisterType((*Area)(nil), "api.area.Area")
}

func init() { proto.RegisterFile("simulation/area/area.proto", fileDescriptor_c8212774c97f163d) }

var fileDescriptor_c8212774c97f163d = []byte{
	// 211 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x2a, 0xce, 0xcc, 0x2d,
	0xcd, 0x49, 0x2c, 0xc9, 0xcc, 0xcf, 0xd3, 0x4f, 0x2c, 0x4a, 0x4d, 0x04, 0x13, 0x7a, 0x05, 0x45,
	0xf9, 0x25, 0xf9, 0x42, 0x1c, 0x89, 0x05, 0x99, 0x7a, 0x20, 0xbe, 0x94, 0x1c, 0x92, 0xaa, 0xe4,
	0xfc, 0xdc, 0x5c, 0x38, 0x05, 0x51, 0xa9, 0x34, 0x8b, 0x91, 0x8b, 0xc5, 0xb1, 0x28, 0x35, 0x51,
	0x88, 0x8f, 0x8b, 0x29, 0x33, 0x45, 0x82, 0x51, 0x81, 0x51, 0x83, 0x25, 0x88, 0x29, 0x33, 0x45,
	0x48, 0x88, 0x8b, 0x25, 0x2f, 0x31, 0x37, 0x55, 0x82, 0x49, 0x81, 0x51, 0x83, 0x33, 0x08, 0xcc,
	0x16, 0xb2, 0xe0, 0xe2, 0x4b, 0x29, 0x2d, 0xc8, 0xc9, 0x4c, 0x4e, 0x2c, 0x49, 0x8d, 0x07, 0x19,
	0x2f, 0xc1, 0xac, 0xc0, 0xac, 0xc1, 0x6d, 0x24, 0xa8, 0x07, 0xb2, 0x0f, 0x6a, 0xae, 0x73, 0x7e,
	0x7e, 0x51, 0x4a, 0x10, 0x2f, 0x5c, 0x21, 0xd8, 0x74, 0x13, 0x2e, 0x9e, 0xe4, 0xfc, 0xbc, 0x92,
	0xa2, 0xfc, 0x1c, 0x88, 0x3e, 0x16, 0x5c, 0xfa, 0xb8, 0xa1, 0xca, 0x40, 0xba, 0x9c, 0xcc, 0xa2,
	0x4c, 0xd2, 0x33, 0x4b, 0x32, 0x4a, 0x93, 0x40, 0x6a, 0xf4, 0x8b, 0x2b, 0xf3, 0x52, 0x8b, 0x52,
	0x2b, 0x60, 0x74, 0x7c, 0x62, 0x4e, 0x41, 0x46, 0xa2, 0x7e, 0x62, 0x41, 0xa6, 0x3e, 0x5a, 0x48,
	0x24, 0xb1, 0x81, 0xfd, 0x66, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x2b, 0x7f, 0x65, 0xb9, 0x23,
	0x01, 0x00, 0x00,
}

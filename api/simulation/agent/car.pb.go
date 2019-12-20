// Code generated by protoc-gen-go. DO NOT EDIT.
// source: simulation/agent/car.proto

package agent

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

type Car struct {
	Status               *CarStatus `protobuf:"bytes,4,opt,name=status,proto3" json:"status,omitempty"`
	Route                *CarRoute  `protobuf:"bytes,7,opt,name=route,proto3" json:"route,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Car) Reset()         { *m = Car{} }
func (m *Car) String() string { return proto.CompactTextString(m) }
func (*Car) ProtoMessage()    {}
func (*Car) Descriptor() ([]byte, []int) {
	return fileDescriptor_5df268eb821aa90f, []int{0}
}

func (m *Car) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Car.Unmarshal(m, b)
}
func (m *Car) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Car.Marshal(b, m, deterministic)
}
func (m *Car) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Car.Merge(m, src)
}
func (m *Car) XXX_Size() int {
	return xxx_messageInfo_Car.Size(m)
}
func (m *Car) XXX_DiscardUnknown() {
	xxx_messageInfo_Car.DiscardUnknown(m)
}

var xxx_messageInfo_Car proto.InternalMessageInfo

func (m *Car) GetStatus() *CarStatus {
	if m != nil {
		return m.Status
	}
	return nil
}

func (m *Car) GetRoute() *CarRoute {
	if m != nil {
		return m.Route
	}
	return nil
}

type CarRoute struct {
	Position             *common.Coord   `protobuf:"bytes,1,opt,name=position,proto3" json:"position,omitempty"`
	Direction            float64         `protobuf:"fixed64,2,opt,name=direction,proto3" json:"direction,omitempty"`
	Speed                float64         `protobuf:"fixed64,3,opt,name=speed,proto3" json:"speed,omitempty"`
	Destination          *common.Coord   `protobuf:"bytes,4,opt,name=destination,proto3" json:"destination,omitempty"`
	Departure            *common.Coord   `protobuf:"bytes,5,opt,name=departure,proto3" json:"departure,omitempty"`
	TransitPoints        []*common.Coord `protobuf:"bytes,6,rep,name=transit_points,json=transitPoints,proto3" json:"transit_points,omitempty"`
	NextTransit          *common.Coord   `protobuf:"bytes,7,opt,name=next_transit,json=nextTransit,proto3" json:"next_transit,omitempty"`
	TotalDistance        float64         `protobuf:"fixed64,8,opt,name=total_distance,json=totalDistance,proto3" json:"total_distance,omitempty"`
	RequiredTime         float64         `protobuf:"fixed64,9,opt,name=required_time,json=requiredTime,proto3" json:"required_time,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *CarRoute) Reset()         { *m = CarRoute{} }
func (m *CarRoute) String() string { return proto.CompactTextString(m) }
func (*CarRoute) ProtoMessage()    {}
func (*CarRoute) Descriptor() ([]byte, []int) {
	return fileDescriptor_5df268eb821aa90f, []int{1}
}

func (m *CarRoute) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CarRoute.Unmarshal(m, b)
}
func (m *CarRoute) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CarRoute.Marshal(b, m, deterministic)
}
func (m *CarRoute) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CarRoute.Merge(m, src)
}
func (m *CarRoute) XXX_Size() int {
	return xxx_messageInfo_CarRoute.Size(m)
}
func (m *CarRoute) XXX_DiscardUnknown() {
	xxx_messageInfo_CarRoute.DiscardUnknown(m)
}

var xxx_messageInfo_CarRoute proto.InternalMessageInfo

func (m *CarRoute) GetPosition() *common.Coord {
	if m != nil {
		return m.Position
	}
	return nil
}

func (m *CarRoute) GetDirection() float64 {
	if m != nil {
		return m.Direction
	}
	return 0
}

func (m *CarRoute) GetSpeed() float64 {
	if m != nil {
		return m.Speed
	}
	return 0
}

func (m *CarRoute) GetDestination() *common.Coord {
	if m != nil {
		return m.Destination
	}
	return nil
}

func (m *CarRoute) GetDeparture() *common.Coord {
	if m != nil {
		return m.Departure
	}
	return nil
}

func (m *CarRoute) GetTransitPoints() []*common.Coord {
	if m != nil {
		return m.TransitPoints
	}
	return nil
}

func (m *CarRoute) GetNextTransit() *common.Coord {
	if m != nil {
		return m.NextTransit
	}
	return nil
}

func (m *CarRoute) GetTotalDistance() float64 {
	if m != nil {
		return m.TotalDistance
	}
	return 0
}

func (m *CarRoute) GetRequiredTime() float64 {
	if m != nil {
		return m.RequiredTime
	}
	return 0
}

type CarStatus struct {
	Age                  string   `protobuf:"bytes,1,opt,name=age,proto3" json:"age,omitempty"`
	Sex                  string   `protobuf:"bytes,2,opt,name=sex,proto3" json:"sex,omitempty"`
	Name                 string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CarStatus) Reset()         { *m = CarStatus{} }
func (m *CarStatus) String() string { return proto.CompactTextString(m) }
func (*CarStatus) ProtoMessage()    {}
func (*CarStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_5df268eb821aa90f, []int{2}
}

func (m *CarStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CarStatus.Unmarshal(m, b)
}
func (m *CarStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CarStatus.Marshal(b, m, deterministic)
}
func (m *CarStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CarStatus.Merge(m, src)
}
func (m *CarStatus) XXX_Size() int {
	return xxx_messageInfo_CarStatus.Size(m)
}
func (m *CarStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_CarStatus.DiscardUnknown(m)
}

var xxx_messageInfo_CarStatus proto.InternalMessageInfo

func (m *CarStatus) GetAge() string {
	if m != nil {
		return m.Age
	}
	return ""
}

func (m *CarStatus) GetSex() string {
	if m != nil {
		return m.Sex
	}
	return ""
}

func (m *CarStatus) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func init() {
	proto.RegisterType((*Car)(nil), "api.car.Car")
	proto.RegisterType((*CarRoute)(nil), "api.car.CarRoute")
	proto.RegisterType((*CarStatus)(nil), "api.car.CarStatus")
}

func init() { proto.RegisterFile("simulation/agent/car.proto", fileDescriptor_5df268eb821aa90f) }

var fileDescriptor_5df268eb821aa90f = []byte{
	// 388 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0x51, 0x6f, 0xd3, 0x30,
	0x10, 0xc7, 0x55, 0xb2, 0x76, 0xcd, 0x6d, 0x9d, 0xe0, 0xc4, 0x83, 0x35, 0x21, 0x54, 0x15, 0x21,
	0x26, 0x24, 0x1a, 0x89, 0x81, 0xe0, 0x99, 0xf2, 0x01, 0x90, 0xd9, 0xd3, 0x5e, 0xa2, 0x5b, 0x72,
	0xea, 0x2c, 0x35, 0xb6, 0xb1, 0x2f, 0x52, 0xf9, 0x14, 0x7c, 0x65, 0x14, 0x27, 0xed, 0x32, 0x34,
	0x9e, 0x7c, 0xf9, 0xdf, 0xef, 0xe2, 0xff, 0x9d, 0x0f, 0x2e, 0xa3, 0x69, 0xda, 0x1d, 0x89, 0x71,
	0xb6, 0xa0, 0x2d, 0x5b, 0x29, 0x2a, 0x0a, 0x6b, 0x1f, 0x9c, 0x38, 0x3c, 0x25, 0x6f, 0xd6, 0x15,
	0x85, 0xcb, 0xd7, 0x23, 0xa8, 0x72, 0x4d, 0x73, 0x3c, 0x7a, 0x70, 0x75, 0x0b, 0xd9, 0x86, 0x02,
	0xbe, 0x87, 0x59, 0x14, 0x92, 0x36, 0xaa, 0x93, 0xe5, 0xe4, 0xea, 0xec, 0x23, 0xae, 0x87, 0x1f,
	0xac, 0x37, 0x14, 0x7e, 0xa6, 0x8c, 0x1e, 0x08, 0x7c, 0x07, 0xd3, 0xe0, 0x5a, 0x61, 0x75, 0x9a,
	0xd0, 0x17, 0x63, 0x54, 0x77, 0x09, 0xdd, 0xe7, 0x57, 0x7f, 0x32, 0x98, 0x1f, 0x34, 0xfc, 0x00,
	0x73, 0xef, 0xa2, 0xe9, 0x8c, 0xa8, 0xc9, 0xb8, 0xb0, 0x77, 0xb3, 0x71, 0x2e, 0xd4, 0xfa, 0x88,
	0xe0, 0x2b, 0xc8, 0x6b, 0x13, 0xb8, 0x4a, 0xfc, 0xb3, 0xe5, 0xe4, 0x6a, 0xa2, 0x1f, 0x04, 0x7c,
	0x09, 0xd3, 0xe8, 0x99, 0x6b, 0x95, 0xa5, 0x4c, 0xff, 0x81, 0xd7, 0x70, 0x56, 0x73, 0x14, 0x63,
	0x53, 0xbb, 0x43, 0x27, 0x4f, 0xdc, 0x32, 0xa6, 0xb0, 0x80, 0xbc, 0x66, 0x4f, 0x41, 0xda, 0xc0,
	0x6a, 0xfa, 0xbf, 0x92, 0x07, 0x06, 0xbf, 0xc2, 0x85, 0x04, 0xb2, 0xd1, 0x48, 0xe9, 0x9d, 0xb1,
	0x12, 0xd5, 0x6c, 0x99, 0x3d, 0x5d, 0xb5, 0x18, 0xc0, 0x1f, 0x89, 0xc3, 0x4f, 0x70, 0x6e, 0x79,
	0x2f, 0xe5, 0xa0, 0x3e, 0x9e, 0xdf, 0x23, 0x83, 0x1d, 0x76, 0xd3, 0x53, 0xf8, 0x16, 0x2e, 0xc4,
	0x09, 0xed, 0xca, 0xda, 0x44, 0x21, 0x5b, 0xb1, 0x9a, 0xa7, 0xa6, 0x17, 0x49, 0xfd, 0x3e, 0x88,
	0xf8, 0x06, 0x16, 0x81, 0x7f, 0xb5, 0x26, 0x70, 0x5d, 0x8a, 0x69, 0x58, 0xe5, 0x89, 0x3a, 0x3f,
	0x88, 0x37, 0xa6, 0xe1, 0xd5, 0x06, 0xf2, 0xe3, 0x7b, 0xe2, 0x73, 0xc8, 0x68, 0xcb, 0xe9, 0x31,
	0x72, 0xdd, 0x85, 0x9d, 0x12, 0x79, 0x9f, 0xc6, 0x9d, 0xeb, 0x2e, 0x44, 0x84, 0x13, 0x4b, 0x0d,
	0xa7, 0x39, 0xe7, 0x3a, 0xc5, 0xdf, 0xbe, 0xdc, 0x7e, 0xde, 0x1a, 0xb9, 0x6f, 0xef, 0x3a, 0xd3,
	0x45, 0xfc, 0x6d, 0x39, 0xf0, 0xfe, 0x70, 0x96, 0xb4, 0xf3, 0xf7, 0x54, 0x90, 0x37, 0xc5, 0xbf,
	0xeb, 0x79, 0x37, 0x4b, 0x2b, 0x77, 0xfd, 0x37, 0x00, 0x00, 0xff, 0xff, 0x84, 0x43, 0x28, 0x5e,
	0xb9, 0x02, 0x00, 0x00,
}
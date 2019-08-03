// Code generated by protoc-gen-go. DO NOT EDIT.
// source: rideshare/rideshare.proto

package rideshare

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	duration "github.com/golang/protobuf/ptypes/duration"
	common "github.com/synerex/synerex_alpha/api/common"
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

type TrafficType int32

const (
	TrafficType_TAXI  TrafficType = 0
	TrafficType_BUS   TrafficType = 1
	TrafficType_TRAIN TrafficType = 2
)

var TrafficType_name = map[int32]string{
	0: "TAXI",
	1: "BUS",
	2: "TRAIN",
}

var TrafficType_value = map[string]int32{
	"TAXI":  0,
	"BUS":   1,
	"TRAIN": 2,
}

func (x TrafficType) String() string {
	return proto.EnumName(TrafficType_name, int32(x))
}

func (TrafficType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_aaf9387a82289ad6, []int{0}
}

type StatusType int32

const (
	StatusType_FREE   StatusType = 0
	StatusType_PICKUP StatusType = 1
	StatusType_RIDE   StatusType = 2
	StatusType_FULL   StatusType = 3
)

var StatusType_name = map[int32]string{
	0: "FREE",
	1: "PICKUP",
	2: "RIDE",
	3: "FULL",
}

var StatusType_value = map[string]int32{
	"FREE":   0,
	"PICKUP": 1,
	"RIDE":   2,
	"FULL":   3,
}

func (x StatusType) String() string {
	return proto.EnumName(StatusType_name, int32(x))
}

func (StatusType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_aaf9387a82289ad6, []int{1}
}

type RideShare struct {
	DepartPoint          *common.Place `protobuf:"bytes,1,opt,name=depart_point,json=departPoint,proto3" json:"depart_point,omitempty"`
	ArrivePoint          *common.Place `protobuf:"bytes,2,opt,name=arrive_point,json=arrivePoint,proto3" json:"arrive_point,omitempty"`
	DepartTime           *common.Time  `protobuf:"bytes,3,opt,name=depart_time,json=departTime,proto3" json:"depart_time,omitempty"`
	ArriveTime           *common.Time  `protobuf:"bytes,4,opt,name=arrive_time,json=arriveTime,proto3" json:"arrive_time,omitempty"`
	NumAdult             uint32        `protobuf:"varint,5,opt,name=num_adult,json=numAdult,proto3" json:"num_adult,omitempty"`
	NumChild             uint32        `protobuf:"varint,6,opt,name=num_child,json=numChild,proto3" json:"num_child,omitempty"`
	Routes               []*Route      `protobuf:"bytes,7,rep,name=routes,proto3" json:"routes,omitempty"`
	AmountPrice          uint32        `protobuf:"varint,8,opt,name=amount_price,json=amountPrice,proto3" json:"amount_price,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *RideShare) Reset()         { *m = RideShare{} }
func (m *RideShare) String() string { return proto.CompactTextString(m) }
func (*RideShare) ProtoMessage()    {}
func (*RideShare) Descriptor() ([]byte, []int) {
	return fileDescriptor_aaf9387a82289ad6, []int{0}
}

func (m *RideShare) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RideShare.Unmarshal(m, b)
}
func (m *RideShare) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RideShare.Marshal(b, m, deterministic)
}
func (m *RideShare) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RideShare.Merge(m, src)
}
func (m *RideShare) XXX_Size() int {
	return xxx_messageInfo_RideShare.Size(m)
}
func (m *RideShare) XXX_DiscardUnknown() {
	xxx_messageInfo_RideShare.DiscardUnknown(m)
}

var xxx_messageInfo_RideShare proto.InternalMessageInfo

func (m *RideShare) GetDepartPoint() *common.Place {
	if m != nil {
		return m.DepartPoint
	}
	return nil
}

func (m *RideShare) GetArrivePoint() *common.Place {
	if m != nil {
		return m.ArrivePoint
	}
	return nil
}

func (m *RideShare) GetDepartTime() *common.Time {
	if m != nil {
		return m.DepartTime
	}
	return nil
}

func (m *RideShare) GetArriveTime() *common.Time {
	if m != nil {
		return m.ArriveTime
	}
	return nil
}

func (m *RideShare) GetNumAdult() uint32 {
	if m != nil {
		return m.NumAdult
	}
	return 0
}

func (m *RideShare) GetNumChild() uint32 {
	if m != nil {
		return m.NumChild
	}
	return 0
}

func (m *RideShare) GetRoutes() []*Route {
	if m != nil {
		return m.Routes
	}
	return nil
}

func (m *RideShare) GetAmountPrice() uint32 {
	if m != nil {
		return m.AmountPrice
	}
	return 0
}

type Route struct {
	TrafficType          TrafficType        `protobuf:"varint,1,opt,name=traffic_type,json=trafficType,proto3,enum=api.rideshare.TrafficType" json:"traffic_type,omitempty"`
	StatusType           StatusType         `protobuf:"varint,2,opt,name=status_type,json=statusType,proto3,enum=api.rideshare.StatusType" json:"status_type,omitempty"`
	TransportName        string             `protobuf:"bytes,3,opt,name=transport_name,json=transportName,proto3" json:"transport_name,omitempty"`
	TransportLine        string             `protobuf:"bytes,4,opt,name=transport_line,json=transportLine,proto3" json:"transport_line,omitempty"`
	Destination          string             `protobuf:"bytes,5,opt,name=destination,proto3" json:"destination,omitempty"`
	DepartPoint          *common.Place      `protobuf:"bytes,6,opt,name=depart_point,json=departPoint,proto3" json:"depart_point,omitempty"`
	ArrivePoint          *common.Place      `protobuf:"bytes,7,opt,name=arrive_point,json=arrivePoint,proto3" json:"arrive_point,omitempty"`
	DepartTime           *common.Time       `protobuf:"bytes,8,opt,name=depart_time,json=departTime,proto3" json:"depart_time,omitempty"`
	ArriveTime           *common.Time       `protobuf:"bytes,9,opt,name=arrive_time,json=arriveTime,proto3" json:"arrive_time,omitempty"`
	AmountTime           *duration.Duration `protobuf:"bytes,10,opt,name=amount_time,json=amountTime,proto3" json:"amount_time,omitempty"`
	AmountPrice          uint32             `protobuf:"varint,11,opt,name=amount_price,json=amountPrice,proto3" json:"amount_price,omitempty"`
	AmountSheets         uint32             `protobuf:"varint,12,opt,name=amount_sheets,json=amountSheets,proto3" json:"amount_sheets,omitempty"`
	AvailableSheets      uint32             `protobuf:"varint,13,opt,name=available_sheets,json=availableSheets,proto3" json:"available_sheets,omitempty"`
	Points               []*common.Point    `protobuf:"bytes,15,rep,name=points,proto3" json:"points,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *Route) Reset()         { *m = Route{} }
func (m *Route) String() string { return proto.CompactTextString(m) }
func (*Route) ProtoMessage()    {}
func (*Route) Descriptor() ([]byte, []int) {
	return fileDescriptor_aaf9387a82289ad6, []int{1}
}

func (m *Route) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Route.Unmarshal(m, b)
}
func (m *Route) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Route.Marshal(b, m, deterministic)
}
func (m *Route) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Route.Merge(m, src)
}
func (m *Route) XXX_Size() int {
	return xxx_messageInfo_Route.Size(m)
}
func (m *Route) XXX_DiscardUnknown() {
	xxx_messageInfo_Route.DiscardUnknown(m)
}

var xxx_messageInfo_Route proto.InternalMessageInfo

func (m *Route) GetTrafficType() TrafficType {
	if m != nil {
		return m.TrafficType
	}
	return TrafficType_TAXI
}

func (m *Route) GetStatusType() StatusType {
	if m != nil {
		return m.StatusType
	}
	return StatusType_FREE
}

func (m *Route) GetTransportName() string {
	if m != nil {
		return m.TransportName
	}
	return ""
}

func (m *Route) GetTransportLine() string {
	if m != nil {
		return m.TransportLine
	}
	return ""
}

func (m *Route) GetDestination() string {
	if m != nil {
		return m.Destination
	}
	return ""
}

func (m *Route) GetDepartPoint() *common.Place {
	if m != nil {
		return m.DepartPoint
	}
	return nil
}

func (m *Route) GetArrivePoint() *common.Place {
	if m != nil {
		return m.ArrivePoint
	}
	return nil
}

func (m *Route) GetDepartTime() *common.Time {
	if m != nil {
		return m.DepartTime
	}
	return nil
}

func (m *Route) GetArriveTime() *common.Time {
	if m != nil {
		return m.ArriveTime
	}
	return nil
}

func (m *Route) GetAmountTime() *duration.Duration {
	if m != nil {
		return m.AmountTime
	}
	return nil
}

func (m *Route) GetAmountPrice() uint32 {
	if m != nil {
		return m.AmountPrice
	}
	return 0
}

func (m *Route) GetAmountSheets() uint32 {
	if m != nil {
		return m.AmountSheets
	}
	return 0
}

func (m *Route) GetAvailableSheets() uint32 {
	if m != nil {
		return m.AvailableSheets
	}
	return 0
}

func (m *Route) GetPoints() []*common.Point {
	if m != nil {
		return m.Points
	}
	return nil
}

func init() {
	proto.RegisterEnum("api.rideshare.TrafficType", TrafficType_name, TrafficType_value)
	proto.RegisterEnum("api.rideshare.StatusType", StatusType_name, StatusType_value)
	proto.RegisterType((*RideShare)(nil), "api.rideshare.RideShare")
	proto.RegisterType((*Route)(nil), "api.rideshare.Route")
}

func init() { proto.RegisterFile("rideshare/rideshare.proto", fileDescriptor_aaf9387a82289ad6) }

var fileDescriptor_aaf9387a82289ad6 = []byte{
	// 591 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x94, 0x5f, 0x6b, 0xdb, 0x3c,
	0x14, 0x87, 0x9b, 0xa4, 0x71, 0xe3, 0xe3, 0xa4, 0xf5, 0xab, 0x77, 0x17, 0x6e, 0x07, 0x23, 0xeb,
	0x18, 0x74, 0xdd, 0xb0, 0xb7, 0x6e, 0xec, 0xa2, 0xb0, 0x8b, 0xfe, 0x1b, 0x84, 0x85, 0x12, 0x94,
	0x04, 0xc6, 0x6e, 0x8c, 0x12, 0x2b, 0x89, 0xc0, 0x96, 0x8d, 0x2c, 0x97, 0xe5, 0xcb, 0xec, 0xf3,
	0xed, 0x63, 0x0c, 0x49, 0x8e, 0x93, 0xa6, 0x14, 0xb6, 0x5e, 0xd9, 0x3a, 0xe7, 0x79, 0xa4, 0x70,
	0xfc, 0x53, 0xe0, 0x50, 0xb0, 0x88, 0xe6, 0x0b, 0x22, 0x68, 0x50, 0xbd, 0xf9, 0x99, 0x48, 0x65,
	0x8a, 0x3a, 0x24, 0x63, 0x7e, 0x55, 0x3c, 0x7a, 0x31, 0x4f, 0xd3, 0x79, 0x4c, 0x03, 0xdd, 0x9c,
	0x14, 0xb3, 0x20, 0x2a, 0x04, 0x91, 0x2c, 0xe5, 0x06, 0x3f, 0xfa, 0x7f, 0x9a, 0x26, 0x49, 0xca,
	0x03, 0xf3, 0x30, 0xc5, 0xe3, 0xdf, 0x75, 0xb0, 0x31, 0x8b, 0xe8, 0x50, 0x6d, 0x81, 0x3e, 0x41,
	0x3b, 0xa2, 0x19, 0x11, 0x32, 0xcc, 0x52, 0xc6, 0xa5, 0x57, 0xeb, 0xd6, 0x4e, 0x9c, 0xb3, 0xff,
	0x7c, 0x75, 0x50, 0xa9, 0x0d, 0x62, 0x32, 0xa5, 0xd8, 0x31, 0xd8, 0x40, 0x51, 0xca, 0x22, 0x42,
	0xb0, 0x3b, 0x5a, 0x5a, 0xf5, 0x47, 0x2d, 0x83, 0x19, 0xeb, 0x03, 0x94, 0x9b, 0x84, 0x92, 0x25,
	0xd4, 0x6b, 0x68, 0xc9, 0xdd, 0x94, 0x46, 0x2c, 0xa1, 0x18, 0x0c, 0xa4, 0xde, 0x95, 0x52, 0x1e,
	0xa4, 0x95, 0xdd, 0xc7, 0x14, 0x03, 0x69, 0xe5, 0x39, 0xd8, 0xbc, 0x48, 0x42, 0x12, 0x15, 0xb1,
	0xf4, 0x9a, 0xdd, 0xda, 0x49, 0x07, 0xb7, 0x78, 0x91, 0x5c, 0xa8, 0xf5, 0xaa, 0x39, 0x5d, 0xb0,
	0x38, 0xf2, 0xac, 0xaa, 0x79, 0xa5, 0xd6, 0xe8, 0x1d, 0x58, 0x22, 0x2d, 0x24, 0xcd, 0xbd, 0xbd,
	0x6e, 0xe3, 0xc4, 0x39, 0x7b, 0xe6, 0xdf, 0x1b, 0xb7, 0x8f, 0x55, 0x13, 0x97, 0x0c, 0x7a, 0x09,
	0x6d, 0x92, 0xa4, 0x05, 0x97, 0x61, 0x26, 0xd8, 0x94, 0x7a, 0x2d, 0xbd, 0x9b, 0x63, 0x6a, 0x03,
	0x55, 0x3a, 0xfe, 0xd5, 0x84, 0xa6, 0x96, 0xd0, 0x17, 0x68, 0x4b, 0x41, 0x66, 0x33, 0x36, 0x0d,
	0xe5, 0x32, 0xa3, 0x7a, 0xcc, 0xfb, 0x67, 0x47, 0x5b, 0x07, 0x8c, 0x0c, 0x32, 0x5a, 0x66, 0x14,
	0x3b, 0x72, 0xbd, 0x40, 0xe7, 0xe0, 0xe4, 0x92, 0xc8, 0x22, 0x37, 0x76, 0x5d, 0xdb, 0x87, 0x5b,
	0xf6, 0x50, 0x13, 0x5a, 0x86, 0xbc, 0x7a, 0x47, 0xaf, 0x61, 0x5f, 0x0a, 0xc2, 0xf3, 0x2c, 0x15,
	0x32, 0xe4, 0xa4, 0x1c, 0xbc, 0x8d, 0x3b, 0x55, 0xf5, 0x96, 0x24, 0x5b, 0x58, 0xcc, 0xb8, 0x19,
	0xf6, 0x26, 0xd6, 0x67, 0x9c, 0xa2, 0xae, 0xfa, 0x86, 0xb9, 0x64, 0x5c, 0xe7, 0x4c, 0xcf, 0xd7,
	0xc6, 0x9b, 0xa5, 0x07, 0x89, 0xb2, 0x9e, 0x94, 0xa8, 0xbd, 0xa7, 0x24, 0xaa, 0xf5, 0xef, 0x89,
	0xb2, 0xff, 0x22, 0x51, 0xe7, 0x50, 0x7e, 0x55, 0xa3, 0x80, 0x56, 0x0e, 0x7d, 0x73, 0xf9, 0xfc,
	0xd5, 0xe5, 0xf3, 0xaf, 0xcb, 0xcb, 0x87, 0xc1, 0xd0, 0xda, 0xdd, 0x4e, 0x89, 0xf3, 0x20, 0x25,
	0xe8, 0x15, 0x74, 0x4a, 0x24, 0x5f, 0x50, 0x2a, 0x73, 0xaf, 0xad, 0x99, 0xd2, 0x1b, 0xea, 0x1a,
	0x7a, 0x03, 0x2e, 0xb9, 0x23, 0x2c, 0x26, 0x93, 0x98, 0xae, 0xb8, 0x8e, 0xe6, 0x0e, 0xaa, 0x7a,
	0x85, 0x5a, 0x7a, 0x86, 0xb9, 0x77, 0xa0, 0x63, 0x7c, 0x7f, 0x88, 0xaa, 0x83, 0x4b, 0xe0, 0xf4,
	0x2d, 0x38, 0x1b, 0x99, 0x43, 0x2d, 0xd8, 0x1d, 0x5d, 0x7c, 0xef, 0xb9, 0x3b, 0x68, 0x0f, 0x1a,
	0x97, 0xe3, 0xa1, 0x5b, 0x43, 0x36, 0x34, 0x47, 0xf8, 0xa2, 0x77, 0xeb, 0xd6, 0x4f, 0x3f, 0x03,
	0xac, 0x23, 0xa6, 0xd8, 0xaf, 0xf8, 0xe6, 0xc6, 0xdd, 0x41, 0x00, 0xd6, 0xa0, 0x77, 0xf5, 0x6d,
	0x3c, 0x70, 0x6b, 0xaa, 0x8a, 0x7b, 0xd7, 0x37, 0x6e, 0x5d, 0xf7, 0xc7, 0xfd, 0xbe, 0xdb, 0xb8,
	0x7c, 0xff, 0xc3, 0x9f, 0x33, 0xb9, 0x28, 0x26, 0xea, 0x37, 0x04, 0xf9, 0x92, 0x53, 0x41, 0x7f,
	0xae, 0x9e, 0x21, 0x89, 0xb3, 0x05, 0x09, 0x48, 0xc6, 0xd6, 0x7f, 0x76, 0x13, 0x4b, 0xcf, 0xf4,
	0xe3, 0x9f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x8c, 0x78, 0x04, 0x40, 0x0a, 0x05, 0x00, 0x00,
}

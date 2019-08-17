// Code generated by protoc-gen-go. DO NOT EDIT.
// source: simulation/area/area.proto

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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type StatusType int32

const (
	StatusType_OK StatusType = 0
	StatusType_NG StatusType = 1
)

var StatusType_name = map[int32]string{
	0: "OK",
	1: "NG",
}

var StatusType_value = map[string]int32{
	"OK": 0,
	"NG": 1,
}

func (x StatusType) String() string {
	return proto.EnumName(StatusType_name, int32(x))
}

func (StatusType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_c8212774c97f163d, []int{0}
}

type MessageType int32

const (
	MessageType_SET_AREA    MessageType = 0
	MessageType_GET_AREA    MessageType = 1
	MessageType_AREA_INFO   MessageType = 2
	MessageType_AREA_STATUS MessageType = 3
)

var MessageType_name = map[int32]string{
	0: "SET_AREA",
	1: "GET_AREA",
	2: "AREA_INFO",
	3: "AREA_STATUS",
}

var MessageType_value = map[string]int32{
	"SET_AREA":    0,
	"GET_AREA":    1,
	"AREA_INFO":   2,
	"AREA_STATUS": 3,
}

func (x MessageType) String() string {
	return proto.EnumName(MessageType_name, int32(x))
}

func (MessageType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_c8212774c97f163d, []int{1}
}

type AreaService struct {
	MessageType          MessageType        `protobuf:"varint,1,opt,name=message_type,json=messageType,proto3,enum=api.area.MessageType" json:"message_type,omitempty"`
	AreaId               uint32             `protobuf:"varint,2,opt,name=area_id,json=areaId,proto3" json:"area_id,omitempty"`
	AreaName             string             `protobuf:"bytes,3,opt,name=area_name,json=areaName,proto3" json:"area_name,omitempty"`
	Coord                *AreaService_Coord `protobuf:"bytes,4,opt,name=coord,proto3" json:"coord,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *AreaService) Reset()         { *m = AreaService{} }
func (m *AreaService) String() string { return proto.CompactTextString(m) }
func (*AreaService) ProtoMessage()    {}
func (*AreaService) Descriptor() ([]byte, []int) {
	return fileDescriptor_c8212774c97f163d, []int{0}
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

func (m *AreaService) GetMessageType() MessageType {
	if m != nil {
		return m.MessageType
	}
	return MessageType_SET_AREA
}

func (m *AreaService) GetAreaId() uint32 {
	if m != nil {
		return m.AreaId
	}
	return 0
}

func (m *AreaService) GetAreaName() string {
	if m != nil {
		return m.AreaName
	}
	return ""
}

func (m *AreaService) GetCoord() *AreaService_Coord {
	if m != nil {
		return m.Coord
	}
	return nil
}

type AreaService_Coord struct {
	StartLat             float32  `protobuf:"fixed32,1,opt,name=start_lat,json=startLat,proto3" json:"start_lat,omitempty"`
	StartLon             float32  `protobuf:"fixed32,2,opt,name=start_lon,json=startLon,proto3" json:"start_lon,omitempty"`
	EndLat               float32  `protobuf:"fixed32,3,opt,name=end_lat,json=endLat,proto3" json:"end_lat,omitempty"`
	EndLon               float32  `protobuf:"fixed32,4,opt,name=end_lon,json=endLon,proto3" json:"end_lon,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AreaService_Coord) Reset()         { *m = AreaService_Coord{} }
func (m *AreaService_Coord) String() string { return proto.CompactTextString(m) }
func (*AreaService_Coord) ProtoMessage()    {}
func (*AreaService_Coord) Descriptor() ([]byte, []int) {
	return fileDescriptor_c8212774c97f163d, []int{0, 0}
}

func (m *AreaService_Coord) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AreaService_Coord.Unmarshal(m, b)
}
func (m *AreaService_Coord) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AreaService_Coord.Marshal(b, m, deterministic)
}
func (m *AreaService_Coord) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AreaService_Coord.Merge(m, src)
}
func (m *AreaService_Coord) XXX_Size() int {
	return xxx_messageInfo_AreaService_Coord.Size(m)
}
func (m *AreaService_Coord) XXX_DiscardUnknown() {
	xxx_messageInfo_AreaService_Coord.DiscardUnknown(m)
}

var xxx_messageInfo_AreaService_Coord proto.InternalMessageInfo

func (m *AreaService_Coord) GetStartLat() float32 {
	if m != nil {
		return m.StartLat
	}
	return 0
}

func (m *AreaService_Coord) GetStartLon() float32 {
	if m != nil {
		return m.StartLon
	}
	return 0
}

func (m *AreaService_Coord) GetEndLat() float32 {
	if m != nil {
		return m.EndLat
	}
	return 0
}

func (m *AreaService_Coord) GetEndLon() float32 {
	if m != nil {
		return m.EndLon
	}
	return 0
}

type AreaInfo struct {
	MessageType          MessageType     `protobuf:"varint,1,opt,name=message_type,json=messageType,proto3,enum=api.area.MessageType" json:"message_type,omitempty"`
	AreaId               uint32          `protobuf:"varint,2,opt,name=area_id,json=areaId,proto3" json:"area_id,omitempty"`
	AreaName             string          `protobuf:"bytes,3,opt,name=area_name,json=areaName,proto3" json:"area_name,omitempty"`
	Coord                *AreaInfo_Coord `protobuf:"bytes,4,opt,name=coord,proto3" json:"coord,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *AreaInfo) Reset()         { *m = AreaInfo{} }
func (m *AreaInfo) String() string { return proto.CompactTextString(m) }
func (*AreaInfo) ProtoMessage()    {}
func (*AreaInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_c8212774c97f163d, []int{1}
}

func (m *AreaInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AreaInfo.Unmarshal(m, b)
}
func (m *AreaInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AreaInfo.Marshal(b, m, deterministic)
}
func (m *AreaInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AreaInfo.Merge(m, src)
}
func (m *AreaInfo) XXX_Size() int {
	return xxx_messageInfo_AreaInfo.Size(m)
}
func (m *AreaInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_AreaInfo.DiscardUnknown(m)
}

var xxx_messageInfo_AreaInfo proto.InternalMessageInfo

func (m *AreaInfo) GetMessageType() MessageType {
	if m != nil {
		return m.MessageType
	}
	return MessageType_SET_AREA
}

func (m *AreaInfo) GetAreaId() uint32 {
	if m != nil {
		return m.AreaId
	}
	return 0
}

func (m *AreaInfo) GetAreaName() string {
	if m != nil {
		return m.AreaName
	}
	return ""
}

func (m *AreaInfo) GetCoord() *AreaInfo_Coord {
	if m != nil {
		return m.Coord
	}
	return nil
}

type AreaInfo_Coord struct {
	StartLat             float32  `protobuf:"fixed32,1,opt,name=start_lat,json=startLat,proto3" json:"start_lat,omitempty"`
	StartLon             float32  `protobuf:"fixed32,2,opt,name=start_lon,json=startLon,proto3" json:"start_lon,omitempty"`
	EndLat               float32  `protobuf:"fixed32,3,opt,name=end_lat,json=endLat,proto3" json:"end_lat,omitempty"`
	EndLon               float32  `protobuf:"fixed32,4,opt,name=end_lon,json=endLon,proto3" json:"end_lon,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AreaInfo_Coord) Reset()         { *m = AreaInfo_Coord{} }
func (m *AreaInfo_Coord) String() string { return proto.CompactTextString(m) }
func (*AreaInfo_Coord) ProtoMessage()    {}
func (*AreaInfo_Coord) Descriptor() ([]byte, []int) {
	return fileDescriptor_c8212774c97f163d, []int{1, 0}
}

func (m *AreaInfo_Coord) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AreaInfo_Coord.Unmarshal(m, b)
}
func (m *AreaInfo_Coord) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AreaInfo_Coord.Marshal(b, m, deterministic)
}
func (m *AreaInfo_Coord) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AreaInfo_Coord.Merge(m, src)
}
func (m *AreaInfo_Coord) XXX_Size() int {
	return xxx_messageInfo_AreaInfo_Coord.Size(m)
}
func (m *AreaInfo_Coord) XXX_DiscardUnknown() {
	xxx_messageInfo_AreaInfo_Coord.DiscardUnknown(m)
}

var xxx_messageInfo_AreaInfo_Coord proto.InternalMessageInfo

func (m *AreaInfo_Coord) GetStartLat() float32 {
	if m != nil {
		return m.StartLat
	}
	return 0
}

func (m *AreaInfo_Coord) GetStartLon() float32 {
	if m != nil {
		return m.StartLon
	}
	return 0
}

func (m *AreaInfo_Coord) GetEndLat() float32 {
	if m != nil {
		return m.EndLat
	}
	return 0
}

func (m *AreaInfo_Coord) GetEndLon() float32 {
	if m != nil {
		return m.EndLon
	}
	return 0
}

type GetArea struct {
	MessageType          MessageType    `protobuf:"varint,1,opt,name=message_type,json=messageType,proto3,enum=api.area.MessageType" json:"message_type,omitempty"`
	AreaId               uint32         `protobuf:"varint,2,opt,name=area_id,json=areaId,proto3" json:"area_id,omitempty"`
	Coord                *GetArea_Coord `protobuf:"bytes,3,opt,name=coord,proto3" json:"coord,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *GetArea) Reset()         { *m = GetArea{} }
func (m *GetArea) String() string { return proto.CompactTextString(m) }
func (*GetArea) ProtoMessage()    {}
func (*GetArea) Descriptor() ([]byte, []int) {
	return fileDescriptor_c8212774c97f163d, []int{2}
}

func (m *GetArea) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetArea.Unmarshal(m, b)
}
func (m *GetArea) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetArea.Marshal(b, m, deterministic)
}
func (m *GetArea) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetArea.Merge(m, src)
}
func (m *GetArea) XXX_Size() int {
	return xxx_messageInfo_GetArea.Size(m)
}
func (m *GetArea) XXX_DiscardUnknown() {
	xxx_messageInfo_GetArea.DiscardUnknown(m)
}

var xxx_messageInfo_GetArea proto.InternalMessageInfo

func (m *GetArea) GetMessageType() MessageType {
	if m != nil {
		return m.MessageType
	}
	return MessageType_SET_AREA
}

func (m *GetArea) GetAreaId() uint32 {
	if m != nil {
		return m.AreaId
	}
	return 0
}

func (m *GetArea) GetCoord() *GetArea_Coord {
	if m != nil {
		return m.Coord
	}
	return nil
}

type GetArea_Coord struct {
	StartLat             float32  `protobuf:"fixed32,1,opt,name=start_lat,json=startLat,proto3" json:"start_lat,omitempty"`
	StartLon             float32  `protobuf:"fixed32,2,opt,name=start_lon,json=startLon,proto3" json:"start_lon,omitempty"`
	EndLat               float32  `protobuf:"fixed32,3,opt,name=end_lat,json=endLat,proto3" json:"end_lat,omitempty"`
	EndLon               float32  `protobuf:"fixed32,4,opt,name=end_lon,json=endLon,proto3" json:"end_lon,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetArea_Coord) Reset()         { *m = GetArea_Coord{} }
func (m *GetArea_Coord) String() string { return proto.CompactTextString(m) }
func (*GetArea_Coord) ProtoMessage()    {}
func (*GetArea_Coord) Descriptor() ([]byte, []int) {
	return fileDescriptor_c8212774c97f163d, []int{2, 0}
}

func (m *GetArea_Coord) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetArea_Coord.Unmarshal(m, b)
}
func (m *GetArea_Coord) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetArea_Coord.Marshal(b, m, deterministic)
}
func (m *GetArea_Coord) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetArea_Coord.Merge(m, src)
}
func (m *GetArea_Coord) XXX_Size() int {
	return xxx_messageInfo_GetArea_Coord.Size(m)
}
func (m *GetArea_Coord) XXX_DiscardUnknown() {
	xxx_messageInfo_GetArea_Coord.DiscardUnknown(m)
}

var xxx_messageInfo_GetArea_Coord proto.InternalMessageInfo

func (m *GetArea_Coord) GetStartLat() float32 {
	if m != nil {
		return m.StartLat
	}
	return 0
}

func (m *GetArea_Coord) GetStartLon() float32 {
	if m != nil {
		return m.StartLon
	}
	return 0
}

func (m *GetArea_Coord) GetEndLat() float32 {
	if m != nil {
		return m.EndLat
	}
	return 0
}

func (m *GetArea_Coord) GetEndLon() float32 {
	if m != nil {
		return m.EndLon
	}
	return 0
}

type AreaStatus struct {
	MessageType          MessageType `protobuf:"varint,1,opt,name=message_type,json=messageType,proto3,enum=api.area.MessageType" json:"message_type,omitempty"`
	AreaId               uint32      `protobuf:"varint,2,opt,name=area_id,json=areaId,proto3" json:"area_id,omitempty"`
	AreaName             string      `protobuf:"bytes,3,opt,name=area_name,json=areaName,proto3" json:"area_name,omitempty"`
	StatusType           StatusType  `protobuf:"varint,4,opt,name=status_type,json=statusType,proto3,enum=api.area.StatusType" json:"status_type,omitempty"`
	Meta                 string      `protobuf:"bytes,5,opt,name=meta,proto3" json:"meta,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *AreaStatus) Reset()         { *m = AreaStatus{} }
func (m *AreaStatus) String() string { return proto.CompactTextString(m) }
func (*AreaStatus) ProtoMessage()    {}
func (*AreaStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_c8212774c97f163d, []int{3}
}

func (m *AreaStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AreaStatus.Unmarshal(m, b)
}
func (m *AreaStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AreaStatus.Marshal(b, m, deterministic)
}
func (m *AreaStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AreaStatus.Merge(m, src)
}
func (m *AreaStatus) XXX_Size() int {
	return xxx_messageInfo_AreaStatus.Size(m)
}
func (m *AreaStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_AreaStatus.DiscardUnknown(m)
}

var xxx_messageInfo_AreaStatus proto.InternalMessageInfo

func (m *AreaStatus) GetMessageType() MessageType {
	if m != nil {
		return m.MessageType
	}
	return MessageType_SET_AREA
}

func (m *AreaStatus) GetAreaId() uint32 {
	if m != nil {
		return m.AreaId
	}
	return 0
}

func (m *AreaStatus) GetAreaName() string {
	if m != nil {
		return m.AreaName
	}
	return ""
}

func (m *AreaStatus) GetStatusType() StatusType {
	if m != nil {
		return m.StatusType
	}
	return StatusType_OK
}

func (m *AreaStatus) GetMeta() string {
	if m != nil {
		return m.Meta
	}
	return ""
}

func init() {
	proto.RegisterEnum("api.area.StatusType", StatusType_name, StatusType_value)
	proto.RegisterEnum("api.area.MessageType", MessageType_name, MessageType_value)
	proto.RegisterType((*AreaService)(nil), "api.area.AreaService")
	proto.RegisterType((*AreaService_Coord)(nil), "api.area.AreaService.Coord")
	proto.RegisterType((*AreaInfo)(nil), "api.area.AreaInfo")
	proto.RegisterType((*AreaInfo_Coord)(nil), "api.area.AreaInfo.Coord")
	proto.RegisterType((*GetArea)(nil), "api.area.GetArea")
	proto.RegisterType((*GetArea_Coord)(nil), "api.area.GetArea.Coord")
	proto.RegisterType((*AreaStatus)(nil), "api.area.AreaStatus")
}

func init() { proto.RegisterFile("simulation/area/area.proto", fileDescriptor_c8212774c97f163d) }

var fileDescriptor_c8212774c97f163d = []byte{
	// 443 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xcc, 0x54, 0x4d, 0x6f, 0xd3, 0x40,
	0x10, 0xed, 0x3a, 0x1f, 0x75, 0xc6, 0x2d, 0x58, 0x2b, 0x50, 0xad, 0x96, 0x43, 0xd4, 0x53, 0x54,
	0x09, 0x47, 0x94, 0x0f, 0x71, 0x0d, 0xa8, 0x44, 0x11, 0x90, 0x4a, 0x4e, 0xb8, 0x70, 0xb1, 0xa6,
	0xf1, 0xd0, 0x5a, 0x8a, 0x77, 0x2d, 0xef, 0x06, 0x91, 0x1b, 0xff, 0x84, 0x13, 0xbf, 0x85, 0xbf,
	0x85, 0x76, 0x52, 0xd7, 0x29, 0x77, 0xaa, 0x5c, 0xec, 0x79, 0xfb, 0x66, 0xe6, 0x3d, 0x3d, 0xad,
	0x16, 0x8e, 0x4d, 0x5e, 0xac, 0x96, 0x68, 0x73, 0xad, 0x86, 0x58, 0x11, 0xf2, 0x27, 0x2e, 0x2b,
	0x6d, 0xb5, 0xf4, 0xb1, 0xcc, 0x63, 0x87, 0x4f, 0x7f, 0x7b, 0x10, 0x8c, 0x2a, 0xc2, 0x19, 0x55,
	0xdf, 0xf3, 0x05, 0xc9, 0xb7, 0x70, 0x50, 0x90, 0x31, 0x78, 0x4d, 0xa9, 0x5d, 0x97, 0x14, 0x89,
	0xbe, 0x18, 0x3c, 0x3a, 0x7f, 0x1a, 0xd7, 0x03, 0xf1, 0xe7, 0x0d, 0x3b, 0x5f, 0x97, 0x94, 0x04,
	0x45, 0x03, 0xe4, 0x11, 0xec, 0xbb, 0x86, 0x34, 0xcf, 0x22, 0xaf, 0x2f, 0x06, 0x87, 0x49, 0xd7,
	0xc1, 0x49, 0x26, 0x4f, 0xa0, 0xc7, 0x84, 0xc2, 0x82, 0xa2, 0x56, 0x5f, 0x0c, 0x7a, 0x89, 0xef,
	0x0e, 0xa6, 0x58, 0x90, 0x7c, 0x01, 0x9d, 0x85, 0xd6, 0x55, 0x16, 0xb5, 0xfb, 0x62, 0x10, 0x9c,
	0x9f, 0x34, 0x42, 0x5b, 0xae, 0xe2, 0xf7, 0xae, 0x25, 0xd9, 0x74, 0x1e, 0x1b, 0xe8, 0x30, 0x76,
	0x8b, 0x8d, 0xc5, 0xca, 0xa6, 0x4b, 0xb4, 0x6c, 0xd4, 0x4b, 0x7c, 0x3e, 0xf8, 0x84, 0x76, 0x8b,
	0xd4, 0x8a, 0x0d, 0xdd, 0x91, 0x5a, 0x39, 0xaf, 0xa4, 0x32, 0x9e, 0x6b, 0x31, 0xd5, 0x25, 0x95,
	0xb9, 0xa9, 0x9a, 0xd0, 0x8a, 0x0d, 0xdd, 0x12, 0x5a, 0x9d, 0xfe, 0xf2, 0xc0, 0x77, 0x8e, 0x26,
	0xea, 0x9b, 0x7e, 0xf0, 0x90, 0xe2, 0xfb, 0x21, 0x45, 0xf7, 0x43, 0x72, 0x96, 0x76, 0x20, 0xa1,
	0x9f, 0x1e, 0xec, 0x8f, 0xc9, 0x3a, 0x47, 0xff, 0x23, 0xa0, 0xe7, 0x75, 0x06, 0x2d, 0xce, 0xe0,
	0xa8, 0xd9, 0x75, 0x2b, 0xba, 0x03, 0x11, 0xfc, 0x11, 0x00, 0x7c, 0x6d, 0x2d, 0xda, 0x95, 0x79,
	0xf0, 0x6b, 0xf2, 0x1a, 0x02, 0xc3, 0xca, 0x1b, 0xb9, 0x36, 0xcb, 0x3d, 0x69, 0xe4, 0x36, 0xb6,
	0x58, 0x0d, 0xcc, 0x5d, 0x2d, 0x25, 0xb4, 0x0b, 0xb2, 0x18, 0x75, 0x78, 0x1d, 0xd7, 0x67, 0xcf,
	0x00, 0x9a, 0x6e, 0xd9, 0x05, 0xef, 0xf2, 0x63, 0xb8, 0xe7, 0xfe, 0xd3, 0x71, 0x28, 0xce, 0x26,
	0x10, 0x6c, 0x59, 0x97, 0x07, 0xe0, 0xcf, 0x2e, 0xe6, 0xe9, 0x28, 0xb9, 0x18, 0x85, 0x7b, 0x0e,
	0x8d, 0x6b, 0x24, 0xe4, 0x21, 0xf4, 0x5c, 0x95, 0x4e, 0xa6, 0x1f, 0x2e, 0x43, 0x4f, 0x3e, 0x86,
	0x80, 0xe1, 0x6c, 0x3e, 0x9a, 0x7f, 0x99, 0x85, 0xad, 0x77, 0x6f, 0xbe, 0xbe, 0xba, 0xce, 0xed,
	0xcd, 0xea, 0x2a, 0x5e, 0xe8, 0x62, 0x68, 0xd6, 0x8a, 0x2a, 0xfa, 0x51, 0xff, 0x53, 0x5c, 0x96,
	0x37, 0x38, 0xc4, 0x32, 0x1f, 0xfe, 0xf3, 0x98, 0x5d, 0x75, 0xf9, 0x21, 0x7b, 0xf9, 0x37, 0x00,
	0x00, 0xff, 0xff, 0xa4, 0xce, 0xeb, 0x19, 0xe6, 0x04, 0x00, 0x00,
}

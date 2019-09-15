// Code generated by protoc-gen-go. DO NOT EDIT.
// source: simulation/participant/participant.proto

package participant

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
	StatusType_OK   StatusType = 0
	StatusType_NG   StatusType = 1
	StatusType_NONE StatusType = 2
)

var StatusType_name = map[int32]string{
	0: "OK",
	1: "NG",
	2: "NONE",
}

var StatusType_value = map[string]int32{
	"OK":   0,
	"NG":   1,
	"NONE": 2,
}

func (x StatusType) String() string {
	return proto.EnumName(StatusType_name, int32(x))
}

func (StatusType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_6631733b0a9fa50e, []int{0}
}

type AgentType int32

const (
	AgentType_PEDESTRIAN AgentType = 0
	AgentType_CAR        AgentType = 1
	AgentType_TRAIN      AgentType = 2
	AgentType_BICYCLE    AgentType = 3
)

var AgentType_name = map[int32]string{
	0: "PEDESTRIAN",
	1: "CAR",
	2: "TRAIN",
	3: "BICYCLE",
}

var AgentType_value = map[string]int32{
	"PEDESTRIAN": 0,
	"CAR":        1,
	"TRAIN":      2,
	"BICYCLE":    3,
}

func (x AgentType) String() string {
	return proto.EnumName(AgentType_name, int32(x))
}

func (AgentType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_6631733b0a9fa50e, []int{1}
}

type ClientType int32

const (
	ClientType_AREA    ClientType = 0
	ClientType_LOG     ClientType = 1
	ClientType_CARAREA ClientType = 2
	ClientType_PEDAREA ClientType = 3
)

var ClientType_name = map[int32]string{
	0: "AREA",
	1: "LOG",
	2: "CARAREA",
	3: "PEDAREA",
}

var ClientType_value = map[string]int32{
	"AREA":    0,
	"LOG":     1,
	"CARAREA": 2,
	"PEDAREA": 3,
}

func (x ClientType) String() string {
	return proto.EnumName(ClientType_name, int32(x))
}

func (ClientType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_6631733b0a9fa50e, []int{2}
}

type DemandType int32

const (
	DemandType_GET DemandType = 0
)

var DemandType_name = map[int32]string{
	0: "GET",
}

var DemandType_value = map[string]int32{
	"GET": 0,
}

func (x DemandType) String() string {
	return proto.EnumName(DemandType_name, int32(x))
}

func (DemandType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_6631733b0a9fa50e, []int{3}
}

type SupplyType int32

const (
	SupplyType_RES_GET SupplyType = 0
)

var SupplyType_name = map[int32]string{
	0: "RES_GET",
}

var SupplyType_value = map[string]int32{
	"RES_GET": 0,
}

func (x SupplyType) String() string {
	return proto.EnumName(SupplyType_name, int32(x))
}

func (SupplyType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_6631733b0a9fa50e, []int{4}
}

type ParticipantInfo struct {
	// participant info
	ClientParticipantId uint64     `protobuf:"varint,1,opt,name=client_participant_id,json=clientParticipantId,proto3" json:"client_participant_id,omitempty"`
	ClientAreaId        uint64     `protobuf:"varint,2,opt,name=client_area_id,json=clientAreaId,proto3" json:"client_area_id,omitempty"`
	ClientAgentId       uint64     `protobuf:"varint,3,opt,name=client_agent_id,json=clientAgentId,proto3" json:"client_agent_id,omitempty"`
	ClientClockId       uint64     `protobuf:"varint,4,opt,name=client_clock_id,json=clientClockId,proto3" json:"client_clock_id,omitempty"`
	ClientType          ClientType `protobuf:"varint,5,opt,name=client_type,json=clientType,proto3,enum=api.participant.ClientType" json:"client_type,omitempty"`
	AreaId              uint32     `protobuf:"varint,6,opt,name=area_id,json=areaId,proto3" json:"area_id,omitempty"`
	AgentType           AgentType  `protobuf:"varint,7,opt,name=agent_type,json=agentType,proto3,enum=api.participant.AgentType" json:"agent_type,omitempty"`
	SupplyType          SupplyType `protobuf:"varint,8,opt,name=supply_type,json=supplyType,proto3,enum=api.participant.SupplyType" json:"supply_type,omitempty"`
	// meta data
	StatusType           StatusType `protobuf:"varint,9,opt,name=status_type,json=statusType,proto3,enum=api.participant.StatusType" json:"status_type,omitempty"`
	Meta                 string     `protobuf:"bytes,10,opt,name=meta,proto3" json:"meta,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *ParticipantInfo) Reset()         { *m = ParticipantInfo{} }
func (m *ParticipantInfo) String() string { return proto.CompactTextString(m) }
func (*ParticipantInfo) ProtoMessage()    {}
func (*ParticipantInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_6631733b0a9fa50e, []int{0}
}

func (m *ParticipantInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ParticipantInfo.Unmarshal(m, b)
}
func (m *ParticipantInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ParticipantInfo.Marshal(b, m, deterministic)
}
func (m *ParticipantInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ParticipantInfo.Merge(m, src)
}
func (m *ParticipantInfo) XXX_Size() int {
	return xxx_messageInfo_ParticipantInfo.Size(m)
}
func (m *ParticipantInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_ParticipantInfo.DiscardUnknown(m)
}

var xxx_messageInfo_ParticipantInfo proto.InternalMessageInfo

func (m *ParticipantInfo) GetClientParticipantId() uint64 {
	if m != nil {
		return m.ClientParticipantId
	}
	return 0
}

func (m *ParticipantInfo) GetClientAreaId() uint64 {
	if m != nil {
		return m.ClientAreaId
	}
	return 0
}

func (m *ParticipantInfo) GetClientAgentId() uint64 {
	if m != nil {
		return m.ClientAgentId
	}
	return 0
}

func (m *ParticipantInfo) GetClientClockId() uint64 {
	if m != nil {
		return m.ClientClockId
	}
	return 0
}

func (m *ParticipantInfo) GetClientType() ClientType {
	if m != nil {
		return m.ClientType
	}
	return ClientType_AREA
}

func (m *ParticipantInfo) GetAreaId() uint32 {
	if m != nil {
		return m.AreaId
	}
	return 0
}

func (m *ParticipantInfo) GetAgentType() AgentType {
	if m != nil {
		return m.AgentType
	}
	return AgentType_PEDESTRIAN
}

func (m *ParticipantInfo) GetSupplyType() SupplyType {
	if m != nil {
		return m.SupplyType
	}
	return SupplyType_RES_GET
}

func (m *ParticipantInfo) GetStatusType() StatusType {
	if m != nil {
		return m.StatusType
	}
	return StatusType_OK
}

func (m *ParticipantInfo) GetMeta() string {
	if m != nil {
		return m.Meta
	}
	return ""
}

type ParticipantDemand struct {
	// demand info
	ClientId uint64 `protobuf:"varint,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	// demand type
	DemandType DemandType `protobuf:"varint,2,opt,name=demand_type,json=demandType,proto3,enum=api.participant.DemandType" json:"demand_type,omitempty"`
	// meta data
	StatusType           StatusType `protobuf:"varint,3,opt,name=status_type,json=statusType,proto3,enum=api.participant.StatusType" json:"status_type,omitempty"`
	Meta                 string     `protobuf:"bytes,4,opt,name=meta,proto3" json:"meta,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *ParticipantDemand) Reset()         { *m = ParticipantDemand{} }
func (m *ParticipantDemand) String() string { return proto.CompactTextString(m) }
func (*ParticipantDemand) ProtoMessage()    {}
func (*ParticipantDemand) Descriptor() ([]byte, []int) {
	return fileDescriptor_6631733b0a9fa50e, []int{1}
}

func (m *ParticipantDemand) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ParticipantDemand.Unmarshal(m, b)
}
func (m *ParticipantDemand) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ParticipantDemand.Marshal(b, m, deterministic)
}
func (m *ParticipantDemand) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ParticipantDemand.Merge(m, src)
}
func (m *ParticipantDemand) XXX_Size() int {
	return xxx_messageInfo_ParticipantDemand.Size(m)
}
func (m *ParticipantDemand) XXX_DiscardUnknown() {
	xxx_messageInfo_ParticipantDemand.DiscardUnknown(m)
}

var xxx_messageInfo_ParticipantDemand proto.InternalMessageInfo

func (m *ParticipantDemand) GetClientId() uint64 {
	if m != nil {
		return m.ClientId
	}
	return 0
}

func (m *ParticipantDemand) GetDemandType() DemandType {
	if m != nil {
		return m.DemandType
	}
	return DemandType_GET
}

func (m *ParticipantDemand) GetStatusType() StatusType {
	if m != nil {
		return m.StatusType
	}
	return StatusType_OK
}

func (m *ParticipantDemand) GetMeta() string {
	if m != nil {
		return m.Meta
	}
	return ""
}

func init() {
	proto.RegisterEnum("api.participant.StatusType", StatusType_name, StatusType_value)
	proto.RegisterEnum("api.participant.AgentType", AgentType_name, AgentType_value)
	proto.RegisterEnum("api.participant.ClientType", ClientType_name, ClientType_value)
	proto.RegisterEnum("api.participant.DemandType", DemandType_name, DemandType_value)
	proto.RegisterEnum("api.participant.SupplyType", SupplyType_name, SupplyType_value)
	proto.RegisterType((*ParticipantInfo)(nil), "api.participant.ParticipantInfo")
	proto.RegisterType((*ParticipantDemand)(nil), "api.participant.ParticipantDemand")
}

func init() {
	proto.RegisterFile("simulation/participant/participant.proto", fileDescriptor_6631733b0a9fa50e)
}

var fileDescriptor_6631733b0a9fa50e = []byte{
	// 499 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x93, 0x41, 0x6f, 0x9b, 0x30,
	0x14, 0xc7, 0x43, 0x48, 0x93, 0xf0, 0xb2, 0x26, 0xcc, 0x53, 0xb5, 0x6c, 0xbd, 0x44, 0xd5, 0x54,
	0x45, 0x1c, 0x88, 0xd4, 0x9d, 0xaa, 0x75, 0x07, 0x4a, 0x50, 0x85, 0x56, 0x91, 0xc8, 0xc9, 0x65,
	0xbb, 0x44, 0x2e, 0x78, 0x2d, 0x1a, 0x01, 0x0b, 0x1c, 0x69, 0xf9, 0x76, 0xfb, 0x22, 0xfb, 0x2e,
	0x93, 0x6d, 0x02, 0xac, 0x5b, 0x2e, 0x3d, 0xf9, 0xf9, 0xf9, 0xf7, 0x7f, 0xfc, 0x9f, 0x1f, 0x86,
	0x69, 0x11, 0x6f, 0x77, 0x09, 0xe1, 0x71, 0x96, 0xce, 0x18, 0xc9, 0x79, 0x1c, 0xc6, 0x8c, 0xa4,
	0xbc, 0x19, 0xdb, 0x2c, 0xcf, 0x78, 0x86, 0x46, 0x84, 0xc5, 0x76, 0x23, 0x7d, 0xf1, 0x5b, 0x87,
	0xd1, 0xb2, 0xde, 0xfb, 0xe9, 0xf7, 0x0c, 0x5d, 0xc1, 0x59, 0x98, 0xc4, 0x34, 0xe5, 0x9b, 0x06,
	0xb9, 0x89, 0xa3, 0xb1, 0x36, 0xd1, 0xa6, 0x1d, 0xfc, 0x46, 0x1d, 0x36, 0x55, 0x11, 0xfa, 0x00,
	0xc3, 0x52, 0x43, 0x72, 0x4a, 0x04, 0xdc, 0x96, 0xf0, 0x2b, 0x95, 0x75, 0x72, 0x4a, 0xfc, 0x08,
	0x5d, 0xc2, 0xe8, 0x40, 0x3d, 0x52, 0x55, 0x53, 0x97, 0xd8, 0x69, 0x89, 0x89, 0xec, 0x5f, 0x5c,
	0x98, 0x64, 0xe1, 0x0f, 0xc1, 0x75, 0x9a, 0x9c, 0x2b, 0xb2, 0x7e, 0x84, 0x6e, 0x60, 0x50, 0x72,
	0x7c, 0xcf, 0xe8, 0xf8, 0x64, 0xa2, 0x4d, 0x87, 0x57, 0xe7, 0xf6, 0xb3, 0x26, 0x6d, 0x57, 0x32,
	0xeb, 0x3d, 0xa3, 0x18, 0xc2, 0x2a, 0x46, 0x6f, 0xa1, 0x77, 0x30, 0xdb, 0x9d, 0x68, 0xd3, 0x53,
	0xdc, 0x25, 0xca, 0xe6, 0x35, 0x80, 0xf2, 0x27, 0xab, 0xf6, 0x64, 0xd5, 0xf7, 0xff, 0x54, 0x95,
	0x66, 0x65, 0x51, 0x83, 0x1c, 0x42, 0xe1, 0xa8, 0xd8, 0x31, 0x96, 0xec, 0x95, 0xb6, 0x7f, 0xc4,
	0xd1, 0x4a, 0x32, 0xca, 0x51, 0x51, 0xc5, 0x52, 0xcd, 0x09, 0xdf, 0x15, 0x4a, 0x6d, 0x1c, 0x53,
	0x4b, 0xa6, 0x54, 0x57, 0x31, 0x42, 0xd0, 0xd9, 0x52, 0x4e, 0xc6, 0x30, 0xd1, 0xa6, 0x06, 0x96,
	0xf1, 0xc5, 0x2f, 0x0d, 0x5e, 0x37, 0x26, 0x35, 0xa7, 0x5b, 0x92, 0x46, 0xe8, 0x1c, 0x8c, 0xf2,
	0xde, 0xaa, 0xa9, 0xf6, 0x55, 0x42, 0x5d, 0x6a, 0x24, 0x31, 0x65, 0xa2, 0x7d, 0xc4, 0x84, 0x2a,
	0xa5, 0x4c, 0x44, 0x55, 0xfc, 0xbc, 0x05, 0xfd, 0x65, 0x2d, 0x74, 0xea, 0x16, 0xac, 0x4b, 0x80,
	0x9a, 0x46, 0x5d, 0x68, 0x2f, 0xbe, 0x98, 0x2d, 0xb1, 0x06, 0x77, 0xa6, 0x86, 0xfa, 0xd0, 0x09,
	0x16, 0x81, 0x67, 0xb6, 0xad, 0x1b, 0x30, 0xaa, 0x91, 0xa0, 0x21, 0xc0, 0xd2, 0x9b, 0x7b, 0xab,
	0x35, 0xf6, 0x9d, 0xc0, 0x6c, 0xa1, 0x1e, 0xe8, 0xae, 0x83, 0x4d, 0x0d, 0x19, 0x70, 0xb2, 0xc6,
	0x8e, 0x1f, 0x98, 0x6d, 0x34, 0x80, 0xde, 0xad, 0xef, 0x7e, 0x75, 0xef, 0x3d, 0x53, 0xb7, 0xae,
	0x01, 0xea, 0xdf, 0x44, 0x54, 0x75, 0xb0, 0xe7, 0x28, 0xe1, 0xfd, 0x42, 0x7c, 0x68, 0x00, 0x3d,
	0xd7, 0xc1, 0x32, 0x2b, 0xa5, 0x4b, 0x6f, 0x2e, 0x37, 0xba, 0x75, 0x06, 0x50, 0x5f, 0x86, 0x10,
	0xdc, 0x79, 0x6b, 0xb3, 0x65, 0xbd, 0x03, 0xa8, 0xc7, 0x2c, 0x14, 0xd8, 0x5b, 0x6d, 0xe4, 0xd1,
	0xed, 0xe7, 0x6f, 0x9f, 0x1e, 0x63, 0xfe, 0xb4, 0x7b, 0xb0, 0xc3, 0x6c, 0x3b, 0x2b, 0xf6, 0x29,
	0xcd, 0xe9, 0xcf, 0xc3, 0xba, 0x21, 0x09, 0x7b, 0x22, 0x33, 0xc2, 0xe2, 0xd9, 0xff, 0xdf, 0xf5,
	0x43, 0x57, 0x3e, 0xe6, 0x8f, 0x7f, 0x02, 0x00, 0x00, 0xff, 0xff, 0xb9, 0x20, 0x66, 0x8d, 0xf8,
	0x03, 0x00, 0x00,
}

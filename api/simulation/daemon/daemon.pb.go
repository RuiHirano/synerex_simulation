// Code generated by protoc-gen-go. DO NOT EDIT.
// source: simulation/daemon/daemon.proto

package daemon

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	agent "github.com/synerex/synerex_alpha/api/simulation/agent"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type OrderType int32

const (
	OrderType_SET_AGENTS   OrderType = 0
	OrderType_CLEAR_AGENTS OrderType = 1
	OrderType_SET_CLOCK    OrderType = 2
	OrderType_START_CLOCK  OrderType = 3
	OrderType_STOP_CLOCK   OrderType = 4
)

var OrderType_name = map[int32]string{
	0: "SET_AGENTS",
	1: "CLEAR_AGENTS",
	2: "SET_CLOCK",
	3: "START_CLOCK",
	4: "STOP_CLOCK",
}

var OrderType_value = map[string]int32{
	"SET_AGENTS":   0,
	"CLEAR_AGENTS": 1,
	"SET_CLOCK":    2,
	"START_CLOCK":  3,
	"STOP_CLOCK":   4,
}

func (x OrderType) String() string {
	return proto.EnumName(OrderType_name, int32(x))
}

func (OrderType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_cf6adfdc580b8b5d, []int{0}
}

type Response struct {
	Ok                   bool     `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
	Err                  string   `protobuf:"bytes,2,opt,name=err,proto3" json:"err,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_cf6adfdc580b8b5d, []int{0}
}

func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (m *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(m, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetOk() bool {
	if m != nil {
		return m.Ok
	}
	return false
}

func (m *Response) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

type OrderMessage struct {
	OrderType OrderType `protobuf:"varint,1,opt,name=order_type,json=orderType,proto3,enum=api.daemon.OrderType" json:"order_type,omitempty"`
	// Types that are valid to be assigned to Message:
	//	*OrderMessage_SetAgentsMessage
	//	*OrderMessage_ClearAgentsMessage
	//	*OrderMessage_SetClockMessage
	//	*OrderMessage_StartClockMessage
	//	*OrderMessage_StopClockMessage
	Message              isOrderMessage_Message `protobuf_oneof:"message"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *OrderMessage) Reset()         { *m = OrderMessage{} }
func (m *OrderMessage) String() string { return proto.CompactTextString(m) }
func (*OrderMessage) ProtoMessage()    {}
func (*OrderMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_cf6adfdc580b8b5d, []int{1}
}

func (m *OrderMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OrderMessage.Unmarshal(m, b)
}
func (m *OrderMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OrderMessage.Marshal(b, m, deterministic)
}
func (m *OrderMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrderMessage.Merge(m, src)
}
func (m *OrderMessage) XXX_Size() int {
	return xxx_messageInfo_OrderMessage.Size(m)
}
func (m *OrderMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_OrderMessage.DiscardUnknown(m)
}

var xxx_messageInfo_OrderMessage proto.InternalMessageInfo

func (m *OrderMessage) GetOrderType() OrderType {
	if m != nil {
		return m.OrderType
	}
	return OrderType_SET_AGENTS
}

type isOrderMessage_Message interface {
	isOrderMessage_Message()
}

type OrderMessage_SetAgentsMessage struct {
	SetAgentsMessage *SetAgentsMessage `protobuf:"bytes,2,opt,name=set_agents_message,json=setAgentsMessage,proto3,oneof"`
}

type OrderMessage_ClearAgentsMessage struct {
	ClearAgentsMessage *ClearAgentsMessage `protobuf:"bytes,3,opt,name=clear_agents_message,json=clearAgentsMessage,proto3,oneof"`
}

type OrderMessage_SetClockMessage struct {
	SetClockMessage *SetClockMessage `protobuf:"bytes,4,opt,name=set_clock_message,json=setClockMessage,proto3,oneof"`
}

type OrderMessage_StartClockMessage struct {
	StartClockMessage *StartClockMessage `protobuf:"bytes,5,opt,name=start_clock_message,json=startClockMessage,proto3,oneof"`
}

type OrderMessage_StopClockMessage struct {
	StopClockMessage *StopClockMessage `protobuf:"bytes,6,opt,name=stop_clock_message,json=stopClockMessage,proto3,oneof"`
}

func (*OrderMessage_SetAgentsMessage) isOrderMessage_Message() {}

func (*OrderMessage_ClearAgentsMessage) isOrderMessage_Message() {}

func (*OrderMessage_SetClockMessage) isOrderMessage_Message() {}

func (*OrderMessage_StartClockMessage) isOrderMessage_Message() {}

func (*OrderMessage_StopClockMessage) isOrderMessage_Message() {}

func (m *OrderMessage) GetMessage() isOrderMessage_Message {
	if m != nil {
		return m.Message
	}
	return nil
}

func (m *OrderMessage) GetSetAgentsMessage() *SetAgentsMessage {
	if x, ok := m.GetMessage().(*OrderMessage_SetAgentsMessage); ok {
		return x.SetAgentsMessage
	}
	return nil
}

func (m *OrderMessage) GetClearAgentsMessage() *ClearAgentsMessage {
	if x, ok := m.GetMessage().(*OrderMessage_ClearAgentsMessage); ok {
		return x.ClearAgentsMessage
	}
	return nil
}

func (m *OrderMessage) GetSetClockMessage() *SetClockMessage {
	if x, ok := m.GetMessage().(*OrderMessage_SetClockMessage); ok {
		return x.SetClockMessage
	}
	return nil
}

func (m *OrderMessage) GetStartClockMessage() *StartClockMessage {
	if x, ok := m.GetMessage().(*OrderMessage_StartClockMessage); ok {
		return x.StartClockMessage
	}
	return nil
}

func (m *OrderMessage) GetStopClockMessage() *StopClockMessage {
	if x, ok := m.GetMessage().(*OrderMessage_StopClockMessage); ok {
		return x.StopClockMessage
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*OrderMessage) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*OrderMessage_SetAgentsMessage)(nil),
		(*OrderMessage_ClearAgentsMessage)(nil),
		(*OrderMessage_SetClockMessage)(nil),
		(*OrderMessage_StartClockMessage)(nil),
		(*OrderMessage_StopClockMessage)(nil),
	}
}

type SetAgentsMessage struct {
	Agents               []*agent.Agent `protobuf:"bytes,1,rep,name=agents,proto3" json:"agents,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *SetAgentsMessage) Reset()         { *m = SetAgentsMessage{} }
func (m *SetAgentsMessage) String() string { return proto.CompactTextString(m) }
func (*SetAgentsMessage) ProtoMessage()    {}
func (*SetAgentsMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_cf6adfdc580b8b5d, []int{2}
}

func (m *SetAgentsMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetAgentsMessage.Unmarshal(m, b)
}
func (m *SetAgentsMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetAgentsMessage.Marshal(b, m, deterministic)
}
func (m *SetAgentsMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetAgentsMessage.Merge(m, src)
}
func (m *SetAgentsMessage) XXX_Size() int {
	return xxx_messageInfo_SetAgentsMessage.Size(m)
}
func (m *SetAgentsMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_SetAgentsMessage.DiscardUnknown(m)
}

var xxx_messageInfo_SetAgentsMessage proto.InternalMessageInfo

func (m *SetAgentsMessage) GetAgents() []*agent.Agent {
	if m != nil {
		return m.Agents
	}
	return nil
}

type ClearAgentsMessage struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ClearAgentsMessage) Reset()         { *m = ClearAgentsMessage{} }
func (m *ClearAgentsMessage) String() string { return proto.CompactTextString(m) }
func (*ClearAgentsMessage) ProtoMessage()    {}
func (*ClearAgentsMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_cf6adfdc580b8b5d, []int{3}
}

func (m *ClearAgentsMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ClearAgentsMessage.Unmarshal(m, b)
}
func (m *ClearAgentsMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ClearAgentsMessage.Marshal(b, m, deterministic)
}
func (m *ClearAgentsMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClearAgentsMessage.Merge(m, src)
}
func (m *ClearAgentsMessage) XXX_Size() int {
	return xxx_messageInfo_ClearAgentsMessage.Size(m)
}
func (m *ClearAgentsMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_ClearAgentsMessage.DiscardUnknown(m)
}

var xxx_messageInfo_ClearAgentsMessage proto.InternalMessageInfo

type SetClockMessage struct {
	GlobalTime           float64  `protobuf:"fixed64,1,opt,name=global_time,json=globalTime,proto3" json:"global_time,omitempty"`
	TimeStep             float64  `protobuf:"fixed64,2,opt,name=time_step,json=timeStep,proto3" json:"time_step,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SetClockMessage) Reset()         { *m = SetClockMessage{} }
func (m *SetClockMessage) String() string { return proto.CompactTextString(m) }
func (*SetClockMessage) ProtoMessage()    {}
func (*SetClockMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_cf6adfdc580b8b5d, []int{4}
}

func (m *SetClockMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetClockMessage.Unmarshal(m, b)
}
func (m *SetClockMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetClockMessage.Marshal(b, m, deterministic)
}
func (m *SetClockMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetClockMessage.Merge(m, src)
}
func (m *SetClockMessage) XXX_Size() int {
	return xxx_messageInfo_SetClockMessage.Size(m)
}
func (m *SetClockMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_SetClockMessage.DiscardUnknown(m)
}

var xxx_messageInfo_SetClockMessage proto.InternalMessageInfo

func (m *SetClockMessage) GetGlobalTime() float64 {
	if m != nil {
		return m.GlobalTime
	}
	return 0
}

func (m *SetClockMessage) GetTimeStep() float64 {
	if m != nil {
		return m.TimeStep
	}
	return 0
}

type StartClockMessage struct {
	StepNum              uint64   `protobuf:"varint,1,opt,name=step_num,json=stepNum,proto3" json:"step_num,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StartClockMessage) Reset()         { *m = StartClockMessage{} }
func (m *StartClockMessage) String() string { return proto.CompactTextString(m) }
func (*StartClockMessage) ProtoMessage()    {}
func (*StartClockMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_cf6adfdc580b8b5d, []int{5}
}

func (m *StartClockMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StartClockMessage.Unmarshal(m, b)
}
func (m *StartClockMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StartClockMessage.Marshal(b, m, deterministic)
}
func (m *StartClockMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StartClockMessage.Merge(m, src)
}
func (m *StartClockMessage) XXX_Size() int {
	return xxx_messageInfo_StartClockMessage.Size(m)
}
func (m *StartClockMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_StartClockMessage.DiscardUnknown(m)
}

var xxx_messageInfo_StartClockMessage proto.InternalMessageInfo

func (m *StartClockMessage) GetStepNum() uint64 {
	if m != nil {
		return m.StepNum
	}
	return 0
}

type StopClockMessage struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StopClockMessage) Reset()         { *m = StopClockMessage{} }
func (m *StopClockMessage) String() string { return proto.CompactTextString(m) }
func (*StopClockMessage) ProtoMessage()    {}
func (*StopClockMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_cf6adfdc580b8b5d, []int{6}
}

func (m *StopClockMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StopClockMessage.Unmarshal(m, b)
}
func (m *StopClockMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StopClockMessage.Marshal(b, m, deterministic)
}
func (m *StopClockMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StopClockMessage.Merge(m, src)
}
func (m *StopClockMessage) XXX_Size() int {
	return xxx_messageInfo_StopClockMessage.Size(m)
}
func (m *StopClockMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_StopClockMessage.DiscardUnknown(m)
}

var xxx_messageInfo_StopClockMessage proto.InternalMessageInfo

func init() {
	proto.RegisterEnum("api.daemon.OrderType", OrderType_name, OrderType_value)
	proto.RegisterType((*Response)(nil), "api.daemon.Response")
	proto.RegisterType((*OrderMessage)(nil), "api.daemon.OrderMessage")
	proto.RegisterType((*SetAgentsMessage)(nil), "api.daemon.SetAgentsMessage")
	proto.RegisterType((*ClearAgentsMessage)(nil), "api.daemon.ClearAgentsMessage")
	proto.RegisterType((*SetClockMessage)(nil), "api.daemon.SetClockMessage")
	proto.RegisterType((*StartClockMessage)(nil), "api.daemon.StartClockMessage")
	proto.RegisterType((*StopClockMessage)(nil), "api.daemon.StopClockMessage")
}

func init() { proto.RegisterFile("simulation/daemon/daemon.proto", fileDescriptor_cf6adfdc580b8b5d) }

var fileDescriptor_cf6adfdc580b8b5d = []byte{
	// 540 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x94, 0x51, 0x6f, 0xda, 0x3e,
	0x14, 0xc5, 0x09, 0xf0, 0xa7, 0xe4, 0xd2, 0x82, 0xb9, 0x7f, 0x26, 0xb1, 0xb6, 0xeb, 0x50, 0x9e,
	0xd0, 0x34, 0x31, 0x89, 0xed, 0x71, 0x2f, 0x94, 0xb2, 0x75, 0x1b, 0x2b, 0x93, 0x93, 0xa7, 0x49,
	0x53, 0x94, 0x52, 0x0b, 0x45, 0x10, 0x1c, 0xc5, 0xe6, 0xa1, 0x9f, 0x75, 0x1f, 0x64, 0xaf, 0x93,
	0x6d, 0x02, 0xd4, 0x14, 0xf6, 0x02, 0xf6, 0xf1, 0xf1, 0xef, 0x5a, 0xe7, 0x5e, 0x80, 0x2b, 0x11,
	0x27, 0xab, 0x45, 0x24, 0x63, 0xbe, 0x7c, 0xf7, 0x10, 0xb1, 0x64, 0xf3, 0xd5, 0x4b, 0x33, 0x2e,
	0x39, 0x42, 0x94, 0xc6, 0x3d, 0xa3, 0x9c, 0x5f, 0xee, 0x78, 0xa3, 0x19, 0x5b, 0x4a, 0xf3, 0x69,
	0x9c, 0xde, 0x5b, 0xa8, 0x52, 0x26, 0x52, 0xbe, 0x14, 0x0c, 0xeb, 0x50, 0xe4, 0xf3, 0xb6, 0xd3,
	0x71, 0xba, 0x55, 0x5a, 0xe4, 0x73, 0x24, 0x50, 0x62, 0x59, 0xd6, 0x2e, 0x76, 0x9c, 0xae, 0x4b,
	0xd5, 0xd2, 0xfb, 0x5d, 0x82, 0xd3, 0x49, 0xf6, 0xc0, 0xb2, 0xef, 0x4c, 0x88, 0x68, 0xc6, 0xf0,
	0x03, 0x00, 0x57, 0xfb, 0x50, 0x3e, 0xa6, 0x4c, 0x5f, 0xad, 0xf7, 0x5f, 0xf4, 0xb6, 0xd5, 0x7b,
	0xda, 0x1d, 0x3c, 0xa6, 0x8c, 0xba, 0x3c, 0x5f, 0xe2, 0x18, 0x50, 0x30, 0x19, 0xea, 0x77, 0x88,
	0x30, 0x31, 0x2c, 0x5d, 0xa7, 0xd6, 0xbf, 0xdc, 0xbd, 0xed, 0x33, 0x39, 0xd0, 0xa6, 0x75, 0xbd,
	0xdb, 0x02, 0x25, 0xc2, 0xd2, 0x90, 0x42, 0x6b, 0xba, 0x60, 0x51, 0x66, 0xf3, 0x4a, 0x9a, 0x77,
	0xb5, 0xcb, 0x1b, 0x2a, 0x9f, 0x4d, 0xc4, 0xe9, 0x9e, 0x8a, 0x5f, 0xa0, 0xa9, 0x5e, 0x38, 0x5d,
	0xf0, 0xe9, 0x7c, 0x03, 0x2c, 0x6b, 0xe0, 0x85, 0xf5, 0xc0, 0xa1, 0xf2, 0x6c, 0x69, 0x0d, 0xf1,
	0x54, 0xc2, 0x09, 0xfc, 0x2f, 0x64, 0x94, 0xd9, 0xb0, 0xff, 0x34, 0xec, 0xd5, 0x13, 0x98, 0xb2,
	0x59, 0xb8, 0xa6, 0xb0, 0x45, 0x9d, 0x9e, 0xe4, 0xa9, 0xc5, 0xab, 0x3c, 0x93, 0x9e, 0xe4, 0xa9,
	0x85, 0x23, 0xc2, 0xd2, 0xae, 0x5d, 0x38, 0x59, 0x23, 0xbc, 0x8f, 0x40, 0xec, 0xc0, 0xb1, 0x0b,
	0x15, 0x13, 0x6b, 0xdb, 0xe9, 0x94, 0xba, 0xb5, 0x3e, 0xd1, 0x05, 0xcc, 0x04, 0x69, 0x27, 0x5d,
	0x9f, 0x7b, 0x2d, 0xc0, 0xfd, 0x78, 0xbd, 0x09, 0x34, 0xac, 0x8c, 0xf0, 0x35, 0xd4, 0x66, 0x0b,
	0x7e, 0x1f, 0x2d, 0x42, 0x19, 0x27, 0x66, 0x68, 0x1c, 0x0a, 0x46, 0x0a, 0xe2, 0x84, 0xe1, 0x05,
	0xb8, 0xea, 0x24, 0x14, 0x92, 0xa5, 0x7a, 0x2a, 0x1c, 0x5a, 0x55, 0x82, 0x2f, 0x59, 0xea, 0xf5,
	0xa0, 0xb9, 0x97, 0x13, 0xbe, 0x84, 0xaa, 0x32, 0x87, 0xcb, 0x55, 0xa2, 0x79, 0x65, 0x7a, 0xa2,
	0xf6, 0x77, 0xab, 0xc4, 0x43, 0x20, 0x76, 0x0e, 0x6f, 0x7e, 0x81, 0xbb, 0x99, 0x4b, 0xac, 0x03,
	0xf8, 0xa3, 0x20, 0x1c, 0x7c, 0x1e, 0xdd, 0x05, 0x3e, 0x29, 0x20, 0x81, 0xd3, 0xe1, 0x78, 0x34,
	0xa0, 0xb9, 0xe2, 0xe0, 0x19, 0xb8, 0xca, 0x31, 0x1c, 0x4f, 0x86, 0xdf, 0x48, 0x11, 0x1b, 0x50,
	0xf3, 0x83, 0x01, 0xcd, 0x85, 0x92, 0x26, 0x04, 0x93, 0x1f, 0xeb, 0x7d, 0xb9, 0xff, 0xa7, 0x08,
	0xae, 0x1f, 0x27, 0x37, 0xba, 0x0b, 0xf8, 0x09, 0xea, 0x9b, 0x54, 0x75, 0x55, 0x3c, 0x3a, 0xe2,
	0xe7, 0xad, 0xdd, 0xd3, 0xfc, 0xb7, 0xe9, 0x15, 0xf0, 0x2b, 0x90, 0x9d, 0x7c, 0x0d, 0xe9, 0x1f,
	0xc3, 0x7d, 0x90, 0x75, 0x03, 0x67, 0x79, 0x57, 0x0c, 0xe8, 0xd8, 0x50, 0x1f, 0xa4, 0xdc, 0x42,
	0x63, 0xdb, 0x0a, 0xc3, 0x39, 0x3e, 0xcf, 0x07, 0x49, 0x2a, 0xa3, 0xbc, 0x49, 0xcf, 0x65, 0x64,
	0x35, 0xf0, 0x10, 0xe7, 0xba, 0xfa, 0xb3, 0x62, 0xc4, 0xfb, 0x8a, 0xfe, 0x7b, 0x7b, 0xff, 0x37,
	0x00, 0x00, 0xff, 0xff, 0x30, 0x9c, 0xcc, 0xd4, 0x2a, 0x05, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// SimDaemonClient is the client API for SimDaemon service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SimDaemonClient interface {
	SetAgentsOrder(ctx context.Context, in *SetAgentsMessage, opts ...grpc.CallOption) (*Response, error)
	ClearAgentsOrder(ctx context.Context, in *ClearAgentsMessage, opts ...grpc.CallOption) (*Response, error)
	SetClockOrder(ctx context.Context, in *SetClockMessage, opts ...grpc.CallOption) (*Response, error)
	StartClockOrder(ctx context.Context, in *StartClockMessage, opts ...grpc.CallOption) (*Response, error)
	StopClockOrder(ctx context.Context, in *StopClockMessage, opts ...grpc.CallOption) (*Response, error)
}

type simDaemonClient struct {
	cc *grpc.ClientConn
}

func NewSimDaemonClient(cc *grpc.ClientConn) SimDaemonClient {
	return &simDaemonClient{cc}
}

func (c *simDaemonClient) SetAgentsOrder(ctx context.Context, in *SetAgentsMessage, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/api.daemon.SimDaemon/SetAgentsOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *simDaemonClient) ClearAgentsOrder(ctx context.Context, in *ClearAgentsMessage, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/api.daemon.SimDaemon/ClearAgentsOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *simDaemonClient) SetClockOrder(ctx context.Context, in *SetClockMessage, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/api.daemon.SimDaemon/SetClockOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *simDaemonClient) StartClockOrder(ctx context.Context, in *StartClockMessage, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/api.daemon.SimDaemon/StartClockOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *simDaemonClient) StopClockOrder(ctx context.Context, in *StopClockMessage, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/api.daemon.SimDaemon/StopClockOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SimDaemonServer is the server API for SimDaemon service.
type SimDaemonServer interface {
	SetAgentsOrder(context.Context, *SetAgentsMessage) (*Response, error)
	ClearAgentsOrder(context.Context, *ClearAgentsMessage) (*Response, error)
	SetClockOrder(context.Context, *SetClockMessage) (*Response, error)
	StartClockOrder(context.Context, *StartClockMessage) (*Response, error)
	StopClockOrder(context.Context, *StopClockMessage) (*Response, error)
}

// UnimplementedSimDaemonServer can be embedded to have forward compatible implementations.
type UnimplementedSimDaemonServer struct {
}

func (*UnimplementedSimDaemonServer) SetAgentsOrder(ctx context.Context, req *SetAgentsMessage) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetAgentsOrder not implemented")
}
func (*UnimplementedSimDaemonServer) ClearAgentsOrder(ctx context.Context, req *ClearAgentsMessage) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClearAgentsOrder not implemented")
}
func (*UnimplementedSimDaemonServer) SetClockOrder(ctx context.Context, req *SetClockMessage) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetClockOrder not implemented")
}
func (*UnimplementedSimDaemonServer) StartClockOrder(ctx context.Context, req *StartClockMessage) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartClockOrder not implemented")
}
func (*UnimplementedSimDaemonServer) StopClockOrder(ctx context.Context, req *StopClockMessage) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopClockOrder not implemented")
}

func RegisterSimDaemonServer(s *grpc.Server, srv SimDaemonServer) {
	s.RegisterService(&_SimDaemon_serviceDesc, srv)
}

func _SimDaemon_SetAgentsOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetAgentsMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SimDaemonServer).SetAgentsOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.daemon.SimDaemon/SetAgentsOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SimDaemonServer).SetAgentsOrder(ctx, req.(*SetAgentsMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _SimDaemon_ClearAgentsOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClearAgentsMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SimDaemonServer).ClearAgentsOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.daemon.SimDaemon/ClearAgentsOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SimDaemonServer).ClearAgentsOrder(ctx, req.(*ClearAgentsMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _SimDaemon_SetClockOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetClockMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SimDaemonServer).SetClockOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.daemon.SimDaemon/SetClockOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SimDaemonServer).SetClockOrder(ctx, req.(*SetClockMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _SimDaemon_StartClockOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartClockMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SimDaemonServer).StartClockOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.daemon.SimDaemon/StartClockOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SimDaemonServer).StartClockOrder(ctx, req.(*StartClockMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _SimDaemon_StopClockOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StopClockMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SimDaemonServer).StopClockOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.daemon.SimDaemon/StopClockOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SimDaemonServer).StopClockOrder(ctx, req.(*StopClockMessage))
	}
	return interceptor(ctx, in, info, handler)
}

var _SimDaemon_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.daemon.SimDaemon",
	HandlerType: (*SimDaemonServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetAgentsOrder",
			Handler:    _SimDaemon_SetAgentsOrder_Handler,
		},
		{
			MethodName: "ClearAgentsOrder",
			Handler:    _SimDaemon_ClearAgentsOrder_Handler,
		},
		{
			MethodName: "SetClockOrder",
			Handler:    _SimDaemon_SetClockOrder_Handler,
		},
		{
			MethodName: "StartClockOrder",
			Handler:    _SimDaemon_StartClockOrder_Handler,
		},
		{
			MethodName: "StopClockOrder",
			Handler:    _SimDaemon_StopClockOrder_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "simulation/daemon/daemon.proto",
}

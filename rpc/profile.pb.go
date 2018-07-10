// Code generated by protoc-gen-go. DO NOT EDIT.
// source: profile.proto

package rpc

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
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

// *
// Used to transmit Mojang or RFC formatted UUIDs as the sole parameter.
type IdRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IdRequest) Reset()         { *m = IdRequest{} }
func (m *IdRequest) String() string { return proto.CompactTextString(m) }
func (*IdRequest) ProtoMessage()    {}
func (*IdRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_profile_f6ac44cad3454f3f, []int{0}
}
func (m *IdRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IdRequest.Unmarshal(m, b)
}
func (m *IdRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IdRequest.Marshal(b, m, deterministic)
}
func (dst *IdRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IdRequest.Merge(dst, src)
}
func (m *IdRequest) XXX_Size() int {
	return xxx_messageInfo_IdRequest.Size(m)
}
func (m *IdRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_IdRequest.DiscardUnknown(m)
}

var xxx_messageInfo_IdRequest proto.InternalMessageInfo

func (m *IdRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

// *
// Stores the parameters for id requests (based on the respective display name and timestamp)
type GetIdRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Timestamp            int64    `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetIdRequest) Reset()         { *m = GetIdRequest{} }
func (m *GetIdRequest) String() string { return proto.CompactTextString(m) }
func (*GetIdRequest) ProtoMessage()    {}
func (*GetIdRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_profile_f6ac44cad3454f3f, []int{1}
}
func (m *GetIdRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetIdRequest.Unmarshal(m, b)
}
func (m *GetIdRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetIdRequest.Marshal(b, m, deterministic)
}
func (dst *GetIdRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetIdRequest.Merge(dst, src)
}
func (m *GetIdRequest) XXX_Size() int {
	return xxx_messageInfo_GetIdRequest.Size(m)
}
func (m *GetIdRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetIdRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetIdRequest proto.InternalMessageInfo

func (m *GetIdRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *GetIdRequest) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

// *
// Represents a profile <-> name mapping at a specified time.
type ProfileId struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	ValidUntil           int64    `protobuf:"varint,5,opt,name=validUntil,proto3" json:"validUntil,omitempty"`
	FirstSeenAt          int64    `protobuf:"varint,6,opt,name=firstSeenAt,proto3" json:"firstSeenAt,omitempty"`
	LastSeenAt           int64    `protobuf:"varint,7,opt,name=lastSeenAt,proto3" json:"lastSeenAt,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ProfileId) Reset()         { *m = ProfileId{} }
func (m *ProfileId) String() string { return proto.CompactTextString(m) }
func (*ProfileId) ProtoMessage()    {}
func (*ProfileId) Descriptor() ([]byte, []int) {
	return fileDescriptor_profile_f6ac44cad3454f3f, []int{2}
}
func (m *ProfileId) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProfileId.Unmarshal(m, b)
}
func (m *ProfileId) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProfileId.Marshal(b, m, deterministic)
}
func (dst *ProfileId) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProfileId.Merge(dst, src)
}
func (m *ProfileId) XXX_Size() int {
	return xxx_messageInfo_ProfileId.Size(m)
}
func (m *ProfileId) XXX_DiscardUnknown() {
	xxx_messageInfo_ProfileId.DiscardUnknown(m)
}

var xxx_messageInfo_ProfileId proto.InternalMessageInfo

func (m *ProfileId) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ProfileId) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ProfileId) GetValidUntil() int64 {
	if m != nil {
		return m.ValidUntil
	}
	return 0
}

func (m *ProfileId) GetFirstSeenAt() int64 {
	if m != nil {
		return m.FirstSeenAt
	}
	return 0
}

func (m *ProfileId) GetLastSeenAt() int64 {
	if m != nil {
		return m.LastSeenAt
	}
	return 0
}

// *
// Represents a complete name history.
type NameHistory struct {
	History              []*NameHistoryEntry `protobuf:"bytes,1,rep,name=history,proto3" json:"history,omitempty"`
	ValidUntil           int64               `protobuf:"varint,2,opt,name=validUntil,proto3" json:"validUntil,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *NameHistory) Reset()         { *m = NameHistory{} }
func (m *NameHistory) String() string { return proto.CompactTextString(m) }
func (*NameHistory) ProtoMessage()    {}
func (*NameHistory) Descriptor() ([]byte, []int) {
	return fileDescriptor_profile_f6ac44cad3454f3f, []int{3}
}
func (m *NameHistory) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NameHistory.Unmarshal(m, b)
}
func (m *NameHistory) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NameHistory.Marshal(b, m, deterministic)
}
func (dst *NameHistory) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NameHistory.Merge(dst, src)
}
func (m *NameHistory) XXX_Size() int {
	return xxx_messageInfo_NameHistory.Size(m)
}
func (m *NameHistory) XXX_DiscardUnknown() {
	xxx_messageInfo_NameHistory.DiscardUnknown(m)
}

var xxx_messageInfo_NameHistory proto.InternalMessageInfo

func (m *NameHistory) GetHistory() []*NameHistoryEntry {
	if m != nil {
		return m.History
	}
	return nil
}

func (m *NameHistory) GetValidUntil() int64 {
	if m != nil {
		return m.ValidUntil
	}
	return 0
}

// *
// Represents a single entry in the name history.
type NameHistoryEntry struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	ChangedToAt          int64    `protobuf:"varint,2,opt,name=changedToAt,proto3" json:"changedToAt,omitempty"`
	ValidUntil           int64    `protobuf:"varint,3,opt,name=validUntil,proto3" json:"validUntil,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NameHistoryEntry) Reset()         { *m = NameHistoryEntry{} }
func (m *NameHistoryEntry) String() string { return proto.CompactTextString(m) }
func (*NameHistoryEntry) ProtoMessage()    {}
func (*NameHistoryEntry) Descriptor() ([]byte, []int) {
	return fileDescriptor_profile_f6ac44cad3454f3f, []int{4}
}
func (m *NameHistoryEntry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NameHistoryEntry.Unmarshal(m, b)
}
func (m *NameHistoryEntry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NameHistoryEntry.Marshal(b, m, deterministic)
}
func (dst *NameHistoryEntry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NameHistoryEntry.Merge(dst, src)
}
func (m *NameHistoryEntry) XXX_Size() int {
	return xxx_messageInfo_NameHistoryEntry.Size(m)
}
func (m *NameHistoryEntry) XXX_DiscardUnknown() {
	xxx_messageInfo_NameHistoryEntry.DiscardUnknown(m)
}

var xxx_messageInfo_NameHistoryEntry proto.InternalMessageInfo

func (m *NameHistoryEntry) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *NameHistoryEntry) GetChangedToAt() int64 {
	if m != nil {
		return m.ChangedToAt
	}
	return 0
}

func (m *NameHistoryEntry) GetValidUntil() int64 {
	if m != nil {
		return m.ValidUntil
	}
	return 0
}

// *
// Stores the parameters for bulk id requests.
type BulkIdRequest struct {
	Names                []string `protobuf:"bytes,1,rep,name=names,proto3" json:"names,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BulkIdRequest) Reset()         { *m = BulkIdRequest{} }
func (m *BulkIdRequest) String() string { return proto.CompactTextString(m) }
func (*BulkIdRequest) ProtoMessage()    {}
func (*BulkIdRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_profile_f6ac44cad3454f3f, []int{5}
}
func (m *BulkIdRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BulkIdRequest.Unmarshal(m, b)
}
func (m *BulkIdRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BulkIdRequest.Marshal(b, m, deterministic)
}
func (dst *BulkIdRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BulkIdRequest.Merge(dst, src)
}
func (m *BulkIdRequest) XXX_Size() int {
	return xxx_messageInfo_BulkIdRequest.Size(m)
}
func (m *BulkIdRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_BulkIdRequest.DiscardUnknown(m)
}

var xxx_messageInfo_BulkIdRequest proto.InternalMessageInfo

func (m *BulkIdRequest) GetNames() []string {
	if m != nil {
		return m.Names
	}
	return nil
}

// *
// Represents a list of bulk id responses.
type BulkIdResponse struct {
	Ids                  []*ProfileId `protobuf:"bytes,1,rep,name=ids,proto3" json:"ids,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *BulkIdResponse) Reset()         { *m = BulkIdResponse{} }
func (m *BulkIdResponse) String() string { return proto.CompactTextString(m) }
func (*BulkIdResponse) ProtoMessage()    {}
func (*BulkIdResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_profile_f6ac44cad3454f3f, []int{6}
}
func (m *BulkIdResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BulkIdResponse.Unmarshal(m, b)
}
func (m *BulkIdResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BulkIdResponse.Marshal(b, m, deterministic)
}
func (dst *BulkIdResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BulkIdResponse.Merge(dst, src)
}
func (m *BulkIdResponse) XXX_Size() int {
	return xxx_messageInfo_BulkIdResponse.Size(m)
}
func (m *BulkIdResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_BulkIdResponse.DiscardUnknown(m)
}

var xxx_messageInfo_BulkIdResponse proto.InternalMessageInfo

func (m *BulkIdResponse) GetIds() []*ProfileId {
	if m != nil {
		return m.Ids
	}
	return nil
}

func init() {
	proto.RegisterType((*IdRequest)(nil), "rpc.IdRequest")
	proto.RegisterType((*GetIdRequest)(nil), "rpc.GetIdRequest")
	proto.RegisterType((*ProfileId)(nil), "rpc.ProfileId")
	proto.RegisterType((*NameHistory)(nil), "rpc.NameHistory")
	proto.RegisterType((*NameHistoryEntry)(nil), "rpc.NameHistoryEntry")
	proto.RegisterType((*BulkIdRequest)(nil), "rpc.BulkIdRequest")
	proto.RegisterType((*BulkIdResponse)(nil), "rpc.BulkIdResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ProfileServiceClient is the client API for ProfileService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ProfileServiceClient interface {
	// *
	// Resolves the profile identifier and correct casing of a given name.
	//
	// When unix epoch (e.g. zero) is passed instead of an actual timestamp, the
	// original user of a name will be resolved (e.g. associations prior to name
	// changing support).
	//
	// If no profile has been associated with the specified name, an unpopulated
	// object is returned instead.
	GetId(ctx context.Context, in *GetIdRequest, opts ...grpc.CallOption) (*ProfileId, error)
	// *
	// Retrieves a complete history of name changes for the profile associated
	// with a given identifier.
	//
	// Names which have been changed to at unix epoch (e.g. zero) refer to the
	// original account name.
	//
	// When no profile with the specified identifier exists, an unpopulated object
	// is returned instead.
	GetNameHistory(ctx context.Context, in *IdRequest, opts ...grpc.CallOption) (*NameHistory, error)
	// *
	// Resolves the profile identifiers and correct casings of multiple names at
	// once.
	//
	// If a name cannot be found, its association will be omitted from the
	// resulting array.
	//
	// Bulk requests do not accept timestamps and will always resolve associations
	// at the current time.
	BulkGetId(ctx context.Context, in *BulkIdRequest, opts ...grpc.CallOption) (*BulkIdResponse, error)
	// *
	// Retrieves a profile based on its associated identifier.
	//
	// If no profile with the specified identifier exists, an unpopulated object
	// is returned instead.
	GetProfile(ctx context.Context, in *IdRequest, opts ...grpc.CallOption) (*Profile, error)
}

type profileServiceClient struct {
	cc *grpc.ClientConn
}

func NewProfileServiceClient(cc *grpc.ClientConn) ProfileServiceClient {
	return &profileServiceClient{cc}
}

func (c *profileServiceClient) GetId(ctx context.Context, in *GetIdRequest, opts ...grpc.CallOption) (*ProfileId, error) {
	out := new(ProfileId)
	err := c.cc.Invoke(ctx, "/rpc.ProfileService/GetId", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileServiceClient) GetNameHistory(ctx context.Context, in *IdRequest, opts ...grpc.CallOption) (*NameHistory, error) {
	out := new(NameHistory)
	err := c.cc.Invoke(ctx, "/rpc.ProfileService/GetNameHistory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileServiceClient) BulkGetId(ctx context.Context, in *BulkIdRequest, opts ...grpc.CallOption) (*BulkIdResponse, error) {
	out := new(BulkIdResponse)
	err := c.cc.Invoke(ctx, "/rpc.ProfileService/BulkGetId", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileServiceClient) GetProfile(ctx context.Context, in *IdRequest, opts ...grpc.CallOption) (*Profile, error) {
	out := new(Profile)
	err := c.cc.Invoke(ctx, "/rpc.ProfileService/GetProfile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProfileServiceServer is the server API for ProfileService service.
type ProfileServiceServer interface {
	// *
	// Resolves the profile identifier and correct casing of a given name.
	//
	// When unix epoch (e.g. zero) is passed instead of an actual timestamp, the
	// original user of a name will be resolved (e.g. associations prior to name
	// changing support).
	//
	// If no profile has been associated with the specified name, an unpopulated
	// object is returned instead.
	GetId(context.Context, *GetIdRequest) (*ProfileId, error)
	// *
	// Retrieves a complete history of name changes for the profile associated
	// with a given identifier.
	//
	// Names which have been changed to at unix epoch (e.g. zero) refer to the
	// original account name.
	//
	// When no profile with the specified identifier exists, an unpopulated object
	// is returned instead.
	GetNameHistory(context.Context, *IdRequest) (*NameHistory, error)
	// *
	// Resolves the profile identifiers and correct casings of multiple names at
	// once.
	//
	// If a name cannot be found, its association will be omitted from the
	// resulting array.
	//
	// Bulk requests do not accept timestamps and will always resolve associations
	// at the current time.
	BulkGetId(context.Context, *BulkIdRequest) (*BulkIdResponse, error)
	// *
	// Retrieves a profile based on its associated identifier.
	//
	// If no profile with the specified identifier exists, an unpopulated object
	// is returned instead.
	GetProfile(context.Context, *IdRequest) (*Profile, error)
}

func RegisterProfileServiceServer(s *grpc.Server, srv ProfileServiceServer) {
	s.RegisterService(&_ProfileService_serviceDesc, srv)
}

func _ProfileService_GetId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServiceServer).GetId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.ProfileService/GetId",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServiceServer).GetId(ctx, req.(*GetIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProfileService_GetNameHistory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServiceServer).GetNameHistory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.ProfileService/GetNameHistory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServiceServer).GetNameHistory(ctx, req.(*IdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProfileService_BulkGetId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BulkIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServiceServer).BulkGetId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.ProfileService/BulkGetId",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServiceServer).BulkGetId(ctx, req.(*BulkIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProfileService_GetProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServiceServer).GetProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.ProfileService/GetProfile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServiceServer).GetProfile(ctx, req.(*IdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ProfileService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.ProfileService",
	HandlerType: (*ProfileServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetId",
			Handler:    _ProfileService_GetId_Handler,
		},
		{
			MethodName: "GetNameHistory",
			Handler:    _ProfileService_GetNameHistory_Handler,
		},
		{
			MethodName: "BulkGetId",
			Handler:    _ProfileService_BulkGetId_Handler,
		},
		{
			MethodName: "GetProfile",
			Handler:    _ProfileService_GetProfile_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "profile.proto",
}

func init() { proto.RegisterFile("profile.proto", fileDescriptor_profile_f6ac44cad3454f3f) }

var fileDescriptor_profile_f6ac44cad3454f3f = []byte{
	// 414 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x53, 0xd1, 0x8a, 0x9b, 0x40,
	0x14, 0x45, 0x6d, 0x76, 0xf1, 0x9a, 0x95, 0xed, 0xb4, 0x05, 0x49, 0x4b, 0x11, 0xa1, 0x10, 0xf6,
	0xc1, 0x82, 0xed, 0x07, 0x74, 0x17, 0x4a, 0xba, 0x2f, 0xa5, 0xb8, 0xed, 0x6b, 0xc1, 0x75, 0x66,
	0xe3, 0x10, 0x75, 0xa6, 0x33, 0x37, 0x81, 0x7c, 0x46, 0x7f, 0xad, 0x5f, 0x54, 0x1c, 0x4d, 0x9c,
	0x98, 0x7d, 0xd3, 0x73, 0xcf, 0x3d, 0xe7, 0x9e, 0x3b, 0x33, 0x70, 0x25, 0x95, 0x78, 0xe2, 0x35,
	0x4b, 0xa5, 0x12, 0x28, 0x88, 0xa7, 0x64, 0xb9, 0x98, 0x97, 0xa2, 0x69, 0x44, 0xdb, 0x43, 0xc9,
	0x5b, 0xf0, 0xef, 0x69, 0xce, 0xfe, 0x6c, 0x99, 0x46, 0x12, 0x82, 0xcb, 0x69, 0xe4, 0xc4, 0xce,
	0xd2, 0xcf, 0x5d, 0x4e, 0x93, 0x2f, 0x30, 0x5f, 0x31, 0x1c, 0xeb, 0x04, 0x5e, 0xb4, 0x45, 0xc3,
	0x06, 0x86, 0xf9, 0x26, 0xef, 0xc0, 0x47, 0xde, 0x30, 0x8d, 0x45, 0x23, 0x23, 0x37, 0x76, 0x96,
	0x5e, 0x3e, 0x02, 0xc9, 0x5f, 0x07, 0xfc, 0x1f, 0xfd, 0x0c, 0xf7, 0x74, 0xaa, 0x7f, 0xd4, 0x73,
	0x2d, 0xbd, 0xf7, 0x00, 0xbb, 0xa2, 0xe6, 0xf4, 0x57, 0x8b, 0xbc, 0x8e, 0x66, 0x46, 0xd0, 0x42,
	0x48, 0x0c, 0xc1, 0x13, 0x57, 0x1a, 0x1f, 0x18, 0x6b, 0x6f, 0x31, 0xba, 0x30, 0x04, 0x1b, 0xea,
	0x14, 0xea, 0xe2, 0x48, 0xb8, 0xec, 0x15, 0x46, 0x24, 0xf9, 0x0d, 0xc1, 0xf7, 0xa2, 0x61, 0xdf,
	0xb8, 0x46, 0xa1, 0xf6, 0xe4, 0x23, 0x5c, 0x56, 0xfd, 0x67, 0xe4, 0xc4, 0xde, 0x32, 0xc8, 0xde,
	0xa4, 0x4a, 0x96, 0xa9, 0x45, 0xf9, 0xda, 0xa2, 0xda, 0xe7, 0x07, 0xd6, 0x64, 0x42, 0x77, 0x3a,
	0x61, 0x52, 0xc1, 0xf5, 0xb4, 0xf9, 0xd9, 0xcd, 0xc5, 0x10, 0x94, 0x55, 0xd1, 0xae, 0x19, 0xfd,
	0x29, 0x6e, 0x71, 0x10, 0xb2, 0xa1, 0x89, 0x93, 0x77, 0xe6, 0xf4, 0x01, 0xae, 0xee, 0xb6, 0xf5,
	0x66, 0x3c, 0xa0, 0xd7, 0x30, 0xeb, 0xa4, 0xb5, 0x49, 0xe2, 0xe7, 0xfd, 0x4f, 0x92, 0x41, 0x78,
	0xa0, 0x69, 0x29, 0x5a, 0xdd, 0x59, 0x7b, 0x9c, 0xea, 0x21, 0x6f, 0x68, 0xf2, 0x1e, 0x4f, 0x29,
	0xef, 0x4a, 0xd9, 0x3f, 0x07, 0xc2, 0x01, 0x7a, 0x60, 0x6a, 0xc7, 0x4b, 0x46, 0x6e, 0x60, 0x66,
	0x6e, 0x03, 0x79, 0x69, 0x1a, 0xec, 0x9b, 0xb1, 0x98, 0x68, 0x90, 0x0c, 0xc2, 0x15, 0x43, 0x7b,
	0xcd, 0x3d, 0x63, 0xec, 0xb8, 0x9e, 0x6e, 0x99, 0x7c, 0x06, 0xbf, 0x1b, 0xb3, 0xf7, 0x20, 0xa6,
	0x7c, 0x92, 0x6e, 0xf1, 0xea, 0x04, 0x1b, 0xa2, 0xdc, 0x00, 0xac, 0x18, 0x0e, 0xce, 0x67, 0x2e,
	0x73, 0x7b, 0xae, 0xbb, 0x04, 0x62, 0x2e, 0xd2, 0x35, 0xc7, 0x6a, 0xfb, 0x98, 0x52, 0x81, 0x1a,
	0x0b, 0x85, 0xa9, 0x46, 0x51, 0x6e, 0x64, 0xf7, 0x4a, 0x94, 0x2c, 0x1f, 0x2f, 0xcc, 0xbb, 0xf8,
	0xf4, 0x3f, 0x00, 0x00, 0xff, 0xff, 0x9b, 0x70, 0xdf, 0x90, 0x3b, 0x03, 0x00, 0x00,
}

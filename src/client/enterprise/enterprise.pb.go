// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.12.3
// source: client/enterprise/enterprise.proto

package enterprise

import (
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type State int32

const (
	State_NONE    State = 0
	State_ACTIVE  State = 1
	State_EXPIRED State = 2
)

// Enum value maps for State.
var (
	State_name = map[int32]string{
		0: "NONE",
		1: "ACTIVE",
		2: "EXPIRED",
	}
	State_value = map[string]int32{
		"NONE":    0,
		"ACTIVE":  1,
		"EXPIRED": 2,
	}
)

func (x State) Enum() *State {
	p := new(State)
	*p = x
	return p
}

func (x State) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (State) Descriptor() protoreflect.EnumDescriptor {
	return file_client_enterprise_enterprise_proto_enumTypes[0].Descriptor()
}

func (State) Type() protoreflect.EnumType {
	return &file_client_enterprise_enterprise_proto_enumTypes[0]
}

func (x State) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use State.Descriptor instead.
func (State) EnumDescriptor() ([]byte, []int) {
	return file_client_enterprise_enterprise_proto_rawDescGZIP(), []int{0}
}

// EnterpriseRecord is the record we store in etcd for a Pachyderm enterprise
// token that has been provided to a Pachyderm cluster
type EnterpriseRecord struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ActivationCode string `protobuf:"bytes,1,opt,name=activation_code,json=activationCode,proto3" json:"activation_code,omitempty"`
	// expires is a timestamp indicating when this activation code will expire.
	Expires *timestamp.Timestamp `protobuf:"bytes,2,opt,name=expires,proto3" json:"expires,omitempty"`
}

func (x *EnterpriseRecord) Reset() {
	*x = EnterpriseRecord{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_enterprise_enterprise_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnterpriseRecord) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnterpriseRecord) ProtoMessage() {}

func (x *EnterpriseRecord) ProtoReflect() protoreflect.Message {
	mi := &file_client_enterprise_enterprise_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnterpriseRecord.ProtoReflect.Descriptor instead.
func (*EnterpriseRecord) Descriptor() ([]byte, []int) {
	return file_client_enterprise_enterprise_proto_rawDescGZIP(), []int{0}
}

func (x *EnterpriseRecord) GetActivationCode() string {
	if x != nil {
		return x.ActivationCode
	}
	return ""
}

func (x *EnterpriseRecord) GetExpires() *timestamp.Timestamp {
	if x != nil {
		return x.Expires
	}
	return nil
}

// TokenInfo contains information about the currently active enterprise token
type TokenInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// expires indicates when the current token expires (unset if there is no
	// current token)
	Expires *timestamp.Timestamp `protobuf:"bytes,1,opt,name=expires,proto3" json:"expires,omitempty"`
}

func (x *TokenInfo) Reset() {
	*x = TokenInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_enterprise_enterprise_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TokenInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TokenInfo) ProtoMessage() {}

func (x *TokenInfo) ProtoReflect() protoreflect.Message {
	mi := &file_client_enterprise_enterprise_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TokenInfo.ProtoReflect.Descriptor instead.
func (*TokenInfo) Descriptor() ([]byte, []int) {
	return file_client_enterprise_enterprise_proto_rawDescGZIP(), []int{1}
}

func (x *TokenInfo) GetExpires() *timestamp.Timestamp {
	if x != nil {
		return x.Expires
	}
	return nil
}

type ActivateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// activation_code is a Pachyderm enterprise activation code. New users can
	// obtain trial activation codes
	ActivationCode string `protobuf:"bytes,1,opt,name=activation_code,json=activationCode,proto3" json:"activation_code,omitempty"`
	// expires is a timestamp indicating when this activation code will expire.
	// This should not generally be set (it's primarily used for testing), and is
	// only applied if it's earlier than the signed expiration time in
	// 'activation_code'.
	Expires *timestamp.Timestamp `protobuf:"bytes,2,opt,name=expires,proto3" json:"expires,omitempty"`
}

func (x *ActivateRequest) Reset() {
	*x = ActivateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_enterprise_enterprise_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ActivateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActivateRequest) ProtoMessage() {}

func (x *ActivateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_client_enterprise_enterprise_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActivateRequest.ProtoReflect.Descriptor instead.
func (*ActivateRequest) Descriptor() ([]byte, []int) {
	return file_client_enterprise_enterprise_proto_rawDescGZIP(), []int{2}
}

func (x *ActivateRequest) GetActivationCode() string {
	if x != nil {
		return x.ActivationCode
	}
	return ""
}

func (x *ActivateRequest) GetExpires() *timestamp.Timestamp {
	if x != nil {
		return x.Expires
	}
	return nil
}

type ActivateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Info *TokenInfo `protobuf:"bytes,1,opt,name=info,proto3" json:"info,omitempty"`
}

func (x *ActivateResponse) Reset() {
	*x = ActivateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_enterprise_enterprise_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ActivateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActivateResponse) ProtoMessage() {}

func (x *ActivateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_client_enterprise_enterprise_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActivateResponse.ProtoReflect.Descriptor instead.
func (*ActivateResponse) Descriptor() ([]byte, []int) {
	return file_client_enterprise_enterprise_proto_rawDescGZIP(), []int{3}
}

func (x *ActivateResponse) GetInfo() *TokenInfo {
	if x != nil {
		return x.Info
	}
	return nil
}

type GetStateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetStateRequest) Reset() {
	*x = GetStateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_enterprise_enterprise_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetStateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetStateRequest) ProtoMessage() {}

func (x *GetStateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_client_enterprise_enterprise_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetStateRequest.ProtoReflect.Descriptor instead.
func (*GetStateRequest) Descriptor() ([]byte, []int) {
	return file_client_enterprise_enterprise_proto_rawDescGZIP(), []int{4}
}

type GetStateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	State          State      `protobuf:"varint,1,opt,name=state,proto3,enum=enterprise.State" json:"state,omitempty"`
	Info           *TokenInfo `protobuf:"bytes,2,opt,name=info,proto3" json:"info,omitempty"`
	ActivationCode string     `protobuf:"bytes,3,opt,name=activation_code,json=activationCode,proto3" json:"activation_code,omitempty"`
}

func (x *GetStateResponse) Reset() {
	*x = GetStateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_enterprise_enterprise_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetStateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetStateResponse) ProtoMessage() {}

func (x *GetStateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_client_enterprise_enterprise_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetStateResponse.ProtoReflect.Descriptor instead.
func (*GetStateResponse) Descriptor() ([]byte, []int) {
	return file_client_enterprise_enterprise_proto_rawDescGZIP(), []int{5}
}

func (x *GetStateResponse) GetState() State {
	if x != nil {
		return x.State
	}
	return State_NONE
}

func (x *GetStateResponse) GetInfo() *TokenInfo {
	if x != nil {
		return x.Info
	}
	return nil
}

func (x *GetStateResponse) GetActivationCode() string {
	if x != nil {
		return x.ActivationCode
	}
	return ""
}

type DeactivateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeactivateRequest) Reset() {
	*x = DeactivateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_enterprise_enterprise_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeactivateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeactivateRequest) ProtoMessage() {}

func (x *DeactivateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_client_enterprise_enterprise_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeactivateRequest.ProtoReflect.Descriptor instead.
func (*DeactivateRequest) Descriptor() ([]byte, []int) {
	return file_client_enterprise_enterprise_proto_rawDescGZIP(), []int{6}
}

type DeactivateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeactivateResponse) Reset() {
	*x = DeactivateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_enterprise_enterprise_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeactivateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeactivateResponse) ProtoMessage() {}

func (x *DeactivateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_client_enterprise_enterprise_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeactivateResponse.ProtoReflect.Descriptor instead.
func (*DeactivateResponse) Descriptor() ([]byte, []int) {
	return file_client_enterprise_enterprise_proto_rawDescGZIP(), []int{7}
}

var File_client_enterprise_enterprise_proto protoreflect.FileDescriptor

var file_client_enterprise_enterprise_proto_rawDesc = []byte{
	0x0a, 0x22, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2f, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x70, 0x72,
	0x69, 0x73, 0x65, 0x2f, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x70, 0x72, 0x69, 0x73, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x70, 0x72, 0x69, 0x73, 0x65,
	0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x71, 0x0a, 0x10, 0x45, 0x6e, 0x74, 0x65, 0x72, 0x70, 0x72, 0x69, 0x73, 0x65, 0x52,
	0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x27, 0x0a, 0x0f, 0x61, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e,
	0x61, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x34,
	0x0a, 0x07, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x65, 0x78, 0x70,
	0x69, 0x72, 0x65, 0x73, 0x22, 0x41, 0x0a, 0x09, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x49, 0x6e, 0x66,
	0x6f, 0x12, 0x34, 0x0a, 0x07, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07,
	0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x22, 0x70, 0x0a, 0x0f, 0x41, 0x63, 0x74, 0x69, 0x76,
	0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x27, 0x0a, 0x0f, 0x61, 0x63,
	0x74, 0x69, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0e, 0x61, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43,
	0x6f, 0x64, 0x65, 0x12, 0x34, 0x0a, 0x07, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x07, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x22, 0x3d, 0x0a, 0x10, 0x41, 0x63, 0x74,
	0x69, 0x76, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x29, 0x0a,
	0x04, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x65, 0x6e,
	0x74, 0x65, 0x72, 0x70, 0x72, 0x69, 0x73, 0x65, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x49, 0x6e,
	0x66, 0x6f, 0x52, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x22, 0x11, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x8f, 0x01, 0x0a, 0x10,
	0x47, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x27, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x11, 0x2e, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x70, 0x72, 0x69, 0x73, 0x65, 0x2e, 0x53, 0x74, 0x61,
	0x74, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x29, 0x0a, 0x04, 0x69, 0x6e, 0x66,
	0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x70,
	0x72, 0x69, 0x73, 0x65, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x04,
	0x69, 0x6e, 0x66, 0x6f, 0x12, 0x27, 0x0a, 0x0f, 0x61, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x61,
	0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x64, 0x65, 0x22, 0x13, 0x0a,
	0x11, 0x44, 0x65, 0x61, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x22, 0x14, 0x0a, 0x12, 0x44, 0x65, 0x61, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2a, 0x2a, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x74,
	0x65, 0x12, 0x08, 0x0a, 0x04, 0x4e, 0x4f, 0x4e, 0x45, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x41,
	0x43, 0x54, 0x49, 0x56, 0x45, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x45, 0x58, 0x50, 0x49, 0x52,
	0x45, 0x44, 0x10, 0x02, 0x32, 0xe6, 0x01, 0x0a, 0x03, 0x41, 0x50, 0x49, 0x12, 0x47, 0x0a, 0x08,
	0x41, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x65, 0x12, 0x1b, 0x2e, 0x65, 0x6e, 0x74, 0x65, 0x72,
	0x70, 0x72, 0x69, 0x73, 0x65, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x70, 0x72, 0x69,
	0x73, 0x65, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x47, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74,
	0x65, 0x12, 0x1b, 0x2e, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x70, 0x72, 0x69, 0x73, 0x65, 0x2e, 0x47,
	0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c,
	0x2e, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x70, 0x72, 0x69, 0x73, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4d,
	0x0a, 0x0a, 0x44, 0x65, 0x61, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x65, 0x12, 0x1d, 0x2e, 0x65,
	0x6e, 0x74, 0x65, 0x72, 0x70, 0x72, 0x69, 0x73, 0x65, 0x2e, 0x44, 0x65, 0x61, 0x63, 0x74, 0x69,
	0x76, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x65, 0x6e,
	0x74, 0x65, 0x72, 0x70, 0x72, 0x69, 0x73, 0x65, 0x2e, 0x44, 0x65, 0x61, 0x63, 0x74, 0x69, 0x76,
	0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x36, 0x5a,
	0x34, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x61, 0x63, 0x68,
	0x79, 0x64, 0x65, 0x72, 0x6d, 0x2f, 0x70, 0x61, 0x63, 0x68, 0x79, 0x64, 0x65, 0x72, 0x6d, 0x2f,
	0x73, 0x72, 0x63, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x2f, 0x65, 0x6e, 0x74, 0x65, 0x72,
	0x70, 0x72, 0x69, 0x73, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_client_enterprise_enterprise_proto_rawDescOnce sync.Once
	file_client_enterprise_enterprise_proto_rawDescData = file_client_enterprise_enterprise_proto_rawDesc
)

func file_client_enterprise_enterprise_proto_rawDescGZIP() []byte {
	file_client_enterprise_enterprise_proto_rawDescOnce.Do(func() {
		file_client_enterprise_enterprise_proto_rawDescData = protoimpl.X.CompressGZIP(file_client_enterprise_enterprise_proto_rawDescData)
	})
	return file_client_enterprise_enterprise_proto_rawDescData
}

var file_client_enterprise_enterprise_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_client_enterprise_enterprise_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_client_enterprise_enterprise_proto_goTypes = []interface{}{
	(State)(0),                  // 0: enterprise.State
	(*EnterpriseRecord)(nil),    // 1: enterprise.EnterpriseRecord
	(*TokenInfo)(nil),           // 2: enterprise.TokenInfo
	(*ActivateRequest)(nil),     // 3: enterprise.ActivateRequest
	(*ActivateResponse)(nil),    // 4: enterprise.ActivateResponse
	(*GetStateRequest)(nil),     // 5: enterprise.GetStateRequest
	(*GetStateResponse)(nil),    // 6: enterprise.GetStateResponse
	(*DeactivateRequest)(nil),   // 7: enterprise.DeactivateRequest
	(*DeactivateResponse)(nil),  // 8: enterprise.DeactivateResponse
	(*timestamp.Timestamp)(nil), // 9: google.protobuf.Timestamp
}
var file_client_enterprise_enterprise_proto_depIdxs = []int32{
	9, // 0: enterprise.EnterpriseRecord.expires:type_name -> google.protobuf.Timestamp
	9, // 1: enterprise.TokenInfo.expires:type_name -> google.protobuf.Timestamp
	9, // 2: enterprise.ActivateRequest.expires:type_name -> google.protobuf.Timestamp
	2, // 3: enterprise.ActivateResponse.info:type_name -> enterprise.TokenInfo
	0, // 4: enterprise.GetStateResponse.state:type_name -> enterprise.State
	2, // 5: enterprise.GetStateResponse.info:type_name -> enterprise.TokenInfo
	3, // 6: enterprise.API.Activate:input_type -> enterprise.ActivateRequest
	5, // 7: enterprise.API.GetState:input_type -> enterprise.GetStateRequest
	7, // 8: enterprise.API.Deactivate:input_type -> enterprise.DeactivateRequest
	4, // 9: enterprise.API.Activate:output_type -> enterprise.ActivateResponse
	6, // 10: enterprise.API.GetState:output_type -> enterprise.GetStateResponse
	8, // 11: enterprise.API.Deactivate:output_type -> enterprise.DeactivateResponse
	9, // [9:12] is the sub-list for method output_type
	6, // [6:9] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_client_enterprise_enterprise_proto_init() }
func file_client_enterprise_enterprise_proto_init() {
	if File_client_enterprise_enterprise_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_client_enterprise_enterprise_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnterpriseRecord); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_client_enterprise_enterprise_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TokenInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_client_enterprise_enterprise_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ActivateRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_client_enterprise_enterprise_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ActivateResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_client_enterprise_enterprise_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetStateRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_client_enterprise_enterprise_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetStateResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_client_enterprise_enterprise_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeactivateRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_client_enterprise_enterprise_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeactivateResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_client_enterprise_enterprise_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_client_enterprise_enterprise_proto_goTypes,
		DependencyIndexes: file_client_enterprise_enterprise_proto_depIdxs,
		EnumInfos:         file_client_enterprise_enterprise_proto_enumTypes,
		MessageInfos:      file_client_enterprise_enterprise_proto_msgTypes,
	}.Build()
	File_client_enterprise_enterprise_proto = out.File
	file_client_enterprise_enterprise_proto_rawDesc = nil
	file_client_enterprise_enterprise_proto_goTypes = nil
	file_client_enterprise_enterprise_proto_depIdxs = nil
}

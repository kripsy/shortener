// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: pkg/api/shortener/v1/service.proto

package shortener

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type SaveURLRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *SaveURLRequest) Reset() {
	*x = SaveURLRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_api_shortener_v1_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveURLRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveURLRequest) ProtoMessage() {}

func (x *SaveURLRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_shortener_v1_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveURLRequest.ProtoReflect.Descriptor instead.
func (*SaveURLRequest) Descriptor() ([]byte, []int) {
	return file_pkg_api_shortener_v1_service_proto_rawDescGZIP(), []int{0}
}

func (x *SaveURLRequest) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type SaveURLResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result        string `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
	IsUniqueError bool   `protobuf:"varint,2,opt,name=is_unique_error,json=isUniqueError,proto3" json:"is_unique_error,omitempty"`
}

func (x *SaveURLResponse) Reset() {
	*x = SaveURLResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_api_shortener_v1_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveURLResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveURLResponse) ProtoMessage() {}

func (x *SaveURLResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_shortener_v1_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveURLResponse.ProtoReflect.Descriptor instead.
func (*SaveURLResponse) Descriptor() ([]byte, []int) {
	return file_pkg_api_shortener_v1_service_proto_rawDescGZIP(), []int{1}
}

func (x *SaveURLResponse) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

func (x *SaveURLResponse) GetIsUniqueError() bool {
	if x != nil {
		return x.IsUniqueError
	}
	return false
}

type GetURLRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *GetURLRequest) Reset() {
	*x = GetURLRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_api_shortener_v1_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetURLRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetURLRequest) ProtoMessage() {}

func (x *GetURLRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_shortener_v1_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetURLRequest.ProtoReflect.Descriptor instead.
func (*GetURLRequest) Descriptor() ([]byte, []int) {
	return file_pkg_api_shortener_v1_service_proto_rawDescGZIP(), []int{2}
}

func (x *GetURLRequest) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type GetURLResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *GetURLResponse) Reset() {
	*x = GetURLResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_api_shortener_v1_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetURLResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetURLResponse) ProtoMessage() {}

func (x *GetURLResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_shortener_v1_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetURLResponse.ProtoReflect.Descriptor instead.
func (*GetURLResponse) Descriptor() ([]byte, []int) {
	return file_pkg_api_shortener_v1_service_proto_rawDescGZIP(), []int{3}
}

func (x *GetURLResponse) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type SaveBatchURLRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UrlBatch []*SaveBatchURLRequest_URLObject `protobuf:"bytes,1,rep,name=url_batch,json=urlBatch,proto3" json:"url_batch,omitempty"`
}

func (x *SaveBatchURLRequest) Reset() {
	*x = SaveBatchURLRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_api_shortener_v1_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveBatchURLRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveBatchURLRequest) ProtoMessage() {}

func (x *SaveBatchURLRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_shortener_v1_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveBatchURLRequest.ProtoReflect.Descriptor instead.
func (*SaveBatchURLRequest) Descriptor() ([]byte, []int) {
	return file_pkg_api_shortener_v1_service_proto_rawDescGZIP(), []int{4}
}

func (x *SaveBatchURLRequest) GetUrlBatch() []*SaveBatchURLRequest_URLObject {
	if x != nil {
		return x.UrlBatch
	}
	return nil
}

type SaveBatchURLResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UrlBatch []*SaveBatchURLResponse_URLObject `protobuf:"bytes,1,rep,name=url_batch,json=urlBatch,proto3" json:"url_batch,omitempty"`
}

func (x *SaveBatchURLResponse) Reset() {
	*x = SaveBatchURLResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_api_shortener_v1_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveBatchURLResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveBatchURLResponse) ProtoMessage() {}

func (x *SaveBatchURLResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_shortener_v1_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveBatchURLResponse.ProtoReflect.Descriptor instead.
func (*SaveBatchURLResponse) Descriptor() ([]byte, []int) {
	return file_pkg_api_shortener_v1_service_proto_rawDescGZIP(), []int{5}
}

func (x *SaveBatchURLResponse) GetUrlBatch() []*SaveBatchURLResponse_URLObject {
	if x != nil {
		return x.UrlBatch
	}
	return nil
}

type GetStatsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Urls  int32 `protobuf:"varint,1,opt,name=urls,proto3" json:"urls,omitempty"`
	Users int32 `protobuf:"varint,2,opt,name=users,proto3" json:"users,omitempty"`
}

func (x *GetStatsResponse) Reset() {
	*x = GetStatsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_api_shortener_v1_service_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetStatsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetStatsResponse) ProtoMessage() {}

func (x *GetStatsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_shortener_v1_service_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetStatsResponse.ProtoReflect.Descriptor instead.
func (*GetStatsResponse) Descriptor() ([]byte, []int) {
	return file_pkg_api_shortener_v1_service_proto_rawDescGZIP(), []int{6}
}

func (x *GetStatsResponse) GetUrls() int32 {
	if x != nil {
		return x.Urls
	}
	return 0
}

func (x *GetStatsResponse) GetUsers() int32 {
	if x != nil {
		return x.Users
	}
	return 0
}

type SaveBatchURLRequest_URLObject struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CorrelationId string `protobuf:"bytes,1,opt,name=correlation_id,json=correlationId,proto3" json:"correlation_id,omitempty"`
	OriginalUrl   string `protobuf:"bytes,2,opt,name=original_url,json=originalUrl,proto3" json:"original_url,omitempty"`
}

func (x *SaveBatchURLRequest_URLObject) Reset() {
	*x = SaveBatchURLRequest_URLObject{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_api_shortener_v1_service_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveBatchURLRequest_URLObject) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveBatchURLRequest_URLObject) ProtoMessage() {}

func (x *SaveBatchURLRequest_URLObject) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_shortener_v1_service_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveBatchURLRequest_URLObject.ProtoReflect.Descriptor instead.
func (*SaveBatchURLRequest_URLObject) Descriptor() ([]byte, []int) {
	return file_pkg_api_shortener_v1_service_proto_rawDescGZIP(), []int{4, 0}
}

func (x *SaveBatchURLRequest_URLObject) GetCorrelationId() string {
	if x != nil {
		return x.CorrelationId
	}
	return ""
}

func (x *SaveBatchURLRequest_URLObject) GetOriginalUrl() string {
	if x != nil {
		return x.OriginalUrl
	}
	return ""
}

type SaveBatchURLResponse_URLObject struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CorrelationId string `protobuf:"bytes,1,opt,name=correlation_id,json=correlationId,proto3" json:"correlation_id,omitempty"`
	ShortUrl      string `protobuf:"bytes,2,opt,name=short_url,json=shortUrl,proto3" json:"short_url,omitempty"`
}

func (x *SaveBatchURLResponse_URLObject) Reset() {
	*x = SaveBatchURLResponse_URLObject{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_api_shortener_v1_service_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveBatchURLResponse_URLObject) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveBatchURLResponse_URLObject) ProtoMessage() {}

func (x *SaveBatchURLResponse_URLObject) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_api_shortener_v1_service_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveBatchURLResponse_URLObject.ProtoReflect.Descriptor instead.
func (*SaveBatchURLResponse_URLObject) Descriptor() ([]byte, []int) {
	return file_pkg_api_shortener_v1_service_proto_rawDescGZIP(), []int{5, 0}
}

func (x *SaveBatchURLResponse_URLObject) GetCorrelationId() string {
	if x != nil {
		return x.CorrelationId
	}
	return ""
}

func (x *SaveBatchURLResponse_URLObject) GetShortUrl() string {
	if x != nil {
		return x.ShortUrl
	}
	return ""
}

var File_pkg_api_shortener_v1_service_proto protoreflect.FileDescriptor

var file_pkg_api_shortener_v1_service_proto_rawDesc = []byte{
	0x0a, 0x22, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65,
	0x6e, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x14, 0x70, 0x6b, 0x67, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x68,
	0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2b, 0x0a, 0x0e, 0x53, 0x61, 0x76, 0x65, 0x55, 0x52, 0x4c, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x03, 0x75, 0x72,
	0x6c, 0x22, 0x5a, 0x0a, 0x0f, 0x53, 0x61, 0x76, 0x65, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1f, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x06, 0x72,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x26, 0x0a, 0x0f, 0x69, 0x73, 0x5f, 0x75, 0x6e, 0x69, 0x71,
	0x75, 0x65, 0x5f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0d,
	0x69, 0x73, 0x55, 0x6e, 0x69, 0x71, 0x75, 0x65, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x22, 0x2a, 0x0a,
	0x0d, 0x47, 0x65, 0x74, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19,
	0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xba, 0x48, 0x04,
	0x72, 0x02, 0x10, 0x01, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x22, 0x2b, 0x0a, 0x0e, 0x47, 0x65, 0x74,
	0x55, 0x52, 0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x19, 0x0a, 0x03, 0x75,
	0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10,
	0x01, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x22, 0xc7, 0x01, 0x0a, 0x13, 0x53, 0x61, 0x76, 0x65, 0x42,
	0x61, 0x74, 0x63, 0x68, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x50,
	0x0a, 0x09, 0x75, 0x72, 0x6c, 0x5f, 0x62, 0x61, 0x74, 0x63, 0x68, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x33, 0x2e, 0x70, 0x6b, 0x67, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x68, 0x6f, 0x72,
	0x74, 0x65, 0x6e, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x61, 0x76, 0x65, 0x42, 0x61, 0x74,
	0x63, 0x68, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x55, 0x52, 0x4c,
	0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x08, 0x75, 0x72, 0x6c, 0x42, 0x61, 0x74, 0x63, 0x68,
	0x1a, 0x5e, 0x0a, 0x09, 0x55, 0x52, 0x4c, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x25, 0x0a,
	0x0e, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x49, 0x64, 0x12, 0x2a, 0x0a, 0x0c, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c,
	0x5f, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xba, 0x48, 0x04, 0x72,
	0x02, 0x10, 0x01, 0x52, 0x0b, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x55, 0x72, 0x6c,
	0x22, 0xc3, 0x01, 0x0a, 0x14, 0x53, 0x61, 0x76, 0x65, 0x42, 0x61, 0x74, 0x63, 0x68, 0x55, 0x52,
	0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x51, 0x0a, 0x09, 0x75, 0x72, 0x6c,
	0x5f, 0x62, 0x61, 0x74, 0x63, 0x68, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x34, 0x2e, 0x70,
	0x6b, 0x67, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72,
	0x2e, 0x76, 0x31, 0x2e, 0x53, 0x61, 0x76, 0x65, 0x42, 0x61, 0x74, 0x63, 0x68, 0x55, 0x52, 0x4c,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x55, 0x52, 0x4c, 0x4f, 0x62, 0x6a, 0x65,
	0x63, 0x74, 0x52, 0x08, 0x75, 0x72, 0x6c, 0x42, 0x61, 0x74, 0x63, 0x68, 0x1a, 0x58, 0x0a, 0x09,
	0x55, 0x52, 0x4c, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x25, 0x0a, 0x0e, 0x63, 0x6f, 0x72,
	0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0d, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64,
	0x12, 0x24, 0x0a, 0x09, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x08, 0x73, 0x68,
	0x6f, 0x72, 0x74, 0x55, 0x72, 0x6c, 0x22, 0x3c, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x53, 0x74, 0x61,
	0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x72,
	0x6c, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x75, 0x72, 0x6c, 0x73, 0x12, 0x14,
	0x0a, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x75,
	0x73, 0x65, 0x72, 0x73, 0x32, 0xf2, 0x02, 0x0a, 0x10, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e,
	0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x56, 0x0a, 0x07, 0x53, 0x61, 0x76,
	0x65, 0x55, 0x52, 0x4c, 0x12, 0x24, 0x2e, 0x70, 0x6b, 0x67, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73,
	0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x61, 0x76, 0x65,
	0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x70, 0x6b, 0x67,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x2e, 0x76,
	0x31, 0x2e, 0x53, 0x61, 0x76, 0x65, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x53, 0x0a, 0x06, 0x47, 0x65, 0x74, 0x55, 0x52, 0x4c, 0x12, 0x23, 0x2e, 0x70, 0x6b,
	0x67, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x2e,
	0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x24, 0x2e, 0x70, 0x6b, 0x67, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74,
	0x65, 0x6e, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x52, 0x4c, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x65, 0x0a, 0x0c, 0x53, 0x61, 0x76, 0x65, 0x42, 0x61,
	0x74, 0x63, 0x68, 0x55, 0x52, 0x4c, 0x12, 0x29, 0x2e, 0x70, 0x6b, 0x67, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x61,
	0x76, 0x65, 0x42, 0x61, 0x74, 0x63, 0x68, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x2a, 0x2e, 0x70, 0x6b, 0x67, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x68, 0x6f, 0x72,
	0x74, 0x65, 0x6e, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x61, 0x76, 0x65, 0x42, 0x61, 0x74,
	0x63, 0x68, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4a, 0x0a,
	0x08, 0x47, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x1a, 0x26, 0x2e, 0x70, 0x6b, 0x67, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x73, 0x68, 0x6f, 0x72,
	0x74, 0x65, 0x6e, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x2f, 0x5a, 0x2d, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6b, 0x72, 0x69, 0x70, 0x73, 0x79, 0x2f, 0x73,
	0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_pkg_api_shortener_v1_service_proto_rawDescOnce sync.Once
	file_pkg_api_shortener_v1_service_proto_rawDescData = file_pkg_api_shortener_v1_service_proto_rawDesc
)

func file_pkg_api_shortener_v1_service_proto_rawDescGZIP() []byte {
	file_pkg_api_shortener_v1_service_proto_rawDescOnce.Do(func() {
		file_pkg_api_shortener_v1_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_api_shortener_v1_service_proto_rawDescData)
	})
	return file_pkg_api_shortener_v1_service_proto_rawDescData
}

var file_pkg_api_shortener_v1_service_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_pkg_api_shortener_v1_service_proto_goTypes = []interface{}{
	(*SaveURLRequest)(nil),                 // 0: pkg.api.shortener.v1.SaveURLRequest
	(*SaveURLResponse)(nil),                // 1: pkg.api.shortener.v1.SaveURLResponse
	(*GetURLRequest)(nil),                  // 2: pkg.api.shortener.v1.GetURLRequest
	(*GetURLResponse)(nil),                 // 3: pkg.api.shortener.v1.GetURLResponse
	(*SaveBatchURLRequest)(nil),            // 4: pkg.api.shortener.v1.SaveBatchURLRequest
	(*SaveBatchURLResponse)(nil),           // 5: pkg.api.shortener.v1.SaveBatchURLResponse
	(*GetStatsResponse)(nil),               // 6: pkg.api.shortener.v1.GetStatsResponse
	(*SaveBatchURLRequest_URLObject)(nil),  // 7: pkg.api.shortener.v1.SaveBatchURLRequest.URLObject
	(*SaveBatchURLResponse_URLObject)(nil), // 8: pkg.api.shortener.v1.SaveBatchURLResponse.URLObject
	(*emptypb.Empty)(nil),                  // 9: google.protobuf.Empty
}
var file_pkg_api_shortener_v1_service_proto_depIdxs = []int32{
	7, // 0: pkg.api.shortener.v1.SaveBatchURLRequest.url_batch:type_name -> pkg.api.shortener.v1.SaveBatchURLRequest.URLObject
	8, // 1: pkg.api.shortener.v1.SaveBatchURLResponse.url_batch:type_name -> pkg.api.shortener.v1.SaveBatchURLResponse.URLObject
	0, // 2: pkg.api.shortener.v1.ShortenerService.SaveURL:input_type -> pkg.api.shortener.v1.SaveURLRequest
	2, // 3: pkg.api.shortener.v1.ShortenerService.GetURL:input_type -> pkg.api.shortener.v1.GetURLRequest
	4, // 4: pkg.api.shortener.v1.ShortenerService.SaveBatchURL:input_type -> pkg.api.shortener.v1.SaveBatchURLRequest
	9, // 5: pkg.api.shortener.v1.ShortenerService.GetStats:input_type -> google.protobuf.Empty
	1, // 6: pkg.api.shortener.v1.ShortenerService.SaveURL:output_type -> pkg.api.shortener.v1.SaveURLResponse
	3, // 7: pkg.api.shortener.v1.ShortenerService.GetURL:output_type -> pkg.api.shortener.v1.GetURLResponse
	5, // 8: pkg.api.shortener.v1.ShortenerService.SaveBatchURL:output_type -> pkg.api.shortener.v1.SaveBatchURLResponse
	6, // 9: pkg.api.shortener.v1.ShortenerService.GetStats:output_type -> pkg.api.shortener.v1.GetStatsResponse
	6, // [6:10] is the sub-list for method output_type
	2, // [2:6] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_pkg_api_shortener_v1_service_proto_init() }
func file_pkg_api_shortener_v1_service_proto_init() {
	if File_pkg_api_shortener_v1_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_api_shortener_v1_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SaveURLRequest); i {
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
		file_pkg_api_shortener_v1_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SaveURLResponse); i {
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
		file_pkg_api_shortener_v1_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetURLRequest); i {
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
		file_pkg_api_shortener_v1_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetURLResponse); i {
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
		file_pkg_api_shortener_v1_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SaveBatchURLRequest); i {
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
		file_pkg_api_shortener_v1_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SaveBatchURLResponse); i {
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
		file_pkg_api_shortener_v1_service_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetStatsResponse); i {
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
		file_pkg_api_shortener_v1_service_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SaveBatchURLRequest_URLObject); i {
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
		file_pkg_api_shortener_v1_service_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SaveBatchURLResponse_URLObject); i {
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
			RawDescriptor: file_pkg_api_shortener_v1_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_api_shortener_v1_service_proto_goTypes,
		DependencyIndexes: file_pkg_api_shortener_v1_service_proto_depIdxs,
		MessageInfos:      file_pkg_api_shortener_v1_service_proto_msgTypes,
	}.Build()
	File_pkg_api_shortener_v1_service_proto = out.File
	file_pkg_api_shortener_v1_service_proto_rawDesc = nil
	file_pkg_api_shortener_v1_service_proto_goTypes = nil
	file_pkg_api_shortener_v1_service_proto_depIdxs = nil
}

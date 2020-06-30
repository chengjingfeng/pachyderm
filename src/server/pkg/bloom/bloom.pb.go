// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.12.3
// source: server/pkg/bloom/bloom.proto

package bloom

import (
	proto "github.com/golang/protobuf/proto"
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

type BloomFilter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NumSubhashes uint32 `protobuf:"varint,1,opt,name=num_subhashes,json=numSubhashes,proto3" json:"num_subhashes,omitempty"`
	// TODO: we could make each bucket signed, which would allow us to do some interesting things,
	// but it might make the design a little confusing.
	// Two BloomFilters with identical hash_length and len(buckets) can be summed
	// to produce a new filter which should be identical to running all operations
	// from the original two filters onto a blank filter.
	// Negative bucket values may be useful if we were to add a background
	// reprocessing stage that would iterate over all existing items and re-add
	// them, but would also track live updates to the set.  Due to live updates,
	// some buckets may need to go negative temporarily - but we would lose the
	// guarantee that removing something that wasn't added to the set is an error.
	// Perhaps better to provide a 'BloomFilterDelta' that can later be combined
	// into an existing filter, while still preserving the invariant that all
	// buckets are positive.
	Buckets []uint32 `protobuf:"varint,2,rep,packed,name=buckets,proto3" json:"buckets,omitempty"`
}

func (x *BloomFilter) Reset() {
	*x = BloomFilter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_server_pkg_bloom_bloom_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BloomFilter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BloomFilter) ProtoMessage() {}

func (x *BloomFilter) ProtoReflect() protoreflect.Message {
	mi := &file_server_pkg_bloom_bloom_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BloomFilter.ProtoReflect.Descriptor instead.
func (*BloomFilter) Descriptor() ([]byte, []int) {
	return file_server_pkg_bloom_bloom_proto_rawDescGZIP(), []int{0}
}

func (x *BloomFilter) GetNumSubhashes() uint32 {
	if x != nil {
		return x.NumSubhashes
	}
	return 0
}

func (x *BloomFilter) GetBuckets() []uint32 {
	if x != nil {
		return x.Buckets
	}
	return nil
}

var File_server_pkg_bloom_bloom_proto protoreflect.FileDescriptor

var file_server_pkg_bloom_bloom_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x62, 0x6c, 0x6f,
	0x6f, 0x6d, 0x2f, 0x62, 0x6c, 0x6f, 0x6f, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05,
	0x62, 0x6c, 0x6f, 0x6f, 0x6d, 0x22, 0x4c, 0x0a, 0x0b, 0x42, 0x6c, 0x6f, 0x6f, 0x6d, 0x46, 0x69,
	0x6c, 0x74, 0x65, 0x72, 0x12, 0x23, 0x0a, 0x0d, 0x6e, 0x75, 0x6d, 0x5f, 0x73, 0x75, 0x62, 0x68,
	0x61, 0x73, 0x68, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0c, 0x6e, 0x75, 0x6d,
	0x53, 0x75, 0x62, 0x68, 0x61, 0x73, 0x68, 0x65, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x62, 0x75, 0x63,
	0x6b, 0x65, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x07, 0x62, 0x75, 0x63, 0x6b,
	0x65, 0x74, 0x73, 0x42, 0x35, 0x5a, 0x33, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x70, 0x61, 0x63, 0x68, 0x79, 0x64, 0x65, 0x72, 0x6d, 0x2f, 0x70, 0x61, 0x63, 0x68,
	0x79, 0x64, 0x65, 0x72, 0x6d, 0x2f, 0x73, 0x72, 0x63, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x62, 0x6c, 0x6f, 0x6f, 0x6d, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_server_pkg_bloom_bloom_proto_rawDescOnce sync.Once
	file_server_pkg_bloom_bloom_proto_rawDescData = file_server_pkg_bloom_bloom_proto_rawDesc
)

func file_server_pkg_bloom_bloom_proto_rawDescGZIP() []byte {
	file_server_pkg_bloom_bloom_proto_rawDescOnce.Do(func() {
		file_server_pkg_bloom_bloom_proto_rawDescData = protoimpl.X.CompressGZIP(file_server_pkg_bloom_bloom_proto_rawDescData)
	})
	return file_server_pkg_bloom_bloom_proto_rawDescData
}

var file_server_pkg_bloom_bloom_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_server_pkg_bloom_bloom_proto_goTypes = []interface{}{
	(*BloomFilter)(nil), // 0: bloom.BloomFilter
}
var file_server_pkg_bloom_bloom_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_server_pkg_bloom_bloom_proto_init() }
func file_server_pkg_bloom_bloom_proto_init() {
	if File_server_pkg_bloom_bloom_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_server_pkg_bloom_bloom_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BloomFilter); i {
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
			RawDescriptor: file_server_pkg_bloom_bloom_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_server_pkg_bloom_bloom_proto_goTypes,
		DependencyIndexes: file_server_pkg_bloom_bloom_proto_depIdxs,
		MessageInfos:      file_server_pkg_bloom_bloom_proto_msgTypes,
	}.Build()
	File_server_pkg_bloom_bloom_proto = out.File
	file_server_pkg_bloom_bloom_proto_rawDesc = nil
	file_server_pkg_bloom_bloom_proto_goTypes = nil
	file_server_pkg_bloom_bloom_proto_depIdxs = nil
}

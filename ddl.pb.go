// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.22.0
// 	protoc        v3.11.4
// source: ddl.proto

package ddl

import (
	proto "github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
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

type IndexOption struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Columns []string `protobuf:"bytes,2,rep,name=columns,proto3" json:"columns,omitempty"`
	Unique  bool     `protobuf:"varint,3,opt,name=unique,proto3" json:"unique,omitempty"`
}

func (x *IndexOption) Reset() {
	*x = IndexOption{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ddl_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IndexOption) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IndexOption) ProtoMessage() {}

func (x *IndexOption) ProtoReflect() protoreflect.Message {
	mi := &file_ddl_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IndexOption.ProtoReflect.Descriptor instead.
func (*IndexOption) Descriptor() ([]byte, []int) {
	return file_ddl_proto_rawDescGZIP(), []int{0}
}

func (x *IndexOption) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *IndexOption) GetColumns() []string {
	if x != nil {
		return x.Columns
	}
	return nil
}

func (x *IndexOption) GetUnique() bool {
	if x != nil {
		return x.Unique
	}
	return false
}

type ColumnOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sequence bool   `protobuf:"varint,1,opt,name=sequence,proto3" json:"sequence,omitempty"`
	Null     bool   `protobuf:"varint,2,opt,name=null,proto3" json:"null,omitempty"`
	Default  string `protobuf:"bytes,3,opt,name=default,proto3" json:"default,omitempty"`
	Size     int32  `protobuf:"varint,4,opt,name=size,proto3" json:"size,omitempty"`
	Type     string `protobuf:"bytes,5,opt,name=type,proto3" json:"type,omitempty"`
}

func (x *ColumnOptions) Reset() {
	*x = ColumnOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ddl_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ColumnOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ColumnOptions) ProtoMessage() {}

func (x *ColumnOptions) ProtoReflect() protoreflect.Message {
	mi := &file_ddl_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ColumnOptions.ProtoReflect.Descriptor instead.
func (*ColumnOptions) Descriptor() ([]byte, []int) {
	return file_ddl_proto_rawDescGZIP(), []int{1}
}

func (x *ColumnOptions) GetSequence() bool {
	if x != nil {
		return x.Sequence
	}
	return false
}

func (x *ColumnOptions) GetNull() bool {
	if x != nil {
		return x.Null
	}
	return false
}

func (x *ColumnOptions) GetDefault() string {
	if x != nil {
		return x.Default
	}
	return ""
}

func (x *ColumnOptions) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *ColumnOptions) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

type TableOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TableName  string         `protobuf:"bytes,1,opt,name=table_name,json=tableName,proto3" json:"table_name,omitempty"`
	PrimaryKey []string       `protobuf:"bytes,2,rep,name=primary_key,json=primaryKey,proto3" json:"primary_key,omitempty"`
	Indexes    []*IndexOption `protobuf:"bytes,3,rep,name=indexes,proto3" json:"indexes,omitempty"`
	Engine     string         `protobuf:"bytes,4,opt,name=engine,proto3" json:"engine,omitempty"`
}

func (x *TableOptions) Reset() {
	*x = TableOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ddl_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TableOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TableOptions) ProtoMessage() {}

func (x *TableOptions) ProtoReflect() protoreflect.Message {
	mi := &file_ddl_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TableOptions.ProtoReflect.Descriptor instead.
func (*TableOptions) Descriptor() ([]byte, []int) {
	return file_ddl_proto_rawDescGZIP(), []int{2}
}

func (x *TableOptions) GetTableName() string {
	if x != nil {
		return x.TableName
	}
	return ""
}

func (x *TableOptions) GetPrimaryKey() []string {
	if x != nil {
		return x.PrimaryKey
	}
	return nil
}

func (x *TableOptions) GetIndexes() []*IndexOption {
	if x != nil {
		return x.Indexes
	}
	return nil
}

func (x *TableOptions) GetEngine() string {
	if x != nil {
		return x.Engine
	}
	return ""
}

var file_ddl_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptor.MessageOptions)(nil),
		ExtensionType: (*TableOptions)(nil),
		Field:         60000,
		Name:          "dev.f110.ddl.table",
		Tag:           "bytes,60000,opt,name=table",
		Filename:      "ddl.proto",
	},
	{
		ExtendedType:  (*descriptor.FieldOptions)(nil),
		ExtensionType: (*ColumnOptions)(nil),
		Field:         60000,
		Name:          "dev.f110.ddl.column",
		Tag:           "bytes,60000,opt,name=column",
		Filename:      "ddl.proto",
	},
}

// Extension fields to descriptor.MessageOptions.
var (
	// optional dev.f110.ddl.TableOptions table = 60000;
	E_Table = &file_ddl_proto_extTypes[0]
)

// Extension fields to descriptor.FieldOptions.
var (
	// optional dev.f110.ddl.ColumnOptions column = 60000;
	E_Column = &file_ddl_proto_extTypes[1]
)

var File_ddl_proto protoreflect.FileDescriptor

var file_ddl_proto_rawDesc = []byte{
	0x0a, 0x09, 0x64, 0x64, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x64, 0x65, 0x76,
	0x2e, 0x66, 0x31, 0x31, 0x30, 0x2e, 0x64, 0x64, 0x6c, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x53, 0x0a, 0x0b, 0x49,
	0x6e, 0x64, 0x65, 0x78, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x63, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x07, 0x63, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x6e, 0x69, 0x71,
	0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x75, 0x6e, 0x69, 0x71, 0x75, 0x65,
	0x22, 0x81, 0x01, 0x0a, 0x0d, 0x43, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x4f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x75, 0x6c, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x6e, 0x75,
	0x6c, 0x6c, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x12, 0x12, 0x0a, 0x04,
	0x73, 0x69, 0x7a, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x22, 0x9b, 0x01, 0x0a, 0x0c, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x4f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x61, 0x62, 0x6c, 0x65,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x72, 0x69, 0x6d, 0x61, 0x72, 0x79, 0x5f,
	0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x72, 0x69, 0x6d, 0x61,
	0x72, 0x79, 0x4b, 0x65, 0x79, 0x12, 0x33, 0x0a, 0x07, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x73,
	0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x66, 0x31, 0x31,
	0x30, 0x2e, 0x64, 0x64, 0x6c, 0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x4f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x07, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x6e,
	0x67, 0x69, 0x6e, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x65, 0x6e, 0x67, 0x69,
	0x6e, 0x65, 0x3a, 0x53, 0x0a, 0x05, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x1f, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xe0, 0xd4, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x66, 0x31, 0x31, 0x30, 0x2e,
	0x64, 0x64, 0x6c, 0x2e, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x52, 0x05, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x3a, 0x54, 0x0a, 0x06, 0x63, 0x6f, 0x6c, 0x75, 0x6d,
	0x6e, 0x12, 0x1d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x18, 0xe0, 0xd4, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x64, 0x65, 0x76, 0x2e, 0x66,
	0x31, 0x31, 0x30, 0x2e, 0x64, 0x64, 0x6c, 0x2e, 0x43, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x4f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x06, 0x63, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x42, 0x1c, 0x5a,
	0x1a, 0x67, 0x6f, 0x2e, 0x66, 0x31, 0x31, 0x30, 0x2e, 0x64, 0x65, 0x76, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x63, 0x2d, 0x64, 0x64, 0x6c, 0x3b, 0x64, 0x64, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_ddl_proto_rawDescOnce sync.Once
	file_ddl_proto_rawDescData = file_ddl_proto_rawDesc
)

func file_ddl_proto_rawDescGZIP() []byte {
	file_ddl_proto_rawDescOnce.Do(func() {
		file_ddl_proto_rawDescData = protoimpl.X.CompressGZIP(file_ddl_proto_rawDescData)
	})
	return file_ddl_proto_rawDescData
}

var file_ddl_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_ddl_proto_goTypes = []interface{}{
	(*IndexOption)(nil),               // 0: dev.f110.ddl.IndexOption
	(*ColumnOptions)(nil),             // 1: dev.f110.ddl.ColumnOptions
	(*TableOptions)(nil),              // 2: dev.f110.ddl.TableOptions
	(*descriptor.MessageOptions)(nil), // 3: google.protobuf.MessageOptions
	(*descriptor.FieldOptions)(nil),   // 4: google.protobuf.FieldOptions
}
var file_ddl_proto_depIdxs = []int32{
	0, // 0: dev.f110.ddl.TableOptions.indexes:type_name -> dev.f110.ddl.IndexOption
	3, // 1: dev.f110.ddl.table:extendee -> google.protobuf.MessageOptions
	4, // 2: dev.f110.ddl.column:extendee -> google.protobuf.FieldOptions
	2, // 3: dev.f110.ddl.table:type_name -> dev.f110.ddl.TableOptions
	1, // 4: dev.f110.ddl.column:type_name -> dev.f110.ddl.ColumnOptions
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	3, // [3:5] is the sub-list for extension type_name
	1, // [1:3] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_ddl_proto_init() }
func file_ddl_proto_init() {
	if File_ddl_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_ddl_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IndexOption); i {
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
		file_ddl_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ColumnOptions); i {
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
		file_ddl_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TableOptions); i {
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
			RawDescriptor: file_ddl_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 2,
			NumServices:   0,
		},
		GoTypes:           file_ddl_proto_goTypes,
		DependencyIndexes: file_ddl_proto_depIdxs,
		MessageInfos:      file_ddl_proto_msgTypes,
		ExtensionInfos:    file_ddl_proto_extTypes,
	}.Build()
	File_ddl_proto = out.File
	file_ddl_proto_rawDesc = nil
	file_ddl_proto_goTypes = nil
	file_ddl_proto_depIdxs = nil
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.9.1
// source: model.proto

package executor

import (
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

type TestCase struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Unique ID for the test case.
	Uid string `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	// Human readable name for the case.
	DisplayName string `protobuf:"bytes,2,opt,name=display_name,json=displayName,proto3" json:"display_name,omitempty"`
	// Markdown description of the test case.
	DescriptionMarkdown string `protobuf:"bytes,3,opt,name=description_markdown,json=descriptionMarkdown,proto3" json:"description_markdown,omitempty"`
	// Labels for selecting the test case.
	Labels map[string]string `protobuf:"bytes,4,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Types that are assignable to TestType:
	//	*TestCase_Skip
	//	*TestCase_Eval
	//	*TestCase_CaptureEval
	TestType isTestCase_TestType `protobuf_oneof:"test_type"`
}

func (x *TestCase) Reset() {
	*x = TestCase{}
	if protoimpl.UnsafeEnabled {
		mi := &file_model_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestCase) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestCase) ProtoMessage() {}

func (x *TestCase) ProtoReflect() protoreflect.Message {
	mi := &file_model_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestCase.ProtoReflect.Descriptor instead.
func (*TestCase) Descriptor() ([]byte, []int) {
	return file_model_proto_rawDescGZIP(), []int{0}
}

func (x *TestCase) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

func (x *TestCase) GetDisplayName() string {
	if x != nil {
		return x.DisplayName
	}
	return ""
}

func (x *TestCase) GetDescriptionMarkdown() string {
	if x != nil {
		return x.DescriptionMarkdown
	}
	return ""
}

func (x *TestCase) GetLabels() map[string]string {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (m *TestCase) GetTestType() isTestCase_TestType {
	if m != nil {
		return m.TestType
	}
	return nil
}

func (x *TestCase) GetSkip() *SkipTest {
	if x, ok := x.GetTestType().(*TestCase_Skip); ok {
		return x.Skip
	}
	return nil
}

func (x *TestCase) GetEval() *EvalTest {
	if x, ok := x.GetTestType().(*TestCase_Eval); ok {
		return x.Eval
	}
	return nil
}

func (x *TestCase) GetCaptureEval() *CaptureEval {
	if x, ok := x.GetTestType().(*TestCase_CaptureEval); ok {
		return x.CaptureEval
	}
	return nil
}

type isTestCase_TestType interface {
	isTestCase_TestType()
}

type TestCase_Skip struct {
	Skip *SkipTest `protobuf:"bytes,5,opt,name=skip,proto3,oneof"`
}

type TestCase_Eval struct {
	Eval *EvalTest `protobuf:"bytes,6,opt,name=eval,proto3,oneof"`
}

type TestCase_CaptureEval struct {
	CaptureEval *CaptureEval `protobuf:"bytes,7,opt,name=capture_eval,json=captureEval,proto3,oneof"`
}

func (*TestCase_Skip) isTestCase_TestType() {}

func (*TestCase_Eval) isTestCase_TestType() {}

func (*TestCase_CaptureEval) isTestCase_TestType() {}

// A test that's skipped.
type SkipTest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Reason the test was skipped.
	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *SkipTest) Reset() {
	*x = SkipTest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_model_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SkipTest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SkipTest) ProtoMessage() {}

func (x *SkipTest) ProtoReflect() protoreflect.Message {
	mi := &file_model_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SkipTest.ProtoReflect.Descriptor instead.
func (*SkipTest) Descriptor() ([]byte, []int) {
	return file_model_proto_rawDescGZIP(), []int{1}
}

func (x *SkipTest) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

// A test that's evaluated and checked against a value.
type EvalTest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Input string `protobuf:"bytes,1,opt,name=input,proto3" json:"input,omitempty"`
	// Types that are assignable to Expect:
	//	*EvalTest_Exact
	Expect isEvalTest_Expect `protobuf_oneof:"expect"`
}

func (x *EvalTest) Reset() {
	*x = EvalTest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_model_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EvalTest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EvalTest) ProtoMessage() {}

func (x *EvalTest) ProtoReflect() protoreflect.Message {
	mi := &file_model_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EvalTest.ProtoReflect.Descriptor instead.
func (*EvalTest) Descriptor() ([]byte, []int) {
	return file_model_proto_rawDescGZIP(), []int{2}
}

func (x *EvalTest) GetInput() string {
	if x != nil {
		return x.Input
	}
	return ""
}

func (m *EvalTest) GetExpect() isEvalTest_Expect {
	if m != nil {
		return m.Expect
	}
	return nil
}

func (x *EvalTest) GetExact() string {
	if x, ok := x.GetExpect().(*EvalTest_Exact); ok {
		return x.Exact
	}
	return ""
}

type isEvalTest_Expect interface {
	isEvalTest_Expect()
}

type EvalTest_Exact struct {
	Exact string `protobuf:"bytes,2,opt,name=exact,proto3,oneof"`
}

func (*EvalTest_Exact) isEvalTest_Expect() {}

// A test that is for undefined behavior, the output is captured.
type CaptureEval struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Input string `protobuf:"bytes,1,opt,name=input,proto3" json:"input,omitempty"`
}

func (x *CaptureEval) Reset() {
	*x = CaptureEval{}
	if protoimpl.UnsafeEnabled {
		mi := &file_model_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CaptureEval) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CaptureEval) ProtoMessage() {}

func (x *CaptureEval) ProtoReflect() protoreflect.Message {
	mi := &file_model_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CaptureEval.ProtoReflect.Descriptor instead.
func (*CaptureEval) Descriptor() ([]byte, []int) {
	return file_model_proto_rawDescGZIP(), []int{3}
}

func (x *CaptureEval) GetInput() string {
	if x != nil {
		return x.Input
	}
	return ""
}

var File_model_proto protoreflect.FileDescriptor

var file_model_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xde, 0x02,
	0x0a, 0x08, 0x54, 0x65, 0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x69, 0x64, 0x12, 0x21, 0x0a, 0x0c,
	0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x31, 0x0a, 0x14, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6d,
	0x61, 0x72, 0x6b, 0x64, 0x6f, 0x77, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x13, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x61, 0x72, 0x6b, 0x64, 0x6f,
	0x77, 0x6e, 0x12, 0x2d, 0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x04, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x15, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x2e, 0x4c, 0x61,
	0x62, 0x65, 0x6c, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c,
	0x73, 0x12, 0x1f, 0x0a, 0x04, 0x73, 0x6b, 0x69, 0x70, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x09, 0x2e, 0x53, 0x6b, 0x69, 0x70, 0x54, 0x65, 0x73, 0x74, 0x48, 0x00, 0x52, 0x04, 0x73, 0x6b,
	0x69, 0x70, 0x12, 0x1f, 0x0a, 0x04, 0x65, 0x76, 0x61, 0x6c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x09, 0x2e, 0x45, 0x76, 0x61, 0x6c, 0x54, 0x65, 0x73, 0x74, 0x48, 0x00, 0x52, 0x04, 0x65,
	0x76, 0x61, 0x6c, 0x12, 0x31, 0x0a, 0x0c, 0x63, 0x61, 0x70, 0x74, 0x75, 0x72, 0x65, 0x5f, 0x65,
	0x76, 0x61, 0x6c, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x43, 0x61, 0x70, 0x74,
	0x75, 0x72, 0x65, 0x45, 0x76, 0x61, 0x6c, 0x48, 0x00, 0x52, 0x0b, 0x63, 0x61, 0x70, 0x74, 0x75,
	0x72, 0x65, 0x45, 0x76, 0x61, 0x6c, 0x1a, 0x39, 0x0a, 0x0b, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x42, 0x0b, 0x0a, 0x09, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x22, 0x24,
	0x0a, 0x08, 0x53, 0x6b, 0x69, 0x70, 0x54, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x22, 0x42, 0x0a, 0x08, 0x45, 0x76, 0x61, 0x6c, 0x54, 0x65, 0x73, 0x74,
	0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x12, 0x16, 0x0a, 0x05, 0x65, 0x78, 0x61, 0x63, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x05, 0x65, 0x78, 0x61, 0x63, 0x74, 0x42, 0x08,
	0x0a, 0x06, 0x65, 0x78, 0x70, 0x65, 0x63, 0x74, 0x22, 0x23, 0x0a, 0x0b, 0x43, 0x61, 0x70, 0x74,
	0x75, 0x72, 0x65, 0x45, 0x76, 0x61, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x42, 0x36, 0x5a,
	0x34, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6a, 0x6f, 0x73, 0x65,
	0x70, 0x68, 0x6c, 0x65, 0x77, 0x69, 0x73, 0x34, 0x32, 0x2f, 0x73, 0x70, 0x65, 0x63, 0x2d, 0x74,
	0x65, 0x73, 0x74, 0x65, 0x72, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x65, 0x72, 0x2f, 0x65, 0x78, 0x65,
	0x63, 0x75, 0x74, 0x6f, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_model_proto_rawDescOnce sync.Once
	file_model_proto_rawDescData = file_model_proto_rawDesc
)

func file_model_proto_rawDescGZIP() []byte {
	file_model_proto_rawDescOnce.Do(func() {
		file_model_proto_rawDescData = protoimpl.X.CompressGZIP(file_model_proto_rawDescData)
	})
	return file_model_proto_rawDescData
}

var file_model_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_model_proto_goTypes = []interface{}{
	(*TestCase)(nil),    // 0: TestCase
	(*SkipTest)(nil),    // 1: SkipTest
	(*EvalTest)(nil),    // 2: EvalTest
	(*CaptureEval)(nil), // 3: CaptureEval
	nil,                 // 4: TestCase.LabelsEntry
}
var file_model_proto_depIdxs = []int32{
	4, // 0: TestCase.labels:type_name -> TestCase.LabelsEntry
	1, // 1: TestCase.skip:type_name -> SkipTest
	2, // 2: TestCase.eval:type_name -> EvalTest
	3, // 3: TestCase.capture_eval:type_name -> CaptureEval
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_model_proto_init() }
func file_model_proto_init() {
	if File_model_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_model_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestCase); i {
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
		file_model_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SkipTest); i {
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
		file_model_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EvalTest); i {
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
		file_model_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CaptureEval); i {
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
	file_model_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*TestCase_Skip)(nil),
		(*TestCase_Eval)(nil),
		(*TestCase_CaptureEval)(nil),
	}
	file_model_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*EvalTest_Exact)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_model_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_model_proto_goTypes,
		DependencyIndexes: file_model_proto_depIdxs,
		MessageInfos:      file_model_proto_msgTypes,
	}.Build()
	File_model_proto = out.File
	file_model_proto_rawDesc = nil
	file_model_proto_goTypes = nil
	file_model_proto_depIdxs = nil
}
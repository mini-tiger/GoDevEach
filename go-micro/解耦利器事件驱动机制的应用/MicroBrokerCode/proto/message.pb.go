// Code generated by protoc-gen-go. DO NOT EDIT.
// source: message.proto

package message

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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

//学生数据体
type Student struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Classes              string   `protobuf:"bytes,2,opt,name=classes,proto3" json:"classes,omitempty"`
	Grade                int32    `protobuf:"varint,3,opt,name=grade,proto3" json:"grade,omitempty"`
	Phone                string   `protobuf:"bytes,4,opt,name=phone,proto3" json:"phone,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Student) Reset()         { *m = Student{} }
func (m *Student) String() string { return proto.CompactTextString(m) }
func (*Student) ProtoMessage()    {}
func (*Student) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{0}
}

func (m *Student) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Student.Unmarshal(m, b)
}
func (m *Student) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Student.Marshal(b, m, deterministic)
}
func (m *Student) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Student.Merge(m, src)
}
func (m *Student) XXX_Size() int {
	return xxx_messageInfo_Student.Size(m)
}
func (m *Student) XXX_DiscardUnknown() {
	xxx_messageInfo_Student.DiscardUnknown(m)
}

var xxx_messageInfo_Student proto.InternalMessageInfo

func (m *Student) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Student) GetClasses() string {
	if m != nil {
		return m.Classes
	}
	return ""
}

func (m *Student) GetGrade() int32 {
	if m != nil {
		return m.Grade
	}
	return 0
}

func (m *Student) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

//请求数据体定义
type StudentRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StudentRequest) Reset()         { *m = StudentRequest{} }
func (m *StudentRequest) String() string { return proto.CompactTextString(m) }
func (*StudentRequest) ProtoMessage()    {}
func (*StudentRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{1}
}

func (m *StudentRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StudentRequest.Unmarshal(m, b)
}
func (m *StudentRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StudentRequest.Marshal(b, m, deterministic)
}
func (m *StudentRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StudentRequest.Merge(m, src)
}
func (m *StudentRequest) XXX_Size() int {
	return xxx_messageInfo_StudentRequest.Size(m)
}
func (m *StudentRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StudentRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StudentRequest proto.InternalMessageInfo

func (m *StudentRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func init() {
	proto.RegisterType((*Student)(nil), "message.Student")
	proto.RegisterType((*StudentRequest)(nil), "message.StudentRequest")
}

func init() { proto.RegisterFile("message.proto", fileDescriptor_33c57e4bae7b9afd) }

var fileDescriptor_33c57e4bae7b9afd = []byte{
	// 137 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcd, 0x4d, 0x2d, 0x2e,
	0x4e, 0x4c, 0x4f, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x87, 0x72, 0x95, 0x92, 0xb9,
	0xd8, 0x83, 0x4b, 0x4a, 0x53, 0x52, 0xf3, 0x4a, 0x84, 0x84, 0xb8, 0x58, 0xf2, 0x12, 0x73, 0x53,
	0x25, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0xc0, 0x6c, 0x21, 0x09, 0x2e, 0xf6, 0xe4, 0x9c, 0xc4,
	0xe2, 0xe2, 0xd4, 0x62, 0x09, 0x26, 0xb0, 0x30, 0x8c, 0x2b, 0x24, 0xc2, 0xc5, 0x9a, 0x5e, 0x94,
	0x98, 0x92, 0x2a, 0xc1, 0xac, 0xc0, 0xa8, 0xc1, 0x1a, 0x04, 0xe1, 0x80, 0x44, 0x0b, 0x32, 0xf2,
	0xf3, 0x52, 0x25, 0x58, 0xc0, 0xaa, 0x21, 0x1c, 0x25, 0x15, 0x2e, 0x3e, 0xa8, 0x25, 0x41, 0xa9,
	0x85, 0xa5, 0xa9, 0xc5, 0x58, 0xed, 0x4a, 0x62, 0x03, 0x3b, 0xcd, 0x18, 0x10, 0x00, 0x00, 0xff,
	0xff, 0x0d, 0x79, 0xd1, 0x66, 0xab, 0x00, 0x00, 0x00,
}

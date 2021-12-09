// Code generated by protoc-gen-go. DO NOT EDIT.
// source: products.proto

package pb

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

type Product struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Product) Reset()         { *m = Product{} }
func (m *Product) String() string { return proto.CompactTextString(m) }
func (*Product) ProtoMessage()    {}
func (*Product) Descriptor() ([]byte, []int) {
	return fileDescriptor_8c6e54f42122eb82, []int{0}
}

func (m *Product) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Product.Unmarshal(m, b)
}
func (m *Product) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Product.Marshal(b, m, deterministic)
}
func (m *Product) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Product.Merge(m, src)
}
func (m *Product) XXX_Size() int {
	return xxx_messageInfo_Product.Size(m)
}
func (m *Product) XXX_DiscardUnknown() {
	xxx_messageInfo_Product.DiscardUnknown(m)
}

var xxx_messageInfo_Product proto.InternalMessageInfo

func (m *Product) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Product) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type GetAllProductsResponse struct {
	Products             []*Product `protobuf:"bytes,1,rep,name=products,proto3" json:"products,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *GetAllProductsResponse) Reset()         { *m = GetAllProductsResponse{} }
func (m *GetAllProductsResponse) String() string { return proto.CompactTextString(m) }
func (*GetAllProductsResponse) ProtoMessage()    {}
func (*GetAllProductsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_8c6e54f42122eb82, []int{1}
}

func (m *GetAllProductsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAllProductsResponse.Unmarshal(m, b)
}
func (m *GetAllProductsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAllProductsResponse.Marshal(b, m, deterministic)
}
func (m *GetAllProductsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAllProductsResponse.Merge(m, src)
}
func (m *GetAllProductsResponse) XXX_Size() int {
	return xxx_messageInfo_GetAllProductsResponse.Size(m)
}
func (m *GetAllProductsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAllProductsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetAllProductsResponse proto.InternalMessageInfo

func (m *GetAllProductsResponse) GetProducts() []*Product {
	if m != nil {
		return m.Products
	}
	return nil
}

type GetAllProductsRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetAllProductsRequest) Reset()         { *m = GetAllProductsRequest{} }
func (m *GetAllProductsRequest) String() string { return proto.CompactTextString(m) }
func (*GetAllProductsRequest) ProtoMessage()    {}
func (*GetAllProductsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8c6e54f42122eb82, []int{2}
}

func (m *GetAllProductsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAllProductsRequest.Unmarshal(m, b)
}
func (m *GetAllProductsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAllProductsRequest.Marshal(b, m, deterministic)
}
func (m *GetAllProductsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAllProductsRequest.Merge(m, src)
}
func (m *GetAllProductsRequest) XXX_Size() int {
	return xxx_messageInfo_GetAllProductsRequest.Size(m)
}
func (m *GetAllProductsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAllProductsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetAllProductsRequest proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Product)(nil), "pb.Product")
	proto.RegisterType((*GetAllProductsResponse)(nil), "pb.GetAllProductsResponse")
	proto.RegisterType((*GetAllProductsRequest)(nil), "pb.GetAllProductsRequest")
}

func init() { proto.RegisterFile("products.proto", fileDescriptor_8c6e54f42122eb82) }

var fileDescriptor_8c6e54f42122eb82 = []byte{
	// 218 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2b, 0x28, 0xca, 0x4f,
	0x29, 0x4d, 0x2e, 0x29, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a, 0x48, 0x52, 0xd2,
	0xe5, 0x62, 0x0f, 0x80, 0x88, 0x0a, 0xf1, 0x71, 0x31, 0x65, 0xa6, 0x48, 0x30, 0x2a, 0x30, 0x6a,
	0x70, 0x06, 0x31, 0x65, 0xa6, 0x08, 0x09, 0x71, 0xb1, 0xe4, 0x25, 0xe6, 0xa6, 0x4a, 0x30, 0x81,
	0x45, 0xc0, 0x6c, 0x25, 0x47, 0x2e, 0x31, 0xf7, 0xd4, 0x12, 0xc7, 0x9c, 0x1c, 0xa8, 0xa6, 0xe2,
	0xa0, 0xd4, 0xe2, 0x82, 0xfc, 0xbc, 0xe2, 0x54, 0x21, 0x75, 0x2e, 0x0e, 0x98, 0xf1, 0x12, 0x8c,
	0x0a, 0xcc, 0x1a, 0xdc, 0x46, 0xdc, 0x7a, 0x05, 0x49, 0x7a, 0x50, 0x75, 0x41, 0x70, 0x49, 0x25,
	0x71, 0x2e, 0x51, 0x74, 0x23, 0x0a, 0x4b, 0x53, 0x8b, 0x4b, 0x8c, 0xa2, 0xb9, 0xf8, 0xa0, 0x42,
	0xc1, 0xa9, 0x45, 0x65, 0x99, 0xc9, 0xa9, 0x42, 0x9e, 0x5c, 0x7c, 0xa8, 0x4a, 0x85, 0x24, 0x41,
	0x66, 0x62, 0xd5, 0x2e, 0x25, 0x85, 0x4d, 0x0a, 0xe2, 0x38, 0x25, 0x06, 0x27, 0xcd, 0x28, 0xf5,
	0xf4, 0xcc, 0x92, 0x8c, 0xd2, 0x24, 0xbd, 0xe4, 0xfc, 0x5c, 0xfd, 0xdc, 0xc4, 0xa2, 0xcc, 0xd2,
	0xe2, 0xec, 0xcc, 0xdc, 0xdc, 0xcc, 0xbc, 0x44, 0xfd, 0xe2, 0xd2, 0x82, 0x82, 0x9c, 0xca, 0xf2,
	0xc4, 0x92, 0xe4, 0x0c, 0xeb, 0x82, 0xa4, 0x24, 0x36, 0x70, 0xe8, 0x18, 0x03, 0x02, 0x00, 0x00,
	0xff, 0xff, 0xe8, 0xdc, 0x89, 0x23, 0x2f, 0x01, 0x00, 0x00,
}

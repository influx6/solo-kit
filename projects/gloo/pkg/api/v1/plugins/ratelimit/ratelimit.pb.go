// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: ratelimit/ratelimit.proto

/*
Package ratelimit is a generated protocol buffer package.

It is generated from these files:
	ratelimit/ratelimit.proto

It has these top-level messages:
	RateLimit
	IngressRateLimit
*/
package ratelimit

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

import bytes "bytes"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type RateLimit_Unit int32

const (
	RateLimit_UNKNOWN RateLimit_Unit = 0
	RateLimit_SECOND  RateLimit_Unit = 1
	RateLimit_MINUTE  RateLimit_Unit = 2
	RateLimit_HOUR    RateLimit_Unit = 3
	RateLimit_DAY     RateLimit_Unit = 4
)

var RateLimit_Unit_name = map[int32]string{
	0: "UNKNOWN",
	1: "SECOND",
	2: "MINUTE",
	3: "HOUR",
	4: "DAY",
}
var RateLimit_Unit_value = map[string]int32{
	"UNKNOWN": 0,
	"SECOND":  1,
	"MINUTE":  2,
	"HOUR":    3,
	"DAY":     4,
}

func (x RateLimit_Unit) String() string {
	return proto.EnumName(RateLimit_Unit_name, int32(x))
}
func (RateLimit_Unit) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_ratelimit_2b7207b1bd97a3f0, []int{0, 0}
}

type RateLimit struct {
	Unit            RateLimit_Unit `protobuf:"varint,1,opt,name=unit,proto3,enum=ratelimit.plugins.gloo.solo.io.RateLimit_Unit" json:"unit,omitempty"`
	RequestsPerUnit uint32         `protobuf:"varint,2,opt,name=requests_per_unit,json=requestsPerUnit,proto3" json:"requests_per_unit,omitempty"`
}

func (m *RateLimit) Reset()         { *m = RateLimit{} }
func (m *RateLimit) String() string { return proto.CompactTextString(m) }
func (*RateLimit) ProtoMessage()    {}
func (*RateLimit) Descriptor() ([]byte, []int) {
	return fileDescriptor_ratelimit_2b7207b1bd97a3f0, []int{0}
}
func (m *RateLimit) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RateLimit.Unmarshal(m, b)
}
func (m *RateLimit) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RateLimit.Marshal(b, m, deterministic)
}
func (dst *RateLimit) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RateLimit.Merge(dst, src)
}
func (m *RateLimit) XXX_Size() int {
	return xxx_messageInfo_RateLimit.Size(m)
}
func (m *RateLimit) XXX_DiscardUnknown() {
	xxx_messageInfo_RateLimit.DiscardUnknown(m)
}

var xxx_messageInfo_RateLimit proto.InternalMessageInfo

func (m *RateLimit) GetUnit() RateLimit_Unit {
	if m != nil {
		return m.Unit
	}
	return RateLimit_UNKNOWN
}

func (m *RateLimit) GetRequestsPerUnit() uint32 {
	if m != nil {
		return m.RequestsPerUnit
	}
	return 0
}

type IngressRateLimit struct {
	AuthrorizedHeader string     `protobuf:"bytes,1,opt,name=authrorized_header,json=authrorizedHeader,proto3" json:"authrorized_header,omitempty"`
	AuthorizedLimits  *RateLimit `protobuf:"bytes,2,opt,name=authorized_limits,json=authorizedLimits" json:"authorized_limits,omitempty"`
	AnonymousLimits   *RateLimit `protobuf:"bytes,3,opt,name=anonymous_limits,json=anonymousLimits" json:"anonymous_limits,omitempty"`
}

func (m *IngressRateLimit) Reset()                    { *m = IngressRateLimit{} }
func (m *IngressRateLimit) String() string            { return proto.CompactTextString(m) }
func (*IngressRateLimit) ProtoMessage()               {}
func (*IngressRateLimit) Descriptor() ([]byte, []int) { return fileDescriptorRatelimit, []int{1} }

func (m *IngressRateLimit) GetAuthrorizedHeader() string {
	if m != nil {
		return m.AuthrorizedHeader
	}
	return ""
}

func (m *IngressRateLimit) GetAuthorizedLimits() *RateLimit {
	if m != nil {
		return m.AuthorizedLimits
	}
	return nil
}

func (m *IngressRateLimit) GetAnonymousLimits() *RateLimit {
	if m != nil {
		return m.AnonymousLimits
	}
	return nil
}

func init() {
	proto.RegisterType((*RateLimit)(nil), "ratelimit.plugins.gloo.solo.io.RateLimit")
	proto.RegisterType((*IngressRateLimit)(nil), "ratelimit.plugins.gloo.solo.io.IngressRateLimit")
	proto.RegisterEnum("ratelimit.plugins.gloo.solo.io.RateLimit_Unit", RateLimit_Unit_name, RateLimit_Unit_value)
}
func (this *RateLimit) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*RateLimit)
	if !ok {
		that2, ok := that.(RateLimit)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Unit != that1.Unit {
		return false
	}
	if this.RequestsPerUnit != that1.RequestsPerUnit {
		return false
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}
func (this *IngressRateLimit) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*IngressRateLimit)
	if !ok {
		that2, ok := that.(IngressRateLimit)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.AuthrorizedHeader != that1.AuthrorizedHeader {
		return false
	}
	if !this.AuthorizedLimits.Equal(that1.AuthorizedLimits) {
		return false
	}
	if !this.AnonymousLimits.Equal(that1.AnonymousLimits) {
		return false
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}

func init() {
	proto.RegisterFile("ratelimit/ratelimit.proto", fileDescriptor_ratelimit_2b7207b1bd97a3f0)
}

var fileDescriptorRatelimit = []byte{
	// 360 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x92, 0xdd, 0x4e, 0xea, 0x40,
	0x10, 0x80, 0x4f, 0xa1, 0x81, 0xc3, 0x92, 0x73, 0x58, 0x36, 0xe7, 0xe2, 0xe8, 0x05, 0x21, 0x5c,
	0xa1, 0x09, 0xbb, 0x11, 0xef, 0x4d, 0x44, 0x30, 0x10, 0xb1, 0x98, 0x4a, 0x35, 0x7a, 0x43, 0x0a,
	0x6c, 0xca, 0x4a, 0xe9, 0xd4, 0xdd, 0xad, 0x89, 0x3e, 0x91, 0xef, 0xe0, 0xdb, 0xf8, 0x02, 0xbe,
	0x82, 0xe9, 0xf2, 0x7b, 0x65, 0xf4, 0x6a, 0x26, 0x33, 0xf3, 0x7d, 0xed, 0x6c, 0x06, 0xed, 0x49,
	0x5f, 0xf3, 0x50, 0x2c, 0x84, 0x66, 0x9b, 0x8c, 0xc6, 0x12, 0x34, 0x90, 0xca, 0x4e, 0x21, 0x4c,
	0x02, 0x11, 0x29, 0x1a, 0x84, 0x00, 0x54, 0x41, 0x08, 0x54, 0xc0, 0xfe, 0xbf, 0x00, 0x02, 0x30,
	0xa3, 0x2c, 0xcd, 0x96, 0x54, 0xed, 0xcd, 0x42, 0x05, 0xd7, 0xd7, 0xbc, 0x9f, 0x82, 0xa4, 0x85,
	0xec, 0x24, 0x12, 0xfa, 0xbf, 0x55, 0xb5, 0xea, 0x7f, 0x9b, 0x94, 0x7e, 0xad, 0xa4, 0x1b, 0x90,
	0x7a, 0x91, 0xd0, 0xae, 0x61, 0xc9, 0x21, 0x2a, 0x4b, 0xfe, 0x98, 0x70, 0xa5, 0xd5, 0x28, 0xe6,
	0x72, 0x64, 0x84, 0x99, 0xaa, 0x55, 0xff, 0xe3, 0x96, 0xd6, 0x8d, 0x2b, 0x2e, 0x53, 0xa2, 0x76,
	0x82, 0xec, 0x34, 0x92, 0x22, 0xca, 0x7b, 0xce, 0x85, 0x33, 0xb8, 0x75, 0xf0, 0x2f, 0x82, 0x50,
	0xee, 0xba, 0x73, 0x36, 0x70, 0xda, 0xd8, 0x4a, 0xf3, 0xcb, 0x9e, 0xe3, 0x0d, 0x3b, 0x38, 0x43,
	0x7e, 0x23, 0xbb, 0x3b, 0xf0, 0x5c, 0x9c, 0x25, 0x79, 0x94, 0x6d, 0x9f, 0xde, 0x61, 0xbb, 0xf6,
	0x61, 0x21, 0xdc, 0x8b, 0x02, 0xc9, 0x95, 0xda, 0x2e, 0xd1, 0x40, 0xc4, 0x4f, 0xf4, 0x4c, 0x82,
	0x14, 0x2f, 0x7c, 0x3a, 0x9a, 0x71, 0x7f, 0xca, 0xa5, 0x59, 0xa9, 0xe0, 0x96, 0x77, 0x3a, 0x5d,
	0xd3, 0x20, 0x37, 0xc8, 0x14, 0x57, 0xd3, 0x66, 0x5b, 0x65, 0xfe, 0xb7, 0xd8, 0x3c, 0xf8, 0xf6,
	0x03, 0xb8, 0x78, 0xeb, 0x30, 0x05, 0x45, 0x86, 0x08, 0xfb, 0x11, 0x44, 0xcf, 0x0b, 0x48, 0xd4,
	0x5a, 0x9b, 0xfd, 0xa9, 0xb6, 0xb4, 0x51, 0x2c, 0xad, 0xad, 0xfe, 0xeb, 0x7b, 0xc5, 0xba, 0x3f,
	0x0f, 0x84, 0x9e, 0x25, 0x63, 0x3a, 0x81, 0x05, 0x4b, 0xa1, 0x86, 0x80, 0x65, 0x9c, 0x0b, 0xcd,
	0x62, 0x09, 0x0f, 0x7c, 0xa2, 0x15, 0x4b, 0x9d, 0x2c, 0x9e, 0x07, 0xcc, 0x8f, 0x05, 0x7b, 0x3a,
	0x62, 0xab, 0x6f, 0x6d, 0x2f, 0x67, 0x9c, 0x33, 0x47, 0x70, 0xfc, 0x19, 0x00, 0x00, 0xff, 0xff,
	0x4f, 0xf9, 0xd4, 0x0f, 0x57, 0x02, 0x00, 0x00,
}

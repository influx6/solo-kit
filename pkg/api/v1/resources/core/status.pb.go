// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: status.proto

package core // import "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

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

type Status_State int32

const (
	// Pending status indicates the resource has not yet been validated
	Status_Pending Status_State = 0
	// Accepted indicates the resource has been validated
	Status_Accepted Status_State = 1
	// Rejected indicates an invalid configuration by the user
	// Rejected resources may be propagated to the xDS server depending on their severity
	Status_Rejected Status_State = 2
)

var Status_State_name = map[int32]string{
	0: "Pending",
	1: "Accepted",
	2: "Rejected",
}
var Status_State_value = map[string]int32{
	"Pending":  0,
	"Accepted": 1,
	"Rejected": 2,
}

func (x Status_State) String() string {
	return proto.EnumName(Status_State_name, int32(x))
}
func (Status_State) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_status_0c7d7c069d2343ab, []int{0, 0}
}

// *
// Status indicates whether a resource has been (in)validated by a reporter in the system.
// Statuses are meant to be read-only by users
type Status struct {
	// State is the enum indicating the state of the resource
	State Status_State `protobuf:"varint,1,opt,name=state,proto3,enum=core.solo.io.Status_State" json:"state,omitempty"`
	// Reason is a description of the error for Rejected resources. If the resource is pending or accepted, this field will be empty
	Reason string `protobuf:"bytes,2,opt,name=reason,proto3" json:"reason,omitempty"`
	// Reference to the reporter who wrote this status
	ReportedBy string `protobuf:"bytes,3,opt,name=reported_by,json=reportedBy,proto3" json:"reported_by,omitempty"`
	// Reference to statuses (by resource-ref string: "Kind.Namespace.Name") of subresources of the parent resource
	SubresourceStatuses  map[string]*Status `protobuf:"bytes,4,rep,name=subresource_statuses,json=subresourceStatuses" json:"subresource_statuses,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *Status) Reset()         { *m = Status{} }
func (m *Status) String() string { return proto.CompactTextString(m) }
func (*Status) ProtoMessage()    {}
func (*Status) Descriptor() ([]byte, []int) {
	return fileDescriptor_status_0c7d7c069d2343ab, []int{0}
}
func (m *Status) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Status.Unmarshal(m, b)
}
func (m *Status) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Status.Marshal(b, m, deterministic)
}
func (dst *Status) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Status.Merge(dst, src)
}
func (m *Status) XXX_Size() int {
	return xxx_messageInfo_Status.Size(m)
}
func (m *Status) XXX_DiscardUnknown() {
	xxx_messageInfo_Status.DiscardUnknown(m)
}

var xxx_messageInfo_Status proto.InternalMessageInfo

func (m *Status) GetState() Status_State {
	if m != nil {
		return m.State
	}
	return Status_Pending
}

func (m *Status) GetReason() string {
	if m != nil {
		return m.Reason
	}
	return ""
}

func (m *Status) GetReportedBy() string {
	if m != nil {
		return m.ReportedBy
	}
	return ""
}

func (m *Status) GetSubresourceStatuses() map[string]*Status {
	if m != nil {
		return m.SubresourceStatuses
	}
	return nil
}

func init() {
	proto.RegisterType((*Status)(nil), "core.solo.io.Status")
	proto.RegisterMapType((map[string]*Status)(nil), "core.solo.io.Status.SubresourceStatusesEntry")
	proto.RegisterEnum("core.solo.io.Status_State", Status_State_name, Status_State_value)
}
func (this *Status) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Status)
	if !ok {
		that2, ok := that.(Status)
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
	if this.State != that1.State {
		return false
	}
	if this.Reason != that1.Reason {
		return false
	}
	if this.ReportedBy != that1.ReportedBy {
		return false
	}
	if len(this.SubresourceStatuses) != len(that1.SubresourceStatuses) {
		return false
	}
	for i := range this.SubresourceStatuses {
		if !this.SubresourceStatuses[i].Equal(that1.SubresourceStatuses[i]) {
			return false
		}
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}

func init() { proto.RegisterFile("status.proto", fileDescriptor_status_0c7d7c069d2343ab) }

var fileDescriptor_status_0c7d7c069d2343ab = []byte{
	// 315 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x51, 0x51, 0x4b, 0xf3, 0x30,
	0x14, 0xfd, 0xda, 0x7d, 0x9b, 0xee, 0x76, 0xc8, 0x88, 0x43, 0xca, 0x1e, 0x74, 0xec, 0x69, 0x08,
	0x6b, 0xe6, 0x44, 0x10, 0x7d, 0x72, 0xe0, 0xbb, 0x74, 0x6f, 0x22, 0xcc, 0xb6, 0xbb, 0xd4, 0xb8,
	0xd9, 0x5b, 0x92, 0x74, 0xd0, 0x7f, 0xe4, 0xef, 0x12, 0xfc, 0x1f, 0x92, 0x64, 0x43, 0xc1, 0xf9,
	0x74, 0xef, 0xc9, 0x39, 0x27, 0x39, 0x37, 0x17, 0x3a, 0x4a, 0x27, 0xba, 0x52, 0x51, 0x29, 0x49,
	0x13, 0xeb, 0x64, 0x24, 0x31, 0x52, 0xb4, 0xa6, 0x48, 0x50, 0xbf, 0x97, 0x53, 0x4e, 0x96, 0xe0,
	0xa6, 0x73, 0x9a, 0xe1, 0xa7, 0x0f, 0xad, 0xb9, 0x35, 0xb1, 0x09, 0x34, 0x8d, 0x1d, 0x43, 0x6f,
	0xe0, 0x8d, 0x8e, 0xa6, 0xfd, 0xe8, 0xa7, 0x3d, 0x72, 0x22, 0x5b, 0x30, 0x76, 0x42, 0x76, 0x02,
	0x2d, 0x89, 0x89, 0xa2, 0x22, 0xf4, 0x07, 0xde, 0xa8, 0x1d, 0x6f, 0x11, 0x3b, 0x83, 0x40, 0x62,
	0x49, 0x52, 0xe3, 0x72, 0x91, 0xd6, 0x61, 0xc3, 0x92, 0xb0, 0x3b, 0x9a, 0xd5, 0xec, 0x19, 0x7a,
	0xaa, 0x4a, 0x25, 0x2a, 0xaa, 0x64, 0x86, 0x0b, 0x97, 0x1a, 0x55, 0xf8, 0x7f, 0xd0, 0x18, 0x05,
	0xd3, 0xf1, 0xfe, 0x97, 0xbf, 0x0d, 0xf3, 0xad, 0xfe, 0xbe, 0xd0, 0xb2, 0x8e, 0x8f, 0xd5, 0x6f,
	0xa6, 0xff, 0x04, 0xe1, 0x5f, 0x06, 0xd6, 0x85, 0xc6, 0x0a, 0x6b, 0x3b, 0x66, 0x3b, 0x36, 0x2d,
	0x3b, 0x87, 0xe6, 0x26, 0x59, 0x57, 0x68, 0xe7, 0x08, 0xa6, 0xbd, 0x7d, 0x01, 0x62, 0x27, 0xb9,
	0xf1, 0xaf, 0xbd, 0xe1, 0x04, 0x9a, 0xf6, 0x23, 0x58, 0x00, 0x07, 0x0f, 0x58, 0x2c, 0x45, 0x91,
	0x77, 0xff, 0xb1, 0x0e, 0x1c, 0xde, 0x65, 0x19, 0x96, 0x1a, 0x97, 0x5d, 0xcf, 0xa0, 0x18, 0x5f,
	0x31, 0x33, 0xc8, 0x9f, 0xdd, 0xbe, 0x7f, 0x9c, 0x7a, 0x8f, 0x57, 0xb9, 0xd0, 0x2f, 0x55, 0x1a,
	0x65, 0xf4, 0xc6, 0xcd, 0xed, 0x63, 0x41, 0xae, 0xae, 0x84, 0xe6, 0xe5, 0x2a, 0xe7, 0x49, 0x29,
	0xf8, 0xe6, 0x82, 0xef, 0x72, 0x2b, 0x6e, 0x82, 0xa4, 0x2d, 0xbb, 0xab, 0xcb, 0xaf, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x85, 0xbb, 0x75, 0x2e, 0xdf, 0x01, 0x00, 0x00,
}

// Code generated by protoc-gen-go.
// source: pwr/pwr.proto
// DO NOT EDIT!

/*
Package pwr is a generated protocol buffer package.

It is generated from these files:
	pwr/pwr.proto

It has these top-level messages:
	PatchHeader
	SyncHeader
	SyncOp
	SignatureHeader
	BlockHash
	CompressionSettings
*/
package pwr

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto.ProtoPackageIsVersion1

type CompressionAlgorithm int32

const (
	CompressionAlgorithm_NONE   CompressionAlgorithm = 0
	CompressionAlgorithm_BROTLI CompressionAlgorithm = 1
	CompressionAlgorithm_GZIP   CompressionAlgorithm = 2
)

var CompressionAlgorithm_name = map[int32]string{
	0: "NONE",
	1: "BROTLI",
	2: "GZIP",
}
var CompressionAlgorithm_value = map[string]int32{
	"NONE":   0,
	"BROTLI": 1,
	"GZIP":   2,
}

func (x CompressionAlgorithm) String() string {
	return proto.EnumName(CompressionAlgorithm_name, int32(x))
}
func (CompressionAlgorithm) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type SyncOp_Type int32

const (
	SyncOp_BLOCK_RANGE SyncOp_Type = 0
	SyncOp_DATA        SyncOp_Type = 1
	// REMOTE_DATA used to be 2 - shouldn't be in the wild, but better not re-use it?
	SyncOp_HEY_YOU_DID_IT SyncOp_Type = 2049
)

var SyncOp_Type_name = map[int32]string{
	0:    "BLOCK_RANGE",
	1:    "DATA",
	2049: "HEY_YOU_DID_IT",
}
var SyncOp_Type_value = map[string]int32{
	"BLOCK_RANGE":    0,
	"DATA":           1,
	"HEY_YOU_DID_IT": 2049,
}

func (x SyncOp_Type) String() string {
	return proto.EnumName(SyncOp_Type_name, int32(x))
}
func (SyncOp_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{2, 0} }

type PatchHeader struct {
	Compression *CompressionSettings `protobuf:"bytes,1,opt,name=compression" json:"compression,omitempty"`
}

func (m *PatchHeader) Reset()                    { *m = PatchHeader{} }
func (m *PatchHeader) String() string            { return proto.CompactTextString(m) }
func (*PatchHeader) ProtoMessage()               {}
func (*PatchHeader) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *PatchHeader) GetCompression() *CompressionSettings {
	if m != nil {
		return m.Compression
	}
	return nil
}

type SyncHeader struct {
	FileIndex int64 `protobuf:"varint,16,opt,name=fileIndex" json:"fileIndex,omitempty"`
}

func (m *SyncHeader) Reset()                    { *m = SyncHeader{} }
func (m *SyncHeader) String() string            { return proto.CompactTextString(m) }
func (*SyncHeader) ProtoMessage()               {}
func (*SyncHeader) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type SyncOp struct {
	Type       SyncOp_Type `protobuf:"varint,1,opt,name=type,enum=io.itch.wharf.pwr.SyncOp_Type" json:"type,omitempty"`
	FileIndex  int64       `protobuf:"varint,2,opt,name=fileIndex" json:"fileIndex,omitempty"`
	BlockIndex int64       `protobuf:"varint,3,opt,name=blockIndex" json:"blockIndex,omitempty"`
	BlockSpan  int64       `protobuf:"varint,4,opt,name=blockSpan" json:"blockSpan,omitempty"`
	Data       []byte      `protobuf:"bytes,5,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *SyncOp) Reset()                    { *m = SyncOp{} }
func (m *SyncOp) String() string            { return proto.CompactTextString(m) }
func (*SyncOp) ProtoMessage()               {}
func (*SyncOp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type SignatureHeader struct {
	Compression *CompressionSettings `protobuf:"bytes,1,opt,name=compression" json:"compression,omitempty"`
}

func (m *SignatureHeader) Reset()                    { *m = SignatureHeader{} }
func (m *SignatureHeader) String() string            { return proto.CompactTextString(m) }
func (*SignatureHeader) ProtoMessage()               {}
func (*SignatureHeader) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *SignatureHeader) GetCompression() *CompressionSettings {
	if m != nil {
		return m.Compression
	}
	return nil
}

type BlockHash struct {
	WeakHash   uint32 `protobuf:"varint,1,opt,name=weakHash" json:"weakHash,omitempty"`
	StrongHash []byte `protobuf:"bytes,2,opt,name=strongHash,proto3" json:"strongHash,omitempty"`
}

func (m *BlockHash) Reset()                    { *m = BlockHash{} }
func (m *BlockHash) String() string            { return proto.CompactTextString(m) }
func (*BlockHash) ProtoMessage()               {}
func (*BlockHash) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type CompressionSettings struct {
	Algorithm CompressionAlgorithm `protobuf:"varint,1,opt,name=algorithm,enum=io.itch.wharf.pwr.CompressionAlgorithm" json:"algorithm,omitempty"`
	Quality   int32                `protobuf:"varint,2,opt,name=quality" json:"quality,omitempty"`
}

func (m *CompressionSettings) Reset()                    { *m = CompressionSettings{} }
func (m *CompressionSettings) String() string            { return proto.CompactTextString(m) }
func (*CompressionSettings) ProtoMessage()               {}
func (*CompressionSettings) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func init() {
	proto.RegisterType((*PatchHeader)(nil), "io.itch.wharf.pwr.PatchHeader")
	proto.RegisterType((*SyncHeader)(nil), "io.itch.wharf.pwr.SyncHeader")
	proto.RegisterType((*SyncOp)(nil), "io.itch.wharf.pwr.SyncOp")
	proto.RegisterType((*SignatureHeader)(nil), "io.itch.wharf.pwr.SignatureHeader")
	proto.RegisterType((*BlockHash)(nil), "io.itch.wharf.pwr.BlockHash")
	proto.RegisterType((*CompressionSettings)(nil), "io.itch.wharf.pwr.CompressionSettings")
	proto.RegisterEnum("io.itch.wharf.pwr.CompressionAlgorithm", CompressionAlgorithm_name, CompressionAlgorithm_value)
	proto.RegisterEnum("io.itch.wharf.pwr.SyncOp_Type", SyncOp_Type_name, SyncOp_Type_value)
}

var fileDescriptor0 = []byte{
	// 414 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xac, 0x52, 0xcf, 0x6f, 0xd3, 0x30,
	0x14, 0xa6, 0x6d, 0x5a, 0xd6, 0x17, 0xb6, 0x05, 0x8f, 0x43, 0x84, 0xd0, 0x84, 0x72, 0x00, 0xb4,
	0x43, 0x90, 0x8a, 0xb4, 0x7b, 0xb2, 0x56, 0x4d, 0xc4, 0xd4, 0x4c, 0x4e, 0x10, 0xda, 0x38, 0x44,
	0x5e, 0xea, 0x25, 0x16, 0x59, 0x1c, 0x1c, 0x8f, 0xd2, 0x23, 0xff, 0x2b, 0x7f, 0x08, 0x8e, 0xbb,
	0xb5, 0x05, 0x2a, 0x4e, 0x3b, 0x44, 0xf2, 0xfb, 0xde, 0xf7, 0xc3, 0xf9, 0x12, 0xd8, 0xaf, 0x17,
	0xe2, 0xbd, 0x7a, 0xdc, 0x5a, 0x70, 0xc9, 0xd1, 0x73, 0xc6, 0x5d, 0x26, 0xb3, 0xc2, 0x5d, 0x14,
	0x44, 0xdc, 0xb8, 0x6a, 0xe1, 0x7c, 0x06, 0xf3, 0x82, 0x28, 0x28, 0xa0, 0x64, 0x4e, 0x05, 0x0a,
	0xc0, 0xcc, 0xf8, 0x6d, 0x2d, 0x68, 0xd3, 0x30, 0x5e, 0xd9, 0x9d, 0xd7, 0x9d, 0x77, 0xe6, 0xe8,
	0x8d, 0xfb, 0x8f, 0xce, 0x3d, 0xdb, 0xb0, 0x62, 0x2a, 0x25, 0xab, 0xf2, 0x06, 0x6f, 0x4b, 0x9d,
	0x13, 0x80, 0x78, 0x59, 0x65, 0xf7, 0xbe, 0xaf, 0x60, 0x78, 0xc3, 0x4a, 0x1a, 0x56, 0x73, 0xfa,
	0xc3, 0xb6, 0x94, 0x6b, 0x0f, 0x6f, 0x00, 0xe7, 0x57, 0x07, 0x06, 0x2d, 0x39, 0xaa, 0xd1, 0x08,
	0x0c, 0xb9, 0xac, 0xa9, 0x4e, 0x3e, 0x18, 0x1d, 0xef, 0x48, 0x5e, 0x11, 0xdd, 0x44, 0xb1, 0xb0,
	0xe6, 0xfe, 0x69, 0xde, 0xfd, 0xcb, 0x1c, 0x1d, 0x03, 0x5c, 0x97, 0x3c, 0xfb, 0xba, 0x5a, 0xf7,
	0xf4, 0x7a, 0x0b, 0x69, 0xd5, 0x7a, 0x8a, 0x6b, 0x52, 0xd9, 0xc6, 0x4a, 0xbd, 0x06, 0x10, 0x02,
	0x63, 0x4e, 0x24, 0xb1, 0xfb, 0x6a, 0xf1, 0x0c, 0xeb, 0xb3, 0x73, 0x0a, 0x46, 0x9b, 0x8e, 0x0e,
	0xc1, 0xf4, 0xcf, 0xa3, 0xb3, 0x8f, 0x29, 0xf6, 0x66, 0xd3, 0x89, 0xf5, 0x04, 0xed, 0x81, 0x31,
	0xf6, 0x12, 0xcf, 0xea, 0xa0, 0x23, 0x38, 0x08, 0x26, 0x97, 0xe9, 0x65, 0xf4, 0x29, 0x1d, 0x87,
	0xe3, 0x34, 0x4c, 0xac, 0x9f, 0x96, 0xf3, 0x05, 0x0e, 0x63, 0x96, 0x57, 0x44, 0xde, 0x09, 0xfa,
	0xe8, 0x7d, 0x4f, 0x61, 0xe8, 0xb7, 0xb7, 0x0e, 0x48, 0x53, 0xa0, 0x97, 0xb0, 0xb7, 0xa0, 0x44,
	0x9f, 0xb5, 0xe7, 0x3e, 0x5e, 0xcf, 0x6d, 0x1f, 0x8d, 0x14, 0xbc, 0xca, 0xf5, 0xb6, 0xab, 0xdf,
	0x6b, 0x0b, 0x71, 0xbe, 0xc3, 0xd1, 0x8e, 0x30, 0x34, 0x81, 0x21, 0x29, 0x73, 0x2e, 0x98, 0x2c,
	0x6e, 0xef, 0xbf, 0xce, 0xdb, 0xff, 0xdf, 0xd3, 0x7b, 0xa0, 0xe3, 0x8d, 0x12, 0xd9, 0xf0, 0xf4,
	0xdb, 0x1d, 0x29, 0x99, 0x5c, 0xea, 0xe8, 0x3e, 0x7e, 0x18, 0x4f, 0x4e, 0xe1, 0xc5, 0x2e, 0x71,
	0x5b, 0xea, 0x2c, 0x9a, 0xb5, 0xf5, 0x02, 0x0c, 0x7c, 0x1c, 0x25, 0xe7, 0xa1, 0x2a, 0x58, 0xa1,
	0xd3, 0xab, 0xf0, 0xc2, 0xea, 0xfa, 0xfd, 0xab, 0x9e, 0x0a, 0xbe, 0x1e, 0xe8, 0x5f, 0xfc, 0xc3,
	0xef, 0x00, 0x00, 0x00, 0xff, 0xff, 0xc3, 0x14, 0xf4, 0x9d, 0xf3, 0x02, 0x00, 0x00,
}

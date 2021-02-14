package nvme

import (
	"github.com/stretchr/testify/assert"
	"math"
	"reflect"
	"testing"
	"unsafe"
)

const (
	expectedOffset = uint64(0xAB<<32 | 0xCD)
	expectedLID    = uint16(0xCA)
	expectedLSP    = uint16(0xE)
	expectedLSI    = uint16(0xA)
)

func TestGetLogCmd_SetDWords(t *testing.T) {
	a := assert.New(t)
	v := make([]byte, maxAdminCmdPageSz)

	const expectedDWords = maxAdminCmdPageSz >> 2

	// 1. create same get-log command
	origin, _ := newGetLogCmd(expectedNSId, expectedOffset, expectedLID, expectedLSP, expectedLSI, v)
	tested, _ := newGetLogCmd(expectedNSId, expectedOffset, expectedLID, expectedLSP, expectedLSI, v)

	tested.SetDWords(^expectedDWords)

	// CDW10 check
	a.NotEqual(origin.CDW10, tested.CDW10)
	a.Equal(origin.CDW10&maskUint16, tested.CDW10&maskUint16) // keep lsp, lid field (lower 16bit)
	a.Equal(umaskUint16, origin.CDW10^tested.CDW10)           // XOR between origin and tested should be 0xFFFF0000

	// CDW11 check
	a.NotEqual(origin.CDW11, tested.CDW11)
	a.Equal(origin.CDW11&umaskUint16, tested.CDW11&umaskUint16) // keep lsi field (upper 16bit)
	a.Equal(maskUint16, origin.CDW11^tested.CDW11)              // XOR between origin and tested should be 0x0000FFFF

	// changed value check
	a.Equal(^expectedDWords, (tested.CDW10>>16)|(tested.CDW11<<16))
	a.NotZero((origin.CDW10 >> 16) | (origin.CDW11 << 16))
	a.NotZero((tested.CDW10 >> 16) | (tested.CDW11 << 16))

	// CDW12/13 check
	a.Equal(origin.CDW12, tested.CDW12)
	a.Equal(origin.CDW13, tested.CDW13)
}

func TestGetLogCmd_SetOffset(t *testing.T) {
	a := assert.New(t)
	v := make([]byte, 4)

	// 1. create same get-log command
	origin, _ := newGetLogCmd(expectedNSId, expectedOffset, expectedLID, expectedLSP, expectedLSI, v)
	tested, _ := newGetLogCmd(expectedNSId, expectedOffset, expectedLID, expectedLSP, expectedLSI, v)

	tested.SetOffset(^expectedOffset)

	// CDW10/11 check
	a.Equal(origin.CDW10, tested.CDW10)
	a.Equal(origin.CDW11, tested.CDW11)

	// CDW12 check
	a.NotEqual(origin.CDW12, tested.CDW12)
	a.Equal(^uint32(0x0), origin.CDW12|tested.CDW12)

	// CDW13 check
	a.NotEqual(origin.CDW13, tested.CDW13)
	a.Equal(^uint32(0x0), origin.CDW13|tested.CDW13)

	// changed value check
	a.Equal(^expectedOffset, uint64(tested.CDW12)<<32|uint64(tested.CDW13))
}

func TestGetLogCmd_SetLSP(t *testing.T) {
	a := assert.New(t)
	v := make([]byte, 4)

	const UMaskLSP = ^(maskUint4 << shiftUint8)

	// 1. create same get-log command
	origin, _ := newGetLogCmd(expectedNSId, expectedOffset, expectedLID, expectedLSP, expectedLSI, v)
	tested, _ := newGetLogCmd(expectedNSId, expectedOffset, expectedLID, expectedLSP, expectedLSI, v)

	// The passwd value is 0xfff1 but 0xfff0 is dirty value. So SetLSP should mask out that dirty
	// value.
	tested.SetLSP(^expectedLSP)

	// CDW11/12/13 check
	a.Equal(origin.CDW11, tested.CDW11)
	a.Equal(origin.CDW12, tested.CDW12)
	a.Equal(origin.CDW13, tested.CDW13)

	// CDW10 Check
	a.NotEqual(origin.CDW10, tested.CDW10)
	a.Equal(maskUint4, (origin.CDW10^tested.CDW10)>>8)
	a.Equal(origin.CDW10&UMaskLSP, tested.CDW10&UMaskLSP)

	// Check set value and masked out dirty bits. "uint32(^expectedLSP) & maskUint4" will be clean
	// 4bit value.
	a.Equal(uint32(^expectedLSP)&maskUint4, (tested.CDW10>>8)&maskUint4)
}

func TestNewGetLogCmd(t *testing.T) {
	a := assert.New(t)

	const (
		expectedSz     = uint32(32)
		expectedDWords = expectedSz >> 2
	)

	tcList := []interface{}{
		make([]byte, expectedSz),
		&struct{ buffer [expectedSz]byte }{},
	}

	for _, tc := range tcList {
		tested, err := newGetLogCmd(expectedNSId, expectedOffset, expectedLID, expectedLSP, expectedLSI, tc)
		a.NotNil(tested)
		a.NoError(err)
		a.Equal(AdminGetLogPage, tested.OpCode)
		a.Equal(uint32(expectedNSId), tested.NSId)
		a.Equal(expectedDWords, (tested.CDW10>>16)|(tested.CDW11<<16))
		a.Equal(expectedOffset, (uint64(tested.CDW12)<<32)|uint64(tested.CDW13))
		a.Equal(expectedLID, uint16(math.MaxUint8&tested.CDW10))
		a.Equal(expectedLSP, uint16((tested.CDW10<<16)>>24))
		a.Equal(expectedLSI, uint16(tested.CDW11>>16))
		a.Equal(expectedSz, tested.DataLength)
		a.Equal(reflect.ValueOf(tc).Pointer(), tested.Data)
	}

	// test invalid case
	tc := struct {
		buffer [expectedSz]byte
	}{}
	tested, err := newGetLogCmd(expectedNSId, expectedOffset, expectedLID, expectedLSP, expectedLSI, tc)
	a.Nil(tested)
	a.Error(err)
}

func TestSMARTSize(t *testing.T) {
	a := assert.New(t)

	a.Equal(uintptr(512), unsafe.Sizeof(smart{}))
}

func TestTelemetryDataBlk_index(t *testing.T) {
	a := assert.New(t)

	a.Equal(0, DataBlock1.index())
	a.Equal(1, DataBlock2.index())
	a.Equal(2, DataBlock3.index())
}

func TestTelemetrySize(t *testing.T) {
	a := assert.New(t)

	a.Equal(uintptr(512), unsafe.Sizeof(telemetry{}))
}

func TestTelemetry_BlockSize(t *testing.T) {
	a := assert.New(t)

	expectedBlk1 := uint32(2 * 512)
	expectedBlk2 := uint32(4 * 512)
	expectedBlk3 := uint32(8 * 512)

	tested := telemetry{
		DataAreaLastBlock: [3]uint16{2, 4, 8},
	}

	a.Equal(expectedBlk1, tested.BlockSize(DataBlock1))
	a.Equal(expectedBlk2, tested.BlockSize(DataBlock2))
	a.Equal(expectedBlk3, tested.BlockSize(DataBlock3))
}

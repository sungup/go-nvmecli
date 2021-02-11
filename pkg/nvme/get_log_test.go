package nvme

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
	"unsafe"
)

func TestSMARTSize(t *testing.T) {
	a := assert.New(t)

	a.Equal(uintptr(512), unsafe.Sizeof(smart{}))
}

func TestNewGetLogCmd(t *testing.T) {
	const (
		expectedNSId   = 1
		expectedDWords = uint32(0xAB<<16 | 0xCD)
		expectedOffset = uint64(0xAB<<32 | 0xCD)
		expectedLID    = uint16(0xCA)
		expectedLSP    = uint16(0xE)
		expectedLSI    = uint16(0xA)
	)

	a := assert.New(t)

	tested := newGetLogCmd(expectedNSId, expectedDWords, expectedOffset, expectedLID, expectedLSP, expectedLSI)
	a.NotNil(tested)
	a.Equal(AdminGetLogPage, tested.OpCode)
	a.Equal(uint32(expectedNSId), tested.NSId)
	a.Equal(expectedDWords, (tested.CDW10>>16)|(tested.CDW11<<16))
	a.Equal(expectedOffset, (uint64(tested.CDW12)<<32)|uint64(tested.CDW13))
	a.Equal(expectedLID, uint16(math.MaxUint8&tested.CDW10))
	a.Equal(expectedLSP, uint16((tested.CDW10<<16)>>24))
	a.Equal(expectedLSI, uint16(tested.CDW11>>16))
}

func TestTelemetrySize(t *testing.T) {
	a := assert.New(t)

	a.Equal(uintptr(512), unsafe.Sizeof(telemetry{}))
}

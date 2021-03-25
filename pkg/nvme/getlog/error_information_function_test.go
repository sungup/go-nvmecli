package getlog

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"unsafe"
)

func TestErrorEntrySize(t *testing.T) {
	a := assert.New(t)

	a.Equal(uintptr(64), unsafe.Sizeof(errorEntry{}))
}

func TestErrStatField_PhaseTag(t *testing.T) {
	a := assert.New(t)

	for _, tc := range []uint16{0x0A, 0x0B} {
		expected := tc & 0x01

		tested := errStatField(tc)

		a.Equal(expected, tested.PhaseTag())
		a.NotEqual(tc, tested.PhaseTag())
	}
}

func TestErrStatField_StatusField(t *testing.T) {
	a := assert.New(t)

	for _, tc := range []uint16{0x0A, 0x0B} {
		expected := tc >> 1

		tested := errStatField(tc)

		a.Equal(expected, tested.StatusField())
		a.NotEqual(expected^(tested.StatusField()<<1), tested.PhaseTag())
	}
}

func TestParamErrLoc_LocationBit(t *testing.T) {
	a := assert.New(t)

	for _, rsv := range []uint16{0x01, 0x02, 0x03} {
		for b0700 := uint16(1); b0700 < 64; b0700++ {
			for expected := uint16(0); expected < 8; expected++ {
				tc := rsv<<11 | expected<<8 | b0700

				tested := paramErrLoc(tc)

				a.Equal(uint8(expected), tested.LocationBit())
			}
		}
	}
}

func TestParamErrLoc_LocationByte(t *testing.T) {
	a := assert.New(t)

	for _, rsv := range []uint16{0x01, 0x02, 0x03} {
		for b1008 := uint16(1); b1008 < 8; b1008++ {
			for expected := uint16(0); expected < 64; expected++ {
				tc := rsv<<11 | b1008<<8 | expected

				tested := paramErrLoc(tc)

				a.Equal(uint8(b1008), tested.LocationBit())
			}
		}
	}
}

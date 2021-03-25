package getlog

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"unsafe"
)

func TestTelemetryDataBlk_index(t *testing.T) {
	a := assert.New(t)

	a.Equal(0, DataBlock1.index())
	a.Equal(1, DataBlock2.index())
	a.Equal(2, DataBlock3.index())
}

func TestTelemetrySize(t *testing.T) {
	a := assert.New(t)

	a.Equal(uintptr(512), unsafe.Sizeof(Telemetry{}))
}

func TestTelemetry_BlockSize(t *testing.T) {
	a := assert.New(t)

	expectedBlk1 := uint32(2 * 512)
	expectedBlk2 := uint32(4 * 512)
	expectedBlk3 := uint32(8 * 512)

	tested := Telemetry{
		DataAreaLastBlock: [3]uint16{2, 4, 8},
	}

	a.Equal(expectedBlk1, tested.BlockSize(DataBlock1))
	a.Equal(expectedBlk2, tested.BlockSize(DataBlock2))
	a.Equal(expectedBlk3, tested.BlockSize(DataBlock3))
}

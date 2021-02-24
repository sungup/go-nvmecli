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
	// TODO implementing here
}

func TestErrStatField_StatusField(t *testing.T) {
	// TODO implementing here
}

func TestParamErrLoc_LocationBit(t *testing.T) {
	// TODO implementing here
}

func TestParamErrLoc_LocationByte(t *testing.T) {
	// TODO implementing here
}

package getlog

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"unsafe"
)

func TestSMARTSize(t *testing.T) {
	a := assert.New(t)

	a.Equal(uintptr(512), unsafe.Sizeof(SMART{}))
}

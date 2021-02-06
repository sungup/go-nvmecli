package nvme

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"unsafe"
)

func TestCtrlIdentifySize(t *testing.T) {
	a := assert.New(t)

	a.Equal(uintptr(4096), unsafe.Sizeof(CtrlIdentify{}))
	a.Equal(uintptr(32), unsafe.Sizeof(PowerStateDescriptor{}))
}

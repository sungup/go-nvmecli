package nvme

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"unsafe"
)

func TestCtrlIdentifySize(t *testing.T) {
	a := assert.New(t)

	a.Equal(uintptr(4096), unsafe.Sizeof(ctrlIdentify{}))
	a.Equal(uintptr(32), unsafe.Sizeof(powerStateDesc{}))
}

func TestNewIdentifyCmd(t *testing.T) {
	const (
		expectedNSId     = 1
		expectedCNTId    = 0xa
		expectedCNS      = cnsNamespace
		expectedNvmSetId = 0xF
	)

	a := assert.New(t)

	tested := newIdentifyCmd(expectedNSId, expectedCNTId, expectedCNS, expectedNvmSetId)
	a.NotNil(tested)
	a.Equal(AdminIdentify, tested.OpCode)
	a.Equal(uint32(expectedCNTId)<<16|uint32(expectedCNS), tested.CDW10)
	a.Equal(uint32(expectedNvmSetId), tested.CDW11)
	a.Equal(uint32(expectedNSId), tested.NSId)
}

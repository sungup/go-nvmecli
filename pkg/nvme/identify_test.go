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
	var (
		a = assert.New(t)

		err    error
		data   interface{}
		tested *AdminCmd
	)
	const (
		expectedNSId     = 1
		expectedCNTId    = 0xa
		expectedCNS      = cnsNamespace
		expectedNvmSetId = 0xF
		expectedTestSize = 3072
	)

	// 1. check invalid dptr data
	data = uint64(0)
	tested, err = newIdentifyCmd(expectedNSId, expectedCNTId, expectedCNS, expectedNvmSetId, data)
	a.Error(err)

	// 2. get data using byte array
	data = make([]uint8, expectedTestSize)
	tested, err = newIdentifyCmd(expectedNSId, expectedCNTId, expectedCNS, expectedNvmSetId, data)
	a.NotNil(tested)
	a.NoError(err)
	a.Equal(AdminIdentify, tested.OpCode)
	a.Equal(uint32(expectedCNTId)<<16|uint32(expectedCNS), tested.CDW10)
	a.Equal(uint32(expectedNvmSetId), tested.CDW11)
	a.Equal(uint32(expectedNSId), tested.NSId)
	a.Equal(uint32(expectedTestSize), tested.DataLength)

	// 3. get data using any structure
	data = &struct {
		raw [expectedTestSize + 10]byte
	}{}
	tested, err = newIdentifyCmd(expectedNSId, expectedCNTId, expectedCNS, expectedNvmSetId, data)
	a.NotNil(tested)
	a.NoError(err)
	a.Equal(AdminIdentify, tested.OpCode)
	a.Equal(uint32(expectedCNTId)<<16|uint32(expectedCNS), tested.CDW10)
	a.Equal(uint32(expectedNvmSetId), tested.CDW11)
	a.Equal(uint32(expectedNSId), tested.NSId)
	a.Equal(uint32(expectedTestSize+10), tested.DataLength)
}

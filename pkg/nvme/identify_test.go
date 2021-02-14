package nvme

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"unsafe"
)

const (
	expectedCNTId    = 0xa
	expectedCNS      = cnsNamespace
	expectedNvmSetId = 0xF
)

func TestCtrlIdentifySize(t *testing.T) {
	a := assert.New(t)

	a.Equal(uintptr(4096), unsafe.Sizeof(ctrlIdentify{}))
	a.Equal(uintptr(32), unsafe.Sizeof(powerStateDesc{}))
}

func TestNewIdentifyCmd(t *testing.T) {
	a := assert.New(t)

	const expectedSz = uint32(32)

	tcList := []interface{}{
		make([]byte, expectedSz),
		&struct{ buffer [expectedSz]byte }{},
	}

	// test valid case
	for _, tc := range tcList {
		tested, err := newIdentifyCmd(expectedNSId, expectedCNTId, expectedCNS, expectedNvmSetId, tc)
		a.NotNil(tested)
		a.NoError(err)
		a.Equal(AdminIdentify, tested.OpCode)
		a.Equal(uint32(expectedCNTId)<<16|uint32(expectedCNS), tested.CDW10)
		a.Equal(uint32(expectedNvmSetId), tested.CDW11)
		a.Equal(uint32(expectedNSId), tested.NSId)
		a.Equal(expectedSz, tested.DataLength)
		a.Equal(reflect.ValueOf(tc).Pointer(), tested.Data)
	}

	// test invalid case
	tc := struct {
		buffer [expectedSz]byte
	}{}

	tested, err := newIdentifyCmd(expectedNSId, expectedCNTId, expectedCNS, expectedNvmSetId, tc)
	a.Nil(tested)
	a.Error(err)
}

// +build with_phys_device

package identify

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"unsafe"
)

func TestGetCtrlIdentify(t *testing.T) {
	a := assert.New(t)

	dev, _ := os.Open(targetDevice)

	// 1. byte array buffer
	buffer := make([]byte, unsafe.Sizeof(CtrlIdentify{}))
	a.NoError(GetCtrlIdentify(dev, buffer))

	// 2. struct data
	identify := CtrlIdentify{}
	a.NoError(GetCtrlIdentify(dev, &identify))

	a.NotEmpty(identify.SN.String())
	a.NotEmpty(identify.MN.String())
	a.NotEmpty(string(identify.FR[:]))
}

func TestParseCtrlIdentify(t *testing.T) {
	a := assert.New(t)

	dev, _ := os.Open(targetDevice)
	buffer := make([]byte, unsafe.Sizeof(CtrlIdentify{}))

	_ = GetCtrlIdentify(dev, buffer)

	// check invalid size error
	identify, err := ParseCtrlIdentify(buffer[1:])
	a.Error(err)
	a.Nil(identify)

	identify, err = ParseCtrlIdentify(buffer)
	a.NoError(err)
	a.NotNil(identify)

	a.NotEmpty(identify.SN.String())
	a.NotEmpty(identify.MN.String())
	a.NotEmpty(string(identify.FR[:]))
}

func TestGetNamespaceIdentify(t *testing.T) {
	a := assert.New(t)

	dev, _ := os.Open(targetDevice)

	// 1. byte array buffer
	buffer := make([]byte, unsafe.Sizeof(NamespaceIdentify{}))
	a.NoError(GetNamespaceIdentify(dev, 1, buffer))

	// 2. struct data
	identify := NamespaceIdentify{}
	a.NoError(GetNamespaceIdentify(dev, 1, &identify))

	a.NotEmpty(identify.EUI64.String())
	a.NotEmpty(identify.NGUID.String())
}

func TestParseNamespaceIdentify(t *testing.T) {
	a := assert.New(t)

	dev, _ := os.Open(targetDevice)
	buffer := make([]byte, unsafe.Sizeof(NamespaceIdentify{}))

	_ = GetNamespaceIdentify(dev, 1, buffer)

	// check invalid size error
	identify, err := ParseNamespaceIdentify(buffer[1:])
	a.Error(err)
	a.Nil(identify)

	identify, err = ParseNamespaceIdentify(buffer)
	a.NoError(err)
	a.NotNil(identify)

	a.NotEmpty(identify.EUI64.String())
	a.NotEmpty(identify.NGUID.String())
}

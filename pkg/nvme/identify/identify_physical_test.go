// +build with_phys_device

package identify

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetCtrlIdentify(t *testing.T) {
	a := assert.New(t)

	dev, _ := os.Open(targetDevice)

	// 1. byte array buffer
	buffer := make([]byte, ctrlIdentifySz)
	a.NoError(GetCtrlIdentify(dev, buffer))

	// 2. struct data
	identify := CtrlIdentify{}
	a.NoError(GetCtrlIdentify(dev, &identify))

	a.NotEmpty(string(identify.SN[:]))
	a.NotEmpty(string(identify.MN[:]))
	a.NotEmpty(string(identify.FR[:]))
}

func TestParseCtrlIdentify(t *testing.T) {
	a := assert.New(t)

	dev, _ := os.Open(targetDevice)
	buffer := make([]byte, ctrlIdentifySz)

	_ = GetCtrlIdentify(dev, buffer)

	// check invalid size error
	identify, err := ParseCtrlIdentify(buffer[1:])
	a.Error(err)
	a.Nil(identify)

	identify, err = ParseCtrlIdentify(buffer)
	a.NoError(err)
	a.NotNil(identify)

	a.NotEmpty(string(identify.SN[:]))
	a.NotEmpty(string(identify.MN[:]))
	a.NotEmpty(string(identify.FR[:]))
}

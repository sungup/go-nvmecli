// +build with_phys_device

package nvme

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCtrlIdentify(t *testing.T) {
	a := assert.New(t)

	dev, _ := os.Open(targetDevice)

	// 1. byte array buffer
	buffer := make([]byte, ctrlIdentifySz)
	a.NoError(CtrlIdentify(dev, buffer))

	// 2. struct data
	identify := ctrlIdentify{}
	a.NoError(CtrlIdentify(dev, &identify))

	a.NotEmpty(string(identify.SN[:]))
	a.NotEmpty(string(identify.MN[:]))
	a.NotEmpty(string(identify.FR[:]))
}

func TestParseCtrlIdentify(t *testing.T) {
	a := assert.New(t)

	dev, _ := os.Open(targetDevice)
	buffer := make([]byte, ctrlIdentifySz)

	_ = CtrlIdentify(dev, buffer)

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

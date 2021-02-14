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

	buffer, err := CtrlIdentify(dev)
	a.NoError(err)
	a.NotNil(buffer)
	a.Len(buffer, ctrlIdentifySz)
}

func TestParseCtrlIdentify(t *testing.T) {
	a := assert.New(t)

	dev, _ := os.Open(targetDevice)

	buffer, _ := CtrlIdentify(dev)

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

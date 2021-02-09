// +build with_phys_device

package nvme

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCtrlIdentify(t *testing.T) {
	a := assert.New(t)

	dev, err := os.Open(targetDevice)
	a.NoError(err)

	buffer, err := CtrlIdentify(dev)
	a.NoError(err)
	a.NotNil(buffer)
	a.Len(buffer, ctrlIdentifySz)
}

func TestParseCtrlIdentify(t *testing.T) {
	a := assert.New(t)

	dev, err := os.Open(targetDevice)
	a.NoError(err)

	buffer, _ := CtrlIdentify(dev)
	identify, err := ParseCtrlIdentify(buffer)
	a.NoError(err)

	a.NotEmpty(string(identify.SN[:]))
	a.NotEmpty(string(identify.MN[:]))
	a.NotEmpty(string(identify.FR[:]))
}

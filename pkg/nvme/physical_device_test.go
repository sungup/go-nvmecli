// +build with_phys_device

package nvme

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const (
	targetDevice = "/dev/nvme0"
)

func TestGetCtrlIdentify(t *testing.T) {
	a := assert.New(t)

	dev, err := os.Open(targetDevice)
	a.NoError(err)

	identify, err := GetCtrlIdentify(dev)
	a.NoError(err)

	t.Log(string(identify.SN[:]))
	t.Log(string(identify.MN[:]))
	t.Log(string(identify.FR[:]))
}

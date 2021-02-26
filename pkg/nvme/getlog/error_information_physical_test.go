// +build with_phys_device

package getlog

import (
	"github.com/stretchr/testify/assert"
	"github.com/sungup/go-nvmecli/pkg/nvme/identify"
	"os"
	"testing"
)

func TestGetELPE(t *testing.T) {
	a := assert.New(t)

	dev, _ := os.Open(targetDevice)

	idCtrl := identify.CtrlIdentify{}
	a.NoError(identify.GetCtrlIdentify(dev, &idCtrl))

	tested, err := getELPE(dev)
	a.NoError(err)
	a.Equal(idCtrl.ELPE.Uint(), uint64(tested))
}

func TestGetErrorInformation(t *testing.T) {
	// TODO re-verify this test code using the log stored NVMe device
	// My Seagate NVMe device didn't collect an error information. So, the following test code
	// should return empty error logs without error.
	a := assert.New(t)

	dev, _ := os.Open(targetDevice)

	tested, err := GetErrorInformation(dev, 65536)
	a.NoError(err)
	a.Empty(tested)
}

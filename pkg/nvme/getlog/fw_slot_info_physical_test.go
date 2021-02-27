// +build with_phys_device

package getlog

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"unsafe"
)

func TestGetFirmwareSlotInfo(t *testing.T) {
	a := assert.New(t)

	dev, _ := os.Open(targetDevice)

	// 1. byte array buffer
	buffer := make([]byte, unsafe.Sizeof(FirmwareSlotInfo{}))
	a.NoError(GetFirmwareSlotInfo(dev, buffer))

	// 2. struct data
	tested := FirmwareSlotInfo{}
	a.NoError(GetFirmwareSlotInfo(dev, &tested))

	a.NotEqual(NoFw, tested.ActiveFwInfo.Active())
}

func TestParseFirmwareSlotInfo(t *testing.T) {
	a := assert.New(t)

	dev, _ := os.Open(targetDevice)
	buffer := make([]byte, unsafe.Sizeof(FirmwareSlotInfo{}))

	_ = GetFirmwareSlotInfo(dev, buffer)

	// check invalid size error
	tested, err := ParseFirmwareSlotInfo(buffer[1:])
	a.Error(err)
	a.Nil(tested)

	// check normal parsing
	tested, err = ParseFirmwareSlotInfo(buffer)
	a.NoError(err)
	a.NotNil(tested)

	testedSlot, testedFw := tested.Active()
	a.NotEqual(NoFw, testedSlot)
	a.NotEmpty(testedFw)
}

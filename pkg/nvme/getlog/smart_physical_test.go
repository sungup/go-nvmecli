// +build with_phys_device

package getlog

import (
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"strconv"
	"testing"
	"unsafe"
)

func TestGetSMART(t *testing.T) {
	a := assert.New(t)

	dev, _ := os.Open(targetDevice)

	// 1. byte array buffer
	buffer := make([]byte, unsafe.Sizeof(SMART{}))

	// TODO Remove following recreated buffer. Currently my environment, if the buffer address has
	//      0xE00 (pointer 3584 and it means the last 512B in 4KB page), IOCTL returns the 0 filled
	//      data into the buffer. To avoid this problem, test code recreate 512B to escape from
	//      0xE00.
	buffer = make([]byte, unsafe.Sizeof(SMART{}))

	a.NoError(GetSMART(dev, buffer))

	// 2. struct data
	tested := SMART{}
	a.NoError(GetSMART(dev, &tested))

	a.NotZero(tested.AvailableSpareThreshold)
}

func TestParseSMART(t *testing.T) {
	a := assert.New(t)

	dev, _ := os.Open(targetDevice)
	buffer := make([]byte, unsafe.Sizeof(SMART{}))

	_ = GetSMART(dev, buffer)

	// check invalid size error
	smart, err := ParseSMART(buffer[1:])
	a.Error(err)
	a.Nil(smart)

	// check normal parsing
	smart, err = ParseSMART(buffer)
	a.NoError(err)
	a.NotNil(smart)

	a.NotZero(smart.AvailableSpareThreshold)
}

func TestGetSMART_ForLast512BDataMissing(t *testing.T) {
	dev, _ := os.Open(targetDevice)

	__getTcPtr := func(data interface{}) uintptr {
		return reflect.ValueOf(data).Pointer()
	}

	for i := 0; i < 128; i++ {
		buffer := make([]byte, unsafe.Sizeof(SMART{})) // 1. create buffer with 512B
		address := __getTcPtr(buffer)                  // 2. get it's pointer address
		_ = GetSMART(dev, buffer)                      // 3. retrieve SMART data into buffer using ioctl
		smart, _ := ParseSMART(buffer)                 // 4. parsing to display AvailableSpareThreshold

		t.Logf(
			"0x%X => 0b%s, AvailableSpareThreshold: %d",
			address,
			strconv.FormatInt(int64(address), 2),
			smart.AvailableSpareThreshold,
		)
	}
}

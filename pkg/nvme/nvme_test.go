package nvme

import (
	"github.com/sungup/go-nvme-ioctl/pkg/ioctl"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

const (
	sizeOfUserIO        = 48
	sizeOfPassthruCmd32 = 72
	sizeOfPassthruCmd64 = 80
	sizeOfAdminCmd      = 72
)

func TestIOCtlCMD(t *testing.T) {
	a := assert.New(t)

	// reference code from nvme-cli project
	expectedNVMeIOCtlId := ioctl.IO('N', 0x40)
	expectedNVMeIOCtlAdminCmd := ioctl.IOWR('N', 0x41, uint64(unsafe.Sizeof(AdminCmd{})))
	expectedNVMeIOCtlSubmitIO := ioctl.IOW('N', 0x42, uint64(unsafe.Sizeof(UserIo{})))
	expectedNVMeIOCtlIOCmd := ioctl.IOWR('N', 0x43, uint64(unsafe.Sizeof(PassthruCmd32{})))
	expectedNVMeIOCtlReset := ioctl.IO('N', 0x44)
	expectedNVMeIOCtlSubSysReset := ioctl.IO('N', 0x45)
	expectedNVMeIOCtlRescan := ioctl.IO('N', 0x46)
	expectedNVMeIOCtlAdminCmd64 := ioctl.IOWR('N', 0x47, uint64(unsafe.Sizeof(PassthruCmd64{})))
	expectedNVMeIOCtlIOCmd64 := ioctl.IOWR('N', 0x48, uint64(unsafe.Sizeof(PassthruCmd64{})))

	a.Equal(expectedNVMeIOCtlId, iocId)
	a.Equal(expectedNVMeIOCtlAdminCmd, iocAdminCmd)
	a.Equal(expectedNVMeIOCtlSubmitIO, iocSubmitIO)
	a.Equal(expectedNVMeIOCtlIOCmd, iocIOCmd)
	a.Equal(expectedNVMeIOCtlReset, iocReset)
	a.Equal(expectedNVMeIOCtlSubSysReset, iocSubSysReset)
	a.Equal(expectedNVMeIOCtlRescan, iocRescan)
	a.Equal(expectedNVMeIOCtlAdminCmd64, iocAdminCmd64)
	a.Equal(expectedNVMeIOCtlIOCmd64, iocIOCmd64)

	// ioctl code from linux c api. please check at test/linux directory.
	expectedNVMeIOCtlId = uint64(0x4e40)
	expectedNVMeIOCtlAdminCmd = uint64(0xc0484e41)
	expectedNVMeIOCtlSubmitIO = uint64(0x40304e42)
	expectedNVMeIOCtlIOCmd = uint64(0xc0484e43)
	expectedNVMeIOCtlReset = uint64(0x4e44)
	expectedNVMeIOCtlSubSysReset = uint64(0x4e45)
	expectedNVMeIOCtlRescan = uint64(0x4e46)
	expectedNVMeIOCtlAdminCmd64 = uint64(0xc0504e47)
	expectedNVMeIOCtlIOCmd64 = uint64(0xc0504e48)

	a.Equal(expectedNVMeIOCtlId, iocId)
	a.Equal(expectedNVMeIOCtlAdminCmd, iocAdminCmd)
	a.Equal(expectedNVMeIOCtlSubmitIO, iocSubmitIO)
	a.Equal(expectedNVMeIOCtlIOCmd, iocIOCmd)
	a.Equal(expectedNVMeIOCtlReset, iocReset)
	a.Equal(expectedNVMeIOCtlSubSysReset, iocSubSysReset)
	a.Equal(expectedNVMeIOCtlRescan, iocRescan)
	a.Equal(expectedNVMeIOCtlAdminCmd64, iocAdminCmd64)
	a.Equal(expectedNVMeIOCtlIOCmd64, iocIOCmd64)
}

func TestDataSizeCheck(t *testing.T) {
	a := assert.New(t)

	a.Equal(sizeOfUserIO, int(unsafe.Sizeof(UserIo{})))
	a.Equal(sizeOfPassthruCmd32, int(unsafe.Sizeof(PassthruCmd32{})))
	a.Equal(sizeOfPassthruCmd64, int(unsafe.Sizeof(PassthruCmd64{})))
	a.Equal(sizeOfAdminCmd, int(unsafe.Sizeof(AdminCmd{})))
}

func Test_getPtr(t *testing.T) {
	a := assert.New(t)

	tcList := []struct {
		obj  interface{}
		ptr  interface{}
		size uint32
	}{
		{UserIo{}, &UserIo{}, sizeOfUserIO},
		{PassthruCmd32{}, &PassthruCmd32{}, sizeOfPassthruCmd32},
		{PassthruCmd64{}, &PassthruCmd64{}, sizeOfPassthruCmd64},
		{AdminCmd{}, &AdminCmd{}, sizeOfAdminCmd},
	}

	for _, tc := range tcList {
		_, _, err := getPtr(tc.obj)
		a.Error(err)

		ptr, sz, err := getPtr(tc.ptr)
		a.NoError(err)
		a.Equal(tc.size, sz)
		a.NotNil(ptr)
	}
}

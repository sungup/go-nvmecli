package nvme

import (
	"fmt"
	"github.com/sungup/go-nvmecli/pkg/ioctl"
	"os"
	"reflect"
	"unsafe"
)

type opcode uint8

const (
	AdminDeleteSQ      = opcode(0x00)
	AdminCreateSQ      = opcode(0x01)
	AdminGetLogPage    = opcode(0x02)
	AdminDeleteCQ      = opcode(0x04)
	AdminCreateCQ      = opcode(0x05)
	AdminIdentify      = opcode(0x06)
	AdminAbortCmd      = opcode(0x08)
	AdminSetFeatures   = opcode(0x09)
	AdminGetFeatures   = opcode(0x0a)
	AdminAsyncEvent    = opcode(0x0c)
	AdminNsMgmt        = opcode(0x0d)
	AdminActivateFW    = opcode(0x10)
	AdminDownloadFW    = opcode(0x11)
	AdminDevSelfTest   = opcode(0x14)
	AdminNsAttach      = opcode(0x15)
	AdminKeepAlive     = opcode(0x18)
	AdminDirectiveSend = opcode(0x19)
	AdminDirectiveRecv = opcode(0x1a)
	AdminVirtualMgmt   = opcode(0x1c)
	AdminNVMeMiSend    = opcode(0x1d)
	AdminNVMeMiRecv    = opcode(0x1e)
	AdminDBBuf         = opcode(0x7C)
	AdminFormatNVM     = opcode(0x80)
	AdminSecuritySend  = opcode(0x81)
	AdminSecurityRecv  = opcode(0x82)
	AdminSanitizeNVM   = opcode(0x84)
	AdminGetLBAStatus  = opcode(0x86)
)

// nvmeCmd interface has two function to set the metadata pointer and data block pointer. To reduce
// code duplication to set the pointer through object -> interface -> reflection (for type checking)
// and pointer assign, each command structure should serve these interface function.
type nvmeCmd interface {
	SetData(data interface{}) error
	SetMeta(data interface{}) error
}

// getPtr convert interface data to it's memory pointer and data length. If input data interface
// is not a pointer of data structure or slice, getPtr will return the error with nil address (0),
// and zero length.
func getPtr(data interface{}) (uintptr, uint32, error) {
	v := reflect.ValueOf(data)
	t := v.Type()

	if t.Kind() == reflect.Ptr && reflect.Indirect(v).Type().Kind() == reflect.Struct {
		// 1. Input data is a pointer of structure
		return v.Pointer(), uint32(reflect.Indirect(v).Type().Size()), nil
	} else if t.Kind() == reflect.Slice {
		// 2. Input data is a slice
		return v.Pointer(), uint32(v.Len()) * uint32(t.Elem().Size()), nil
	}

	return 0, 0, fmt.Errorf("input data is not a pointer of structure or a byte slice")
}

// UserIO is a structure to send normal io command to a NVMe device.
type UserIO struct {
	OpCode  opcode
	Flags   uint8
	Control uint16
	NBlocks uint16
	_       uint16
	Meta    uintptr
	Data    uintptr
	SLBA    uintptr
	DSMgmt  uint32
	RefTag  uint32
	AppTag  uint16
	AppMask uint16
}

// SetData set the memory pointer of data block.
func (c *UserIO) SetData(data interface{}) error {
	if ptr, _, err := getPtr(data); err == nil {
		c.Data = ptr
		return nil
	} else {
		return err
	}
}

// SetMeta set the memory pointer of metadata block.
func (c *UserIO) SetMeta(data interface{}) error {
	if ptr, _, err := getPtr(data); err == nil {
		c.Meta = ptr
		return nil
	} else {
		return err
	}
}

// PassthruCmd is a base structure for admin and passthru command about the NVMe device. Following
// Passthru32, Passthru64, and AdminCmd has same structure format without the TimeoutMSec and Result
// field. So, 3 structures share this basic structure to manipulate as an interface class.
type PassthruCmd struct {
	OpCode     opcode
	Flags      uint8
	_          uint16 // Reserved
	NSId       uint32
	CDW2       uint32
	CDW3       uint32
	Meta       uintptr
	Data       uintptr
	MetaLength uint32
	DataLength uint32
	CDW10      uint32
	CDW11      uint32
	CDW12      uint32
	CDW13      uint32
	CDW14      uint32
	CDW15      uint32
}

// SetData set the memory pointer and it's size of data block.
func (c *PassthruCmd) SetData(data interface{}) error {
	if ptr, size, err := getPtr(data); err == nil {
		c.Data, c.DataLength = ptr, size
		return nil
	} else {
		return err
	}
}

// SetMeta set the memory pointer and it's size of metadata block.
func (c *PassthruCmd) SetMeta(data interface{}) error {
	if ptr, size, err := getPtr(data); err == nil {
		c.Meta, c.MetaLength = ptr, size
		return nil
	} else {
		return err
	}
}

// PassthruCmd32 is a passthru command which is same with nvme_passthru_cmd in linux kernel.
// (please check /include/uapi/linux/nvme_ioctl.h)
type PassthruCmd32 struct {
	PassthruCmd

	TimeoutMSec uint32
	Result      uint32
}

// PassthruCmd64 is a passthru command which is same with nvme_passthru_cmd64 in linux kernel.
// (please check /include/uapi/linux/nvme_ioctl.h)
type PassthruCmd64 struct {
	PassthruCmd

	TimeoutMSec uint32
	_           uint32
	Result      uint64
}

// AdminCmd is an admin command which is same with nvme_admin_cmd in linux kernel. But different
// from linux kernel, AdminCmd is inherited (not aliased) structure from PassthruCmd because of
// the go language's aliasing rule.
type AdminCmd struct {
	PassthruCmd

	TimeoutMSec uint32
	Result      uint32
}

const (
	iocNVMeType = 'N' << ioctl.TypeShift

	iocId          = ioctl.IOCNone | iocNVMeType | (0x40 << ioctl.NrShift)
	iocAdminCmd    = ioctl.IOCInOut | iocNVMeType | (0x41 << ioctl.NrShift) | uint64(unsafe.Sizeof(AdminCmd{})<<ioctl.SizeShift)
	iocSubmitIO    = ioctl.IOCIn | iocNVMeType | (0x42 << ioctl.NrShift) | uint64(unsafe.Sizeof(UserIO{})<<ioctl.SizeShift)
	iocIOCmd       = ioctl.IOCInOut | iocNVMeType | (0x43 << ioctl.NrShift) | uint64(unsafe.Sizeof(PassthruCmd32{})<<ioctl.SizeShift)
	iocReset       = ioctl.IOCNone | iocNVMeType | (0x44 << ioctl.NrShift)
	iocSubSysReset = ioctl.IOCNone | iocNVMeType | (0x45 << ioctl.NrShift)
	iocRescan      = ioctl.IOCNone | iocNVMeType | (0x46 << ioctl.NrShift)
	iocAdminCmd64  = ioctl.IOCInOut | iocNVMeType | (0x47 << ioctl.NrShift) | uint64(unsafe.Sizeof(PassthruCmd64{})<<ioctl.SizeShift)
	iocIOCmd64     = ioctl.IOCInOut | iocNVMeType | (0x48 << ioctl.NrShift) | uint64(unsafe.Sizeof(PassthruCmd64{})<<ioctl.SizeShift)
)

// IOCtlAdminCmd issues an received admin command.
func IOCtlAdminCmd(file *os.File, cmd *AdminCmd) error {
	return ioctl.Submit(file, uintptr(iocAdminCmd), uintptr(unsafe.Pointer(cmd)))
}

package nvme

import (
	"fmt"
	"github.com/sungup/go-nvme-ioctl/pkg/ioctl"
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

type nvmeCmd interface {
	SetData(data interface{}) error
	SetMeta(data interface{}) error
}

func getPtr(data interface{}) (uintptr, uint32, error) {
	v := reflect.ValueOf(data)
	if v.Type().Kind() != reflect.Ptr {
		return 0, 0, fmt.Errorf("input data is not pointer")
	}

	return v.Pointer(), uint32(reflect.Indirect(v).Type().Size()), nil
}

type UserIo struct {
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

func (c *PassthruCmd) SetData(data interface{}) error {
	if ptr, size, err := getPtr(data); err == nil {
		c.Data, c.DataLength = ptr, size
		return nil
	} else {
		return err
	}
}

func (c *PassthruCmd) SetMeta(data interface{}) error {
	if ptr, size, err := getPtr(data); err == nil {
		c.Meta, c.MetaLength = ptr, size
		return nil
	} else {
		return err
	}
}

type PassthruCmd32 struct {
	PassthruCmd

	TimeoutMSec uint32
	Result      uint32
}

type PassthruCmd64 struct {
	PassthruCmd

	TimeoutMSec uint32
	_           uint32
	Result      uint64
}

type AdminCmd struct {
	PassthruCmd

	TimeoutMSec uint32
	Result      uint32
}

const (
	iocNVMeType = 'N' << ioctl.TypeShift

	iocId          = ioctl.IOCNone | iocNVMeType | (0x40 << ioctl.NrShift)
	iocAdminCmd    = ioctl.IOCInOut | iocNVMeType | (0x41 << ioctl.NrShift) | uint64(unsafe.Sizeof(AdminCmd{})<<ioctl.SizeShift)
	iocSubmitIO    = ioctl.IOCIn | iocNVMeType | (0x42 << ioctl.NrShift) | uint64(unsafe.Sizeof(UserIo{})<<ioctl.SizeShift)
	iocIOCmd       = ioctl.IOCInOut | iocNVMeType | (0x43 << ioctl.NrShift) | uint64(unsafe.Sizeof(PassthruCmd32{})<<ioctl.SizeShift)
	iocReset       = ioctl.IOCNone | iocNVMeType | (0x44 << ioctl.NrShift)
	iocSubSysReset = ioctl.IOCNone | iocNVMeType | (0x45 << ioctl.NrShift)
	iocRescan      = ioctl.IOCNone | iocNVMeType | (0x46 << ioctl.NrShift)
	iocAdminCmd64  = ioctl.IOCInOut | iocNVMeType | (0x47 << ioctl.NrShift) | uint64(unsafe.Sizeof(PassthruCmd64{})<<ioctl.SizeShift)
	iocIOCmd64     = ioctl.IOCInOut | iocNVMeType | (0x48 << ioctl.NrShift) | uint64(unsafe.Sizeof(PassthruCmd64{})<<ioctl.SizeShift)
)

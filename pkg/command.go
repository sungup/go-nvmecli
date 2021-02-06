package pkg

type NVMePassthruCmd32 struct {
	OpCode      uint8
	Flags       uint8
	_           uint16 // Reserved
	NSId        uint32
	CDW2        uint32
	CDW3        uint32
	Meta        uintptr
	Data        uintptr
	MetaLength  uint32
	DataLength  uint32
	CDW10       uint32
	CDW11       uint32
	CDW12       uint32
	CDW13       uint32
	CDW14       uint32
	CDW15       uint32
	TimeoutMSec uint32
	Result      uint32
}

type NVMePassthruCmd64 struct {
	OpCode      uint8
	Flags       uint8
	_           uint16 // Reserved
	NSId        uint32
	CDW2        uint32
	CDW3        uint32
	Meta        uintptr
	Data        uintptr
	MetaLength  uint32
	DataLength  uint32
	CDW10       uint32
	CDW11       uint32
	CDW12       uint32
	CDW13       uint32
	CDW14       uint32
	CDW15       uint32
	TimeoutMSec uint32
	_           uint32
	Result      uint64
}

type NVMeAdminCmd NVMePassthruCmd32

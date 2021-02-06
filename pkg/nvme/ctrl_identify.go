package nvme

import (
	"github.com/sungup/go-nvme-ioctl/pkg/ioctl"
	"os"
	"unsafe"
)

type PowerStateDescriptor struct {
	MP       uint16
	_        [1]byte
	MxpsNops byte
	ENLAT    uint32
	EXLAT    uint32
	RRT      uint8
	RRL      uint8
	RWT      uint8
	RWL      uint8
	IDLP     uint16
	IPS      uint8
	_        [1]byte
	ACTP     uint16
	ApwAps   uint8
	_        [9]byte
}

type CtrlIdentify struct {
	VID   uint16
	SSVID uint16
	SN    [20]byte
	MN    [40]byte
	FR    [8]byte
	RAB   uint8
	IEEE  [3]byte
	CMIC  byte
	MDTS  uint8

	CNTLID uint16
	VER    uint32
	RTD3R  uint32
	RTD3E  uint32
	OAES   uint32

	CTRATT uint32

	RRLS      uint16
	_         [9]byte // Reserved
	CNTRLTYPE byte

	FGUID [2]uint64 // TODO make 128byte long int structure
	CRDT1 uint16
	CRDT2 uint16
	CRDT3 uint16
	_     [106]byte // Reserved
	_     [16]byte  // Refer to the NVMe Management Interface Specification for definition.

	// Admin Command Set Attributes & optional Controller Capabilities
	OACS uint16
	ACL  uint8
	AERL uint8

	FRMW  uint8
	LPA   uint8
	ELPE  uint8
	NPSS  uint8
	AVSCC uint8

	APSTA   uint8
	WCTEMP  uint16
	CCTEMP  uint16
	MTFA    uint16
	HMPRE   uint32
	HMMIN   uint32
	TNVMCAP [2]uint64 // TODO make 128byte long int structure
	UNVMCAP [2]uint64 // TODO make 128byte long int structure

	RPMBS uint32
	EDSTT uint16
	DSTO  uint8
	FWUG  uint8

	KAS   uint16
	HCTMA uint16
	MNTMT uint16
	MXTMT uint16

	SNICAP    uint32
	HMMINDS   uint32
	HMMAXD    uint16
	NSETIDMAX uint16
	ENDGIDMAX uint16
	ANATT     uint8

	ANACAP    uint8
	ANAGRPMAX uint32
	NANAGRPID uint32
	PELS      uint32
	_         [156]byte

	// NVM Command Set Attributes
	SQES   uint8
	CQES   uint8
	MAXCMD uint16
	NN     uint32

	ONCS  uint16
	FUSES uint16
	FNA   uint8

	VWC  uint8
	AWUN uint16

	AWUPF uint16
	NVSCC uint8
	NWPC  uint8

	ACWU uint16
	_    [2]byte

	SQLS   uint32
	MNAN   uint32
	_      [224]byte // Reserved
	SUBNQN [256]byte
	_      [768]byte // Reserved
	_      [256]byte // Refer to the NVMe over Fabrics specification.

	// Power State Descriptors
	PSD [32]PowerStateDescriptor

	// Vendor Specific
	_ [1024]byte
}

func newIdentifyCmd13(namespaceId, cdw10, cdw11 uint32, data interface{}) (*AdminCmd, error) {
	cmd := AdminCmd{
		PassthruCmd: PassthruCmd{
			OpCode: AdminIdentify,
			NSId:   namespaceId,
			CDW10:  cdw10,
			CDW11:  cdw11,
		},
		TimeoutMSec: 0,
		Result:      0,
	}

	err := cmd.SetData(data)

	return &cmd, err
}

func newIdentifyCmd(namespaceId, cdw10 uint32, data interface{}) (*AdminCmd, error) {
	return newIdentifyCmd13(namespaceId, cdw10, 0, data)
}

func GetCtrlIdentify(file *os.File) (*CtrlIdentify, error) {
	// TODO reimplementing here
	identify := CtrlIdentify{}

	cmd, _ := newIdentifyCmd(0, 1, &identify)

	if err := ioctl.Submit(file, uintptr(iocAdminCmd), uintptr(unsafe.Pointer(cmd))); err == nil {
		return &identify, nil
	} else {
		return nil, err
	}
}

package identify

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/sungup/go-nvmecli/pkg/nvme"
	"github.com/sungup/go-nvmecli/pkg/nvme/types"
	"github.com/sungup/go-nvmecli/pkg/utils"
	"math"
	"os"
	"unsafe"
)

const (
	// Controller Identify data size is 4096B
	ctrlIdentifySz = 4096
)

const (
	cnsNamespace    = uint16(0x00)
	cnsController   = uint16(0x01)
	cnsActiveNSList = uint16(0x02)
	cnsNSDescList   = uint16(0x03)
	cnsNVMSetList   = uint16(0x04)
)

// newIdentifyCmd generates an AdminCmd structure to retrieve the NVMe's identify related structure.
// The cntid and cns will be set on CDW10 and nvmSetId also set on CDW11.
func newIdentifyCmd(nsid uint32, cntid, cns, nvmSetId uint16, v interface{}) (*nvme.AdminCmd, error) {
	cmd := nvme.AdminCmd{
		PassthruCmd: nvme.PassthruCmd{
			OpCode: nvme.AdminIdentify,
			NSId:   nsid,
			CDW10:  uint32(cntid)<<16 | uint32(cns),
			CDW11:  uint32(nvmSetId),
		},
		TimeoutMSec: 0,
		Result:      0,
	}

	if err := cmd.SetData(v); err != nil {
		return nil, err
	} else {
		return &cmd, nil
	}

}

// GetCtrlIdentify fills v interface with the controller identify data from an NVMe device.
func GetCtrlIdentify(file *os.File, v interface{}) error {
	if cmd, err := newIdentifyCmd(0, 0, cnsController, 0, v); err != nil {
		return err
	} else {
		return nvme.IOCtlAdminCmd(file, cmd)
	}
}

// powerStateDesc is an entry of structure for Power State Descriptor in NVMe SPEC. To prevent
// externally creation, make it private structure
type powerStateDesc struct {
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

// CtrlIdentify is an structure for the controller identify information of an NVMe device.
type CtrlIdentify struct {
	// Controller Capabilities and Features
	VID   types.VID
	SSVID types.SSVID
	SN    types.SN
	MN    types.MN
	FR    [8]byte
	RAB   types.Uint8
	IEEE  types.IEEE
	CMIC  types.Hex8
	MDTS  types.Uint8

	CNTLID uint16
	VER    uint32
	RTD3R  uint32
	RTD3E  uint32
	OAES   uint32

	CTRATT uint32

	RRLS      uint16
	_         [9]byte // Reserved
	CNTRLTYPE byte

	FGUID types.Uint128
	CRDT1 uint16
	CRDT2 uint16
	CRDT3 uint16
	_     [106]byte // Reserved
	_     [16]byte  // Refer to the NVMe Management Interface Specification for definition.

	// Admin Command Set Attributes & optional Controller Capabilities
	OACS uint16
	ACL  types.Uint8
	AERL types.Uint8

	FRMW  types.Uint8
	LPA   types.Uint8
	ELPE  types.Uint8
	NPSS  types.Uint8
	AVSCC types.Uint8

	APSTA   types.Uint8
	WCTEMP  uint16
	CCTEMP  uint16
	MTFA    uint16
	HMPRE   uint32
	HMMIN   uint32
	TNVMCAP types.Uint128
	UNVMCAP types.Uint128

	RPMBS uint32
	EDSTT uint16
	DSTO  types.Uint8
	FWUG  types.Uint8

	KAS   uint16
	HCTMA uint16
	MNTMT uint16
	MXTMT uint16

	SNICAP    uint32
	HMMINDS   uint32
	HMMAXD    uint16
	NSETIDMAX uint16
	ENDGIDMAX uint16
	ANATT     types.Uint8

	ANACAP    types.Uint8
	ANAGRPMAX uint32
	NANAGRPID uint32
	PELS      uint32
	_         [156]byte

	// NVM Command Set Attributes
	SQES   types.Uint8
	CQES   types.Uint8
	MAXCMD uint16
	NN     uint32

	ONCS  uint16
	FUSES uint16
	FNA   types.Uint8

	VWC  types.Uint8
	AWUN uint16

	AWUPF uint16
	NVSCC types.Uint8
	NWPC  types.Uint8

	ACWU uint16
	_    [2]byte

	SQLS   uint32
	MNAN   uint32
	_      [224]byte // Reserved
	SUBNQN [256]byte
	_      [768]byte // Reserved
	_      [256]byte // Refer to the NVMe over Fabrics specification.

	// Power State Descriptors
	PSD [32]powerStateDesc

	// Vendor Specific
	_ [1024]byte
}

// ParseCtrlIdentify creates an CtrlIdentify object from byte slice to convert easy usable
// controller identify data format.
func ParseCtrlIdentify(raw []byte) (*CtrlIdentify, error) {
	if len(raw) != int(unsafe.Sizeof(CtrlIdentify{})) {
		return nil, fmt.Errorf("unexpected identify raw data size: %d", len(raw))
	}

	i := CtrlIdentify{}

	if err := binary.Read(bytes.NewReader(raw), utils.SystemEndian, &i); err == nil {
		return &i, nil
	} else {
		return nil, err
	}
}

// GetNamespaceIdentify fills v interface with the namespace identify data from a namespace of NVMe
// device.
func GetNamespaceIdentify(file *os.File, nsid uint32, v interface{}) error {
	if cmd, err := newIdentifyCmd(nsid, 0, cnsNamespace, 0, v); err != nil {
		return err
	} else {
		return nvme.IOCtlAdminCmd(file, cmd)
	}
}

// relativePerf is the relative performance of the LBA format indicated relative to other LBA
// formats supported by the controller.
type relativePerf uint8

const (
	NSBestPerf     = relativePerf(0b00)
	NSBetterPerf   = relativePerf(0b01)
	NSGoodPerf     = relativePerf(0b10)
	NSDegradedPerf = relativePerf(0b11)
)

// lbaFormat is a structure supported by the controller.
type lbaFormat uint32

// RelativePerformance returns relativePerf of the LBA format.
func (l lbaFormat) RelativePerformance() relativePerf {
	return relativePerf(uint8(l>>24) & 0x03)
}

// LBADataSize returns the LBA data size supported. The value is reported in terms of a power of
// two (2^n), and cannot be smaller than 9 (512B). If returned value is 0h, target LBA format is
// not supported/used or is not currently available.
func (l lbaFormat) LBADataSize() int {
	return int(l >> 16 & math.MaxUint8)
}

// MetadataSize returns the number of metadata bytes provided per LBA based on the LBA data size.
// If there is no metadata supported, then this function returns 0.
func (l lbaFormat) MetadataSize() int {
	return int(l & math.MaxUint16)
}

// NamespaceIdentify is an structure for identify information for an namespace in NVMe device.
type NamespaceIdentify struct {
	NSZE     uint64        // [07:00]    M: Namespace Size
	NCAP     uint64        // [15:08]    M: Namespace Capacity
	NUSE     uint64        // [23:16]    M: Namespace Utilization
	NSFEAT   uint8         // [24]       M: Namespace Features
	NLBAF    uint8         // [25]       M: Number of LBA Formats
	FLBAS    uint8         // [26]       M: Formatted LBA Size
	MC       uint8         // [27]       M: Metadata Capabilities
	DPC      uint8         // [28]       M: End-to-End Protection Capabilities
	DPS      uint8         // [29]       M: End-to-End Protection Type Settings
	NMIC     uint8         // [30]       O: Namespace Multi-path I/O and Namespace Sharing Capabilities
	RESCAP   uint8         // [31]       O: Reservation Capabilities
	FPI      uint8         // [32]       O: Format Progress Indicator
	DLFEAT   uint8         // [33]       O: Deallocate Logical Block Features
	NAWUN    uint16        // [35:34]    O: Namespace Atomic Write Unit Normal
	NAWUPF   uint16        // [37:36]    O: Namespace Atomic Write Unit Power Fail
	NACWU    uint16        // [39:38]    O: Namespace Atomic Compare & Write Unit
	NABSN    uint16        // [41:40]    O: Namespace Atomic Boundary Size Normal
	NABO     uint16        // [43:42]    O: Namespace Atomic Boundary Offset
	NABSPF   uint16        // [45:44]    O: Namespace Atomic Boundary Size Power Fail
	NOIOB    uint16        // [47:46]    O: Namespace Optimal I/O Boundary
	NVMCAP   types.Uint128 // [63:48]    O: NVM Capacity
	NPWG     uint16        // [65:64]    O: Namespace Preferred Write Granularity
	NPWA     uint16        // [67:66]    O: Namespace Preferred Write Alignment
	NPDG     uint16        // [69:68]    O: Namespace Preferred Deallocate Granularity
	NPDA     uint16        // [71:70]    O: Namespace Preferred Deallocate Alignment
	NOWS     uint16        // [73:72]    O: Namespace Optimal Write Size
	_        [18]byte      // [91:74]    reserved
	ANAGRPID uint32        // [95:92]    O: ANA Group Identifier
	_        [3]byte       // [98:96]    reserved
	NSATTR   uint8         // [99]       O: Namespace Attributes
	NVMSETID uint16        // [101:100]  O: NVM Set Identifier
	ENDGID   uint16        // [103:102]  O: Endurance Group Identifier
	NGUID    types.NGUID   // [119:104]  O: Namespace Globally Unique Identifier
	EUI64    types.EUI64   // [127:120]  O: IEEE Extended Unique Identifier
	LBAF     [16]lbaFormat // [191:128]  M/O: LBA Format 0~15 Support. Only 0 is mandatory
	_        [192]byte     // [383:192]  reserved
	Vendor   [3712]byte    // [4095:384] O: Vendor Specific
}

// ParseNamespaceIdentify create an NamespaceIdentify object from byte slice to convert easy usable
// namespace identify data format.
func ParseNamespaceIdentify(raw []byte) (*NamespaceIdentify, error) {
	if len(raw) != int(unsafe.Sizeof(NamespaceIdentify{})) {
		return nil, fmt.Errorf("unexpected identify raw data size: %d", len(raw))
	}

	i := NamespaceIdentify{}

	if err := binary.Read(bytes.NewReader(raw), utils.SystemEndian, &i); err == nil {
		return &i, nil
	} else {
		return nil, err
	}
}

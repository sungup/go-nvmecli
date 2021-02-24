package identify

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/sungup/go-nvmecli/pkg/nvme"
	"github.com/sungup/go-nvmecli/pkg/nvme/types"
	"github.com/sungup/go-nvmecli/pkg/utils"
	"os"
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

// GetCtrlIdentify returns 4096B byte slice which contains the controller identify data from an
// NVMe device. If you want the parsed controller identify data, call the ParseCtrlIdentify using
// the returned byte slice.
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

// CtrlIdentify is an structure for the controller identify information of an NVMe device. To
// prevent externally creation, make it private structure.
type CtrlIdentify struct {
	// Controller Capabilities and Features
	VID   types.Hex16
	SSVID types.Hex16
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

	FGUID types.Uint128
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
	TNVMCAP types.Uint128
	UNVMCAP types.Uint128

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
	PSD [32]powerStateDesc

	// Vendor Specific
	_ [1024]byte
}

// ParseCtrlIdentify creates an CtrlIdentify object to manipulate each value from outsize of this
// package.
func ParseCtrlIdentify(raw []byte) (*CtrlIdentify, error) {
	if len(raw) != ctrlIdentifySz {
		return nil, fmt.Errorf("unexpected identify raw data size: %d", len(raw))
	}

	i := CtrlIdentify{}

	if err := binary.Read(bytes.NewReader(raw), utils.SystemEndian, &i); err == nil {
		return &i, nil
	} else {
		return nil, err
	}
}

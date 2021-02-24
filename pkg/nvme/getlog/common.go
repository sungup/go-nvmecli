package getlog

import (
	"github.com/sungup/go-nvme-ioctl/pkg/nvme"
	"math"
	"os"
)

const (
	logPageErrorInfo      = uint8(0x01)
	logPageSMART          = uint8(0x02)
	logPageFWSlot         = uint8(0x03)
	logPageChangedNsList  = uint8(0x04)
	logPageCommandSupport = uint8(0x05)
	logPageDevSelfTest    = uint8(0x06)
	logPageTelemetryHost  = uint8(0x07)
	logPageTelemetryCtrl  = uint8(0x08)
	logPageEndurGrpInfo   = uint8(0x09)
	logPagePredLatNVMSet  = uint8(0x0A)
	logPagePredLatEvt     = uint8(0x0B)
	logPageAsyncNsAccess  = uint8(0x0C)
	logPagePersistEvtLog  = uint8(0x0D)
	logPageLBAStatusInfo  = uint8(0x0E)
	logPageEndurGrpEvt    = uint8(0x0F)

	getLogSMARTSz             = 512
	getLogTelemetryHeaderSz   = 512
	getLogTelemetryBlkSzShift = 9

	maskUint4   = uint32(1<<4 - 1)
	maskUint8   = uint32(math.MaxUint8)
	maskUint16  = uint32(math.MaxUint16)
	maskUint32  = uint64(math.MaxUint32)
	umaskUint16 = ^maskUint16
	shiftUint8  = 8
	shiftUint16 = 16
	shiftUint32 = 32
)

type getLogCmd struct {
	nvme.AdminCmd
}

// SetDWords changes the dwords fields (NUMDL and NUMDU).
func (l *getLogCmd) SetDWords(dwords uint32) {
	l.CDW10 = (dwords << shiftUint16) | (l.CDW10 & maskUint16)
	l.CDW11 = (dwords >> shiftUint16) | (l.CDW11 & umaskUint16)
}

// SetDWords changes the offset fields (LPOL and LPOU).
func (l *getLogCmd) SetOffset(offset uint64) {
	l.CDW12 = uint32(offset >> shiftUint32)
	l.CDW13 = uint32(offset & maskUint32)
}

// SetLSP change the 4bit Log Specific Identifier.
func (l *getLogCmd) SetLSP(lsp uint8) {
	const UMaskLSP = ^(maskUint4 << shiftUint8)

	l.CDW10 = (l.CDW10 & UMaskLSP) | (uint32(lsp)&maskUint4)<<shiftUint8
}

// newGetLogCmd generate an AdminCmd structure to retrieve the NVMe's log pages. To issue a get-log
// command, the size and offset should be set for the large size log data like the Telemetry.
// However the dwords is a value to retrieve from log page and not same with the size of the dptr.
// If user issue get-log command larger than specified size of the real log data, get-log commands
// will return with undefined results beyond the end of the log page.
// Host software should clear the RAE bit to '0' for log pages that are not used with Asynchronous
// Events.
func newGetLogCmd(nsid uint32, offset uint64, lid, lsp uint8, lsi uint16, v interface{}) (*getLogCmd, error) {
	cmd := getLogCmd{
		nvme.AdminCmd{
			PassthruCmd: nvme.PassthruCmd{
				OpCode: nvme.AdminGetLogPage,
				NSId:   nsid,
				CDW10:  (uint32(lsp)&maskUint4)<<shiftUint8 | uint32(lid)&maskUint8,
				CDW11:  uint32(lsi) << shiftUint16,
				CDW12:  uint32(offset >> shiftUint32), // Log Page Offset Lower
				CDW13:  uint32(offset & maskUint32),   // Log Page Offset Upper
			},
			TimeoutMSec: 0,
			Result:      0,
		},
	}

	if err := cmd.SetData(v); err != nil {
		return nil, err
	} else {
		cmd.SetDWords(cmd.DataLength >> 2)
		return &cmd, nil
	}
}

// --------------------------------- //
// LID XXh: Undefined Vendor Command //
// --------------------------------- //

// GetVendorCMD retrieve a log data for the vendor specific command. The vendorID is an aliased
// parameter about the lid (Log Page Identifier).
func GetVendorCMD(file *os.File, nsid uint32, vendorID, lsp uint8, lsi uint16, v interface{}) error {
	if cmd, err := newGetLogCmd(nsid, 0, vendorID, lsp, lsi, v); err != nil {
		return err
	} else {
		return nvme.IOCtlAdminCmd(file, &cmd.AdminCmd)
	}
}

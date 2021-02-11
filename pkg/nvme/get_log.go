package nvme

import "math"

const (
	logPageErrorInfo      = uint16(0x01)
	logPageSMART          = uint16(0x02)
	logPageFWSlot         = uint16(0x03)
	logPageChangedNsList  = uint16(0x04)
	logPageCommandSupport = uint16(0x05)
	logPageDevSelfTest    = uint16(0x06)
	logPageTelemetryHost  = uint16(0x07)
	logPageTelemetryCtrl  = uint16(0x08)
	logPageEndurGrpInfo   = uint16(0x09)
	logPagePredLatNVMSet  = uint16(0x0A)
	logPagePredLatEvt     = uint16(0x0B)
	logPageAsyncNsAccess  = uint16(0x0C)
	logPagePersistEvtLog  = uint16(0x0D)
	logPageLBAStatusInfo  = uint16(0x0E)
	logPageEndurGrpEvt    = uint16(0x0F)
)

// newGetLogCmd generate an AdminCmd structure to retrieve the NVMe's log pages. To issue a get-log
// command, the size and offset should be set for the large size log data like the telemetry.
// However the dwords is a value to retrieve from log page and not same with the size of the dptr.
// If user issue get-log command larger than specified size of the real log data, get-log commands
// will return with undefined results beyond the end of the log page.
// Host software should clear the RAE bit to '0' for log pages that are not used with Asynchronous
// Events.
func newGetLogCmd(nsid, dwords uint32, offset uint64, lid, lsp, lsi uint16, dptr interface{}) (*AdminCmd, error) {
	const (
		MaskUint4   = uint32(1<<4 - 1)
		MaskUint8   = uint32(math.MaxUint8)
		MaskUint16  = uint32(math.MaxUint16)
		MaskUint32  = uint64(math.MaxUint32)
		ShiftUint8  = 8
		ShiftUint16 = 16
		ShiftUint32 = 32
	)

	cmd := AdminCmd{
		PassthruCmd: PassthruCmd{
			OpCode: AdminGetLogPage,
			NSId:   nsid,
			CDW10:  dwords<<ShiftUint16 | (uint32(lsp)&MaskUint4)<<ShiftUint8 | uint32(lid)&MaskUint8,
			CDW11:  uint32(lsi)<<ShiftUint16 | dwords>>ShiftUint16,
			CDW12:  uint32(offset >> ShiftUint32), // Log Page Offset Lower
			CDW13:  uint32(offset & MaskUint32),   // Log Page Offset Upper
		},
		TimeoutMSec: 0,
		Result:      0,
	}

	err := cmd.SetData(dptr)

	return &cmd, err
}

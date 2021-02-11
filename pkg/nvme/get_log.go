package nvme

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/sungup/go-nvme-ioctl/pkg/utils"
	"math"
	"os"
)

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

	getLogSMARTSz = 512
)

// newGetLogCmd generate an AdminCmd structure to retrieve the NVMe's log pages. To issue a get-log
// command, the size and offset should be set for the large size log data like the telemetry.
// However the dwords is a value to retrieve from log page and not same with the size of the dptr.
// If user issue get-log command larger than specified size of the real log data, get-log commands
// will return with undefined results beyond the end of the log page.
// Host software should clear the RAE bit to '0' for log pages that are not used with Asynchronous
// Events.
func newGetLogCmd(nsid, dwords uint32, offset uint64, lid, lsp, lsi uint16) *AdminCmd {
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

	return &cmd
}

// GetLogSMART will retrieve SMART data from NVMe device.
func GetLogSMART(file *os.File) ([]byte, error) {
	return ioctlAdminCmd(
		file,
		getLogSMARTSz,
		func() *AdminCmd { return newGetLogCmd(0, getLogSMARTSz>>2, 0, logPageSMART, 0, 0) },
	)
}

// TODO rebuild SMART structure because of 16B unsigned integer parsing
type smart struct {
	CriticalWarning         uint8    // 00
	CompositeTemperature    [2]uint8 // TODO need the 2B custom integer because uint16 occur unaligned error
	AvailableSpare          uint8    // 03
	AvailableSpareThreshold uint8    // 04
	PercentageUsed          uint8    // 05
	EnduranceGrpCritWarn    uint8    // 06
	_                       [25]byte

	DataUnitsRead       [2]uint64 // TODO need the 16B integer
	DataUnitsWritten    [2]uint64 // TODO need the 16B integer
	HostReadCommands    [2]uint64 // TODO need the 16B integer
	HostWrittenCommands [2]uint64 // TODO need the 16B integer

	ControllerBusyTime [2]uint64 // TODO need the 16B integer
	PowerCycles        [2]uint64 // TODO need the 16B integer
	PowerOnHours       [2]uint64 // TODO need the 16B integer
	UnsafeShutdowns    [2]uint64 // TODO need the 16B integer
	IntegrityErrors    [2]uint64 // TODO need the 16B integer
	ErrorLogEntries    [2]uint64 // TODO need the 16B integer

	WarningCompositeTempTime     uint32
	CriticalCompositeTempTime    uint32
	TemperatureSensor            [8]uint16
	ThermalMgmtTemp1TransCount   uint32
	ThermalMgmtTemp2TransCount   uint32
	TotalTimeForThermalMgmtTemp1 uint32
	TotalTimeForThermalMgmtTemp2 uint32

	_ [280]byte // reserved
}

// ParseSMART creates an SMART object to get the device health information.
func ParseSMART(raw []byte) (*smart, error) {
	if len(raw) != getLogSMARTSz {
		return nil, fmt.Errorf("unexpected SMART raw data size: %d", len(raw))
	}

	s := smart{}

	if err := binary.Read(bytes.NewReader(raw), utils.SystemEndian, &s); err == nil {
		return &s, nil
	} else {
		return nil, err
	}
}

func GetLogTelemetryHostInit(file *os.File, dataBlock uint) {
	// TODO implementing here
}

func GetLogTelemetryCtrlInit(file *os.File, dataBlock uint) {
	// TODO implementing here
}

// telemetry is a header of the telemetry page. To retrieve the telemetry log, host SW should call
// ioctl with only header to generate and check the data block size. After that, host SW can
// retrieve the telemetry log with one more ioctl command.
type telemetry struct {
	Identifier byte
	_          [4]byte // reserved
	IEEE       [3]byte

	DataAreaLastBlock [3]uint16

	_ [368]byte //reserved

	CtrlInitiativeAvailable    uint8
	CtrlInitiativeGenerationNo uint8
	ReasonIdentifier           [128]byte
}

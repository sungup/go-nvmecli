package nvme

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/sungup/go-nvme-ioctl/pkg/ioctl"
	"github.com/sungup/go-nvme-ioctl/pkg/utils"
	"math"
	"os"
	"unsafe"
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
	AdminCmd
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
func (l *getLogCmd) SetLSP(lsp uint16) {
	const UMaskLSP = ^(maskUint4 << shiftUint8)

	l.CDW10 = (l.CDW10 & UMaskLSP) | (uint32(lsp)&maskUint4)<<shiftUint8
}

// newGetLogCmd generate an AdminCmd structure to retrieve the NVMe's log pages. To issue a get-log
// command, the size and offset should be set for the large size log data like the telemetry.
// However the dwords is a value to retrieve from log page and not same with the size of the dptr.
// If user issue get-log command larger than specified size of the real log data, get-log commands
// will return with undefined results beyond the end of the log page.
// Host software should clear the RAE bit to '0' for log pages that are not used with Asynchronous
// Events.
func newGetLogCmd(nsid uint32, offset uint64, lid, lsp, lsi uint16, v interface{}) (*getLogCmd, error) {
	cmd := getLogCmd{
		AdminCmd{
			PassthruCmd: PassthruCmd{
				OpCode: AdminGetLogPage,
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

// ----------------------------------- //
// LID 02h: SMART / Health Information //
// ----------------------------------- //

// GetLogSMART will retrieve SMART data from NVMe device.
func GetLogSMART(file *os.File, v interface{}) error {
	if cmd, err := newGetLogCmd(0, 0, logPageSMART, 0, 0, v); err != nil {
		return err
	} else {
		return ioctlAdminCmd(file, func() *AdminCmd { return &cmd.AdminCmd })
	}
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

// ---------------------------------------------------------------------------- //
// LID 07h: Telemetry Host-Initiative, LID 08h: Telemetry Controller-Initiative //
// ---------------------------------------------------------------------------- //
type telemetryDataBlk int

const (
	DataBlock1 = telemetryDataBlk(1)
	DataBlock2 = telemetryDataBlk(2)
	DataBlock3 = telemetryDataBlk(3)
)

// index function change the telemetryDataBlk macro to index of telemetry.DataBlockLast's index
func (b telemetryDataBlk) index() int {
	return int(b) - 1
}

// getLogTelemetry retrieve telemetry data from NVMe device. Host-initiated and Ctrl-initiated
// telemetry has same format except lsp field, so this function receive the lid to determine the
// get-log Log Identifier and the lsp to create telemetry data for the Host-initiated telemetry.
func getLogTelemetry(file *os.File, block telemetryDataBlk, lid uint16, lsp uint16) ([]byte, error) {
	var (
		header *telemetry
		err    error
	)

	const request = uintptr(iocAdminCmd)

	buffer := make([]byte, maxAdminCmdPageSz)
	cmd, _ := newGetLogCmd(0, 0, lid, lsp, 0, buffer)
	cmdPtr := uintptr(unsafe.Pointer(cmd))

	// 1. get telemetry header logs with lsp value
	cmd.SetDWords(getLogTelemetryHeaderSz >> 2)

	if err = ioctl.Submit(file, request, cmdPtr); err != nil {
		return nil, err
	}

	if header, err = ParseTelemetryHeader(buffer); err != nil {
		return nil, err
	}

	// 2. calculate retrieving data size, create return data and resize buffer
	dataSz := header.BlockSize(block)
	offset := uint32(getLogTelemetryHeaderSz)
	fetchSz := maxAdminCmdPageSz

	cmd.SetDWords(fetchSz >> 2)
	cmd.SetLSP(0x0)

	data := make([]byte, 0, getLogTelemetryHeaderSz+dataSz)
	data = append(data, buffer[:getLogTelemetryHeaderSz]...)

	// 3. retrieving telemetry log by basic page size
	for offset < dataSz {
		// 3-1. update offset and dwords of getLogCmd
		if dataSz < offset+fetchSz {
			fetchSz = dataSz - offset
			cmd.SetDWords(fetchSz >> 2)
		}
		cmd.SetOffset(uint64(offset))

		// 3-2. resend ioctl to device
		if err := ioctl.Submit(file, request, cmdPtr); err != nil {
			return nil, err
		}

		data = append(data, buffer[:fetchSz]...)

		// 3-3. move next offset
		offset += fetchSz
	}

	return data, nil
}

// GetLogTelemetryHostInit retrieves the host-initiated  telemetry data from NVMe device. If host sw
// call with create=true, GetLogTelemetryHostInit will recreate the host-initiated telemetry data
// before gathering telemetry data
func GetLogTelemetryHostInit(file *os.File, block telemetryDataBlk, create bool) ([]byte, error) {
	var lsp uint16 = 0x00
	if create {
		lsp = 0x01
	}

	return getLogTelemetry(file, block, logPageTelemetryHost, lsp)
}

// GetLogTelemetryCtrlInit retrieves the controller-initiated telemetry data from NVMe device.
func GetLogTelemetryCtrlInit(file *os.File, block telemetryDataBlk) ([]byte, error) {
	return getLogTelemetry(file, block, logPageTelemetryCtrl, 0x0)
}

// telemetry is a header of the telemetry page. To retrieve the telemetry log, host SW should call
// ioctl with only header to generate and check the data block size. After that, host SW can
// retrieve the telemetry log with one more ioctl command. Host-initiated log and Controller-
// initiated log have same format.
type telemetry struct {
	Identifier byte    // [00]
	_          [4]byte // [04:01] reserved
	IEEE       [3]byte // [07:05]

	DataAreaLastBlock [3]uint16 // [13:08]

	_ [368]byte // [381:14] reserved

	CtrlInitiativeAvailable    uint8     // [382]
	CtrlInitiativeGenerationNo uint8     // [383]
	ReasonIdentifier           [128]byte // [511:384]
}

// BlockSize returns the Byte unit each data block size calculating the DataAreaLastBlock.
func (t *telemetry) BlockSize(block telemetryDataBlk) uint32 {
	return uint32(t.DataAreaLastBlock[block.index()]) << getLogTelemetryBlkSzShift
}

// ParseTelemetryHeader parses the telemetry's header information from raw data. If the size of raw
// data is under 512B, this function raise an error.
func ParseTelemetryHeader(raw []byte) (*telemetry, error) {
	if len(raw) < getLogTelemetryHeaderSz {
		return nil, fmt.Errorf("unexpected Telemetry header size: %d", len(raw))
	}

	t := telemetry{}
	if err := binary.Read(bytes.NewReader(raw), utils.SystemEndian, &t); err == nil {
		return &t, nil
	} else {
		return nil, err
	}
}

// --------------------------------- //
// LID XXh: Undefined Vendor Command //
// --------------------------------- //

// GetLogVendorCMD retrieve a log data for the vendor specific command. The vendorID is an aliased
// parameter about the lid (Log Page Identifier).
func GetLogVendorCMD(file *os.File, nsid uint32, vendorID, lsp, lsi uint16, v interface{}) error {
	if cmd, err := newGetLogCmd(nsid, 0, vendorID, lsp, lsi, v); err != nil {
		return err
	} else {
		return ioctlAdminCmd(file, func() *AdminCmd { return &cmd.AdminCmd })
	}
}

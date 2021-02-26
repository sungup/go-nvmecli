package getlog

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/sungup/go-nvmecli/pkg/nvme"
	"github.com/sungup/go-nvmecli/pkg/nvme/types"
	"github.com/sungup/go-nvmecli/pkg/utils"
	"os"
)

// ---------------------------------------------------------------------------- //
// LID 07h: Telemetry Host-Initiative, LID 08h: Telemetry Controller-Initiative //
// ---------------------------------------------------------------------------- //

type telemetryDataBlk int

const (
	DataBlock1 = telemetryDataBlk(1)
	DataBlock2 = telemetryDataBlk(2)
	DataBlock3 = telemetryDataBlk(3)

	maxTelemetryPageSz = uint32(4096)

	telemetryHeaderSz   = uint32(512)
	telemetryBlkSzShift = 9
)

// index function change the telemetryDataBlk macro to index of Telemetry.DataBlockLast's index
func (b telemetryDataBlk) index() int {
	return int(b) - 1
}

// getLogTelemetry retrieve telemetry data from NVMe device. Host-initiated and Ctrl-initiated
// telemetry has same format except lsp field, so this function receive the lid to determine the
// get-log Log Identifier and the lsp to create telemetry data for the Host-initiated telemetry.
func getLogTelemetry(file *os.File, block telemetryDataBlk, lid, lsp uint8) ([]byte, error) {
	var (
		header *Telemetry
		err    error
	)

	buffer := make([]byte, maxTelemetryPageSz)
	cmd, _ := newGetLogCmd(0, 0, lid, lsp, 0, buffer)

	// 1. get Telemetry header logs with lsp value
	cmd.SetDWords(telemetryHeaderSz >> 2)

	if err = nvme.IOCtlAdminCmd(file, &cmd.AdminCmd); err != nil {
		return nil, err
	}

	if header, err = ParseTelemetryHeader(buffer); err != nil {
		return nil, err
	}

	// 2. calculate retrieving data size, create return data and resize buffer
	dataSz := header.BlockSize(block)
	fetchSz := maxTelemetryPageSz

	cmd.SetDWords(fetchSz >> 2)
	cmd.SetLSP(0x0)

	data := make([]byte, 0, telemetryHeaderSz+dataSz)
	data = append(data, buffer[:telemetryHeaderSz]...)

	// 3. retrieving Telemetry log by basic page size
	for offset := telemetryHeaderSz; offset < dataSz; offset += fetchSz {
		// 3-1. update offset and dwords of getLogCmd
		if dataSz < offset+fetchSz {
			fetchSz = dataSz - offset
			cmd.SetDWords(fetchSz >> 2)
		}
		cmd.SetOffset(uint64(offset))

		// 3-2. resend ioctl to device
		if err := nvme.IOCtlAdminCmd(file, &cmd.AdminCmd); err != nil {
			return nil, err
		}

		data = append(data, buffer[:fetchSz]...)
	}

	return data, nil
}

// GetTelemetryHostInit retrieves the host-initiated telemetry data from NVMe device. If host sw
// call with create=true, GetTelemetryHostInit will recreate the host-initiated telemetry data
// before gathering telemetry data
func GetTelemetryHostInit(file *os.File, block telemetryDataBlk, create bool) ([]byte, error) {
	var lsp uint8 = 0x00
	if create {
		lsp = 0x01
	}

	return getLogTelemetry(file, block, logPageTelemetryHost, lsp)
}

// GetTelemetryCtrlInit retrieves the controller-initiated telemetry data from NVMe device.
func GetTelemetryCtrlInit(file *os.File, block telemetryDataBlk) ([]byte, error) {
	return getLogTelemetry(file, block, logPageTelemetryCtrl, 0x0)
}

// Telemetry is a header of the telemetry page. To retrieve the telemetry log, host SW should call
// ioctl with only header to generate and check the data block size. After that, host SW can
// retrieve the telemetry log with one more ioctl command. Host-initiated log and Controller-
// initiated log have same format.
type Telemetry struct {
	Identifier byte       // [00]
	_          [4]byte    // [04:01] reserved
	IEEE       types.IEEE // [07:05]

	DataAreaLastBlock [3]uint16 // [13:08]

	_ [368]byte // [381:14] reserved

	CtrlInitiativeAvailable    uint8     // [382]
	CtrlInitiativeGenerationNo uint8     // [383]
	ReasonIdentifier           [128]byte // [511:384]
}

// BlockSize returns the Byte unit each data block size calculating the DataAreaLastBlock.
func (t *Telemetry) BlockSize(block telemetryDataBlk) uint32 {
	return uint32(t.DataAreaLastBlock[block.index()]) << telemetryBlkSzShift
}

// ParseTelemetryHeader parses the Telemetry's header information from raw data. If the size of raw
// data is under 512B, this function raise an error.
func ParseTelemetryHeader(raw []byte) (*Telemetry, error) {
	if len(raw) < int(telemetryHeaderSz) {
		return nil, fmt.Errorf("unexpected Telemetry header size: %d", len(raw))
	}

	t := Telemetry{}
	if err := binary.Read(bytes.NewReader(raw), utils.SystemEndian, &t); err == nil {
		return &t, nil
	} else {
		return nil, err
	}
}

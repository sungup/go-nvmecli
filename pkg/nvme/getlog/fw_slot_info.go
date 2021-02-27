package getlog

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/sungup/go-nvmecli/pkg/nvme"
	"github.com/sungup/go-nvmecli/pkg/utils"
	"os"
	"unsafe"
)

// GetFirmwareSlotInfo will retrieve Firmware Slot Information (03h) from NVMe device.
func GetFirmwareSlotInfo(file *os.File, v interface{}) error {
	if cmd, err := newGetLogCmd(0, 0, logPageFWSlot, 0, 0, v); err != nil {
		return err
	} else {
		return nvme.IOCtlAdminCmd(file, &cmd.AdminCmd)
	}
}

// fwSlotNo is type to indexing the firmware slot number from 1-based to 0 based.
type fwSlotNo int

const (
	NoFw = fwSlotNo(0)
	FRS1 = fwSlotNo(1)
	FRS2 = fwSlotNo(2)
	FRS3 = fwSlotNo(3)
	FRS4 = fwSlotNo(4)
	FRS5 = fwSlotNo(5)
	FRS6 = fwSlotNo(6)
	FRS7 = fwSlotNo(7)
)

// index convert the 1-based number to the 0-based slice index.
func (f fwSlotNo) index() int {
	return int(f) - 1
}

// fwRev is a slice of 8 byte string contains revision information.
type fwRev [8]byte

// String of fwRev returns a string revision info from slice
func (f fwRev) String() string {
	const trimCutSet = " \x00"
	return string(bytes.Trim(f[:], trimCutSet))
}

// activeFwInfo contains a specified information about the active firmware revision.
type activeFwInfo byte

// Active returns the firmware slot number that is currently activated.
func (f activeFwInfo) Active() fwSlotNo {
	return fwSlotNo(int(f) & 0b111)
}

// Next returns the firmware slot number that is going to be activated at the next Controller Level
// Reset command. If Next returns 0, there is no a next active firmware after reset command, because
// slot number is on the 1-based index
func (f activeFwInfo) Next() fwSlotNo {
	return fwSlotNo(int(f) >> 4 & 0b111)
}

// FirmwareSlotInfo is a log page structure contains the active firmware information and the
// currently installed firmware revisions.
type FirmwareSlotInfo struct {
	ActiveFwInfo activeFwInfo
	_            [7]byte // reserved

	fwRevision [7]fwRev
	_          [448]byte // reserved
}

// Slot returns a revision at the slot number. If slot number is NoFw, It will returns empty string.
func (f *FirmwareSlotInfo) Slot(slot fwSlotNo) string {
	if slot != NoFw {
		return f.fwRevision[slot.index()].String()
	} else {
		return ""
	}
}

// Active returns the slot number and the firmware revision that is currently activated.
//goland:noinspection GoExportedFuncWithUnexportedType
func (f *FirmwareSlotInfo) Active() (fwSlotNo, string) {
	active := f.ActiveFwInfo.Active()

	return active, f.Slot(active)
}

// Next returns the firmware slot number and firmware revision that is going to be activated at the
// next Controller Level Reset command.
//goland:noinspection GoExportedFuncWithUnexportedType
func (f *FirmwareSlotInfo) Next() (fwSlotNo, string) {
	next := f.ActiveFwInfo.Next()

	return next, f.Slot(next)
}

// ParseFirmwareSlotInfo parses the firmware activation information from raw data. If the size of
// raw data is under 512B (size of FirmwareSlotInfo), this function raises an error.
func ParseFirmwareSlotInfo(raw []byte) (*FirmwareSlotInfo, error) {
	if len(raw) < int(unsafe.Sizeof(FirmwareSlotInfo{})) {
		return nil, fmt.Errorf("unexpected firmware slot information data size: %d", len(raw))
	}

	// make clone structure with assignable member
	fw := struct {
		AFI activeFwInfo
		_   [7]byte
		FRS [7]fwRev
		_   [448]byte
	}{}

	if err := binary.Read(bytes.NewReader(raw), utils.SystemEndian, &fw); err == nil {
		return &FirmwareSlotInfo{
			ActiveFwInfo: fw.AFI,
			fwRevision:   fw.FRS,
		}, nil
	} else {
		return nil, err
	}
}

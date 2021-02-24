package getlog

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/sungup/go-nvme-ioctl/pkg/nvme"
	"github.com/sungup/go-nvme-ioctl/pkg/nvme/types"
	"github.com/sungup/go-nvme-ioctl/pkg/utils"
	"os"
)

// ----------------------------------- //
// LID 02h: GetSMART / Health Information //
// ----------------------------------- //

type SMART struct {
	CriticalWarning         uint8        // [00]
	CompositeTemperature    types.Kelvin // [01:02]
	AvailableSpare          uint8        // [03]
	AvailableSpareThreshold uint8        // [04]
	PercentageUsed          uint8        // [05]
	EnduranceGrpCritWarn    uint8        // [06]
	_                       [25]byte     // reserved

	DataUnitsRead       types.Uint128
	DataUnitsWritten    types.Uint128
	HostReadCommands    types.Uint128
	HostWrittenCommands types.Uint128

	ControllerBusyTime types.Uint128
	PowerCycles        types.Uint128
	PowerOnHours       types.Uint128
	UnsafeShutdowns    types.Uint128
	IntegrityErrors    types.Uint128
	ErrorLogEntries    types.Uint128

	WarningCompositeTempTime     uint32
	CriticalCompositeTempTime    uint32
	TemperatureSensor            [8]types.Kelvin
	ThermalMgmtTemp1TransCount   uint32
	ThermalMgmtTemp2TransCount   uint32
	TotalTimeForThermalMgmtTemp1 uint32
	TotalTimeForThermalMgmtTemp2 uint32

	_ [280]byte // reserved
}

// GetSMART will retrieve GetSMART data from NVMe device.
func GetSMART(file *os.File, v interface{}) error {
	if cmd, err := newGetLogCmd(0, 0, logPageSMART, 0, 0, v); err != nil {
		return err
	} else {
		return nvme.IOCtlAdminCmd(file, &cmd.AdminCmd)
	}
}

// ParseSMART creates an GetSMART object to get the device health information.
func ParseSMART(raw []byte) (*SMART, error) {
	if len(raw) != getLogSMARTSz {
		return nil, fmt.Errorf("unexpected GetSMART raw data size: %d", len(raw))
	}

	s := SMART{}

	if err := binary.Read(bytes.NewReader(raw), utils.SystemEndian, &s); err == nil {
		return &s, nil
	} else {
		return nil, err
	}
}

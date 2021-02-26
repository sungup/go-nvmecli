package getlog

import (
	"fmt"
	"github.com/sungup/go-nvmecli/pkg/nvme"
	"github.com/sungup/go-nvmecli/pkg/nvme/identify"
	"math"
	"os"
	"unsafe"
)

// -------------------------- //
// LID 01h: Error Information //
// -------------------------- //

func getELPE(file *os.File) (uint32, error) {
	idCtrl := identify.CtrlIdentify{}
	err := identify.GetCtrlIdentify(file, &idCtrl)

	return uint32(idCtrl.ELPE.Uint()), err
}

// GetErrorInformation will retrieve all NVMe error log entries from NVMe device
func GetErrorInformation(file *os.File, latest uint32) ([]errorEntry, error) {
	const (
		errEntrySz     = uint32(unsafe.Sizeof(errorEntry{}))
		unitErrInfoCnt = 4096 / errEntrySz
	)

	// 1. get identify from the identify.GetCtrlIdentify
	maxEntry, err := getELPE(file)
	if err != nil || maxEntry == 0 {
		if err == nil {
			err = fmt.Errorf("ELPE (Error Log Page Entries) is 0")
		}

		return nil, fmt.Errorf("getting ELPE failed or unsupported: %v", err)
	} else if latest < maxEntry {
		maxEntry = latest
	}

	// 2. make Buffer, Return-Data and Get-Log command.
	fetchCnt := unitErrInfoCnt

	errors := make([]errorEntry, 0, maxEntry)
	buffer := make([]errorEntry, unitErrInfoCnt)
	cmd, _ := newGetLogCmd(0, 0, logPageErrorInfo, 0, 0, buffer)

	// 3. get entries each unitErrInfoCnt error logs.
	for index := uint32(0); index < maxEntry; index += unitErrInfoCnt {
		// 3-1. update offset and dwords of getLogCmd
		if maxEntry < index+unitErrInfoCnt {
			fetchCnt = maxEntry - index
			cmd.SetDWords(fetchCnt * errEntrySz >> 2)
		}
		cmd.SetOffset(uint64(index * errEntrySz))

		// 3-2. send ioctl to device
		if err := nvme.IOCtlAdminCmd(file, &cmd.AdminCmd); err != nil {
			return nil, err
		}

		errors = append(errors, buffer[:fetchCnt]...)
	}

	// 4. returns the valid error logs has ErrorCount > 1.
	count := 0
	for ; errors[count].ErrorCount > 0; count++ {
	}

	return errors[:count], nil
}

type errStatField uint16

// StatusField returns bit[15:1] information for the command that completed.
func (e errStatField) StatusField() uint16 {
	return uint16(e) >> 1
}

// PhaseTag returns bit[0] indicate the Phase Tag posted for the command.
func (e errStatField) PhaseTag() uint16 {
	return uint16(e) & 0x01
}

type paramErrLoc uint16

// LocationBit returns the bit of the command parameters that error is associated with.
func (p paramErrLoc) LocationBit() uint8 {
	return uint8(p>>8) & 0b111
}

// LocationByte returns the byte of the command parameters that error is associated with.
func (p paramErrLoc) LocationByte() uint8 {
	return uint8(math.MaxUint8 & uint16(p))
}

type trType uint8

const (
	NoTransport        = trType(0)
	RDMATransport      = trType(1)
	FibreChTransport   = trType(2)
	TCPTransport       = trType(3)
	IntraHostTransport = trType(0xFE)
)

type errorEntry struct {
	ErrorCount uint64

	SqID      uint16
	CommandID uint16

	StatusField      errStatField
	ParamErrLocation paramErrLoc

	LBA       uint64
	Namespace uint32

	VendorSpecificInfo uint8
	TransportType      trType

	_ [2]byte

	CommandSpecific   uint64
	TransportSpecific uint16

	_ [22]byte
}

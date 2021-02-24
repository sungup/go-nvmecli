package getlog

import (
	"fmt"
	"math"
	"os"
)

// -------------------------- //
// LID 01h: Error Information //
// -------------------------- //

// GetErrorInformation will retrieve all NVMe error log entries from NVMe device
func GetErrorInformation(file *os.File, latest uint16) ([]errorEntry, error) {
	// TODO 1. get identify from the identify.GetCtrlIdentify

	// TODO 2. make Buffer, Return-Data and Get-Log command.

	// TODO 3. get entries each 64 error logs.

	// TODO 4. returns the valid error logs has ErrorCount > 1.

	// TODO implementing here
	return nil, fmt.Errorf("unimplemented function")
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
	return uint8(p>>8) & 0b11
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

type trTypeInfo uint16

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

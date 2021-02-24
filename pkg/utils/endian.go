package utils

import (
	"encoding/binary"
	"unsafe"
)

var (
	SystemEndian  binary.ByteOrder
	ReverseEndian binary.ByteOrder
)

func init() {
	buf := [2]byte{}
	*(*uint16)(unsafe.Pointer(&buf[0])) = uint16(0xABCD)

	switch buf {
	case [2]byte{0xCD, 0xAB}:
		SystemEndian = binary.LittleEndian
		ReverseEndian = binary.BigEndian
	case [2]byte{0xAB, 0xCD}:
		SystemEndian = binary.BigEndian
		ReverseEndian = binary.LittleEndian
	default:
		panic("could not determine system endian")
	}
}

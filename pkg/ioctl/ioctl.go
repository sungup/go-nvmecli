package ioctl

import (
	"fmt"
	"os"
	"syscall"
)

const (
	nrBits   = 8
	typeBits = 8

	sizeBits = 14
	dirBits  = 2

	nrMask   = (1 << nrBits) - 1
	typeMask = (1 << typeBits) - 1
	sizeMask = (1 << sizeBits) - 1
	dirMask  = (1 << dirBits) - 1

	NrShift   = 0
	TypeShift = NrShift + nrBits
	SizeShift = TypeShift + typeBits
	DirShift  = SizeShift + sizeBits

	iocNone  = uint64(0)
	iocWrite = uint64(1)
	iocRead  = uint64(2)
	iocWrRd  = iocWrite | iocRead

	IOCNone  = iocNone << DirShift
	IOCIn    = iocWrite << DirShift
	IOCOut   = iocRead << DirShift
	IOCInOut = iocWrRd << DirShift
)

func ioc(dir, typ, nr, size uint64) uint64 {
	return (dir << DirShift) |
		(typ << TypeShift) |
		(nr << NrShift) |
		(size << SizeShift)
}

// create ioctl numbers
func IO(typ, nr uint64) uint64         { return ioc(iocNone, typ, nr, 0) }
func IOR(typ, nr, size uint64) uint64  { return ioc(iocRead, typ, nr, size) }
func IOW(typ, nr, size uint64) uint64  { return ioc(iocWrite, typ, nr, size) }
func IOWR(typ, nr, size uint64) uint64 { return ioc(iocRead|iocWrite, typ, nr, size) }

// decode ioctl numbers
func iocDir(nr uint64) uint64  { return nr >> DirShift & dirMask }
func iocType(nr uint64) uint64 { return nr >> TypeShift & typeMask }
func iocNr(nr uint64) uint64   { return nr >> NrShift & nrMask }
func iocSize(nr uint64) uint64 { return nr >> SizeShift & sizeMask }

// submit ioctl command
func Submit(f *os.File, request, data uintptr) error {
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), request, data); errno != 0 {
		return os.NewSyscallError("ioctl", fmt.Errorf("%d", errno))
	}

	return nil
}

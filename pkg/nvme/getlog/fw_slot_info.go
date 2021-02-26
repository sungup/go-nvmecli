package getlog

type fwRev [8]byte

func (f fwRev) String() string {
	return string(f[:])
}

type activeFwInfo byte

func (f activeFwInfo) Next() int {
	return int(f) >> 4 & 0b111
}

func (f activeFwInfo) Active() int {
	return int(f) & 0b111
}

type FirmwareSlotInfo struct {
	ActiveFwInfo activeFwInfo
	_            [7]byte // reserved

	FwRevision [7]fwRev
	_          [448]byte // reserved
}

func (f *FirmwareSlotInfo) Active() (int, string) {
	active := f.ActiveFwInfo.Active()

	return active, f.FwRevision[active-1].String()
}

func (f *FirmwareSlotInfo) Next() (int, string) {
	next := f.ActiveFwInfo.Next()

	if next != 0 {
		return next, f.FwRevision[next-1].String()
	} else {
		return next, ""
	}
}

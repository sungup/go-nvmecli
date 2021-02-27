package getlog

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	tcFwList = []struct {
		tested   fwRev
		expected string
	}{
		{tested: fwRev{'F', 'w', 'S', 'l', 'o', 't', '0', '1'}, expected: "FwSlot01"},
		{tested: fwRev{'F', 'w', 'S', 'l', 'o', 't', '0', '2'}, expected: "FwSlot02"},
		{tested: fwRev{'F', 'w', 'S', 'l', 'o', 't', '0', '3'}, expected: "FwSlot03"},
		{tested: fwRev{'F', 'w', 'S', 'l', 'o', 't', '0', '4'}, expected: "FwSlot04"},
		{tested: fwRev{'F', 'w', 'S', 'l', 'o', 't', '0', '5'}, expected: "FwSlot05"},
		{tested: fwRev{'F', 'w', 'S', 'l', 'o', 't', '0', '6'}, expected: "FwSlot06"},
		{tested: fwRev{'F', 'w', 'S', 'l', 'o', 't', '0', '7'}, expected: "FwSlot07"},
	}
)

func TestFwSlotNo_index(t *testing.T) {
	a := assert.New(t)

	for i := 1; i < 8; i++ {
		tested := fwSlotNo(i)
		expected := i - 1

		a.Equal(expected, tested.index())
	}
}

func TestFwRev_String(t *testing.T) {
	a := assert.New(t)

	type tcType struct {
		tested   fwRev
		expected string
	}

	tcList := []tcType{
		{tested: fwRev{'H', 'E', 'L', 'L', 'O', ' ', ' ', ' '}, expected: "HELLO"},
		{tested: fwRev{'H', 'E', 'L', 'L', 'O', ' ', 'W', 'O'}, expected: "HELLO WO"},
		{tested: fwRev{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '}, expected: ""},
		{tested: fwRev{'H', 'E', 'L', 'L', 'O'}, expected: "HELLO"},
		{tested: fwRev{}, expected: ""},
	}

	for _, tc := range tcList {
		a.Equal(tc.expected, tc.tested.String())
	}
}

func TestActiveFwInfo_Active(t *testing.T) {
	a := assert.New(t)

	expectedList := []fwSlotNo{FRS1, FRS2, FRS3, FRS4, FRS5, FRS6, FRS7}

	for active := 1; active < 8; active++ {
		for next := 0; next < 8; next++ {
			afi := activeFwInfo(next<<4 | active)

			a.Contains(expectedList, afi.Active())
		}
	}
}

func TestActiveFwInfo_Next(t *testing.T) {
	a := assert.New(t)

	expectedList := []fwSlotNo{NoFw, FRS1, FRS2, FRS3, FRS4, FRS5, FRS6, FRS7}

	for active := 1; active < 8; active++ {
		for next := 0; next < 8; next++ {
			afi := activeFwInfo(next<<4 | active)

			a.Contains(expectedList, afi.Next())
		}
	}
}

func TestFirmwareSlotInfo_Slot(t *testing.T) {
	a := assert.New(t)

	tested := FirmwareSlotInfo{}

	for i, tc := range tcFwList {
		tested.fwRevision[i] = tc.tested
	}

	for _, slot := range []fwSlotNo{FRS1, FRS2, FRS3, FRS4, FRS5, FRS6, FRS7} {
		a.Equal(tcFwList[slot.index()].expected, tested.Slot(slot))
	}
}

func TestFirmwareSlotInfo_Active(t *testing.T) {
	a := assert.New(t)

	tested := FirmwareSlotInfo{}
	expectedList := []fwSlotNo{FRS1, FRS2, FRS3, FRS4, FRS5, FRS6, FRS7}

	for i, tc := range tcFwList {
		tested.fwRevision[i] = tc.tested
	}

	for active := 1; active < 8; active++ {
		for next := 0; next < 8; next++ {
			tested.ActiveFwInfo = activeFwInfo(next<<4 | active)

			testedSlot, testedFw := tested.Active()

			a.Contains(expectedList, testedSlot)
			a.Equal(tcFwList[testedSlot.index()].expected, testedFw)
		}
	}
}

func TestFirmwareSlotInfo_Next(t *testing.T) {

	a := assert.New(t)

	tested := FirmwareSlotInfo{}
	expectedList := []fwSlotNo{NoFw, FRS1, FRS2, FRS3, FRS4, FRS5, FRS6, FRS7}

	for i, tc := range tcFwList {
		tested.fwRevision[i] = tc.tested
	}

	for active := 1; active < 8; active++ {
		for next := 0; next < 8; next++ {
			tested.ActiveFwInfo = activeFwInfo(next<<4 | active)

			testedSlot, testedFw := tested.Next()

			a.Contains(expectedList, testedSlot)

			if testedSlot != NoFw {
				a.Equal(tcFwList[testedSlot.index()].expected, testedFw)
			} else {
				a.Empty(testedFw)
			}
		}
	}
}

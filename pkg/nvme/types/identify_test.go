package types

import (
	"bytes"
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"github.com/sungup/go-nvmecli/pkg/utils"
	"testing"
)

func TestVID_String(t *testing.T) {
	a := assert.New(t)

	expected := "ABCDh"
	tested := VID(0xabcd)

	a.Equal(expected, tested.String())
}

func TestSSVID_String(t *testing.T) {
	a := assert.New(t)

	expected := "1234h"
	tested := SSVID(0x1234)

	a.Equal(expected, tested.String())
}

func TestSN_String(t *testing.T) {
	a := assert.New(t)

	expected := "HELLO SN"
	tested := SN{'H', 'E', 'L', 'L', 'O', ' ', 'S', 'N', ' ', ' '}

	a.Equal(expected, tested.String())
}

func TestMN_String(t *testing.T) {
	a := assert.New(t)

	expected := "Hello MN"
	tested := MN{'H', 'e', 'l', 'l', 'o', ' ', 'M', 'N', ' ', ' '}

	a.Equal(expected, tested.String())
}

func TestIEEE_String(t *testing.T) {
	// Reference: 7.10.3 IEEE OUI Identifier (IEEE); p303, NVM-Express-1.4a
	a := assert.New(t)

	expected := "ABCDEFh"
	tested := IEEE{0xef, 0xcd, 0xab}

	a.Equal(expected, tested.String())
}

func TestEUI64_OUI(t *testing.T) {
	// Reference: 7.10.4 IEEE Extended Unique Identifier (EUI64); p303, NVM-Express-1.4a
	a := assert.New(t)

	expected := "ABCDEFh"
	tested := EUI64{
		Oui: [3]Hex8{0xab, 0xcd, 0xef},
		Ext: [5]Hex8{0x01, 0x23, 0x45, 0x67, 0x89},
	}

	a.Equal(expected, tested.OUI())
}

func TestEUI64_Extension(t *testing.T) {
	// Reference: 7.10.4 IEEE Extended Unique Identifier (EUI64); p303-304, NVM-Express-1.4a
	a := assert.New(t)

	expected := "0123456789h"
	tested := EUI64{
		Oui: [3]Hex8{0xab, 0xcd, 0xef},
		Ext: [5]Hex8{0x01, 0x23, 0x45, 0x67, 0x89},
	}

	a.Equal(expected, tested.Extension())
}

func TestEUI64_String(t *testing.T) {
	a := assert.New(t)

	expected := "EUI64{OUI: ABCDEFh, Extension: 0123456789h}"
	tested := EUI64{
		Oui: [3]Hex8{0xab, 0xcd, 0xef},
		Ext: [5]Hex8{0x01, 0x23, 0x45, 0x67, 0x89},
	}

	a.Equal(expected, tested.String())
}

func TestNGUIDFormatCheck(t *testing.T) {
	// Reference: 7.10.5 Namespace Globally Unique Identifier (NGUID); p304-305, NVM-Express-1.4a
	a := assert.New(t)

	expectedOUI := "ABCDEFh"
	expectedExt := "0123456789h"

	tested := NGUID{
		Vendor: [8]Hex8{0xFE, 0xDC, 0xBA, 0x98, 0x76, 0x54, 0x32, 0x10},
		EUI64: EUI64{
			Oui: [3]Hex8{0xab, 0xcd, 0xef},
			Ext: [5]Hex8{0x01, 0x23, 0x45, 0x67, 0x89},
		},
	}

	// check other function working fine
	a.Equal(expectedOUI, tested.OUI())
	a.Equal(expectedExt, tested.Extension())

	// check byte orders start from Vendor
	buffer := new(bytes.Buffer)
	a.NoError(binary.Write(buffer, utils.SystemEndian, tested))

	for i, tc := range tested.Vendor {
		a.Equal(buffer.Bytes()[i], byte(tc))
	}

	for i, tc := range tested.Oui {
		a.Equal(buffer.Bytes()[8+i], byte(tc))
	}

	for i, tc := range tested.Ext {
		a.Equal(buffer.Bytes()[8+3+i], byte(tc))
	}
}

func TestNGUID_VendorID(t *testing.T) {
	// Reference: 7.10.5 Namespace Globally Unique Identifier (NGUID); p304-305, NVM-Express-1.4a
	a := assert.New(t)

	expectedVdr := "FEDCBA9876543210h"

	tested := NGUID{
		Vendor: [8]Hex8{0xFE, 0xDC, 0xBA, 0x98, 0x76, 0x54, 0x32, 0x10},
		EUI64: EUI64{
			Oui: [3]Hex8{0xab, 0xcd, 0xef},
			Ext: [5]Hex8{0x01, 0x23, 0x45, 0x67, 0x89},
		},
	}

	a.Equal(expectedVdr, tested.VendorID())
}

func TestNGUID_String(t *testing.T) {
	// Reference: 7.10.5 Namespace Globally Unique Identifier (NGUID); p304-305, NVM-Express-1.4a
	a := assert.New(t)

	expectedVdr := "NGUID{Vendor: FEDCBA9876543210h, OUI: ABCDEFh, Extension: 0123456789h}"

	tested := NGUID{
		Vendor: [8]Hex8{0xFE, 0xDC, 0xBA, 0x98, 0x76, 0x54, 0x32, 0x10},
		EUI64: EUI64{
			Oui: [3]Hex8{0xab, 0xcd, 0xef},
			Ext: [5]Hex8{0x01, 0x23, 0x45, 0x67, 0x89},
		},
	}

	a.Equal(expectedVdr, tested.String())
}

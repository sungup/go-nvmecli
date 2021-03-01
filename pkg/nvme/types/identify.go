package types

import (
	"bytes"
	"fmt"
)

const (
	trimSet = " \x00"
)

// VID is the PCI Vendor ID defined in the Identify Controller Data Structure. This value is
// assigned by the PCI SIG. VID is in little endian format.
// Reference: 7.10.1 PCI Vendor ID (VID) and PCI Subsystem Vendor ID (SSVID); p302, NVM-Express-1.4a
type VID uint16

// String convert VID to string
func (i VID) String() string {
	return fmt.Sprintf("%04Xh", uint16(i))
}

// SSVID is the PCI Subsystem Vendor ID defined in the Identify Controller Data Structure. This
// value is assigned by the PCI SIG. SSVID is in little endian format.
// Reference: 7.10.1 PCI Vendor ID (VID) and PCI Subsystem Vendor ID (SSVID); p302, NVM-Express-1.4a
type SSVID uint16

// String convert VID to string
func (i SSVID) String() string {
	return fmt.Sprintf("%04Xh", uint16(i))
}

// SN is the Serial Number defined in the Identify Controller Data Structure. The values are ASCII
// string and in big endian format.
// Reference: 7.10.2. Serial Number(SN) and Model Number (MN); p303, NVM-Express-1.4a
type SN [20]byte

// String convert SN to trimmed string
func (i SN) String() string {
	return string(bytes.Trim(i[:], trimSet))
}

// MN is the Model Number defined in the Identify Controller Data Structure. The values are ASCII
// string and in big endian format.
// Reference: 7.10.2. Serial Number(SN) and Model Number (MN); p303, NVM-Express-1.4a
type MN [40]byte

// String convert MN to trimmed string
func (i MN) String() string {
	return string(bytes.Trim(i[:], trimSet))
}

// IEEE is the IEEE OUI Identifier defined in the Identify Controller data structure. IEEE OUI is
// in little endian format.
// Reference: 7.10.3 IEEE OUI Identifier (IEEE); p303, NVM-Express-1.4a
type IEEE [3]Hex8

// String convert OUI identifier to string.
func (i IEEE) String() string {
	// Big endian but only 3 byte data structure
	return fmt.Sprintf("%02X%02X%02Xh", i[2], i[1], i[0])
}

// EUI64 is the IEEE Extended Unique Identifier defined in the Identify Namespace data structure.
// Different from IEEE OUI, EUI64 is in big endian format. (IEEE OUI is in little endian format)
// Reference: 7.10.4 IEEE Extended Unique Identifier (EUI64); p303-304, NVM-Express-1.4a
//
// Don't access Oui and Ext directly because these member have been exported to access from binary
// decoder.
type EUI64 struct {
	Oui [3]Hex8
	Ext [5]Hex8
}

// OUI convert the oui identifier in EUI 64 to string.
func (i EUI64) OUI() string {
	return fmt.Sprintf("%02X%02X%02Xh", i.Oui[0], i.Oui[1], i.Oui[2])
}

// Extension convert the extension identifier in EUI64 to string
func (i EUI64) Extension() string {
	// Little endian
	return fmt.Sprintf("%02X%02X%02X%02X%02Xh", i.Ext[0], i.Ext[1], i.Ext[2], i.Ext[3], i.Ext[4])
}

// String convert EUI64 to string
func (i EUI64) String() string {
	return fmt.Sprintf("EUI64{OUI: %s, Extension: %s}", i.OUI(), i.Extension())
}

// NGUID is the Namespace Globally Unique Identifier defined in the Identify Namespace data
// structure. It is composed of an EUI64 (OUI + extension) and a vendor specific extension
// identifier. Also, it is in big endian format.
// Reference: 7.10.5 Namespace Globally Unique Identifier (NGUID); p304-305, NVM-Express-1.4a
//
// Don't access Vendor directly because this member has been exported to access from binary decoder.
type NGUID struct {
	Vendor [8]Hex8
	EUI64
}

// VendorID convert the vendor specific identifier in NGUID to string
func (i NGUID) VendorID() string {
	// Little endian
	return fmt.Sprintf(
		"%02X%02X%02X%02X%02X%02X%02X%02Xh",
		i.Vendor[0], i.Vendor[1], i.Vendor[2], i.Vendor[3],
		i.Vendor[4], i.Vendor[5], i.Vendor[6], i.Vendor[7],
	)
}

// String convert NGUID to string
func (i NGUID) String() string {
	return fmt.Sprintf("NGUID{Vendor: %s, OUI: %s, Extension: %s}", i.VendorID(), i.OUI(), i.Extension())
}

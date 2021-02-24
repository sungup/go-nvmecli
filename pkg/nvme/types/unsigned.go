package types

// All data type in this package should use byte array because avoiding the byte alignment bugs.

type Hex8 byte
type Hex16 [2]byte
type Hex32 [4]byte
type Hex64 [8]byte
type Hex128 [2]Hex64

type Uint8 byte
type Uint16 [2]byte
type Uint32 [4]byte
type Uint64 [8]byte
type Uint128 [2]Uint64

func (u Uint8) Uint() uint64 { return uint64(u) }

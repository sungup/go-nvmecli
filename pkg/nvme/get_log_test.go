package nvme

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestNewGetLogCmd(t *testing.T) {
	var (
		a = assert.New(t)

		err    error
		data   interface{}
		tested *AdminCmd
	)

	const (
		expectedNSId   = 1
		expectedDWords = uint32(0xAB<<16 | 0xCD)
		expectedOffset = uint64(0xAB<<32 | 0xCD)
		expectedLID    = uint16(0xCA)
		expectedLSP    = uint16(0xE)
		expectedLSI    = uint16(0xA)

		expectedTestSize = 3072
	)

	// 1. check invalid dptr data
	data = uint64(0)
	tested, err = newGetLogCmd(expectedNSId, expectedDWords, expectedOffset, expectedLID, expectedLSP, expectedLSI, data)
	a.Error(err)

	// 2. get data using byte array
	data = make([]uint8, expectedTestSize)
	tested, err = newGetLogCmd(expectedNSId, expectedDWords, expectedOffset, expectedLID, expectedLSP, expectedLSI, data)
	a.NotNil(tested)
	a.NoError(err)
	a.Equal(AdminGetLogPage, tested.OpCode)
	a.Equal(uint32(expectedNSId), tested.NSId)
	a.Equal(uint32(expectedTestSize), tested.DataLength)
	a.Equal(expectedDWords, (tested.CDW10>>16)|(tested.CDW11<<16))
	a.Equal(expectedOffset, (uint64(tested.CDW12)<<32)|uint64(tested.CDW13))
	a.Equal(expectedLID, uint16(math.MaxUint8&tested.CDW10))
	a.Equal(expectedLSP, uint16((tested.CDW10<<16)>>24))
	a.Equal(expectedLSI, uint16(tested.CDW11>>16))

	data = &struct {
		raw [expectedTestSize + 10]byte
	}{}
	tested, err = newGetLogCmd(expectedNSId, expectedDWords, expectedOffset, expectedLID, expectedLSP, expectedLSI, data)
	a.NotNil(tested)
	a.NoError(err)
	a.Equal(AdminGetLogPage, tested.OpCode)
	a.Equal(uint32(expectedNSId), tested.NSId)
	a.Equal(uint32(expectedTestSize+10), tested.DataLength)
	a.Equal(expectedDWords, (tested.CDW10>>16)|(tested.CDW11<<16))
	a.Equal(expectedOffset, (uint64(tested.CDW12)<<32)|uint64(tested.CDW13))
	a.Equal(expectedLID, uint16(math.MaxUint8&tested.CDW10))
	a.Equal(expectedLSP, uint16((tested.CDW10<<16)>>24))
	a.Equal(expectedLSI, uint16(tested.CDW11>>16))
}

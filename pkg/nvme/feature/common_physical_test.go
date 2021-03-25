// +build with_phys_device

package feature

import (
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetFeature(t *testing.T) {
	a := assert.New(t)

	const (
		expectedFID = FIDTimeStamp
		expectedSEL = SELCurrent
	)

	dev, _ := os.Open(targetDevice)
	buffer := [8]byte{}

	// get feature of timestamp will return not 0 value
	a.NoError(GetFeature(dev, expectedNSId, expectedFID, 0, expectedSEL, buffer[:]))
	a.NotZero(binary.BigEndian.Uint64(buffer[:]))
}

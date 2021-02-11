// +build with_phys_device

package nvme

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetLogSMART(t *testing.T) {
	a := assert.New(t)

	dev, err := os.Open(targetDevice)
	a.NoError(err)

	buffer, err := GetLogSMART(dev)
	a.NoError(err)
	a.NotNil(buffer)
	a.Len(buffer, getLogSMARTSz)
}

func TestParseSMART(t *testing.T) {
	a := assert.New(t)

	dev, err := os.Open(targetDevice)

	a.NoError(err)

	buffer, _ := GetLogSMART(dev)
	smart, err := ParseSMART(buffer)
	a.NoError(err)

	a.NotZero(smart.AvailableSpareThreshold)
	t.Log(smart.ErrorLogEntries)
}

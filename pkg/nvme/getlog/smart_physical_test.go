// +build with_phys_device

package getlog

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetSMART(t *testing.T) {
	a := assert.New(t)

	dev, _ := os.Open(targetDevice)

	// 1. byte array buffer
	buffer := make([]byte, getLogSMARTSz)
	a.NoError(GetSMART(dev, buffer))

	// 2. struct data
	tested := SMART{}
	a.NoError(GetSMART(dev, &tested))

	a.NotZero(tested.AvailableSpareThreshold)
}

func TestParseSMART(t *testing.T) {
	a := assert.New(t)

	dev, _ := os.Open(targetDevice)
	buffer := make([]byte, getLogSMARTSz)

	_ = GetSMART(dev, buffer)

	// check invalid size error
	smart, err := ParseSMART(buffer[1:])
	a.Error(err)
	a.Nil(smart)

	// check normal parsing
	smart, err = ParseSMART(buffer)
	a.NoError(err)
	a.NotNil(smart)

	a.NotZero(smart.AvailableSpareThreshold)
}

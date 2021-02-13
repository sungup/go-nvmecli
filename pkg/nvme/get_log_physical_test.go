// +build with_phys_device

package nvme

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetLogSMART(t *testing.T) {
	a := assert.New(t)

	dev, _ := os.Open(targetDevice)

	buffer, err := GetLogSMART(dev)
	a.NoError(err)
	a.NotNil(buffer)
	a.Len(buffer, getLogSMARTSz)
}

func TestParseSMART(t *testing.T) {
	a := assert.New(t)

	dev, _ := os.Open(targetDevice)

	buffer, _ := GetLogSMART(dev)
	smart, err := ParseSMART(buffer)
	a.NoError(err)

	a.NotZero(smart.AvailableSpareThreshold)
}

func TestTelemetryHeader(t *testing.T) {
	a := assert.New(t)

	dev, _ := os.Open(targetDevice)

	buffer, err := ioctlAdminCmd(
		dev,
		getLogTelemetryHeaderSz,
		func() *AdminCmd {
			cmd := newGetLogCmd(0, getLogTelemetryHeaderSz>>2, 0, logPageTelemetryHost, 0x01, 0)
			return &cmd.AdminCmd
		},
	)

	a.NotEmpty(buffer)
	a.NoError(err)

	// make blocking telemetry checking interface because my Seagate NVMe device doesn't
	// support telemetry feature.
	/*
		telemetry, err := ParseTelemetryHeader(buffer)
		a.NoError(err)
		a.Equal(logPageTelemetryHost, uint16(telemetry.Identifier))

		t.Log(string(telemetry.ReasonIdentifier[:]))
		t.Log(telemetry.DataAreaLastBlock)
	*/
}

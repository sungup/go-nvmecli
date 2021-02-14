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

	// 1. byte array buffer
	buffer := make([]byte, getLogSMARTSz)
	a.NoError(GetLogSMART(dev, buffer))

	// 2. struct data
	tested := smart{}
	a.NoError(GetLogSMART(dev, &tested))

	a.NotZero(tested.AvailableSpareThreshold)
}

func TestParseSMART(t *testing.T) {
	a := assert.New(t)

	dev, _ := os.Open(targetDevice)
	buffer := make([]byte, getLogSMARTSz)

	_ = GetLogSMART(dev, buffer)

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

func TestGetLogTelemetry(t *testing.T) {
	/* nothing to do, because ctrl/host init function will call it with their specific param */
}

func TestGetLogTelemetryHostInit(t *testing.T) {
	// TODO re-verify this test code using the telemetry support NVME device
	// My Seagate NVMe device doesn't support telemetry, but there is no error to call the get-log
	// telemetry operation. So, total length of telemetry data is 512B for the telemetry header.
	/*
		a := assert.New(t)

		dev, _ := os.Open(targetDevice)

		buffer, err := GetLogTelemetryHostInit(dev, DataBlock3, true)
		a.NoError(err)
		a.NotNil(buffer)
		a.Len(buffer, 512)
	*/
}

func TestGetLogTelemetryCtrlInit(t *testing.T) {
	// TODO re-verify this test code using the telemetry support NVME device
	// My Seagate NVMe device doesn't support telemetry, but there is no error to call the get-log
	// telemetry operation. So, total length of telemetry data is 512B for the telemetry header.
	/*
		a := assert.New(t)

		dev, _ := os.Open(targetDevice)

		buffer, err := GetLogTelemetryCtrlInit(dev, DataBlock3)
		a.NoError(err)
		a.NotNil(buffer)
		a.Len(buffer, 512)
	*/
}

func TestParseTelemetryHeader(t *testing.T) {
	// TODO re-verify this test code using the telemetry support NVME device
	// My Seagate NVMe device doesn't support telemetry, but there is no error to call the get-log
	// telemetry operation. So, total length of telemetry data is 512B for the telemetry header.
	/*
		a := assert.New(t)

		dev, _ := os.Open(targetDevice)

		buffer, _ := GetLogTelemetryCtrlInit(dev, DataBlock3)

		// check invalid size error
		tested, err := ParseTelemetryHeader(buffer[1:])
		a.Error(err)
		a.Nil(tested)

		// check normal parsing
		tested, err = ParseTelemetryHeader(buffer)
		a.NoError(err)
		a.NotNil(tested)

			a.Equal(logPageTelemetryHost, uint16(tested.Identifier))

			t.Log(string(tested.ReasonIdentifier[:]))
			t.Log(tested.DataAreaLastBlock)
	*/
}

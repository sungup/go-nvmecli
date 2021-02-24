// +build with_phys_device

package getlog

import "testing"

func TestGetLogTelemetry(t *testing.T) {
	/* nothing to do, because ctrl/host init function will call it with their specific param */
}

func TestGetTelemetryHostInit(t *testing.T) {
	// TODO re-verify this test code using the Telemetry support NVME device
	// My Seagate NVMe device doesn't support Telemetry, but there is no error to call the get-log
	// Telemetry operation. So, total length of Telemetry data is 512B for the Telemetry header.
	/*
		a := assert.New(t)

		dev, _ := os.Open(targetDevice)

		buffer, err := GetTelemetryHostInit(dev, DataBlock3, true)
		a.NoError(err)
		a.NotNil(buffer)
		a.Len(buffer, 512)
	*/
}

func TestGetTelemetryCtrlInit(t *testing.T) {
	// TODO re-verify this test code using the Telemetry support NVME device
	// My Seagate NVMe device doesn't support Telemetry, but there is no error to call the get-log
	// Telemetry operation. So, total length of Telemetry data is 512B for the Telemetry header.
	/*
		a := assert.New(t)

		dev, _ := os.Open(targetDevice)

		buffer, err := GetTelemetryCtrlInit(dev, DataBlock3)
		a.NoError(err)
		a.NotNil(buffer)
		a.Len(buffer, 512)
	*/
}

func TestParseTelemetryHeader(t *testing.T) {
	// TODO re-verify this test code using the Telemetry support NVME device
	// My Seagate NVMe device doesn't support Telemetry, but there is no error to call the get-log
	// Telemetry operation. So, total length of Telemetry data is 512B for the Telemetry header.
	/*
		a := assert.New(t)

		dev, _ := os.Open(targetDevice)

		buffer, _ := GetTelemetryCtrlInit(dev, DataBlock3)

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

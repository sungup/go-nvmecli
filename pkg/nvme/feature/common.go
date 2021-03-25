package feature

import (
	"github.com/sungup/go-nvmecli/pkg/nvme"
	"os"
)

type sel uint32

const (
	SELCurrent    = sel(0x00) << 8
	SELDefault    = sel(0x01) << 8
	SELSaved      = sel(0x02) << 8
	SELSupportCap = sel(0x03) << 8

	FIDArbitration                 = uint8(0x001)
	FIDPowerManagement             = uint8(0x002)
	FIDLBARangeType                = uint8(0x003)
	FIDTemperatureThreshold        = uint8(0x004)
	FIDErrorRecovery               = uint8(0x005)
	FIDVolatileWriteCache          = uint8(0x006)
	FIDNumberOfQueues              = uint8(0x007)
	FIDInterruptCoalescing         = uint8(0x008)
	FIDInterruptVectorConf         = uint8(0x009)
	FIDWriteAtomicityNormal        = uint8(0x00a)
	FIDAsyncEventConf              = uint8(0x00b)
	FIDAutoPowerStateTransition    = uint8(0x00c)
	FIDHostMemoryBuffer            = uint8(0x00d)
	FIDTimeStamp                   = uint8(0x00e)
	FIDKeepAliveTimer              = uint8(0x00f)
	FIDHostCtrlThermalManagement   = uint8(0x010)
	FIDNonOPPowerStateConf         = uint8(0x011)
	FIDReadRecoveryLevelConf       = uint8(0x012)
	FIDPredictableLatModeConf      = uint8(0x013)
	FIDPredictableLatModeWin       = uint8(0x014)
	FIDLBAStatusInfoReportInterval = uint8(0x015)
	FIDHostBehaviorSupport         = uint8(0x016)
	FIDSanitizeConf                = uint8(0x017)
	FIDEnduranceGroupEventConf     = uint8(0x018)

	selFieldUMask = ^(uint32(0x007) << 8)
)

type getFeatureCmd struct {
	nvme.AdminCmd
}

// SEL() change the SEL fields
func (f *getFeatureCmd) SEL(sel sel) {
	f.CDW10 = f.CDW10&selFieldUMask | uint32(sel)
}

// newGetFeatureCmd generate an AdminCmd structure to retrieve the NVMe's feature pages.
func newGetFeatureCmd(nsid uint32, fid uint8, specific uint32, sel sel, v interface{}) (*getFeatureCmd, error) {
	// TODO check common FID doesn't need specific field
	cmd := getFeatureCmd{
		nvme.AdminCmd{
			PassthruCmd: nvme.PassthruCmd{
				OpCode: nvme.AdminGetFeatures,
				NSId:   nsid,
				CDW10:  uint32(sel) | uint32(fid),
				CDW11:  specific,
			},
			TimeoutMSec: 0,
			Result:      0,
		},
	}

	if v != nil {
		if err := cmd.SetData(v); err != nil {
			return nil, err
		}
	}

	return &cmd, nil
}

// GetFeatureCMD retrieve a feature data.
func GetFeature(file *os.File, nsid uint32, fid uint8, specific uint32, sel sel, v interface{}) error {
	if cmd, err := newGetFeatureCmd(nsid, fid, specific, sel, v); err != nil {
		return err
	} else {
		return nvme.IOCtlAdminCmd(file, &cmd.AdminCmd)
	}
}

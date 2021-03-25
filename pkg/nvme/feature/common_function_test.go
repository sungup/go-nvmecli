package feature

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetFeatureCmd_SEL(t *testing.T) {
	a := assert.New(t)

	const (
		expectedSpecific       = 0xEFCDAB89
		expectedFID      uint8 = 0x0CF
		expectedSEL            = SELSupportCap
	)

	v := make([]byte, 4)

	// 1. create same get-feature command
	origin, _ := newGetFeatureCmd(expectedNSId, expectedFID, expectedSpecific, expectedSEL, v)
	tested, _ := newGetFeatureCmd(expectedNSId, expectedFID, expectedSpecific, expectedSEL, v)

	tested.SEL(SELDefault)

	// CDW10 check
	a.NotEqual(origin.CDW10, tested.CDW10)
	a.Equal(origin.CDW10&selFieldUMask, tested.CDW10&selFieldUMask)

	// CDW11/12/13 check
	a.Equal(origin.CDW11, tested.CDW11)
	a.Equal(origin.CDW12, tested.CDW12)
	a.Equal(origin.CDW13, tested.CDW13)
}

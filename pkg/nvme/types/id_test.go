package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIEEE_String(t *testing.T) {
	a := assert.New(t)

	expected := "ABCDEFh"
	tested := IEEE{0xef, 0xcd, 0xab}

	a.Equal(expected, tested.String())
}

package types

import "fmt"

type IEEE [3]Hex8

func (i IEEE) String() string {
	return fmt.Sprintf("%X%X%Xh", i[2], i[1], i[0])
}

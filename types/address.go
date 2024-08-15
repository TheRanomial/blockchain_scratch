package types

import (
	"encoding/hex"
	"fmt"
)

type Address [20]uint8

func (a Address) ToSlice() []byte {

	temp:=make([]byte,20)

	for i:=0;i<20;i++{
		temp[i]=a[i]
	}
	return temp
}
func (a Address) String() string {
	return hex.EncodeToString(a.ToSlice())
}

func AddressFromBytes(b []byte) Address{
	if len(b) != 20 {
		msg := fmt.Sprintf("given bytes with length %d should be 20", len(b))
		panic(msg)
	}

	var value [20]uint8

	for i:=0;i<20;i++{
		value[i]=b[i]
	}

	return Address(value)
}
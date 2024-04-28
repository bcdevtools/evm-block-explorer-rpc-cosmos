package utils

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
)

func UnpackAbiString(bz []byte, _ string) (string, error) {
	unpacked, err := abiArgsSingleString.Unpack(bz)
	if err != nil {
		return "", err
	}
	return unpacked[0].(string), nil
}

var abiTypeString abi.Type
var abiArgsSingleString abi.Arguments

func init() {
	var err error
	abiTypeString, err = abi.NewType("string", "string", nil)
	if err != nil {
		panic(err)
	}

	abiArgsSingleString = abi.Arguments{
		abi.Argument{
			Name: "content",
			Type: abiTypeString,
		},
	}
}

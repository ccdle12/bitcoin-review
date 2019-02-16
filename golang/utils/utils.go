package utils

import (
	"errors"
	"math/big"
)

// ConvStrBigInt will take a string representation of a large number and
// convert it to a *big.Int.
func ConvStrBigInt(n string) (*big.Int, error) {
	x, success := new(big.Int).SetString(n, 10)
	if !success {
		return nil, errors.New("unable to convert string representation to big int")
	}

	return x, nil
}

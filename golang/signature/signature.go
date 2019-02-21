package signature

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
)

// Signature is a struct that holds the R and Q value for generating DER
// formatted signatures.
type Signature struct {
	R *big.Int
	S *big.Int
}

// GenerateDERSig will generate a signature given the R and S values in the
// Signature.
func (sig *Signature) GenerateDERSig() []byte {
	// Signature Format:
	// <DER><Length of signature><Marker for R><Length of R><R Value><Marker for S><Length of S><S Value>

	// Convert R and S to bytes.
	rbin := sig.R.Bytes()
	sbin := sig.S.Bytes()

	result := []byte{}

	// Check if rbin has a high bit.
	if rbin[0] > 128 {
		rbin = append([]byte{0x00}, rbin...)
	}

	// Get the length of R converted to bytes.
	rLen := convIntToTrimmedByte(len(rbin))

	// Append the mark 0x02 marker for r.
	result = append(result, []byte{0x02}...)

	// Append the length of r.
	result = append(result, rLen...)

	// Append rbin.
	result = append(result, rbin...)
	fmt.Printf("result after append: %x\n", result)

	// Check if sbin has a high bit.
	if sbin[0] > 128 {
		sbin = append([]byte{0x00}, sbin...)
	}

	// Get the length of S converted to bytes.
	sLen := convIntToTrimmedByte(len(sbin))

	// Append s values.
	result = append(result, []byte{0x02}...)
	result = append(result, sLen...)
	result = append(result, sbin...)

	// Get the length of the signature converted to bytes.
	sigLen := convIntToTrimmedByte(len(result))

	// Prepend the signature length and the DER type.
	result = append(sigLen, result...)
	result = append([]byte{0x30}, result...)

	return result
}

func convIntToTrimmedByte(i int) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(i))

	// Remove the leading 0s from the byte slice.
	b = bytes.Trim(b, "\x00")

	return b
}

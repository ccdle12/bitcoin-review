package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"golang.org/x/crypto/ripemd160"
	"hash"
	"math/big"
)

var (
	bigRadix       = big.NewInt(58)
	bigZero        = big.NewInt(0)
	base58Alphabet = []string{
		"1", "2", "3", "4", "5", "6", "7", "8", "9", "A",
		"B", "C", "D", "E", "F", "G", "H", "J", "K", "L",
		"M", "N", "P", "Q", "R", "S", "T", "U", "V", "W",
		"X", "Y", "Z", "a", "b", "c", "d", "e", "f", "g",
		"h", "i", "j", "k", "m", "n", "o", "p", "q", "r",
		"s", "t", "u", "v", "w", "x", "y", "z",
	}
	base58Table = [256]byte{
		255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255,
		255, 0, 1, 2, 3, 4, 5, 6,
		7, 8, 255, 255, 255, 255, 255, 255,
		255, 9, 10, 11, 12, 13, 14, 15,
		16, 255, 17, 18, 19, 20, 21, 255,
		22, 23, 24, 25, 26, 27, 28, 29,
		30, 31, 32, 255, 255, 255, 255, 255,
		255, 33, 34, 35, 36, 37, 38, 39,
		40, 41, 42, 43, 255, 44, 45, 46,
		47, 48, 49, 50, 51, 52, 53, 54,
		55, 56, 57, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255,
		255, 255, 255, 255, 255, 255, 255, 255,
	}
)

// EncodeBase58 will encode a slice of bytes in Base58 Encoded format.
func EncodeBase58(b []byte) string {
	var result string

	// Convert b to a Big Endian number.
	num := new(big.Int).SetBytes(b)

	// Use the remainder to find the next base58 character and concatenate to
	// the result.
	for num.Cmp(bigZero) > 0 {
		mod := new(big.Int)
		num.DivMod(num, bigRadix, mod)
		result = base58Alphabet[mod.Int64()] + result
	}

	// Replace leading 0s with 1s.
	for _, i := range b {
		if i != 0 {
			break
		}
		result = base58Alphabet[0] + result
	}

	return result
}

// DecodeBase58 will decode a string as bytes.
func DecodeBase58(s string) ([]byte, error) {
	if len(s) == 0 {
		return nil, errors.New("cannot pass empty string to decode")
	}

	result := big.NewInt(0)
	multi := big.NewInt(1)
	tmpBig := new(big.Int)

	// Decrement over the input string s.
	for i := len(s) - 1; i >= 0; i-- {
		tmp := base58Table[s[i]]

		// Catch invalid base58 characters.
		if tmp == 255 {
			return nil, errors.New("invalid base58 character")
		}

		tmpBig.SetInt64(int64(tmp))
		tmpBig.Mul(multi, tmpBig)
		result.Add(result, tmpBig)
		multi.Mul(multi, bigRadix)
	}

	// Assign the number of the result as bytes to tmpBytes.
	tmpBytes := result.Bytes()

	// Replace leading 1s.
	var numZeros int
	for numZeros = 0; numZeros < len(s); numZeros++ {
		if s[numZeros] != '1' {
			break
		}
	}

	// Create the output of bytes.
	length := numZeros + len(tmpBytes)
	output := make([]byte, length)

	// Copy the tmpBytes (result) to the output, after the leading zeroes.
	copy(output[numZeros:], tmpBytes)

	return output, nil
}

// ConvStrBigInt will take a string representation of a large number and
// convert it to a *big.Int.
func ConvStrBigInt(n string) (*big.Int, error) {
	x, success := new(big.Int).SetString(n, 10)
	if !success {
		return nil, errors.New("unable to convert string representation to big int")
	}

	return x, nil
}

// generateHash will will receive a bytes buffer and perform the hash function that is also passed as parameter.
func generateHash(b []byte, hasher hash.Hash) []byte {
	hasher.Write(b)

	// Returns the hashed result as a []byte.
	return hasher.Sum(nil)
}

// Hash160 generates the hash ripemd160(sha256(b)).
func Hash160(b []byte) []byte {
	// generate the hash of sha256 and then hash the result as ripemd160.
	return generateHash(generateHash(b, sha256.New()), ripemd160.New())
}

// DoubleSHA256 generates the hash sha256(sha256(b)).
func DoubleSHA256(b []byte) []byte {
	return generateHash(generateHash(b, sha256.New()), sha256.New())
}

// ConvIntStrToBigInt will convert the string representations of the int string
// to a big int.
func ConvIntStrToBigInt(s string) (*big.Int, error) {
	sInt, success := new(big.Int).SetString(s, 10)
	if !success {
		return nil, errors.New("unable to convert s parameter to big Int")
	}

	return sInt, nil
}

// ReverseStr will take a string as an argument and then reverse it, returning a string.
func ReverseStr(s string) string {
	var result string
	for _, r := range s {
		result = string(r) + result
	}

	return result
}

// ReadVarint will take a *byte.Buffer as an argument and return an int of the varint.
func ReadVarint(stream *bytes.Buffer) int {
	// Read the varint to find out the length of the script sig
	varintByte := stream.Next(1)
	varintBuf := make([]byte, 3)
	varintBuf = append(varintBuf, varintByte...)
	varintNum := binary.BigEndian.Uint32(varintBuf)

	var varint uint64
	switch varintNum {
	case 0xff:
		// Read the next 8 bytes.
		varintByte := stream.Next(8)
		varintBuf := make([]byte, 8)
		varintBuf = append(varintBuf, varintByte...)
		varint = binary.LittleEndian.Uint64(varintBuf)
	case 0xfe:
		// Read the next 4 bytes.
		varintByte := stream.Next(4)
		varintBuf := make([]byte, 4)
		varintBuf = append(varintBuf, varintByte...)
		varint = binary.LittleEndian.Uint64(varintBuf)
	case 0xfd:
		// Read the next 2 bytes.
		varintByte := stream.Next(2)
		varintBuf := make([]byte, 2)
		varintBuf = append(varintBuf, varintByte...)
		varint = binary.LittleEndian.Uint64(varintBuf)
	default:
		varint = uint64(varintNum)
	}

	return int(varint)
}

// EncodeVarint will take an int and encode it as a 1 byte varint.
func EncodeVarint(n int) ([]byte, error) {
	var varint []byte
	if n < 0xfd {
		varint = make([]byte, 8)
		binary.LittleEndian.PutUint32(varint, uint32(n))
		// TODO: Think of a way to consistenly us this as a function.
		// Removes trailing null characters.
		varint = bytes.Trim(varint, "\x00")
		return varint, nil

	} else if n < 0x10000 {
		varint = make([]byte, 8)
		varint = append(varint, 0xfd)
		binary.LittleEndian.PutUint32(varint, uint32(n))
		varint = bytes.Trim(varint, "\x00")
		return varint, nil

	} else if n < 0x100000000 {
		varint = make([]byte, 8)
		varint = append(varint, 0xfe)
		binary.LittleEndian.PutUint32(varint, uint32(n))
		varint = bytes.Trim(varint, "\x00")
		return varint, nil

		// TODO: Need to fix this overflow with n.
		// } else if n < 0x10000000000000000 {
		// varint = make([]byte, 8)
		// varint = append(varint, 0xff)
		// binary.LittleEndian.PutUint32(varint, uint32(n))
		// return varint, nil
	}

	return nil, errors.New("failed to encode varint")
}

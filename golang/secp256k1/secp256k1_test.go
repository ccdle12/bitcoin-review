package secp256k1

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"
)

func TestPrintConsts(t *testing.T) {
	// Generate a curve.
	secp256k1 := New()

	fmt.Printf("%v\n", secp256k1.P)
}

// TestHexDecodeString will test that we can convert the P value of secp256k1
// to bytes and then back into hex.
func TestHexDecodeString(t *testing.T) {
	decoded, err := hex.DecodeString(p)
	if err != nil {
		t.Fatalf("could not decode string")
	}

	// Check that when we encode it back, it is still the P string.
	encodedStr := hex.EncodeToString(decoded)
	if strings.ToUpper(encodedStr) != p {
		t.Fatalf("decoding failed, encodedStr: %s does not equal P: %s", encodedStr, p)
	}
}

// TestBigIntString will test that we can take a string and create a Big Int.
func TestBigIntString(t *testing.T) {
	pFromStr, convSuccesss := new(big.Int).SetString(p, 16)
	if !convSuccesss {
		t.Fatalf("failed to convert P to big Int: %v", convSuccesss)
	}

	decoded, err := hex.DecodeString(p)
	if err != nil {
		t.Fatalf("could not decode string")
	}

	p := new(big.Int).SetBytes(decoded)
	fmt.Printf("Printing p as a big int: %v\n", p)

	// Check that both ways can create the same p value.
	if p.String() != pFromStr.String() {
		t.Fatalf("p converted to string: %s and p from string: %s, failed to match",
			p.String(),
			pFromStr.String())
	}
}

// TestPrivKeyCollision is a quick sanity check, to make sure the private keys
// are not colliding.
func TestPrivKeyCollision(t *testing.T) {
	// Generate a curve.
	secp256k1 := New()

	// Generate two random private keys, they should NEVER match.
	privKey, err := rand.Int(rand.Reader, secp256k1.N)
	if err != nil {
		t.Fatalf("failed to generate a private key")
	}
	fmt.Printf("Private Key: %v\n", privKey)

	privKey2, err := rand.Int(rand.Reader, secp256k1.N)
	if err != nil {
		t.Fatalf("failed to generate a private key2 ")
	}
	fmt.Printf("Private Key 2: %v\n", privKey2)

	// Compare private keys.
	if privKey == privKey2 {
		t.Fatal("failed to generate different private keys")
	}
}

// TestImplemenationScalarMultiply will be an implementation test for Scalar
// Multiplication meaning we are programming a function to generate a Public
// Key given a Private Key.
func TestImplemenationScalarMultiply(t *testing.T) {
	// Generate a curve.
	secp256k1 := New()

	// Generate a random number below N (the private key).
	privKey, err := rand.Int(rand.Reader, secp256k1.N)
	if err != nil {
		t.Fatalf("failed to generate a private key")
	}
	fmt.Printf("Private Key: %v\n", privKey)

	// Convert the private key to big-endian byte slice.
	k := privKey.Bytes()
	fmt.Printf("k slice in bytes: %v\n", k)

	// Assign Bx, By, and Bz as the base.
	Bx := secp256k1.Gx
	By := secp256k1.Gy
	Bz := new(big.Int).SetInt64(1)

	// x, y, z will be used for point doubling.
	x := Bx
	y := By
	z := Bz

	// Loop over the bytes of the secret k.
	// Uses the double and add algorithm.
	for _, byte := range k {
		for bitNum := 0; bitNum < 8; bitNum++ {
			x, y, z = secp256k1.jacobianDouble(x, y, z)

			if byte&0x80 == 0x80 {
				x, y, z = secp256k1.JacobianAdd(Bx, By, Bz, x, y, z)
			}
			// TODO: (ccdle12) need to look intowhy we need to shift the byte
			// by 1.
			byte <<= 1
		}
	}

	// Convert the jacobian back to the affine.
	xOut, yOut := secp256k1.AffineFromJacobian(x, y, z)
	fmt.Printf("Public Key: xOut: %v, yOut: %v \n", xOut, yOut)

	// Check that the output of x and y are valid points on the curve.
	valid := secp256k1.IsOnCurve(xOut, yOut)
	fmt.Printf("Is public key valid: %v\n", valid)
	if !valid {
		t.Fatalf("xOut: %v and yOut: %v are not valid for the curve", xOut, yOut)
	}
}

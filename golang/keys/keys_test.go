package keys

import (
	"fmt"
	"github.com/ccdle12/bitcoin-review/golang/secp256k1"
	"github.com/ccdle12/bitcoin-review/golang/utils"
	"math/big"
	"testing"
)

// TestGenPrivKey will test that we can genreate a Private Key and ensure it is
// below N.
func TestGenPrivKey(t *testing.T) {
	// Generate a curve.
	curve := secp256k1.New()

	// Generate the private key.
	privateKey, err := generatePrivateKey(curve)
	if err != nil {
		t.Fatalf("failed to generate a private key: %s", err)
	}

	// Ensure the random number (private key is below N).
	c := privateKey.secret.Cmp(curve.N)
	if c == 1 {
		t.Fatalf("private key: %v should not be greater than n: %v",
			privateKey, curve.N)
	}
}

// TestGenKeyPair will test that we can use the key pair constructor to create
// a Private and Public Key Pair. The constructor will already check if the
// Public Key is valid, but for sanities sake we will check it again.
func TestGenKeyPair(t *testing.T) {
	keys, err := New()
	if err != nil {
		t.Fatalf("test key gen pair failed: %v\n", err)
	}

	// Sanity Check for Public Key Pair.
	x := keys.PublicKey.X
	y := keys.PublicKey.Y

	// Check that the Public Key points x and y are valid points on the curve.
	valid := keys.Curve.IsOnCurve(x, y)
	if !valid {
		t.Fatalf("x: %v and y: %v are not valid for the curve", x, y)
	}
}

// TestUncompressedSec will test that we can generate an uncompressed sec
// Public Key.
func TestUncompressedSec(t *testing.T) {
	// Use existing pub key.
	// x:= 4009715469895962904302745416817721540571577912364644137838095050706137667860
	// y:= 32025336288095498019218993550383068707359510270784983226210884843871535451292
	x, err := utils.ConvStrBigInt("4009715469895962904302745416817721540571577912364644137838095050706137667860")
	if err != nil {
		t.Fatalf("unable to convert x value to big int")
	}

	y, err := utils.ConvStrBigInt("32025336288095498019218993550383068707359510270784983226210884843871535451292")
	if err != nil {
		t.Fatalf("unable to convert y value to big int")
	}

	// Create a Public Key.
	publicKey := &PublicKey{X: x, Y: y}

	// Convert big Ints to big endian bytes := secX, secY.
	xBytes := publicKey.X.Bytes()
	yBytes := publicKey.Y.Bytes()

	fmt.Printf("SEC: %x\n", xBytes)
	fmt.Printf("SEC: %x\n", yBytes)

	// Created expected for uncompressed, prepend b'x04' to the (secX + secY).
	expectedUncompressed := append(xBytes, yBytes...)
	expectedUncompressed = append([]byte{0x04}, expectedUncompressed...)

	// Generate a sec formatted uncompressed key.
	secUncompressed := generateUncompressedSec(publicKey)
	fmt.Printf("Uncompressed SEC: %x\n", expectedUncompressed)
	fmt.Printf("Generated Uncompressed SEC: %x\n", secUncompressed)

	// Make sure we generate the expected uncompressed sec.
	for i := range secUncompressed {
		if secUncompressed[i] != expectedUncompressed[i] {
			t.Fatalf("failed to generate the expected uncompressed sec pub key."+
				"Expected: %x, Received: %x\n", expectedUncompressed,
				secUncompressed)
		}
	}
}

// TestCompressedSec will test that we can generate an uncompressed sec
// Public Key.
func TestCompressedSec(t *testing.T) {
	// Use existing pub key: (even Y)
	// x:= 4009715469895962904302745416817721540571577912364644137838095050706137667860
	// y:= 32025336288095498019218993550383068707359510270784983226210884843871535451292
	x, err := utils.ConvStrBigInt("4009715469895962904302745416817721540571577912364644137838095050706137667860")
	if err != nil {
		t.Fatalf("unable to convert x value to big int")
	}

	y, err := utils.ConvStrBigInt("32025336288095498019218993550383068707359510270784983226210884843871535451292")
	if err != nil {
		t.Fatalf("unable to convert y value to big int")
	}

	// Create a Public Key.
	publicKey := &PublicKey{X: x, Y: y}

	// Convert big Ints to big endian bytes := secX, secY.
	xBytes := publicKey.X.Bytes()

	// Check if Y is odd or even.
	yMarker := y.Mod(y, big.NewInt(2))

	// Comparison should return 0 meaning they are the same value.
	if yMarker.Cmp(big.NewInt(0)) != 0 {
		t.Fatalf("yMarker should have been even")
	}

	// Created expected for uncompressed, prepend b'x02' to the (secX).
	expectedCompressed := append([]byte{0x02}, xBytes...)

	// Generate a sec formatted uncompressed key.
	secCompressed := generateCompressedSec(publicKey)
	fmt.Printf("Compressed SEC: %x\n", expectedCompressed)
	fmt.Printf("Generated Compressed SEC: %x\n", secCompressed)

	// Make sure we generate the expected uncompressed sec.
	for i := range secCompressed {
		if secCompressed[i] != expectedCompressed[i] {
			t.Fatalf("failed to generate the expected compressed sec pub key."+
				"Expected: %x, Received: %x\n", expectedCompressed,
				secCompressed)
		}
	}

}

func TestOddCompressedSec(t *testing.T) {
	// Use an existing pub key: (odd Y)
	// x := 43733605778270459583874364812384261459365992207657902102567152558096696733127
	// y := 64346778444748414606606796249150556060624935198788845168028978963277938956739
	x, err := utils.ConvStrBigInt("43733605778270459583874364812384261459365992207657902102567152558096696733127")
	if err != nil {
		t.Fatalf("unable to convert x value to big int")
	}

	y, err := utils.ConvStrBigInt("64346778444748414606606796249150556060624935198788845168028978963277938956739")
	if err != nil {
		t.Fatalf("unable to convert y value to big int")
	}

	// Check if Y is odd or even.
	yMarker := new(big.Int).Mod(y, big.NewInt(2))

	// Comparison should not return 0 meaning they are not the same value,
	// indicating the yMarker is odd.
	if yMarker.Cmp(big.NewInt(0)) == 0 {
		t.Fatalf("yMarker should have been odd")
	}

	// Create a Public Key.
	publicKey := &PublicKey{X: x, Y: y}
	fmt.Printf("SEC: Pubkey y: %v\n", publicKey.Y)

	// Convert big Ints to big endian bytes := secX, secY.
	xBytes := publicKey.X.Bytes()

	// Created expected for uncompressed, prepend b'x03' to the (secX) since
	// the Y value is odd.
	expectedCompressed := append([]byte{0x03}, xBytes...)

	// Generate a sec formatted uncompressed key.
	secCompressed := generateCompressedSec(publicKey)
	fmt.Printf("Compressed SEC: %x\n", expectedCompressed)
	fmt.Printf("Generated Compressed SEC: %x\n", secCompressed)

	// Make sure we generate the expected uncompressed sec.
	for i := range secCompressed {
		if secCompressed[i] != expectedCompressed[i] {
			t.Fatalf("failed to generate the expected compressed sec pub key."+
				"Expected: %x, Received: %x\n", expectedCompressed,
				secCompressed)
		}
	}
}

// TestGenTestnetAddress will test that we can generate a testnet compatible
// address.
func TestGenTestnetAddress(t *testing.T) {
	// Use an existing pub key: (odd Y)
	// x := 43651727216793576570341989570883305974491642311510342469928224726666590034225
	// y := 109857391791750504773247734335453148952192151977881622854599464318335318347795
	x, err := utils.ConvStrBigInt("43651727216793576570341989570883305974491642311510342469928224726666590034225")
	if err != nil {
		t.Fatalf("unable to convert x value to big int")
	}

	y, err := utils.ConvStrBigInt("109857391791750504773247734335453148952192151977881622854599464318335318347795")
	if err != nil {
		t.Fatalf("unable to convert y value to big int")
	}

	publicKey := &PublicKey{X: x, Y: y}
	secCompressed := generateCompressedSec(publicKey)

	address := GenerateTestnetAddress(secCompressed)

	expected := "mo24iC138ffpdWiFsH8y7dq6v5CDD1UbiT"

	if address != expected {
		t.Fatalf("failed to generate correct address, expeted: %v, received: %v", expected, address)
	}
}

// TestGenMainnetAddress will test that we can generate a mainnet compatible
// address.
func TestGenMainnetAddress(t *testing.T) {
	// Use an existing pub key: (odd Y)
	// x := 43651727216793576570341989570883305974491642311510342469928224726666590034225
	// y := 109857391791750504773247734335453148952192151977881622854599464318335318347795
	x, err := utils.ConvStrBigInt("43651727216793576570341989570883305974491642311510342469928224726666590034225")
	if err != nil {
		t.Fatalf("unable to convert x value to big int")
	}

	y, err := utils.ConvStrBigInt("109857391791750504773247734335453148952192151977881622854599464318335318347795")
	if err != nil {
		t.Fatalf("unable to convert y value to big int")
	}

	publicKey := &PublicKey{X: x, Y: y}
	secCompressed := generateCompressedSec(publicKey)

	address := GenerateMainnetAddress(secCompressed)

	expected := "18W7R8v4KeEZrQEe9iAbHicn45bWNn2QBe"

	if address != expected {
		t.Fatalf("failed to generate correct address, expeted: %v, received: %v", expected, address)
	}
}

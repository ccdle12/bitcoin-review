package utils

import (
	"fmt"
	"testing"
)

// TestConvStrBigInt will test that we can take a string representation of a
// large number and conert it to a big.Int.
func TestConvStrBigInt(t *testing.T) {
	x, err := ConvStrBigInt("4009715469895962904302745416817721540571577912364644137838095050706137667860")

	if err != nil {
		t.Fatalf("should have been able to convert to big int")
	}

	expected := "4009715469895962904302745416817721540571577912364644137838095050706137667860"
	if x.String() != expected {
		t.Fatalf("unable to convert the x: %v to the correct represenation: %v\n", x.String(), expected)
	}
}

// TestDecodeBase58 will test that we can decode a Base58 encoded address.
func TestDecodeBase58(t *testing.T) {
	// Hex String representation of the decoded address.
	// Expected the Hexadcimal String of the pre-encoding of address (mo24iC138ffpdWiFsH8y7dq6v5CDD1UbiT)
	expected := "6f524a4c9f658b9e482c40669096d93f2a6d96de523106c664"
	address := "mo24iC138ffpdWiFsH8y7dq6v5CDD1UbiT"

	result, err := DecodeBase58(address)
	if err != nil {
		t.Fatalf("failed to decode a base58 address: %v", err)
	}
	resultHex := fmt.Sprintf("%x", string(result))

	if resultHex != expected {
		t.Fatalf("failed to decode the address "+
			"expected: %v received: %v", expected, resultHex)
	}

	// Hex String representation of the decoded address.
	// Expected is the pre-encoding of address (18W7R8v4KeEZrQEe9iAbHicn45bWNn2QBe)
	expected = "00524a4c9f658b9e482c40669096d93f2a6d96de52eb64c4fd"
	address = "18W7R8v4KeEZrQEe9iAbHicn45bWNn2QBe"

	result, err = DecodeBase58(address)
	if err != nil {
		t.Fatalf("failed to decode a base58 address: %v", err)
	}
	resultHex = fmt.Sprintf("%x", string(result))

	if resultHex != expected {
		t.Fatalf("failed to decode the address "+
			"expected: %v received: %v", expected, resultHex)
	}
}

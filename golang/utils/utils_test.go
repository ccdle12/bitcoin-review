package utils

import (
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

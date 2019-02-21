package signature

import (
	"fmt"
	"github.com/ccdle12/bitcoin-review/golang/utils"
	"testing"
)

// TestGenerateDERSignature will test that we can generate a valid DER
// formatted signature.
func TestParseSignature(t *testing.T) {
	r, err := utils.ConvIntStrToBigInt("65768643913645672968978426589689987237850374542483501912088659345491159391021")
	if err != nil {
		t.Fatalf("unable to convert r")
	}

	s, err := utils.ConvIntStrToBigInt("55618899300744280687710599871980893657541124572884031214465422719409044157728")
	if err != nil {
		t.Fatalf("unable to convert s")
	}

	sig := &Signature{R: r, S: s}
	if sig == nil {
		t.Fatalf("unable to create sig object")
	}

	// Expected DER signature, from the (r,s).
	expected := "30450221009167bbb944c67d650cab2f3d5cbd06c2391977de478832c50d4af00b0a2f9b2d02207af72e71cecf43022204cce257af9625f799d5a90ed904c42b903771dd217520"
	der := sig.GenerateDERSig()

	derStr := fmt.Sprintf("%x", string(der))

	fmt.Printf("TEST: %x", der)
	if derStr != expected {
		t.Fatalf("der signature: %v, does not match expected: %v", derStr, expected)
	}
}

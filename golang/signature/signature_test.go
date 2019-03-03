package signature

import (
	// "crypto/ecdsa"
	"fmt"
	"github.com/ccdle12/bitcoin-review/golang/secp256k1"
	"github.com/ccdle12/bitcoin-review/golang/utils"
	"math/big"
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

// TestValidSignature will test that we can validate a signature.
func TestValidSignature(t *testing.T) {
	px, err := utils.ConvHexStrToBigInt("04519fac3d910ca7e7138f7013706f619fa8f033e6ec6e09370ea38cee6a7574")
	if err != nil {
		t.Fatalf(err.Error())
	}
	py, err := utils.ConvHexStrToBigInt("82b51eab8c27c66e26c858a079bcdf4f1ada34cec420cafc7eac1a42216fb6c4")
	if err != nil {
		t.Fatalf(err.Error())
	}

	z1, err := utils.ConvHexStrToBigInt("bc62d4b80d9e36da29c16c5d4d9f11731f36052c72401a76c23c0fb5a9b74423")
	if err != nil {
		t.Fatalf(err.Error())
	}
	r1, err := utils.ConvHexStrToBigInt("37206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c6")
	if err != nil {
		t.Fatalf(err.Error())
	}
	s1, err := utils.ConvHexStrToBigInt("8ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec")
	if err != nil {
		t.Fatalf(err.Error())
	}

	curve := secp256k1.New()

	n2 := new(big.Int).Sub(curve.N, big.NewInt(2))

	// u = z * pow(s, N-2, N) % N
	s1Inverse := new(big.Int).Exp(s1, n2, curve.N)
	fmt.Printf("u1: %v\n", s1Inverse)
	u1 := new(big.Int).Mul(z1, s1Inverse)
	u1.Mod(u1, curve.N)
	fmt.Printf("sig: %v\n", u1)

	// v1 := new(big.Int).Exp(s1, new(big.Int).Sub(curve.N, big.NewInt(2)), curve.N)
	// v1 := new(big.Int).Exp()
	v1 := new(big.Int).Mul(r1, s1Inverse)
	v1.Mod(v1, curve.N)
	fmt.Printf("sig: %v\n", v1)

	x2, y2, z2 := curve.ScalarMult(u1.Bytes())
	// aX2, aY2 := curve.AffineFromJacobian(x2, y2, z2)
	// fmt.Printf("sig x2: %v\n %v\n %v\n", x2, y2, z2)

	x3, y3, z3 := curve.GenericScalarMult(px, py, v1.Bytes())
	// aX3, aY3 := curve.AffineFromJacobian(x3, y3, z3)
	fmt.Printf("sig x3: %v\n %v\n %v\n", x3, y3, z3)

	// TODO: Use simple add
	// resultX, _ := curve.SimpleAdd(aX2, aY2, aX3, aY3)

	// TODO: Now we need to add the two points together.
	x4, y4, z4 := curve.JacobianAdd(x2, y2, z2, x3, y3, z3)
	resultX, _ := curve.AffineFromJacobian(x4, y4, z4)

	eq := resultX.Cmp(r1)
	fmt.Printf("result: resultx %v\n", resultX)
	fmt.Printf("result: r1 %v\n", r1)
	fmt.Printf("eq: %v\n", eq)
}

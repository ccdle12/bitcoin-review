package keys

import (
	"github.com/ccdle12/bitcoin-review/golang/secp256k1"
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

package keys

import (
	"crypto/rand"
	"errors"
	"github.com/ccdle12/bitcoin-review/golang/secp256k1"
	"math/big"
)

// Keys contains a Private Key and Public Key.
type Keys struct {
	Curve      *secp256k1.Secp256k1 // TODO: replace this with Curve interface
	PrivateKey *PrivateKey
	PublicKey  *PublicKey
}

// New is the constructor for creating a key pair. It will generate a Private
// Key and Public Key pair.
// TODO: Pass a curve as an argument construct.
func New() (*Keys, error) {
	// Generate a curve.
	secp256k1 := secp256k1.New()

	// 1. Generate a private key and assign.
	privateKey, err := generatePrivateKey(secp256k1)
	if err != nil {
		return nil, err
	}
	// 2. Generate a public key and assign.
	publicKey, err := generatePublicKey(secp256k1, privateKey)
	if err != nil {
		return nil, err
	}

	return &Keys{secp256k1, privateKey, publicKey}, err
}

// PrivateKeys is the struct to hold Private Key information.
type PrivateKey struct {
	secret *big.Int
}

// PublicKey is the struct that holds Public Key information.
type PublicKey struct {
	X *big.Int
	Y *big.Int
}

// generatePrivateKey will generate a new Private Key.
// TODO: Update the parameter to use the Curve interface
func generatePrivateKey(curve *secp256k1.Secp256k1) (*PrivateKey, error) {
	// Generate the Private Key secret.
	secret, err := rand.Int(rand.Reader, curve.N)
	if err != nil {
		return nil, errors.New("failed to generate a private key")
	}

	return &PrivateKey{secret: secret}, nil
}

// generatePublicKey will generate a new Public Key.
// TODO: Update the parameter to use the Curve interface
func generatePublicKey(curve *secp256k1.Secp256k1, pk *PrivateKey) (*PublicKey, error) {

	jx, jy, jz := curve.ScalarMult(pk.secret.Bytes())
	x, y := curve.AffineFromJacobian(jx, jy, jz)

	validPoint := curve.IsOnCurve(x, y)
	if !validPoint {
		return nil, errors.New("the public key generated is not on the curve and therefore invalid")
	}

	return &PublicKey{X: x, Y: y}, nil
}

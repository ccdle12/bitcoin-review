package secp256k1

import (
	"errors"
	"math/big"
)

const (
	// Keep them all as strings for now.
	p = "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F" // Prime Modulo of the Field.
	a = "0000000000000000000000000000000000000000000000000000000000000000" // Part of the equation for secp256k1. y^2 = x^3 + ax + b.
	b = "0000000000000000000000000000000000000000000000000000000000000007" // Part of the equation for secp256k1. y^2 = x^3 + ax + b
	// g = "0479BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8" // The publicly known base point (generator point).
	gx = "79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798" // X co-ordinate of the base point (generator point)
	gy = "483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8" // X co-ordinate of the base point (generator point)
	n  = "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141" // Number of points in the field of secp256k1.
)

// Secp256k1 is the implementation of this s.
type Secp256k1 struct {
	P  *big.Int
	A  *big.Int
	B  *big.Int
	Gx *big.Int
	Gy *big.Int
	N  *big.Int
}

// New is the constructor for the s Secp256k.
// TODO: make this implement a Curve interface.
func New() *Secp256k1 {
	P, _ := convHexStrToBigInt(p)
	A, _ := convHexStrToBigInt(a)
	B, _ := convHexStrToBigInt(b)
	Gx, _ := convHexStrToBigInt(gx)
	Gy, _ := convHexStrToBigInt(gy)
	N, _ := convHexStrToBigInt(n)

	return &Secp256k1{P, A, B, Gx, Gy, N}
}

// convHexStrToBigInt will convert the constants of the s, that are in
// String representations of Hex to bigInts.
func convHexStrToBigInt(sParam string) (*big.Int, error) {
	sInt, success := new(big.Int).SetString(sParam, 16)
	if !success {
		return nil, errors.New("unable to convert s parameter to big Int")
	}

	return sInt, nil
}

// TODO: EXPLAIN
func (s *Secp256k1) jacobianAdd(x1, y1, z1, x2, y2, z2 *big.Int) (*big.Int, *big.Int, *big.Int) {
	// See http://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#addition-add-2007-bl
	z1z1 := new(big.Int).Mul(z1, z1)
	z1z1.Mod(z1z1, s.P)
	z2z2 := new(big.Int).Mul(z2, z2)
	z2z2.Mod(z2z2, s.P)

	u1 := new(big.Int).Mul(x1, z2z2)
	u1.Mod(u1, s.P)

	u2 := new(big.Int).Mul(x2, z1z1)
	u2.Mod(u2, s.P)

	h := new(big.Int).Sub(u2, u1)
	// Check if h is less than 0.
	if h.Sign() == -1 {
		h.Add(h, s.P)
	}
	i := new(big.Int).Lsh(h, 1)
	i.Mul(i, i)
	j := new(big.Int).Mul(h, i)

	s1 := new(big.Int).Mul(y1, z2)
	s1.Mul(s1, z2z2)
	s1.Mod(s1, s.P)
	s2 := new(big.Int).Mul(y2, z1)
	s2.Mul(s2, z1z1)
	s2.Mod(s2, s.P)
	r := new(big.Int).Sub(s2, s1)
	// Check if r is less than 0.
	if r.Sign() == -1 {
		r.Add(r, s.P)
	}
	r.Lsh(r, 1)
	v := new(big.Int).Mul(u1, i)

	x3 := new(big.Int).Set(r)
	x3.Mul(x3, x3)
	x3.Sub(x3, j)
	x3.Sub(x3, v)
	x3.Sub(x3, v)
	x3.Mod(x3, s.P)

	y3 := new(big.Int).Set(r)
	v.Sub(v, x3)
	y3.Mul(y3, v)
	s1.Mul(s1, j)
	s1.Lsh(s1, 1)
	y3.Sub(y3, s1)
	y3.Mod(y3, s.P)

	z3 := new(big.Int).Add(z1, z2)
	z3.Mul(z3, z3)
	z3.Sub(z3, z1z1)
	if z3.Sign() == -1 {
		z3.Add(z3, s.P)
	}
	z3.Sub(z3, z2z2)
	if z3.Sign() == -1 {
		z3.Add(z3, s.P)
	}
	z3.Mul(z3, h)
	z3.Mod(z3, s.P)

	return x3, y3, z3

}

// TODO: EXPLAIN
func (s *Secp256k1) jacobianDouble(x, y, z *big.Int) (*big.Int, *big.Int, *big.Int) {
	// See http://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#doubling-dbl-2009-l
	a := new(big.Int).Mul(x, x) // X1^2
	b := new(big.Int).Mul(y, y) // Y1^2
	c := new(big.Int).Mul(b, b) // B^2

	d := new(big.Int).Add(x, b) //X1 + B
	d.Mul(d, d)                 // (X1 +B)^2
	d.Sub(d, a)                 //  (X1 + B)^2 - A
	d.Sub(d, c)                 //  (X1 + B)^2 - C
	d.Mul(d, big.NewInt(2))     // 2 * ((X1 + B)^2 - A - C)

	e := new(big.Int).Mul(big.NewInt(3), a) // 3 * A
	f := new(big.Int).Mul(e, e)             // E^2

	x3 := new(big.Int).Mul(big.NewInt(2), d) // 2 * D
	x3.Sub(f, x3)                            // F - 2 * D
	x3.Mod(x3, s.P)

	y3 := new(big.Int).Sub(d, x3)                  // D - X3
	y3.Mul(e, y3)                                  // E * (D-X3)
	y3.Sub(y3, new(big.Int).Mul(big.NewInt(8), c)) // E * (D-X3) - 8 * C
	y3.Mod(y3, s.P)

	z3 := new(big.Int).Mul(y, z) // Y1 * Z1
	z3.Mul(big.NewInt(2), z3)
	z3.Mod(z3, s.P)

	return x3, y3, z3
}

// TODO: EXPLAIN
func (s *Secp256k1) AffineFromJacobian(x, y, z *big.Int) (*big.Int, *big.Int) {
	zinv := new(big.Int).ModInverse(z, s.P)
	zinvsq := new(big.Int).Mul(zinv, zinv)

	xOut := new(big.Int).Mul(x, zinvsq)
	xOut.Mod(xOut, s.P)

	zinvsq.Mul(zinvsq, zinv)

	yOut := new(big.Int).Mul(y, zinvsq)
	yOut.Mod(yOut, s.P)

	return xOut, yOut
}

// IsOnCurve is a function to check whether the x,y co-ordinates satisfy the
// curve.
func (s *Secp256k1) IsOnCurve(x, y *big.Int) bool {
	// EQUATION: secp256k1 = y^2 = x^3 + b
	// After inputting x and y, if both sides of the equation are satisfied,
	// then the Point is a valid point on the curve.
	y2 := new(big.Int).Mul(y, y) // y^2
	y2.Mod(y2, s.P)              // y^2 % P

	x3 := new(big.Int).Mul(x, x) // x^2
	x3.Mul(x3, x)                // x^3

	x3.Add(x3, s.B) // x^3 + B
	x3.Mod(x3, s.P) // (x^3 + B) % P

	return x3.Cmp(y2) == 0
}

// ScalarMult is the open function for scalar multiplication on a curve.
func (s *Secp256k1) ScalarMult(k []byte) (*big.Int, *big.Int, *big.Int) {
	// Assign Bx, By, and Bz as the base.
	Bx := s.Gx
	By := s.Gy
	Bz := new(big.Int).SetInt64(1)

	// x, y, z will be used for point doubling.
	x := Bx
	y := By
	z := Bz

	// Loop over the bytes of the secret k.
	// Uses the double and add algorithm.
	for _, byte := range k {
		for bitNum := 0; bitNum < 8; bitNum++ {
			x, y, z = s.jacobianDouble(x, y, z)

			if byte&0x80 == 0x80 {
				x, y, z = s.jacobianAdd(Bx, By, Bz, x, y, z)
			}
			// TODO: (ccdle12) need to look intowhy we need to shift the byte
			// by 1.
			// byte <<= 1
		}
	}

	return x, y, z

	// Convert the jacobian back to the affine.
	// xOut, yOut := secp256k1.affineFromJacobian(x, y, z)
	// fmt.Printf("Public Key: xOut: %v, yOut: %v \n", xOut, yOut)

	// Check that the output of x and y are valid points on the curve.
	// valid := secp256k1.IsOnCurve(xOut, yOut)
	// fmt.Printf("Is public key valid: %v\n", valid)
	// if !valid {
	// 	t.Fatalf("xOut: %v and yOut: %v are not valid for the curve", xOut, yOut)
	// }
}

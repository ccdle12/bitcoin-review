extern crate num_integer;
extern crate ramp;
use self::ramp::Int;

/// Secp256k1 constants.
const P: &str = "115792089237316195423570985008687907853269984665640564039457584007908834671663"; // Prime Modulo of the Field.
const A: &str = "0000000000000000000000000000000000000000000000000000000000000000"; // Part of the equation for secp256k1. y^2 = x^3 + ax + b.
const B: &str = "0000000000000000000000000000000000000000000000000000000000000007"; // Part of the equation for secp256k1. y^2 = x^3 + ax + b.
const GX: &str = "79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798"; // X co-ordinate of the base point (generator point).
const GY: &str = "483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8"; // Y co-ordinate of the base point (generator point).
                                                                                     // NOTE: base 16, need to change all of it to the same format.
const N: &str = "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141"; // Number of points in the field of secp256k1.

// get_p will return the constant P as a Int.
pub fn get_p() -> Int {
    return Int::from_str_radix(
        "115792089237316195423570985008687907853269984665640564039457584007908834671663",
        10,
    )
    .unwrap();
}

// get_b will return the constant B as a Int.
pub fn get_b() -> Int {
    return Int::from_str_radix(
        "0000000000000000000000000000000000000000000000000000000000000007",
        10,
    )
    .unwrap();
}

// get_n will return the constant N as a Int.
pub fn get_n() -> Int {
    return Int::from_str_radix(
        "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141",
        16,
    )
    .unwrap();
}

// is_on_curve checks whether two x,y co-ordinates satisfy the curve secp256k1.
pub fn is_on_curve(x: &Int, y: &Int) -> bool {
    // EQUATION: secp256k1 = y^2 = x^3 + b
    let p = get_p();
    let b = get_b();

    let y1 = y.pow(2) % &p;
    let x1 = (x.pow(3) + &b) % &p;

    if y1 == x1 {
        return true;
    }

    return false;
}

#[cfg(test)]
mod tests {
    #[test]
    fn test_big_int_addition() {
        let a = super::Int::from(10);
        let b = super::Int::from(5);

        let c = &a + &b;

        let expected = super::Int::from(15);

        assert_eq!(expected, c);
    }

    #[test]
    fn test_get_number_p() {
        let p = super::get_p();

        assert_eq!(
            "115792089237316195423570985008687907853269984665640564039457584007908834671663",
            p.to_string()
        );
    }

    #[test]
    fn test_is_on_curve() {
        // Pre-generated x,y co-ordinates, they should be valid on secp256k1.
        let x = super::Int::from_str_radix(
            "38691711181538895418159153243459074037181781406247798684369213484306497014884",
            10,
        )
        .unwrap();

        let y = super::Int::from_str_radix(
            "52700585418916410786512382566155139688549611752353407161567662246516672768510",
            10,
        )
        .unwrap();

        assert!(super::is_on_curve(&x, &y))
    }
}

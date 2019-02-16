extern crate num_integer;
extern crate ramp;
extern crate rand;

use self::ramp::Int;
use self::ramp::RandomInt;
use secp256k1::{get_n, is_on_curve};

// PrivateKey is a struct that holds the generated secret.
pub struct PrivateKey {
    secret: Int,
}

// PublicKey is a struct that holds the generated Public Key.
pub struct PublicKey {
    x: Int,
    y: Int,
}

// generate_private_key will generate a secret and initialise a Private Key object.
pub fn generate_private_key() -> PrivateKey {
    // Generate a random number below n.
    let mut rng = rand::thread_rng();
    let secret = rng.gen_uint_below(&get_n());

    PrivateKey { secret }
}

// generate_public_key will generate a public key given the secret as a scalar.
// pub fn generate_public_key(secret: &Int) -> PublicKey {
//     
// }

#[cfg(test)]
mod tests {
    use std::collections::hash_map::Entry;
    use std::collections::HashMap;

    #[test]
    fn test_private_key_generate() {
        // Generate a private key.
        // NOTE: should pass in a curve object.
        let private_key = super::generate_private_key();
        println!("Private Key below: {}", private_key.secret);

        // Check the private key is below the number N.
        assert!(private_key.secret < super::get_n())
    }

    #[test]
    // Naive test for private key collision.
    fn test_private_key_collision() {
        // Create Hash Map.
        let mut cache: HashMap<String, bool> = HashMap::new();

        // Loop from 0 to 100.
        for i in 0..100 {
            // Generate a private key.
            let mut private_key = super::generate_private_key();

            // Check the generated key does not exist in the cache.
            match cache.entry(private_key.secret.to_string()) {
                Entry::Occupied(_) => panic!(
                    "test_private_key_collision failed: there was a
                                             collision with private keys generated"
                ),
                Entry::Vacant(_) => {}
            };

            // Add the generated private key to the cache.
            cache.insert(private_key.secret.to_string(), true);
        }
    }

    #[test]
    fn test_public_key_generate() {
        // Generate a private key.
        let private_key = super::generate_private_key();

        println!("{}", private_key.secret);

        // TEST:
        // Get the byte size of the private key secret.
        // pub fn bit_length(&self) -> u32
        //
        // Loop over the size of the bit_length and call the function on each ith
        // pub fn bit(&self, bit: u32) -> bool
        //
        // Implement double and add algorithm
        // https://crypto.stackexchange.com/questions/3907/how-does-one-calculate-the-scalar-multiplication-on-elliptic-curves

        // Generate a public key, using the private key.
        // let public_key = super::generate_public_key(&private_key.secret);

        // Check that the public key is valid.
        // let valid_key = super::is_on_curve(&public_key.x, &public_key.y);

        // assert!(valid_key);
    }

}

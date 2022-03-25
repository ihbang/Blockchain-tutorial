use crypto::{digest::Digest, sha2::Sha256};
use primitive_types::U256;
use std::ops::Shl;

use crate::block::Block;

const DIFFICULTY: u16 = 6;
const NONCE_MAX: u64 = std::u64::MAX;

pub struct ProofOfWork<'a> {
    block: &'a Block,
    target: U256,
}

impl<'a> ProofOfWork<'a> {
    pub fn new(block: &'a Block) -> Self {
        let target: U256 = U256::from(1u16);
        let target = target.shl(256 - 4 * DIFFICULTY);

        ProofOfWork { block, target }
    }

    pub fn run(&self) -> (u64, Vec<u8>) {
        let mut hash;

        println!(
            "Mining the block containing \"{}\"",
            String::from_utf8_lossy(self.block.data().as_slice())
        );
        for nonce in 0..NONCE_MAX {
            hash = self.calculate_hash(nonce);
            print!("\r{}", hex::encode(hash));

            let hash_val = U256::from(hash);

            if self.target > hash_val {
                println!("\n");
                return (nonce, Vec::from(hash));
            }
        }
        (NONCE_MAX, vec![0; 32])
    }

    pub fn validate(&self) -> bool {
        let hash = self.calculate_hash(self.block.nonce());

        self.target > U256::from(hash)
    }

    fn calculate_hash(&self, nonce: u64) -> [u8; 32] {
        let mut hasher = Sha256::new();

        let hash_data = self.block.timestamp().as_secs().to_string();
        hasher.input_str(&hash_data);

        let hash_data = self.block.data().as_slice();
        hasher.input(hash_data);
        hasher.input(self.block.prev_block_hash());
        hasher.input(&nonce.to_ne_bytes());

        let mut hash = [0; 32];
        hasher.result(&mut hash);
        hash
    }
}

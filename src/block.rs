use std::time::{Duration, SystemTime};

use crate::proofofwork::ProofOfWork;

pub struct Block {
    timestamp: Duration,
    data: Vec<u8>,
    prev_block_hash: [u8; 32],
    hash: [u8; 32],
    nonce: u64,
}

impl Block {
    pub fn new(data: &str, prev_block_hash: &[u8; 32]) -> Self {
        let timestamp = SystemTime::now()
            .duration_since(SystemTime::UNIX_EPOCH)
            .unwrap();
        let data = Vec::from(data.as_bytes());

        let mut block = Block {
            timestamp,
            data,
            prev_block_hash: *prev_block_hash,
            hash: [0; 32],
            nonce: 0,
        };

        let pow = ProofOfWork::new(&block);
        let (nonce, hash) = pow.run();
        block.hash = hash;
        block.nonce = nonce;
        block
    }

    pub fn timestamp(&self) -> &Duration {
        &self.timestamp
    }

    pub fn data(&self) -> &Vec<u8> {
        &self.data
    }

    pub fn prev_block_hash(&self) -> &[u8; 32] {
        &self.prev_block_hash
    }

    pub fn hash(&self) -> &[u8; 32] {
        &self.hash
    }

    pub fn nonce(&self) -> u64 {
        self.nonce
    }
}

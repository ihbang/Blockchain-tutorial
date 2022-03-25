use serde::{Deserialize, Serialize};
use std::time::{Duration, SystemTime};

use crate::proofofwork::ProofOfWork;

#[derive(Serialize, Deserialize)]
pub struct Block {
    timestamp: Duration,
    data: Vec<u8>,
    prev_block_hash: Vec<u8>,
    hash: Vec<u8>,
    nonce: u64,
}

impl Block {
    pub fn new(data: &str, prev_block_hash: &[u8]) -> Self {
        let timestamp = SystemTime::now()
            .duration_since(SystemTime::UNIX_EPOCH)
            .unwrap();
        let data = Vec::from(data.as_bytes());

        let mut block = Block {
            timestamp,
            data,
            prev_block_hash: prev_block_hash.to_vec(),
            hash: Vec::new(),
            nonce: 0,
        };

        let pow = ProofOfWork::new(&block);
        let (nonce, hash) = pow.run();
        block.hash = hash;
        block.nonce = nonce;
        block
    }

    pub fn serialize(&self) -> Vec<u8> {
        bincode::serialize(&self).unwrap()
    }

    pub fn deserialize(data: &[u8]) -> Block {
        bincode::deserialize(data).unwrap()
    }

    pub fn timestamp(&self) -> &Duration {
        &self.timestamp
    }

    pub fn data(&self) -> &Vec<u8> {
        &self.data
    }

    pub fn prev_block_hash(&self) -> &[u8] {
        &self.prev_block_hash
    }

    pub fn hash(&self) -> &[u8] {
        &self.hash
    }

    pub fn nonce(&self) -> u64 {
        self.nonce
    }
}

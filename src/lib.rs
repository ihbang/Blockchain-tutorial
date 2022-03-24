use hex;
use sha2::{
    digest::{consts::U32, generic_array::GenericArray},
    Digest, Sha256,
};
use std::time::{Duration, SystemTime};

pub struct Block {
    timestamp: Duration,
    data: Vec<u8>,
    prev_block_hash: GenericArray<u8, U32>,
    hash: GenericArray<u8, U32>,
}

impl Block {
    pub fn new(data: &str, prev_block_hash: &GenericArray<u8, U32>) -> Block {
        let timestamp = SystemTime::now()
            .duration_since(SystemTime::UNIX_EPOCH)
            .unwrap();
        let data = Vec::from(data.as_bytes());

        let mut hasher = Sha256::new();

        let hash_data = timestamp.as_secs().to_string();
        hasher.update(hash_data.as_bytes());

        let hash_data = data.as_slice();
        hasher.update(hash_data);

        let hash_data = prev_block_hash.as_slice();
        hasher.update(hash_data);

        let hash = hasher.finalize();

        Block {
            timestamp,
            data,
            prev_block_hash: prev_block_hash.clone(),
            hash,
        }
    }

    pub fn data(&self) -> String {
        String::from_utf8_lossy(self.data.as_slice()).to_string()
    }

    pub fn prev_block_hash(&self) -> String {
        hex::encode(self.prev_block_hash.as_slice())
    }

    pub fn hash(&self) -> String {
        hex::encode(self.hash.as_slice())
    }
}

pub struct BlockChain {
    pub blocks: Vec<Block>,
}

impl BlockChain {
    pub fn new() -> BlockChain {
        let mut blocks = Vec::new();

        blocks.push(BlockChain::new_genesis_block());

        BlockChain { blocks }
    }

    pub fn add_block(&mut self, data: &str) {
        let prev_block = self.blocks.last().unwrap();
        let new_block = Block::new(data, &prev_block.hash);

        self.blocks.push(new_block);
    }

    fn new_genesis_block() -> Block {
        Block::new("Genesis Block", &Default::default())
    }
}

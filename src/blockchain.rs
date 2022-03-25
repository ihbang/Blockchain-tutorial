use crate::block::Block;

pub struct BlockChain {
    pub blocks: Vec<Block>,
}

impl BlockChain {
    pub fn new() -> Self {
        let mut blocks = Vec::new();

        blocks.push(BlockChain::new_genesis_block());

        BlockChain { blocks }
    }

    pub fn add_block(&mut self, data: &str) {
        let prev_block = self.blocks.last().unwrap();
        let new_block = Block::new(data, prev_block.hash());

        self.blocks.push(new_block);
    }

    fn new_genesis_block() -> Block {
        Block::new("Genesis Block", &[0; 32])
    }
}

use crate::block::Block;
use crate::blockchain::BLOCKS_TREE;

pub struct BlockchainIterator<'a> {
    pub current_hash: Vec<u8>,
    pub db: &'a sled::Db,
}

impl<'a> Iterator for BlockchainIterator<'a> {
    type Item = Block;

    fn next(&mut self) -> Option<Self::Item> {
        let tree = self.db.open_tree(BLOCKS_TREE).expect("db tree open failed");

        let block = tree.get(&self.current_hash).expect("block get failed");

        match block {
            Some(b) => {
                let block = Block::deserialize(b.as_ref());
                self.current_hash = block.prev_block_hash().to_vec();
                return Some(block);
            }
            None => return None,
        }
    }
}

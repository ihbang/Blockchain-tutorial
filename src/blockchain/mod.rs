use crate::block::Block;
use crate::blockchain::iterator::BlockchainIterator;

pub mod iterator;

const DB_PATH: &str = "./db";
const BLOCKS_TREE: &str = "blocks";

pub struct Blockchain {
    tip: Vec<u8>,
    db: sled::Db,
}

impl<'a> Blockchain {
    pub fn new() -> Self {
        let db = sled::open(DB_PATH).expect("db open failed");
        let tree = db.open_tree(BLOCKS_TREE).expect("db tree open failed");

        if let Ok(Some(data)) = tree.get(b"l") {
            let tip = data.as_ref().to_vec();
            Blockchain { tip, db }
        } else {
            let genesis = Blockchain::new_genesis_block();
            tree.insert(genesis.hash(), genesis.serialize())
                .expect("insert failed");
            tree.insert(b"l", genesis.hash()).expect("insert failed");

            let tip = genesis.hash().to_vec();

            Blockchain { tip, db }
        }
    }

    pub fn add_block(&mut self, data: &str) {
        let tree = self.db.open_tree(BLOCKS_TREE).expect("db tree open failed");

        let tip = tree.get(b"l").unwrap().unwrap();
        let tip = tip.as_ref();

        let new_block = Block::new(data, tip);

        tree.insert(new_block.hash(), new_block.serialize())
            .expect("insert failed");
        tree.insert(b"l", new_block.hash()).expect("insert failed");
        self.tip = new_block.hash().to_vec();
    }

    pub fn iter(&'a self) -> BlockchainIterator<'a> {
        BlockchainIterator {
            current_hash: self.tip.clone(),
            db: &self.db,
        }
    }

    fn new_genesis_block() -> Block {
        Block::new("Genesis Block", &[0; 32])
    }
}

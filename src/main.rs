use blockchain_tutorial::BlockChain;

fn main() {
    let mut bc = BlockChain::new();

    bc.add_block("Send 1 BTC to Ivan");
    bc.add_block("Send 2 more BTC to Ivan");

    for block in bc.blocks {
        println!("Prev. hash: {}", block.prev_block_hash());
        println!("Data: {}", block.data());
        println!("Hash: {}\n", block.hash());
    }
}

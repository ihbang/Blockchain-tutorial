use blockchain_tutorial::blockchain::Blockchain;
use blockchain_tutorial::proofofwork::ProofOfWork;

fn main() {
    let mut bc = Blockchain::new();

    bc.add_block("Send 1 BTC to Ivan");
    bc.add_block("Send 2 more BTC to Ivan");

    for block in bc.iter() {
        println!("Prev. hash: {}", hex::encode(block.prev_block_hash()));
        println!("Data: {}", String::from_utf8_lossy(block.data().as_slice()));
        println!("Hash: {}", hex::encode(block.hash()));

        let pow = ProofOfWork::new(&block);
        println!("PoW: {}\n", pow.validate());
    }
}

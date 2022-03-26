use clap::{arg, Command};

use crate::blockchain::Blockchain;
use crate::proofofwork::ProofOfWork;

pub fn parse() -> clap::ArgMatches {
    Command::new("rustchain")
        .about(
            "Simple blockchain implemented in Rust with reference to \"Building Blockchain in Go\"",
        )
        .subcommand_required(true)
        .arg_required_else_help(true)
        .subcommand(
            Command::new("addblock")
                .about("Add new block")
                .arg(arg!(-d --data <STRING> "Data to be added in the block"))
                .arg_required_else_help(true),
        )
        .subcommand(Command::new("printchain").about("Print all blocks in blockchain"))
        .get_matches()
}

pub fn run(args: clap::ArgMatches) {
    let mut bc = Blockchain::new();
    match args.subcommand() {
        Some(("addblock", sub_matches)) => {
            bc.add_block(sub_matches.value_of("data").unwrap());
        }
        Some(("printchain", _)) => {
            for block in bc.iter() {
                println!("Prev. hash: {}", hex::encode(block.prev_block_hash()));
                println!("Data: {}", String::from_utf8_lossy(block.data().as_slice()));
                println!("Hash: {}", hex::encode(block.hash()));

                let pow = ProofOfWork::new(&block);
                println!("PoW: {}\n", pow.validate());
            }
        }
        _ => unreachable!(),
    }
}

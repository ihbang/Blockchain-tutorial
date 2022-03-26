use blockchain_tutorial::cli;

fn main() {
    let args = cli::parse();
    cli::run(args);
}

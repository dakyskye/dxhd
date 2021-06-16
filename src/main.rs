use std::io::{self, BufRead};

fn main() {
    println!("Hello, parser 2!");


    // match io::stdin().read_line(&mut parse_line) {
    //     Ok(_) => {}
    //     Err(error) => {
    //         println!("Error reading line: {}", error)
    //     }
    // }

    // let stdin = io::stdin();
    // let mut iterator = stdin.lock().lines();
    // let parse_line = iterator.next().unwrap().unwrap();

    let args: Vec<String> = std::env::args().collect();
    let output = args.join(" ");

    // cargo run some arguments i want to pass

    // # super + ctrl + {d, n}
    // scw {"d","n"}

    println!("We'll parse: {:?}", output);
}


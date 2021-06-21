// # shift + print
// flameshot printscreen

// shift + {0-9}
// set-brightness {0-9}

// shift + -

mod parser;

fn main() {
    println!("Hello, parser 2!");

    let args: Vec<String> = std::env::args().collect();
    let output = args[1..].join(" ");

    println!("We'll parse: {:?}", output);
    
    let tokens = crate::parser::tokenizer::tokenize(&output);
    println!("Tried to tokenize: {:?}", tokens);

    let split = crate::parser::lexer::split_till_plus(&tokens);
    println!("Tried to tokenize: {:?}", split);

    let lexed = crate::parser::lexer::lex(&tokens);
    println!("Lexed:");
    println!("{:?}", lexed);
}
    

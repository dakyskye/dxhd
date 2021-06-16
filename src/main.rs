use std::option;

// # shift + print
// flameshot printscreen

// shift + {0-9}
// set-brightness {0-9}

// shift + -


fn main() {
    println!("Hello, parser 2!");

    let args: Vec<String> = std::env::args().collect();
    let output = args[1..].join(" ");

    println!("We'll parse: {:?}", output);

    let tokens = tokenize(&output);
    println!("Tried to tokenize: {:?}", tokens);

    let lexed = split_till_plus(&tokens);
    println!("Tried to tokenize: {:?}", lexed);
}
    
fn split_till_plus(vec: &Vec<Token>) -> Result<Vec<Vec<Token>>, String> {
    let mut split: Vec<Vec<Token>> = Vec::new();
    let mut option_depth: i32 = 0;
    let index: i32 = 0;

    for (idx, value) in vec.iter().enumerate() {
        match *value {
            Token::Plus => {
                if option_depth == 0 {
                    split.push(vec[0..idx].to_vec());
                    match split_till_plus(&vec[idx+1..].to_vec()) {
                        Ok(mut series) => {
                            split.append(&mut series);
                            return Ok(split)
                        }
                        Err(e) => return Err(e)
                    }
                }
            }
            _ => {}
        }
    }
    split.push(vec[..].to_vec());
    Ok(split)
}

fn lex(vec: &Vec<Token>) -> Result<Vec<LexNode>, String> {
    let mut result: Vec<LexNode> = Vec::new();

    Ok(result)
}

#[derive(Debug, Clone)]
struct LexNode {
    content: Vec<LexNode>,
    of_type: LexItem
}

#[derive(Debug, Clone)]
enum LexItem {
    Range,
    Text,
    Option
}


#[derive(Debug, Clone)]
enum Token {
    RangeStart,
    RangeSeparator,
    Comma,
    RangeEnd,
    Plus,
    Text(String),
    Whitespace
}


fn tokenize(input: &String) -> Vec<Token> {
    let mut result: Vec<Token> = Vec::new();
    let mut it = input.chars().peekable();
    let mut text = String::new();

    fn push_text(text: &mut String, vec: &mut Vec<Token>) {
        if text.len() != 0 {
            vec.push(Token::Text(String::from(&*text)));
            text.clear()
        }
    }

    while let Some(&c) = it.peek() {
        match c {
            '{' => {
                push_text(&mut text, &mut result);
                result.push(Token::RangeStart)
            }
            '}' => {
                push_text(&mut text, &mut result);
                result.push(Token::RangeEnd)
            }
            '+' => {
                push_text(&mut text, &mut result);
                result.push(Token::Plus)
            }
            ' ' => {
                push_text(&mut text, &mut result);
                result.push(Token::Whitespace)
            }
            ',' => {
                push_text(&mut text, &mut result);
                result.push(Token::Comma)
            }
            '-' => {
                push_text(&mut text, &mut result);
                result.push(Token::RangeSeparator)
            }
            a => {
                text.push(a);
            }
        }
        it.next();
    }
    if text.len() != 0 {
        result.push(Token::Text(String::from(text)));
    }
    return result
}
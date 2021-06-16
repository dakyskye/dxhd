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
    println!("Tried to Lex: {:?}", tokens);
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


fn tokenize(input: &String) -> Result<Vec<Token>, String> {
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
    return Ok(result)
}
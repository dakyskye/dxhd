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

    let split = split_till_plus(&tokens);
    println!("Tried to tokenize: {:?}", split);

    let lexed = lex(&tokens);
    println!("Lexed:");
    println!("{:?}", lexed);
}
    
fn split_till_plus(vec: &Vec<Token>) -> Result<Vec<Vec<Token>>, String> {
    let mut split: Vec<Vec<Token>> = Vec::new();
    let mut option_depth: i32 = 0;

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
            Token::RangeStart   => option_depth += 1,
            Token::RangeEnd     => option_depth -= 1,
            _ => {}
        }
    }
    split.push(vec[..].to_vec());
    Ok(split)
}

fn lex(vec: &Vec<Token>) -> Result<Vec<LexNode>, String> {
    let mut result: Vec<LexNode> = Vec::new();

    let split_result = split_till_plus(vec);

    let parts = match split_result {
        Ok(parts) => parts,
        Err(error) => return Err(error)
    };

    for (_, part) in parts.iter().enumerate() {
        match lex_part(part) {
            Ok(node) => result.push(node),
            Err(err) => return Err(err)
        }
    }

    Ok(result)
}

fn lex_part(vec: &Vec<Token>) -> Result<LexNode, String> {
    if vec.len() == 0 {
        return Err(String::from("lex_part: vector size is 0"));
    }
    let token = vec.first().unwrap();
    match token {
        Token::Text(str) => {
            let content = if vec.len() > 1 {
                    match lex_part(&vec[1..].to_vec()) {
                    Ok(content) => Some([content].to_vec()),
                    Err(err) => return Err(err)
                }
            } else {
                None
            };
            return Ok(LexNode{
                of_type: LexItem::Text(String::from(str)),
                content: content
            })
        }
        Token::RangeStart => {
            // Wrong check! Go deep down til you find range end at same depth
            // Use that instead
            if *(vec.last().unwrap()) == Token::RangeEnd {
                let closure = match lex_closure(&vec[1..vec.len()-1].to_vec()) {
                    Ok(closure) => closure,
                    Err(err) => return Err(err)
                };
                return Ok(LexNode{
                    of_type: LexItem::Closure,
                    content: Some([closure].to_vec())
                })
            } else {
                return Err(String::from("No matching ending brace (}) to a starting brace ({)"))
            }
        }
        _ => return Err(String::from("Bad expression!"))
    }
}

fn lex_closure(vec: &Vec<Token>) -> Result<LexNode, String>
{
    fn split_till_comma(vec: &Vec<Token>) -> Vec<Vec<Token>> {
        let mut split = Vec::new();
        return split;
    }
    Err(String::from("Unknown error"))
}

#[derive(Debug, Clone)]
struct LexNode {
    content: Option<Vec<LexNode>>,
    of_type: LexItem
}

#[derive(Debug, Clone)]
enum LexItem {
    Text(String),
    Closure
}


#[derive(Debug, Clone, PartialEq)]
enum Token {
    RangeStart,
    RangeSeparator,
    Comma,
    RangeEnd,
    Plus,
    Text(String),
    // Whitespace
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
                // result.push(Token::Whitespace)
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

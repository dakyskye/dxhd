use crate::parser::tokenizer::Token;

pub fn split_till_plus(vec: &Vec<Token>) -> Result<Vec<Vec<Token>>, String> {
    let mut split: Vec<Vec<Token>> = Vec::new();
    let mut option_depth: i32 = 0;

    for (idx, value) in vec.iter().enumerate() {
        match *value {
            Token::Plus => {
                if option_depth == 0 {
                    split.push(vec[0..idx].to_vec());
                    // Recursion step
                    match split_till_plus(&vec[idx+1..].to_vec()) {
                        Ok(mut series) => {
                            split.append(&mut series);
                            return Ok(split)
                        }
                        Err(e) => return Err(e)
                    }
                }
            }
            Token::OptionStart   => option_depth += 1,
            Token::OptionEnd     => option_depth -= 1,
            _ => {}
        }
    }
    split.push(vec[..].to_vec());
    Ok(split)
}

pub fn lex(vec: &Vec<Token>) -> Result<Vec<LexNode>, String> {
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
        Token::OptionStart => {
            // Wrong check! Go deep down til you find range end at same depth
            // Use that instead
            if *(vec.last().unwrap()) == Token::OptionEnd {
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

    fn split_till_comma(vec: &Vec<Token>) -> Result<Vec<Vec<Token>>, String> {
        let mut option_depth = 0;
        let mut split = Vec::new();

        for (idx, value) in vec.iter().enumerate() {
            match *value {
                Token::Comma => {
                    if option_depth == 0 {
                        split.push(vec[0..idx].to_vec());
                        match split_till_comma(&vec[idx+1..].to_vec()) {
                            Ok(mut series) => {
                                split.append(&mut series);
                                return Ok(split);
                            }
                            Err(e) => return Err(e)
                        }
                    }
                }
                Token::OptionStart   => option_depth += 1,
                Token::OptionEnd     => option_depth -= 1,
                _ => {}
            }
        }
        split.push(vec.to_vec());
        Ok(split)
    }
    Err(String::from("Unknown error"))
}

#[derive(Debug, Clone)]
pub struct LexNode {
    content: Option<Vec<LexNode>>,
    of_type: LexItem
}

#[derive(Debug, Clone)]
pub enum LexItem {
    Text(String),
    Closure
}


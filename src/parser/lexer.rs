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
                let option = match lex_option(&vec[1..vec.len()-1].to_vec()) {
                    Ok(option) => option,
                    Err(err) => return Err(err)
                };
                return Ok(option)
            } else {
                return Err(String::from("No matching ending brace (}) to a starting brace ({)"))
            }
        }
        _ => return Err(String::from("Bad expression!"))
    }
}

fn lex_option(vec: &Vec<Token>) -> Result<LexNode, String>
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

    let options = match split_till_comma(vec) {
        Ok(options) => options,
        Err(e) => return Err(e)
    };

    if options.len() == 1 {
        match try_find_range(&options[0][..]) {
            Some(range) => return Ok(range),
            None => ()
        };
    }

    Err(String::from("Unknown error"))
}

fn try_find_range(slice: &[Token]) -> Option<LexNode> {
    if slice.len() != 3 {
        return None
    }

    match slice {
        [Token::Text(a), Token::RangeSeparator, Token::Text(b)] => Some(LexNode{
            of_type: LexItem::Range,
            content: Some(vec![
                LexNode{
                    of_type: LexItem::Text(String::from(a)),
                    content: None
                },
                LexNode{
                    of_type: LexItem::Text(String::from(b)),
                    content: None
                },
            ])
        }),
        _ => None
    }
}

#[derive(Debug, Clone, PartialEq)]
pub struct LexNode {
    content: Option<Vec<LexNode>>,
    of_type: LexItem
}

#[derive(Debug, Clone, PartialEq)]
pub enum LexItem {
    Text(String),
    Option,
    Range,
}

#[cfg(test)]
mod tests {
    use crate::parser::tokenizer::tokenize;

    use super::*;

    #[test]
    fn test_text_single() {
        let nodes = lex(&tokenize(&String::from("x")));

        let expected = [LexNode{of_type: LexItem::Text(String::from("x")), content: None}];

        assert!(nodes.is_err() == false);

        let unwrapped = nodes.unwrap();
        assert_eq!(unwrapped, expected);
    }

    #[test]
    fn test_text_multiple() {
        let nodes = lex(&tokenize(&String::from("x y")));

        let expected = [
            LexNode{
                of_type: LexItem::Text(String::from("x")),
                content: Some(vec![
                    LexNode{
                        of_type: LexItem::Text(String::from("y")),
                        content: None
                    }
                ])
            }
        ];

        assert!(nodes.is_err() == false);

        let unwrapped = nodes.unwrap();
        assert_eq!(unwrapped, expected);
    }

    #[test]
    fn test_text_addition() {
        let nodes = lex(&tokenize(&String::from("x + y")));

        let expected = [
            LexNode{of_type: LexItem::Text(String::from("x")), content: None},
            LexNode{of_type: LexItem::Text(String::from("y")), content: None}
        ];

        assert!(nodes.is_err() == false);

        let unwrapped = nodes.unwrap();
        assert_eq!(unwrapped, expected);
    }

    #[test]
    fn test_option() {
        let nodes = lex(&tokenize(&String::from("{a,b,c}")));

        let expected = [
            LexNode{
                of_type: LexItem::Option,
                content: Some(vec![
                    LexNode{
                        of_type: LexItem::Text(String::from("a")),
                        content: None
                    },
                    LexNode{
                        of_type: LexItem::Text(String::from("b")),
                        content: None
                    },
                    LexNode{
                        of_type: LexItem::Text(String::from("c")),
                        content: None
                    },
                ])
            }
        ];

        assert!(nodes.is_err() == false);

        let unwrapped = nodes.unwrap();
        assert_eq!(unwrapped, expected);
    }

    #[test]
    fn test_range() {
        let nodes = lex(&tokenize(&String::from("{a-z}")));

        let expected = [
            LexNode{
                of_type: LexItem::Range,
                content: Some(vec![
                    LexNode{
                        of_type: LexItem::Text(String::from("a")),
                        content: None
                    },
                    LexNode{
                        of_type: LexItem::Text(String::from("z")),
                        content: None
                    },
                ])
            }
        ];

        assert!(nodes.is_err() == false);

        let unwrapped = nodes.unwrap();
        assert_eq!(unwrapped, expected);
    }
}
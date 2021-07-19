#[derive(Debug, Clone, PartialEq)]
pub enum Token {
    OptionStart,
    RangeSeparator,
    Comma,
    OptionEnd,
    Plus,
    Text(String),
    Delimiters,
    //Whitespace
}

pub fn tokenize(input: &String) -> Vec<Token> {
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
                result.push(Token::OptionStart)
            }
            '}' => {
                push_text(&mut text, &mut result);
                result.push(Token::OptionEnd)
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
            '\\' =>{
                push_text(&mut text, &mut result);
                result.push(Token::Delimiters)
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

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_blank() {
        let tokens = tokenize(&String::from(""));
        assert!(tokens.is_empty())
    }

    #[test]
    fn test_text() {
        let tokens = tokenize(&String::from("x"));

        let expected = [Token::Text(String::from("x"))];

        assert_eq!(tokens, expected);
    }

    #[test]
    fn test_addition() {
        let tokens = tokenize(&String::from("x + y"));

        let expected = [
            Token::Text(String::from("x")),
            Token::Plus,
            Token::Text(String::from("y"))
        ];

        assert_eq!(tokens, expected);
    }

    #[test]
    fn test_option() {
        let tokens = tokenize(&String::from("{x, y}"));

        let expected = [
            Token::OptionStart,
            Token::Text(String::from("x")),
            Token::Comma,
            Token::Text(String::from("y")),
            Token::OptionEnd
        ];

        assert_eq!(tokens, expected);
    }

    #[test]
    fn test_range() {
        let tokens = tokenize(&String::from("{0-9}"));

        let expected = [
            Token::OptionStart,
            Token::Text(String::from("0")),
            Token::RangeSeparator,
            Token::Text(String::from("9")),
            Token::OptionEnd
        ];

        assert_eq!(tokens, expected);
    }
}

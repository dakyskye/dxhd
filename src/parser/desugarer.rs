use crate::parser::lexer::LexNode;
use crate::parser::lexer::LexItem;

// The desugarer does one thing
// 1. Split ranges into options. E.g. `{a-c}` becomes `{a, b, c}`
#[allow(non_snake_case)]
pub fn desugar(vec: &Vec<LexNode>) -> Result<Vec<LexNode>, String> {
    let mut result: Vec<LexNode> = Vec::new();
    for(_, value) in vec.iter().enumerate() {
        let local:Result<LexNode, String> = match value {
            // Case we're looking for: Ranges
            LexNode { content: Some(children), of_type: LexItem::Range } => {
                match split_range(&children[..]) {
                    Ok(node) => Ok(node),
                    Err(err) => return Err(err)
                }
            },
            // Types without children: Themselves
            LexNode{
                of_type: _,
                content: None
            } => Ok(value.clone()),
            // Types with children: Themselves with children desugared.
            LexNode{
                of_type: T,
                content: Some(children)
            } => match desugar(&children) {
                Ok(nodes) => Ok(LexNode{
                    of_type: T.clone(),
                    content: Some(nodes)
                }),
                Err(err) => return Err(err)
            }
        };
        match local {
            Err(err) => return Err(err),
            Ok(vec) => result.push(vec)
        };
    }
    Ok(result)
}

fn split_range(vec: &[LexNode]) -> Result<LexNode, String>{
    match vec {
        [
            LexNode{
                of_type: LexItem::Text(a),
                content: None
            },
            LexNode{
                of_type: LexItem::Text(b),
                content: None
            },
        ] => {
            if a.len() != 1 || b.len() != 1 {
                return Err(format!("Range options must be a single character! Left: {} || Right: {}", a, b))
            } else
            {
                let start: char = a.chars().next().unwrap();
                let end: char = b.chars().next().unwrap();
                // start..=end - Inclusive on both ends due to `=`
                let mut result: Vec<LexNode> = Vec::new();
                for n in start..=end {
                    result.push(LexNode{
                        of_type: LexItem::Option,
                        content: Some(vec![
                            LexNode{
                                of_type: LexItem::Text(String::from(n)),
                                content: None
                            }
                        ])
                    })
                }
                return Ok(LexNode{
                            of_type: LexItem::OptionGroup,
                            content: Some(result)
                        });
            }
        }
        _ => return Err(format!("Invalid format for range: {:?}", vec))
    }
}

#[cfg(test)]
mod tests {
    use crate::parser::tokenizer::tokenize;
    use crate::parser::lexer::lex;
    use super::*;

    #[test]
    fn test_range_expands_to_options() {
        let result = desugar(&lex(&tokenize(&String::from("{0-3}"))).unwrap());

        let expected = vec![LexNode{
            of_type: LexItem::OptionGroup,
            content: Some(vec![
                LexNode{
                    of_type: LexItem::Option,
                    content: Some(vec![
                        LexNode{
                            of_type: LexItem::Text(String::from("0")),
                            content: None
                        }])
                },
                LexNode{
                    of_type: LexItem::Option,
                    content: Some(vec![
                        LexNode{
                            of_type: LexItem::Text(String::from("1")),
                            content: None
                        }])
                },
                LexNode{
                    of_type: LexItem::Option,
                    content: Some(vec![
                        LexNode{
                            of_type: LexItem::Text(String::from("2")),
                            content: None
                        }])
                },
                LexNode{
                    of_type: LexItem::Option,
                    content: Some(vec![
                        LexNode{
                            of_type: LexItem::Text(String::from("3")),
                            content: None
                        }])
                },
            ])
        }];

        assert_eq!(result.unwrap(), expected);
    }
}
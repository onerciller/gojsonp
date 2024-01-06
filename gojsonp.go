package gojsonp

import (
	"github.com/onerciller/gojsonp/parser"
	"github.com/onerciller/gojsonp/token"
)

// DecodeJson function to convert JSON string to map.
// It uses the tokenizer to convert the JSON string into tokens.
// It uses the parser to convert the tokens into AST nodes.
func DecodeJson(data string) (map[string]interface{}, error) {
	tokens := token.Tokenizer(data)
	return parser.AstToMap(tokens)
}

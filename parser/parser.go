package parser

import (
	"fmt"
	"github.com/onerciller/gojsonp/token"
	"strconv"
)

// AstNode struct to represent an AST node.
type AstNode struct {
	Type  token.Type
	Value interface{}
}

// Parser function to iterate over tokens and generate AST nodes.
func parser(tokens []token.Token) ([]*AstNode, error) {
	var ast []*AstNode

	for _, tk := range tokens {
		astNode, err := parseValue(tk)
		if err != nil {
			return nil, err
		}
		if astNode != nil {
			ast = append(ast, astNode)
		}

	}

	return ast, nil
}

// AstToMap function to convert AST to map.
// It iterates over the AST nodes and converts them into a map.
// It uses the key to store the key of the key-value pair.
func AstToMap(tokens []token.Token) (map[string]interface{}, error) {
	ast, err := parser(tokens)
	if err != nil {
		return nil, err
	}
	result := make(map[string]interface{})
	var key string
	for _, node := range ast {
		switch node.Type {
		case token.String:
			if key == "" {
				key = node.Value.(string)
			} else {
				result[key] = node.Value.(string)
				key = ""
			}
		case token.Number:
			if key == "" {
				key = node.Value.(string)
			} else {
				result[key] = node.Value.(float64)
				key = ""
			}
		case token.Boolean:
			if key == "" {
				key = node.Value.(string)
			} else {
				result[key] = node.Value.(bool)
				key = ""
			}
		case token.Null:
			if key == "" {
				key = node.Value.(string)
			} else {
				result[key] = nil
				key = ""
			}
		}
	}
	return result, nil
}

// parseValue converts a token to an AST node.
func parseValue(tk token.Token) (*AstNode, error) {
	switch tk.Type {
	case token.RightBrace, token.RightBracket, token.Comma, token.Colon, token.LeftBrace, token.LeftBracket, token.EOF:
		return nil, nil
	case token.String:
		return &AstNode{Type: tk.Type, Value: tk.Val}, nil
	case token.Number:
		number, err := strconv.ParseFloat(tk.Val, 64)
		if err != nil {
			return nil, err
		}
		return &AstNode{Type: tk.Type, Value: number}, nil
	case token.Boolean:
		boolean, err := strconv.ParseBool(tk.Val)
		if err != nil {
			return nil, err
		}
		return &AstNode{Type: tk.Type, Value: boolean}, nil
	case token.Null:
		return &AstNode{Type: tk.Type, Value: nil}, nil
	default:
		return nil, fmt.Errorf("unexpected token type: %s", tk.Type)
	}
}

package token

import (
	"bytes"
	"fmt"
	"unicode"
)

//Overview of the Code Structure
//Our parser is structured around the token package, comprising several key components:
//Token Types: Define the different elements in a JSON-like structure.
//Token Struct: Represents a single token with a specific type and value.
//Tokenizer Function: Converts a string into a series of tokens.
//Supporting Functions: Assist in the tokenization process.

// Type is a string alias representing different token types in the tokenizer.
type Type string

// Constants for token types. These represent different elements in a JSON-like structure.
const (
	// ILLEGAL unknown token character
	ILLEGAL Type = "ILLEGAL"

	// EOF end of file
	EOF Type = "EOF"

	// LeftBrace Structural characters in JSON
	LeftBrace    Type = "{"
	RightBrace   Type = "}"
	LeftBracket  Type = "["
	RightBracket Type = "]"
	Comma        Type = ","
	Colon        Type = ":"
	Quote        Type = "\""

	// String literals
	// Data types in JSON
	String  Type = "STRING"
	Number  Type = "NUMBER"
	Boolean Type = "BOOLEAN"
	Null    Type = "NULL"
)

// String method for Type to get the string representation of the Type.
func (t Type) String() string {
	return string(t)
}

// Token represents a single token with its Type and value.
type Token struct {

	// Type of the token (e.g., String, Number, LeftBrace)
	Type Type

	// The actual value of the token as a string
	Val string
}

// Tokenizer takes a string input and tokenizes it into a slice of Token.
// The tokens are used by the parser to build the AST
// The tokenizer is a simple state machine that iterates over the input string and returns a slice of tokens
// lexer is another name for tokenizer

// Tokenizer It is a simple state machine iterating over the input and categorizing characters into tokens.
// Example Input: `{"name": "John"}`
// Example Output: [{Type: LeftBrace, Val: "{"}, {Type: String, Val: "name"}, ...]
func Tokenizer(input []byte) []Token {
	current := 0
	var tokens []Token
	stack := NewStack()
	var prevTokenType = ILLEGAL
	for current < len(input) {
		char := input[current]

		// Determine token type based on the current character
		currentTokenType := determineTokenType(char, input, current)

		// Skip whitespace
		if unicode.IsSpace(rune(char)) {
			current++
			continue
		}

		// Handle illegal token sequences
		if !isValidSequences(prevTokenType, currentTokenType) {
			errorToken := Token{
				Type: ILLEGAL,
				Val:  fmt.Sprintf("Invalid token sequence"),
			}
			tokens = append(tokens, errorToken)
			return tokens
		}

		// Switch based on the current character to determine token type
		switch currentTokenType {

		// Example case: '{' is tokenized as {Type: LeftBrace, Val: "{"}
		case LeftBrace:
			stack.Push(LeftBrace)
			tokens = append(tokens, Token{Type: LeftBrace, Val: string(char)})
		// Example case: '}' is tokenized as {Type: RightBrace, Val: "}"}
		case RightBrace:
			if stack.Peek() == LeftBrace {
				stack.Pop()
			}
			tokens = append(tokens, Token{Type: RightBrace, Val: string(char)})
		// Example case: '[' is tokenized as {Type: LeftBracket, Val: "["}
		case LeftBracket:
			stack.Push(LeftBracket)
			tokens = append(tokens, Token{Type: LeftBracket, Val: string(char)})
		// Example case: ']' is tokenized as {Type: RightBracket, Val: "]"}
		case RightBracket:
			if stack.Peek() == LeftBracket {
				stack.Pop()
			}
			tokens = append(tokens, Token{Type: RightBracket, Val: string(char)})
		// Example case: ',' is tokenized as {Type: Comma, Val: ","}
		case Comma:
			tokens = append(tokens, Token{Type: Comma, Val: string(char)})
		// Example case: ':' is tokenized as {Type: Colon, Val: ":"}
		case Colon:
			tokens = append(tokens, Token{Type: Colon, Val: string(char)})
		// Example case: '"' is tokenized as {Type: LeftQuote, Val: '"'}
		case Quote:
			current++ // skip opening quote: '"'
			start := current

			// iterate until we find the closing quote: '"'
			for current < len(input) && input[current] != '"' {
				current++
			}

			// check quote is closed
			if current < len(input) {
				value := input[start:current]
				tokens = append(tokens, Token{Type: String, Val: string(value)})
			} else {
				tokens = append(tokens, Token{Type: ILLEGAL, Val: "Unclosed string literal"})
				return tokens
			}
		default:
			if unicode.IsDigit(rune(char)) {
				start := current

				// iterate until we find the closing quote: '"'
				for current < len(input) && isDigit(input[current]) {
					current++
				}

				// example not valid digit:
				if current != len(input) && !isTerminatingCharacter(input[current]) {
					tokens = append(tokens, Token{Type: ILLEGAL, Val: "Invalid number format"})
					return tokens
				} else {
					prevTokenType = Number
					value := input[start:current]
					tokens = append(tokens, Token{Type: Number, Val: string(value)})
				}
				continue
			} else if char == 't' || char == 'f' {
				if isBoolean(input, current) {
					var length int
					if bytes.Equal(input[current:current+4], []byte("true")) {
						length = 4
					} else {
						length = 5
					}

					prevTokenType = Boolean
					value := input[current : current+length] // true or false
					tokens = append(tokens, Token{Type: Boolean, Val: string(value)})

					current += length
					continue

				} else {
					tokens = append(tokens, Token{Type: ILLEGAL, Val: "Invalid boolean literal"})
					return tokens
				}
			} else if char == 'n' {
				if isNull(input, current) {
					if bytes.Equal(input[current:current+4], []byte("null")) {
						value := input[current : current+4] // null

						tokens = append(tokens, Token{Type: Null, Val: string(value)})
						current += 4
						prevTokenType = Null
						continue
					}
				} else {
					current++
					prevTokenType = Null
					continue
				}
			}
		}
		prevTokenType = currentTokenType
		current++
	}

	if len(stack.TokenTypes) > 0 {
		tokens = append(tokens, Token{Type: ILLEGAL, Val: "Unclosed token"})
		return tokens
	}

	tokens = append(tokens, Token{Type: EOF, Val: ""})

	return tokens
}

// isTerminatingCharacter checks if a character is a valid terminating character for a number.
// Valid terminating characters are ',', '}', ']', and EOF.
func isTerminatingCharacter(c byte) bool {
	return c == ',' || c == '}' || c == ']' || string(c) == EOF.String()
}

// isDigit checks if a byte is a digit (0-9).
func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// isBoolean checks if a substring represents a boolean value ('true' or 'false').
// Example: For input "true", it returns true.
func isBoolean(input []byte, index int) bool {
	trueLiteral := []byte("true")
	falseLiteral := []byte("false")

	if len(input)-index >= len(trueLiteral) {

		if !bytes.Equal(input[index:index+len(trueLiteral)], trueLiteral) {
			return false
		}

		afterTrueLiteral := input[index+len(trueLiteral)]
		if afterTrueLiteral == ',' || afterTrueLiteral == '}' || afterTrueLiteral == ']' {
			return true
		}
	}

	if len(input)-index >= len(falseLiteral) {

		if !bytes.Equal(input[index:index+len(falseLiteral)], falseLiteral) {
			return false
		}

		afterTrueLiteral := input[index+len(trueLiteral)]
		if afterTrueLiteral == ',' || afterTrueLiteral == '}' || afterTrueLiteral == ']' {
			return true
		}
	}

	return false
}

// isNull checks if a substring represents a null value ('null').
func isNull(input []byte, index int) bool {
	nullLiteral := []byte("null")
	if len(input)-index >= len(nullLiteral) {

		if !bytes.Equal(input[index:index+len(nullLiteral)], nullLiteral) {
			return false
		}

		afterTrueLiteral := input[index+len(nullLiteral)]
		if afterTrueLiteral == ',' || afterTrueLiteral == '}' || afterTrueLiteral == ']' {
			return true
		}
	}

	return false
}

// determineTokenType returns the type of token based on the input character
func determineTokenType(char byte, input []byte, currentIndex int) Type {
	switch char {
	case '{':
		return LeftBrace
	case '}':
		return RightBrace
	case '[':
		return LeftBracket
	case ']':
		return RightBracket
	case ',':
		return Comma
	case ':':
		return Colon
	case '"':
		return Quote // or String, if you're immediately recognizing the string token
	default:
		if unicode.IsDigit(rune(char)) {
			return Number
		} else if char == 't' || char == 'f' {
			if isBoolean(input, currentIndex) {
				return Boolean
			}
		} else if char == 'n' {
			if isNull(input, currentIndex) {
				return Null
			}
		}
	}
	return ILLEGAL
}

// Stack is a simple stack implementation for token types
type Stack struct {
	TokenTypes []Type
}

func NewStack() *Stack {
	return &Stack{}
}

// Push adds a token type to the stack
func (s *Stack) Push(t Type) {
	s.TokenTypes = append(s.TokenTypes, t)
}

// Peek returns the top token type in the stack
func (s *Stack) Peek() Type {
	return s.TokenTypes[len(s.TokenTypes)-1]
}

// Pop removes the top token type from the stack
func (s *Stack) Pop() {
	s.TokenTypes = s.TokenTypes[:len(s.TokenTypes)-1]
}

// isValidSequences checks if the current token is a valid next token for the previous token.
// For example, a colon (:) can only be followed by a string, number, boolean, null, left brace, or left bracket.
// Example: For input ":", it returns true.
func isValidSequences(prevToken, currentToken Type) bool {
	validSequences := map[Type][]Type{
		ILLEGAL:      {String, Number, Boolean, Null, LeftBrace, LeftBracket, Quote},
		LeftBrace:    {String, Number, Boolean, Null, LeftBrace, LeftBracket, RightBrace, RightBracket, Quote},
		RightBrace:   {Comma, EOF},
		LeftBracket:  {String, Number, Boolean, Null, LeftBrace, LeftBracket, RightBracket, Quote},
		RightBracket: {Comma, EOF},
		Comma:        {String, Number, Boolean, Null, LeftBrace, LeftBracket, Quote},
		Colon:        {String, Number, Boolean, Null, LeftBrace, LeftBracket, Quote},
		String:       {Comma, RightBrace, RightBracket, Colon},
		Number:       {Comma, RightBrace, RightBracket},
		Boolean:      {Comma, RightBrace, RightBracket},
		Null:         {Comma, RightBrace, RightBracket},
		Quote:        {Comma, RightBrace, RightBracket, Colon},
	}

	// Check if the currentToken is in the valid next tokens for prevToken.
	if nextTokens, ok := validSequences[prevToken]; ok {
		return ContainsInArrays(nextTokens, currentToken)
	}
	return false
}

// ContainsInArrays checks if the specified array contains the given value.
func ContainsInArrays(arr []Type, val Type) bool {
	for _, a := range arr {
		if a == val {
			return true
		}
	}
	return false
}

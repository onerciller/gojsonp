package token

import (
	"reflect"
	"strconv"
	"testing"
)

// TestTokenizer tests the Tokenizer function for various cases.
func TestTokenizer(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "Empty string",
			input: "",
			expected: []Token{
				{Type: EOF, Val: ""},
			},
		},
		{
			name:  "Braces and brackets",
			input: "{}[]",
			expected: []Token{
				{Type: LeftBrace, Val: "{"},
				{Type: RightBrace, Val: "}"},
				{Type: ILLEGAL, Val: "Invalid token sequence"},
			},
		},
		{
			name:  "Comma and colon",
			input: ",:",
			expected: []Token{
				{Type: ILLEGAL, Val: "Invalid token sequence"},
			},
		},

		{
			name:  "String literal",
			input: `{"name": "John"}`,
			expected: []Token{
				{Type: LeftBrace, Val: "{"},
				{Type: String, Val: "name"},
				{Type: Colon, Val: ":"},
				{Type: String, Val: "John"},
				{Type: RightBrace, Val: "}"},
				{Type: EOF, Val: ""},
			},
		},

		{
			name:  "string literal with empty value",
			input: `{"name": ""}`,
			expected: []Token{
				{Type: LeftBrace, Val: "{"},
				{Type: String, Val: "name"},
				{Type: Colon, Val: ":"},
				{Type: String, Val: ""},
				{Type: RightBrace, Val: "}"},
				{Type: EOF, Val: ""},
			},
		},

		{
			name:  "string literal with whitespace",
			input: `{"name": "John Smith"}`,
			expected: []Token{
				{Type: LeftBrace, Val: "{"},
				{Type: String, Val: "name"},
				{Type: Colon, Val: ":"},
				{Type: String, Val: "John Smith"},
				{Type: RightBrace, Val: "}"},
				{Type: EOF, Val: ""},
			},
		},

		{
			name:  "string literal with number",
			input: `{"name": 123}`,
			expected: []Token{
				{Type: LeftBrace, Val: "{"},
				{Type: String, Val: "name"},
				{Type: Colon, Val: ":"},
				{Type: Number, Val: strconv.Itoa(123)},
				{Type: RightBrace, Val: "}"},
				{Type: EOF, Val: ""},
			},
		},

		{
			name:  "boolean true literal invalid trues",
			input: `{"isActive": trues}`,
			expected: []Token{
				{Type: LeftBrace, Val: "{"},
				{Type: String, Val: "isActive"},
				{Type: Colon, Val: ":"},
				{Type: ILLEGAL, Val: "Invalid token sequence"},
			},
		},

		{
			name:  "boolean false literal invalid falsee",
			input: `{"isActive": falsee}`,
			expected: []Token{
				{Type: LeftBrace, Val: "{"},
				{Type: String, Val: "isActive"},
				{Type: Colon, Val: ":"},
				{Type: ILLEGAL, Val: "Invalid token sequence"},
			},
		},

		{
			name:  "null literal",
			input: `{"isActive":null}`,
			expected: []Token{
				{Type: LeftBrace, Val: "{"},
				{Type: String, Val: "isActive"},
				{Type: Colon, Val: ":"},
				{Type: Null, Val: "null"},
				{Type: RightBrace, Val: "}"},
				{Type: EOF, Val: ""},
			},
		},

		{
			name:  "null literal invalid nulls",
			input: `{"isActive": nulls}`,
			expected: []Token{
				{Type: LeftBrace, Val: "{"},
				{Type: String, Val: "isActive"},
				{Type: Colon, Val: ":"},
				{Type: ILLEGAL, Val: "Invalid token sequence"},
			},
		},

		{
			name:  "null literal invalid nul",
			input: `{"isActive": nul}`,
			expected: []Token{
				{Type: LeftBrace, Val: "{"},
				{Type: String, Val: "isActive"},
				{Type: Colon, Val: ":"},
				{Type: ILLEGAL, Val: "Invalid token sequence"},
			},
		},

		{
			name:  "Unclosed string literal",
			input: `{"name": "Alice`,
			expected: []Token{
				{Type: LeftBrace, Val: "{"},
				{Type: String, Val: "name"},
				{Type: Colon, Val: ":"},
				{Type: ILLEGAL, Val: "Unclosed string literal"},
			},
		},

		{
			name:  "Incomplete JSON structure",
			input: `{"name": "Alice", "age"`,
			expected: []Token{
				{Type: LeftBrace, Val: "{"},
				{Type: String, Val: "name"},
				{Type: Colon, Val: ":"},
				{Type: String, Val: "Alice"},
				{Type: Comma, Val: ","},
				{Type: String, Val: "age"},
				{Type: ILLEGAL, Val: "Unclosed token"},
			},
		},
		{
			name:  "Invalid Number Format",
			input: `{"age": 123abc}`,
			expected: []Token{
				{Type: LeftBrace, Val: "{"},
				{Type: String, Val: "age"},
				{Type: Colon, Val: ":"},
				{Type: ILLEGAL, Val: "Invalid number format"},
			},
		},
		{
			name:  "Misplaced comma",
			input: `{"name": "Alice",, "age": 30}`,
			expected: []Token{
				{Type: LeftBrace, Val: "{"},
				{Type: String, Val: "name"},
				{Type: Colon, Val: ":"},
				{Type: String, Val: "Alice"},
				{Type: Comma, Val: ","},
				{Type: ILLEGAL, Val: "Invalid token sequence"},
			},
		},

		{
			name:  "Misplaced colon",
			input: `{"name":: "Alice", "age": 30}`,
			expected: []Token{
				{Type: LeftBrace, Val: "{"},
				{Type: String, Val: "name"},
				{Type: Colon, Val: ":"},
				{Type: ILLEGAL, Val: "Invalid token sequence"},
			},
		},
		{
			name:  "Extra Characters After Close",
			input: `{"name": "Alice"}x`,
			expected: []Token{
				{Type: LeftBrace, Val: "{"},
				{Type: String, Val: "name"},
				{Type: Colon, Val: ":"},
				{Type: String, Val: "Alice"},
				{Type: RightBrace, Val: "}"},
				{Type: ILLEGAL, Val: "Invalid token sequence"},
			},
		},
		{
			name:  "No Colon Separator",
			input: `{"name" "Alice"}`,
			expected: []Token{
				{Type: LeftBrace, Val: "{"},
				{Type: String, Val: "name"},
				{Type: ILLEGAL, Val: "Invalid token sequence"},
			},
		},
		{
			name:  "Malformed structure",
			input: `{name: "Alice", "age": 30}`,
			expected: []Token{
				{Type: LeftBrace, Val: "{"},
				{Type: ILLEGAL, Val: "Invalid token sequence"},
			},
		},
		{
			name:  "Numbers as Strings",
			input: `{"age": "30"}`,
			expected: []Token{
				{Type: LeftBrace, Val: "{"},
				{Type: String, Val: "age"},
				{Type: Colon, Val: ":"},
				{Type: String, Val: "30"},
				{Type: RightBrace, Val: "}"},
				{Type: EOF, Val: ""},
			},
		},

		{
			name:  "Strings as Number",
			input: `{"age": 30 and thirty}`,
			expected: []Token{
				{Type: LeftBrace, Val: "{"},
				{Type: String, Val: "age"},
				{Type: Colon, Val: ":"},
				{Type: ILLEGAL, Val: "Invalid number format"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Tokenizer(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Test %s failed. Expected %#v\n, got %#v'\n", tc.name, tc.expected, result)
			}
		})
	}
}

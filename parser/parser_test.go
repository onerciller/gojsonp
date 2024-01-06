package parser

import (
	"github.com/onerciller/gojsonp/token"
	"reflect"
	"testing"
)

// TestParser tests the parser function for various token inputs.
func TestParser(t *testing.T) {
	tests := []struct {
		name    string
		tokens  []token.Token
		want    []*AstNode
		wantErr bool
	}{
		{
			name: "String token",
			tokens: []token.Token{
				{Type: token.String, Val: "hello"},
			},
			want: []*AstNode{
				{Type: token.String, Value: "hello"},
			},
			wantErr: false,
		},
		{
			name: "Number token",
			tokens: []token.Token{
				{Type: token.Number, Val: "123"},
			},
			want: []*AstNode{
				{Type: token.Number, Value: float64(123)},
			},
			wantErr: false,
		},
		{
			name: "Boolean token",
			tokens: []token.Token{
				{Type: token.Boolean, Val: "true"},
			},
			want: []*AstNode{
				{Type: token.Boolean, Value: true},
			},
			wantErr: false,
		},
		{
			name: "Null token",
			tokens: []token.Token{
				{Type: token.Null, Val: "null"},
			},
			want: []*AstNode{
				{Type: token.Null, Value: nil},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser(tt.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("parser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestAstToMap tests the AstToMap function for converting AST nodes to a map.
func TestAstToMap(t *testing.T) {
	tests := []struct {
		name    string
		tokens  []token.Token
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "Simple key-value pair",
			tokens: []token.Token{
				{Type: token.String, Val: "key"},
				{Type: token.String, Val: "value"},
			},
			want: map[string]interface{}{
				"key": "value",
			},
			wantErr: false,
		},

		{
			name: "Multiple key-value pairs",
			tokens: []token.Token{
				{Type: token.String, Val: "key1"},
				{Type: token.String, Val: "value1"},
				{Type: token.String, Val: "key2"},
				{Type: token.String, Val: "value2"},
			},
			want: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AstToMap(tt.tokens)
			if (err != nil) != tt.wantErr {
				t.Errorf("AstToMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AstToMap() got = %v, want %v", got, tt.want)
			}
		})
	}
}

package lex

import (
	"bytes"
	"testing"
)

func TestLexerNormalizeLines(t *testing.T) {
	for _, tc := range []struct{ input, want string }{
		{"", ""}, {"a", "a\n"}, {"a\n", "a\n"}, {"a\r", "a\n"},
		{"a\r\nb\rc\n", "a\nb\nc\n"}, {"\r\r\n\n", "\n\n\n"},
		{"\x00", "\uFFFD\n"}, {"a\x00\x00b\r\nc\x00\rd", "a\uFFFD\uFFFDb\nc\uFFFD\nd\n"},
		{"中文\r\n日文\x00\r", "中文\n日文\uFFFD\n"}, {"\xff\r\n\xfe\x00", "\xff\n\xfe\uFFFD\n"},
	} {
		t.Run(tc.input, func(t *testing.T) {
			lexer := NewLexer([]byte(tc.input))
			var lines [][]byte
			for line := lexer.NextLine(); line != nil; line = lexer.NextLine() {
				lines = append(lines, line)
			}
			if got := string(bytes.Join(lines, nil)); got != tc.want {
				t.Fatalf("want %q, got %q", tc.want, got)
			}
			if lexer.NextLine() != nil {
				t.Fatal("expected end of input")
			}
		})
	}
}

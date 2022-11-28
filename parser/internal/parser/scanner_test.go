package parser

import (
	"bufio"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestScanner(t *testing.T) {
	tsts := []struct {
		Name  string
		Input string
		Want  []token
	}{
		{"empty", "", nil},

		{"emptycomment", "/**/", nil},
		{"comment", "/* abc */", nil},
		{"multicomment", "/* abc\ndef */", nil},
		{"precomment", "/* abc */%", []token{{'%', yySymType{P: Position{Offset: 9, Column: 9}, PreWS: []rune("/* abc */")}}}},

		{"headingws", " %", []token{{'%', yySymType{P: Position{Offset: 1, Column: 1}, PreWS: []rune{' '}}}}},
		{"trailingws", "% ", []token{{Tok: '%'}}},

		{"operator", ":", []token{{Tok: ':'}}},
		{"slash", "/", []token{{Tok: '/'}}},
		{"rslasharrow", "/->", []token{{Tok: RSLASHARROW}}},
		{"rslashdash", "/-", []token{{Tok: RSLASHDASH}}},

		{"dash", "-", []token{{Tok: '-'}}},
		{"ldashbracket", "-[", []token{{Tok: LDASHBRACKET}}},
		{"ldashslash", "-/", []token{{Tok: LDASHSLASH}}},
		{"rarrow", "->", []token{{Tok: RARROW}}},

		{"pipe", "|", []token{{Tok: '|'}}},
		{"dpipe", "||", []token{{Tok: DPIPE}}},

		{"lt", "<", []token{{Tok: '<'}}},
		{"larrowbracket", "<-[", []token{{Tok: LARROWBRACKET}}},
		{"larrowslash", "<-/", []token{{Tok: LARROWSLASH}}},
		{"larrow", "<-", []token{{Tok: LARROW}}},
		{"lteq", "<=", []token{{Tok: LTEQ}}},
		{"ltgt", "<>", []token{{Tok: LTGT}}},

		{"gt", ">", []token{{Tok: '>'}}},
		{"gteq", ">=", []token{{Tok: GTEQ}}},

		{"rbracket", "]", []token{{Tok: ']'}}},
		{"rbracketarrow", "]->", []token{{Tok: RBRACKETARROW}}},
		{"rbracketdash", "]-", []token{{Tok: RBRACKETDASH}}},

		{"string_literal", "'abc'", []token{{STRING_LITERAL, yySymType{S: "'abc'"}}}},
		{"string_literal_escape", "'abc''def'", []token{{STRING_LITERAL, yySymType{S: "'abc''def'"}}}},
		{"string_literal_escape_end", "'abc'''", []token{{STRING_LITERAL, yySymType{S: "'abc'''"}}}},

		{"quoted_identifier", `"abc"`, []token{{QUOTED_IDENTIFIER, yySymType{S: `"abc"`}}}},
		{"quoted_identifier_escape", `"abc""def"`, []token{{QUOTED_IDENTIFIER, yySymType{S: `"abc""def"`}}}},
		{"quoted_identifier_escape_end", `"abc"""`, []token{{QUOTED_IDENTIFIER, yySymType{S: `"abc"""`}}}},

		{"dot", ".", []token{{Tok: '.'}}},
		{"dot_unsigned_decimal", ".42", []token{{UNSIGNED_DECIMAL, yySymType{S: ".42"}}}},

		{"unsigned_decimal_first", "42.", []token{{UNSIGNED_DECIMAL, yySymType{S: "42."}}}},
		{"unsigned_decimal_both", "42.4711", []token{{UNSIGNED_DECIMAL, yySymType{S: "42.4711"}}}},
		{"unsigned_integer", "42", []token{{UNSIGNED_INTEGER, yySymType{S: "42"}}}},

		{"unquoted_identifier", "abc", []token{{UNQUOTED_IDENTIFIER, yySymType{S: "abc"}}}},
		{"keyword_ucase", "CREATE", []token{{Tok: CREATE}}},
		{"keyword_lcase", "create", []token{{Tok: CREATE}}},
		{"keyword_ccase", "Create", []token{{Tok: CREATE}}},

		{"multiple", "CREATE- /", []token{{Tok: CREATE}, {'-', yySymType{P: Position{Offset: 6, Column: 6}}}, {'/', yySymType{P: Position{Offset: 8, Column: 8}, PreWS: []rune{' '}}}}},
		{"multiplelines", "CREATE\n - /*\nabc*//", []token{
			{Tok: CREATE},
			{'-', yySymType{P: Position{Offset: 8, Line: 1, Column: 1}, PreWS: []rune("\n ")}},
			{'/', yySymType{P: Position{Offset: 18, Line: 2, Column: 5}, PreWS: []rune(" /*\nabc*/")}}},
		},
	}
	for _, tst := range tsts {
		tst := tst
		t.Run(tst.Name, func(t *testing.T) {
			t.Parallel()

			l := newScanner(bufio.NewReader(strings.NewReader(tst.Input)))
			var got []token

			for {
				var tok token
				tok.Tok = l.Lex(&tok.LVal)
				if tok.Tok == bad {
					t.Fatalf("Lex failed: %v", l.errs)
				} else if tok.Tok == eof {
					break
				}
				got = append(got, tok)
			}

			if diff := cmp.Diff(tst.Want, got, cmpopts.IgnoreUnexported(yySymType{})); diff != "" {
				t.Errorf("Lex: +got, -want:\n%s", diff)
			}
		})
	}
}

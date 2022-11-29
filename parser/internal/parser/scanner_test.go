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
		Want  []testToken
	}{
		{"empty", "", nil},

		{"emptycomment", "/**/", nil},
		{"comment", "/* abc */", nil},
		{"multicomment", "/* abc\ndef */", nil},
		{"precomment", "/* abc */%", []testToken{{'%', yySymType{L: &lexValue{P: Position{Offset: 9, Column: 9}, PreWS: []rune("/* abc */")}}}}},

		{"headingws", " %", []testToken{{'%', yySymType{L: &lexValue{P: Position{Offset: 1, Column: 1}, PreWS: []rune{' '}}}}}},
		{"trailingws", "% ", []testToken{{Tok: '%'}}},

		{"operator", ":", []testToken{{Tok: ':'}}},
		{"slash", "/", []testToken{{Tok: '/'}}},
		{"rslasharrow", "/->", []testToken{{Tok: RSLASHARROW}}},
		{"rslashdash", "/-", []testToken{{Tok: RSLASHDASH}}},

		{"dash", "-", []testToken{{Tok: '-'}}},
		{"ldashbracket", "-[", []testToken{{Tok: LDASHBRACKET}}},
		{"ldashslash", "-/", []testToken{{Tok: LDASHSLASH}}},
		{"rarrow", "->", []testToken{{Tok: RARROW}}},

		{"pipe", "|", []testToken{{Tok: '|'}}},
		{"dpipe", "||", []testToken{{Tok: DPIPE}}},

		{"lt", "<", []testToken{{Tok: '<'}}},
		{"larrowbracket", "<-[", []testToken{{Tok: LARROWBRACKET}}},
		{"larrowslash", "<-/", []testToken{{Tok: LARROWSLASH}}},
		{"larrow", "<-", []testToken{{Tok: LARROW}}},
		{"lteq", "<=", []testToken{{Tok: LTEQ}}},
		{"ltgt", "<>", []testToken{{Tok: LTGT}}},

		{"gt", ">", []testToken{{Tok: '>'}}},
		{"gteq", ">=", []testToken{{Tok: GTEQ}}},

		{"rbracket", "]", []testToken{{Tok: ']'}}},
		{"rbracketarrow", "]->", []testToken{{Tok: RBRACKETARROW}}},
		{"rbracketdash", "]-", []testToken{{Tok: RBRACKETDASH}}},

		{"string_literal", "'abc'", []testToken{{STRING_LITERAL, yySymType{L: &lexValue{S: "'abc'"}}}}},
		{"string_literal_escape", "'abc''def'", []testToken{{STRING_LITERAL, yySymType{L: &lexValue{S: "'abc''def'"}}}}},
		{"string_literal_escape_end", "'abc'''", []testToken{{STRING_LITERAL, yySymType{L: &lexValue{S: "'abc'''"}}}}},

		{"quoted_identifier", `"abc"`, []testToken{{QUOTED_IDENTIFIER, yySymType{L: &lexValue{S: `"abc"`}}}}},
		{"quoted_identifier_escape", `"abc""def"`, []testToken{{QUOTED_IDENTIFIER, yySymType{L: &lexValue{S: `"abc""def"`}}}}},
		{"quoted_identifier_escape_end", `"abc"""`, []testToken{{QUOTED_IDENTIFIER, yySymType{L: &lexValue{S: `"abc"""`}}}}},

		{"dot", ".", []testToken{{Tok: '.'}}},
		{"dot_unsigned_decimal", ".42", []testToken{{UNSIGNED_DECIMAL, yySymType{L: &lexValue{S: ".42"}}}}},

		{"unsigned_decimal_first", "42.", []testToken{{UNSIGNED_DECIMAL, yySymType{L: &lexValue{S: "42."}}}}},
		{"unsigned_decimal_both", "42.4711", []testToken{{UNSIGNED_DECIMAL, yySymType{L: &lexValue{S: "42.4711"}}}}},
		{"unsigned_integer", "42", []testToken{{UNSIGNED_INTEGER, yySymType{L: &lexValue{S: "42"}}}}},

		{"unquoted_identifier", "abc", []testToken{{UNQUOTED_IDENTIFIER, yySymType{L: &lexValue{S: "abc"}}}}},
		{"keyword_ucase", "CREATE", []testToken{{Tok: CREATE}}},
		{"keyword_lcase", "create", []testToken{{Tok: CREATE}}},
		{"keyword_ccase", "Create", []testToken{{Tok: CREATE}}},

		{"multiple", "CREATE- /", []testToken{{Tok: CREATE}, {'-', yySymType{L: &lexValue{P: Position{Offset: 6, Column: 6}}}}, {'/', yySymType{L: &lexValue{P: Position{Offset: 8, Column: 8}, PreWS: []rune{' '}}}}}},
		{"multiplelines", "CREATE\n - /*\nabc*//", []testToken{
			{Tok: CREATE},
			{'-', yySymType{L: &lexValue{P: Position{Offset: 8, Line: 1, Column: 1}, PreWS: []rune("\n ")}}},
			{'/', yySymType{L: &lexValue{P: Position{Offset: 18, Line: 2, Column: 5}, PreWS: []rune(" /*\nabc*/")}}}},
		},
	}
	for _, tst := range tsts {
		tst := tst
		t.Run(tst.Name, func(t *testing.T) {
			t.Parallel()

			l := newScanner(bufio.NewReader(strings.NewReader(tst.Input)))
			var got []testToken

			for {
				var tok testToken
				tok.Tok = l.Lex(&tok.LVal)
				if tok.Tok == bad {
					t.Fatalf("Lex failed: %v", l.errs)
				} else if tok.Tok == eof {
					break
				}
				if tok.LVal.L.P == (Position{}) && tok.LVal.L.S == "" {
					// Canonicalize so we don't have to add pointless &lexValue{}s to tst.Want.
					tok.LVal.L = nil
				}
				got = append(got, tok)
			}

			if diff := cmp.Diff(tst.Want, got, cmpopts.IgnoreUnexported(yySymType{})); diff != "" {
				t.Errorf("Lex: +got, -want:\n%s", diff)
			}
		})
	}
}

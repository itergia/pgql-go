package parser

import "testing"

func TestParse(t *testing.T) {
	l := &sliceLexer{
		toks: []token{
			{SELECT, yySymType{}},
			{TRUE, yySymType{}},
			{FROM, yySymType{}},
			{MATCH, yySymType{}},
			{'(', yySymType{}},
			{')', yySymType{}},
			{';', yySymType{}},
		},
	}
	if ret := yyParse(l); ret != 0 {
		t.Fatalf("yyParse failed: %d, %s", ret, l.errs)
	}
}

type sliceLexer struct {
	toks []token
	errs []string
}

func (l *sliceLexer) Lex(lval *yySymType) int {
	if len(l.toks) == 0 {
		return eof
	}

	t := l.toks[0]
	l.toks = l.toks[1:]
	*lval = t.LVal
	return t.Tok
}

func (l *sliceLexer) Error(e string) {
	l.errs = append(l.errs, e)
}

type token struct {
	Tok  int
	LVal yySymType
}

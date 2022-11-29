// Package parser contains a goyacc-based parser for PGQL.
//
//go:generate go run golang.org/x/tools/cmd/goyacc pgql.y
package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/itergia/pgql-go/ast"
)

// Position indicates a scanner position in the input stream.
type Position struct {
	// Offset is the zero-based UTF-8 codepoint offset from the start
	// of the stream.
	Offset int

	// Line is the zero-based line number, as separated by NL.
	Line int

	// Column is the zero-based column number on the current line.
	Column int
}

const (
	// eof is the token value for EOF in goyacc.
	eof = 0

	// bad is the token value for error in goyacc.
	bad = 1
)

type RuneReader interface {
	ReadRune() (rune, int, error)
}

func init() {
	// The terse error messages are simply "syntax error" which is
	// never good enough.
	yyErrorVerbose = true
}

type Statements struct {
	Stmts []ast.Stmt
}

func Parse(r RuneReader) (*Statements, error) {
	pc := parserContext{scanner: newScanner(r)}
	var yy yyParserImpl
	if yy.Parse(&pc) != 0 || len(pc.errs) > 0 {
		// For yy.lval.P to work for reporting the faulty token, there
		// cannot be any error-recovery terms.
		errs := append(pc.errs, pc.Errors()...)
		if len(errs) == 0 {
			errs = []error{errors.New("parsing failed without further information")}
		}
		return nil, parseError{errs: errs, pos: yy.lval.L.P}
	}

	return &Statements{Stmts: pc.stmts}, nil
}

type parserContext struct {
	*scanner
	errs  []error
	stmts []ast.Stmt
}

func (pc *parserContext) Error(e string) {
	pc.errs = append(pc.errs, errors.New(e))
}

func (pc *parserContext) Stmts(ss []ast.Stmt) {
	pc.stmts = ss
}

type parseError struct {
	errs []error
	pos  Position
}

func (e parseError) Error() string {
	var sb strings.Builder
	if e.pos != (Position{}) {
		fmt.Fprintf(&sb, "at %d:%d: ", e.pos.Line+1, e.pos.Column+1)
	}
	if len(e.errs) == 0 {
		return "<nil>"
	} else if len(e.errs) == 1 {
		return sb.String() + e.errs[0].Error()
	}

	fmt.Fprintf(&sb, "%d errors", len(e.errs))
	for _, err := range e.errs {
		sb.WriteString("; ")
		sb.WriteString(err.Error())
	}

	return sb.String()
}

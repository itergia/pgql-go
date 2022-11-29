package parser

import (
	"bufio"
	"io"

	"github.com/itergia/pgql-go/parser/internal/parser"
)

type Statements = parser.Statements

// Parse parses the given UTF-8 stream as a list of PGQL statements,
// separated by semicolons. If the reader is a *bufio.Reader, it is
// used directly, otherwise a new bufio.Reader is created, which means
// the parser may read more data than needed from r.
func Parse(r io.Reader) (*Statements, error) {
	if br, ok := r.(parser.RuneReader); ok {
		return parser.Parse(br)
	}

	return parser.Parse(bufio.NewReader(r))
}

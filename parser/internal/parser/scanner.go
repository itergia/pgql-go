package parser

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

// scanner splits a Reader into tokens.
type scanner struct {
	r    RuneReader
	la   []rune
	errs []error
	pos  Position
}

// newScanner creates a new scanner using the Reader.
func newScanner(r RuneReader) *scanner {
	return &scanner{r: r}
}

// Errors returns the errors encountered by the scanner.
func (s *scanner) Errors() []error {
	return s.errs
}

// Lex implements yyLexer and returns the next token.
func (s *scanner) Lex(lval *yySymType) int {
	v := lexValue{P: s.pos}
	lval.L = &v

	for {
		switch s.peekRune() {
		case eof:
			return eof

		case ' ', '\t', '\n':
			// Optimized common cases, so we don't need IsSpace.
			v.PreWS = append(v.PreWS, s.readRune())
			v.P = s.pos

		case ':', '?', ';', ',', '{', '}', '(', ')', '=', '+', '*', '%':
			return int(s.readRune())

		case '/':
			s.readRune()
			r := s.peekRune()
			if r == bad {
				return bad
			}
			switch r {
			case '*': // Block comment.
				v.PreWS = append(v.PreWS, '/', s.readRune())
				ss, err := s.readWhile(func(r rune) bool { return r != '*' })
				if err != nil {
					s.errs = append(s.errs, err)
					return bad
				}
				v.PreWS = append(v.PreWS, []rune(ss)...)
				r := s.readRune()
				r2 := s.readRune()
				if r2 == bad {
					return bad
				} else if r2 != '/' {
					s.errs = append(s.errs, fmt.Errorf("unexpected character %q, expected %q (asterisk not allowed in comments)", r2, '/'))
					return bad
				}
				v.PreWS = append(v.PreWS, r, r2)
				v.P = s.pos

			case '-':
				s.readRune()
				r = s.peekRune()
				if r == bad {
					return bad
				}
				switch r {
				case '>':
					s.readRune()
					return RSLASHARROW

				default:
					return RSLASHDASH
				}

			default:
				return '/'
			}

		case '-':
			s.readRune()
			r := s.peekRune()
			if r == bad {
				return bad
			}
			switch r {
			case '[':
				s.readRune()
				return LDASHBRACKET

			case '/':
				s.readRune()
				return LDASHSLASH

			case '>':
				s.readRune()
				return RARROW

			default:
				return '-'
			}

		case '|':
			s.readRune()
			r := s.peekRune()
			if r == bad {
				return bad
			}
			switch r {
			case '|':
				s.readRune()
				return DPIPE

			default:
				return '|'
			}

		case '<':
			s.readRune()
			r := s.peekRune()
			if r == bad {
				return bad
			}
			switch r {
			case '-':
				s.readRune()
				r = s.peekRune()
				if r == bad {
					return bad
				}
				switch r {
				case '[':
					s.readRune()
					return LARROWBRACKET

				case '/':
					s.readRune()
					return LARROWSLASH

				default:
					return LARROW
				}

			case '=':
				s.readRune()
				return LTEQ

			case '>':
				s.readRune()
				return LTGT

			default:
				return '<'
			}

		case '>':
			s.readRune()
			r := s.peekRune()
			if r == bad {
				return bad
			}
			switch r {
			case '=':
				s.readRune()
				return GTEQ

			default:
				return '>'
			}

		case ']':
			s.readRune()
			r := s.peekRune()
			if r == bad {
				return bad
			}
			switch r {
			case '-':
				s.readRune()
				r = s.peekRune()
				if r == bad {
					return bad
				}
				switch r {
				case '>':
					s.readRune()
					return RBRACKETARROW

				default:
					return RBRACKETDASH
				}

			default:
				return ']'
			}

		case '\'':
			ss, err := s.readQuoted(s.readRune())
			if err != nil {
				s.errs = append(s.errs, err)
				return bad
			}
			v.S = ss
			return STRING_LITERAL

		case '"':
			ss, err := s.readQuoted(s.readRune())
			if err != nil {
				s.errs = append(s.errs, err)
				return bad
			}
			v.S = ss
			return QUOTED_IDENTIFIER

		case '.':
			s.readRune()
			ss, err := s.readWhile(func(r rune) bool { return unicode.IsDigit(r) })
			if err == io.EOF {
				return '.'
			} else if err != nil {
				s.errs = append(s.errs, err)
				return bad
			}
			if ss != "" {
				v.S = "." + ss
				return UNSIGNED_DECIMAL
			}
			return '.'

		default:
			r := s.readRune()
			switch {
			case unicode.IsDigit(r):
				ss, err := s.readWhile(func(r rune) bool { return unicode.IsDigit(r) })
				if err != nil {
					s.errs = append(s.errs, err)
					return bad
				}
				ss = string([]rune{r}) + ss
				switch s.peekRune() {
				case '.':
					s.readRune()
					ss2, err := s.readWhile(func(r rune) bool { return unicode.IsDigit(r) })
					if err == io.EOF {
						// Decimals may end in a period.
					} else if err != nil {
						s.errs = append(s.errs, err)
						return bad
					}
					v.S = ss + "." + ss2
					return UNSIGNED_DECIMAL

				default:
					v.S = ss
					return UNSIGNED_INTEGER
				}

			case unicode.IsLetter(r):
				ss, err := s.readWhile(func(r rune) bool { return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' })
				if err != nil {
					s.errs = append(s.errs, err)
					return bad
				}
				ss = string([]rune{r}) + ss
				tok, ok := keywords[strings.ToUpper(ss)]
				if ok {
					return tok
				}
				v.S = ss
				return UNQUOTED_IDENTIFIER

			case unicode.IsSpace(r):
				// Ignore.
				v.PreWS = append(v.PreWS, s.readRune())
				v.P = s.pos

			default:
				s.errs = append(s.errs, fmt.Errorf("unexpected character %q", r))
				return bad
			}
		}
	}
}

// readQuoted reads a quoted string or identifier until the end. The
// returned string includes the surrounding quotes.
func (s *scanner) readQuoted(quote rune) (string, error) {
	buf := string([]rune{quote})
	for {
		ss, err := s.readWhile(func(r rune) bool { return r != quote })
		if err != nil {
			return "", err
		}
		buf += ss

		r := s.readRune()
		switch s.peekRune() {
		case quote:
			buf += string([]rune{r, s.readRune()})
			// An escaped quote. Continue reading.

		default:
			buf += string([]rune{r})
			return buf, nil
		}
	}
}

// peekRune returns the next rune, without removing it.
func (s *scanner) peekRune() rune {
	if len(s.la) > 0 {
		return s.la[0]
	}

	r, _, err := s.r.ReadRune()
	if err == io.EOF {
		return eof
	} else if err != nil {
		s.errs = append(s.errs, err)
		return bad
	}
	s.la = append(s.la, r)
	return r
}

// readRune returns the next rune and removes it.
func (s *scanner) readRune() rune {
	if len(s.la) > 0 {
		r := s.la[0]
		s.la = s.la[1:]
		s.pos.Offset++
		if r == '\n' {
			s.pos.Line++
			s.pos.Column = 0
		} else {
			s.pos.Column++
		}
		return r
	}

	r, _, err := s.r.ReadRune()
	if err == io.EOF {
		return eof
	} else if err != nil {
		s.errs = append(s.errs, err)
		return bad
	}
	s.pos.Offset++
	if r == '\n' {
		s.pos.Line++
		s.pos.Column = 0
	} else {
		s.pos.Column++
	}
	return r
}

// readWhile returns the runes where f() returns true.
func (s *scanner) readWhile(f func(rune) bool) (string, error) {
	if len(s.la) > 0 {
		for i, r := range s.la {
			if !f(r) {
				ret := string(s.la[:i])
				s.la = s.la[i:]
				return ret, nil
			}
			s.pos.Offset++
			if r == '\n' {
				s.pos.Line++
				s.pos.Column = 0
			} else {
				s.pos.Column++
			}
		}
	}

	for {
		r, _, err := s.r.ReadRune()
		if err == io.EOF && len(s.la) > 0 {
			ret := string(s.la)
			s.la = nil
			return ret, nil
		} else if err != nil {
			return "", err
		}
		s.la = append(s.la, r)
		if !f(r) {
			i := len(s.la) - 1
			ret := string(s.la[:i])
			s.la = s.la[i:]
			return ret, nil
		}
		s.pos.Offset++
		if r == '\n' {
			s.pos.Line++
			s.pos.Column = 0
		} else {
			s.pos.Column++
		}
	}
}

// keywords is the list of reserved identifiers. They are case-insensitive.
var keywords = map[string]int{
	"ALL": ALL,
	"AND": AND, "OR": OR, "NOT": NOT,
	"ANY": ANY,
	"ARE": ARE,
	"AS":  AS,
	"ASC": ASC, "DESC": DESC,
	"BETWEEN": BETWEEN,
	"BY":      BY,
	"CASE":    CASE, "WHEN": WHEN, "THEN": THEN, "ELSE": ELSE, "END": END,
	"CAST":    CAST,
	"COLUMNS": COLUMNS,
	"COST":    COST,
	"COUNT":   COUNT, "MIN": MIN, "MAX": MAX, "AVG": AVG, "SUM": SUM, "ARRAY_AGG": ARRAY_AGG, "LISTAGG": LISTAGG,
	"CREATE": CREATE,
	"DATE":   DATE, "TIME": TIME, "TIMESTAMP": TIMESTAMP,
	"DESTINATION": DESTINATION,
	"DISTINCT":    DISTINCT,
	"DROP":        DROP,
	"EDGE":        EDGE,
	"EXCEPT":      EXCEPT,
	"EXISTS":      EXISTS,
	"EXTRACT":     EXTRACT,
	"FROM":        FROM,
	"GRAPH":       GRAPH,
	"GROUP":       GROUP,
	"HAVING":      HAVING,
	"IN":          IN,
	"INSERT":      INSERT, "UPDATE": UPDATE, "DELETE": DELETE,
	"INTERVAL": INTERVAL,
	"INTO":     INTO,
	"IS":       IS,
	"KEY":      KEY,
	"LABEL":    LABEL,
	"LABELS":   LABELS,
	"LIMIT":    LIMIT, "OFFSET": OFFSET,
	"MATCH":      MATCH,
	"NO":         NO,
	"NULL":       NULL,
	"ON":         ON,
	"ONE":        ONE,
	"ORDER":      ORDER,
	"PATH":       PATH,
	"PER":        PER,
	"PREFIX":     PREFIX,
	"PROPERTIES": PROPERTIES,
	"PROPERTY":   PROPERTY,
	"REFERENCES": REFERENCES,
	"ROW":        ROW,
	"SELECT":     SELECT,
	"SET":        SET,
	"SHORTEST":   SHORTEST, "CHEAPEST": CHEAPEST,
	"SOURCE": SOURCE,
	"STRING": STRING, "BOOLEAN": BOOLEAN, "INTEGER": INTEGER, "INT": INT, "LONG": LONG, "FLOAT": FLOAT, "DOUBLE": DOUBLE,
	"STEP":   STEP,
	"TABLES": TABLES,
	"TOP":    TOP,
	"TRUE":   TRUE, "FALSE": FALSE,
	"VERTEX": VERTEX,
	"WHERE":  WHERE,
	"WITH":   WITH,
	"YEAR":   YEAR, "MONTH": MONTH, "DAY": DAY, "HOUR": HOUR, "MINUTE": MINUTE, "SECOND": SECOND, "TIMEZONE_HOUR": TIMEZONE_HOUR, "TIMEZONE_MINUTE": TIMEZONE_MINUTE,
	"ZONE": ZONE,
}

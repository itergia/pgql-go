package token

import (
	"fmt"
	"strings"
)

// UnquoteIdentifier turns a quoted identifier token into an
// identifier. The function panics if the input is not surrounded by
// double-quotes.
func UnquoteIdentifier(s string) string {
	if s[0] != '"' || s[len(s)-1] != '"' {
		panic(fmt.Errorf("PGQL quoted identifier without surrounding quotes: %q", s))
	}

	return strings.Replace(s[1:len(s)-1], `""`, `"`, -1)
}

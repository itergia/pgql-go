package token

import "testing"

func TestUnquoteIdentifier(t *testing.T) {
	tsts := []struct {
		Input string
		Want  string
	}{
		{`"a"`, `a`},
		{`"abc"`, `abc`},
		{`"a"""`, `a"`},
		{`"""a"`, `"a`},
		{`"a""b"`, `a"b`},
		{`"a""""b"`, `a""b`},
	}
	for _, tst := range tsts {
		t.Run(tst.Input, func(t *testing.T) {
			got := UnquoteIdentifier(tst.Input)
			if got != tst.Want {
				t.Errorf("UnquoteIdentifier: got %q, want %q", got, tst.Want)
			}
		})
	}
}

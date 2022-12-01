package parser

import (
	"errors"

	"github.com/itergia/pgql-go/ast"
)

// reportError reports a non-nil error to the lexer. Returns zero on
// nil and one on error.
func reportError(lex yyLexer, err error) int {
	if err == nil {
		return 0
	}

	lex.Error(err.Error())
	return 1
}

// checkPathPattern validates the AllPathPattern production.
func checkPathPattern(pats []*ast.PathPattern) error {
	pat := pats[0]
	if pat.Cardinality == ast.AllCardinality && pat.Metric == ast.NoMetric {
		if e := pat.Es[0]; e.Quantity == nil || e.Quantity.Max == nil {
			return errors.New("an ALL pattern must have an upper bound quantifier")
		}
	}

	return nil
}

// checkModifyQuerySimple validates the ModifyQuery production.
func checkModifyQuerySimple(stmt *ast.ModifyStmt) error {
	if stmt.From == nil {
		if len(stmt.PathMacros) > 0 {
			return errors.New("cannot have PATH if FROM is missing")
		}
		// Mods cannot be empty, by grammar.
		if _, ok := stmt.Mods[0].(*ast.InsertClause); len(stmt.Mods) != 1 || !ok {
			return errors.New("expected a single INSERT if FROM is missing")
		}
		if stmt.Where != nil {
			return errors.New("cannot have WHERE if FROM is missing")
		}
		if len(stmt.GroupBy) > 0 {
			return errors.New("cannot have GROUP BY if FROM is missing")
		}
		if stmt.Having != nil {
			return errors.New("cannot have HAVING if FROM is missing")
		}
		if len(stmt.OrderBy) > 0 {
			return errors.New("cannot have ORDER BY if FROM is missing")
		}
		if stmt.Limit != nil {
			return errors.New("cannot have LIMIT if FROM is missing")
		}
		if stmt.Offset != nil {
			return errors.New("cannot have OFFSET if FROM is missing")
		}
	}

	return nil
}

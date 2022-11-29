package parser

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/itergia/pgql-go/ast"
)

func TestParse(t *testing.T) {
	tsts := []struct {
		Name string
		Toks []testToken
		Want []ast.Stmt
	}{
		{
			"create",
			testToks(kw(CREATE), kw(PROPERTY), kw(GRAPH), id("mygraph"), kw(VERTEX), kw(TABLES), kw('('), id("atbl"), kw(')'), kw(';')),
			[]ast.Stmt{&ast.CreateStmt{GraphName: &ast.QIdent{Names: []*ast.Ident{{Name: "mygraph"}}}, VertexTables: []*ast.VertexTableDecl{{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl"}}}}}}},
		},
		{
			"createWithEdge",
			testToks(kw(CREATE), kw(PROPERTY), kw(GRAPH), id("mygraph"), kw(VERTEX), kw(TABLES), kw('('), id("atbl"), kw(')'), kw(EDGE), kw(TABLES), kw('('), id("atbl2"), kw(SOURCE), id("atbl3"), kw(DESTINATION), id("atbl4"), kw(')'), kw(';')),
			[]ast.Stmt{&ast.CreateStmt{GraphName: &ast.QIdent{Names: []*ast.Ident{{Name: "mygraph"}}}, VertexTables: []*ast.VertexTableDecl{{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl"}}}}}, EdgeTables: []*ast.EdgeTableDecl{{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl2"}}}, Source: &ast.VertexTableRef{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl3"}}}}, Dest: &ast.VertexTableRef{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl4"}}}}}}}},
		},
		{
			"createWithSchema",
			testToks(kw(CREATE), kw(PROPERTY), kw(GRAPH), id("asch"), kw('.'), id("mygraph"), kw(VERTEX), kw(TABLES), kw('('), id("atbl"), kw(')'), kw(';')),
			[]ast.Stmt{&ast.CreateStmt{GraphName: &ast.QIdent{Names: []*ast.Ident{{Name: "asch"}, {Name: "mygraph"}}}, VertexTables: []*ast.VertexTableDecl{{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl"}}}}}}},
		},
		{
			"createWithDualVertex",
			testToks(kw(CREATE), kw(PROPERTY), kw(GRAPH), id("mygraph"), kw(VERTEX), kw(TABLES), kw('('), id("atbl"), kw(','), id("atbl2"), kw(')'), kw(';')),
			[]ast.Stmt{&ast.CreateStmt{GraphName: &ast.QIdent{Names: []*ast.Ident{{Name: "mygraph"}}}, VertexTables: []*ast.VertexTableDecl{{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl"}}}}, {TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl2"}}}}}}},
		},
		{
			"createWithDualEdge",
			testToks(kw(CREATE), kw(PROPERTY), kw(GRAPH), id("mygraph"), kw(VERTEX), kw(TABLES), kw('('), id("atbl"), kw(')'), kw(EDGE), kw(TABLES), kw('('), id("atbl2"), kw(SOURCE), id("atbl3"), kw(DESTINATION), id("atbl4"), kw(','), id("atbl5"), kw(SOURCE), id("atbl6"), kw(DESTINATION), id("atbl7"), kw(')'), kw(';')),
			[]ast.Stmt{&ast.CreateStmt{GraphName: &ast.QIdent{Names: []*ast.Ident{{Name: "mygraph"}}}, VertexTables: []*ast.VertexTableDecl{{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl"}}}}}, EdgeTables: []*ast.EdgeTableDecl{{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl2"}}}, Source: &ast.VertexTableRef{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl3"}}}}, Dest: &ast.VertexTableRef{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl4"}}}}}, {TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl5"}}}, Source: &ast.VertexTableRef{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl6"}}}}, Dest: &ast.VertexTableRef{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl7"}}}}}}}},
		},
		{
			"createWithAlias",
			testToks(kw(CREATE), kw(PROPERTY), kw(GRAPH), id("mygraph"), kw(VERTEX), kw(TABLES), kw('('), id("atbl"), kw(AS), id("atbl2"), kw(')'), kw(';')),
			[]ast.Stmt{&ast.CreateStmt{GraphName: &ast.QIdent{Names: []*ast.Ident{{Name: "mygraph"}}}, VertexTables: []*ast.VertexTableDecl{{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl"}}}, TableAlias: &ast.Ident{Name: "atbl2"}}}}},
		},
		{
			"createWithKey",
			testToks(kw(CREATE), kw(PROPERTY), kw(GRAPH), id("mygraph"), kw(VERTEX), kw(TABLES), kw('('), id("atbl"), kw(KEY), kw('('), id("acol"), kw(')'), kw(')'), kw(';')),
			[]ast.Stmt{&ast.CreateStmt{GraphName: &ast.QIdent{Names: []*ast.Ident{{Name: "mygraph"}}}, VertexTables: []*ast.VertexTableDecl{{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl"}}}, Keys: []*ast.Ident{{Name: "acol"}}}}}},
		},
		{
			"createWithLabel",
			testToks(kw(CREATE), kw(PROPERTY), kw(GRAPH), id("mygraph"), kw(VERTEX), kw(TABLES), kw('('), id("atbl"), kw(LABEL), id("albl"), kw(')'), kw(';')),
			[]ast.Stmt{&ast.CreateStmt{GraphName: &ast.QIdent{Names: []*ast.Ident{{Name: "mygraph"}}}, VertexTables: []*ast.VertexTableDecl{{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl"}}}, Label: &ast.Ident{Name: "albl"}}}}},
		},
		{
			"createWithNoProperties",
			testToks(kw(CREATE), kw(PROPERTY), kw(GRAPH), id("mygraph"), kw(VERTEX), kw(TABLES), kw('('), id("atbl"), kw(NO), kw(PROPERTIES), kw(')'), kw(';')),
			[]ast.Stmt{&ast.CreateStmt{GraphName: &ast.QIdent{Names: []*ast.Ident{{Name: "mygraph"}}}, VertexTables: []*ast.VertexTableDecl{{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl"}}}, Props: &ast.PropsClause{None: true}}}}},
		},
		{
			"createWithEdgeAlias",
			testToks(kw(CREATE), kw(PROPERTY), kw(GRAPH), id("mygraph"), kw(VERTEX), kw(TABLES), kw('('), id("atbl"), kw(')'), kw(EDGE), kw(TABLES), kw('('), id("atbl2"), kw(AS), id("atbl5"), kw(SOURCE), id("atbl3"), kw(DESTINATION), id("atbl4"), kw(')'), kw(';')),
			[]ast.Stmt{&ast.CreateStmt{GraphName: &ast.QIdent{Names: []*ast.Ident{{Name: "mygraph"}}}, VertexTables: []*ast.VertexTableDecl{{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl"}}}}}, EdgeTables: []*ast.EdgeTableDecl{{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl2"}}}, TableAlias: &ast.Ident{Name: "atbl5"}, Source: &ast.VertexTableRef{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl3"}}}}, Dest: &ast.VertexTableRef{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl4"}}}}}}}},
		},
		{
			"createWithEdgeKey",
			testToks(kw(CREATE), kw(PROPERTY), kw(GRAPH), id("mygraph"), kw(VERTEX), kw(TABLES), kw('('), id("atbl"), kw(')'), kw(EDGE), kw(TABLES), kw('('), id("atbl2"), kw(KEY), kw('('), id("acol"), kw(')'), kw(SOURCE), id("atbl3"), kw(DESTINATION), id("atbl4"), kw(')'), kw(';')),
			[]ast.Stmt{&ast.CreateStmt{GraphName: &ast.QIdent{Names: []*ast.Ident{{Name: "mygraph"}}}, VertexTables: []*ast.VertexTableDecl{{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl"}}}}}, EdgeTables: []*ast.EdgeTableDecl{{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl2"}}}, Keys: []*ast.Ident{{Name: "acol"}}, Source: &ast.VertexTableRef{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl3"}}}}, Dest: &ast.VertexTableRef{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl4"}}}}}}}},
		},
		{
			"createWithEdgeLabel",
			testToks(kw(CREATE), kw(PROPERTY), kw(GRAPH), id("mygraph"), kw(VERTEX), kw(TABLES), kw('('), id("atbl"), kw(')'), kw(EDGE), kw(TABLES), kw('('), id("atbl2"), kw(SOURCE), id("atbl3"), kw(DESTINATION), id("atbl4"), kw(LABEL), id("albl"), kw(')'), kw(';')),
			[]ast.Stmt{&ast.CreateStmt{GraphName: &ast.QIdent{Names: []*ast.Ident{{Name: "mygraph"}}}, VertexTables: []*ast.VertexTableDecl{{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl"}}}}}, EdgeTables: []*ast.EdgeTableDecl{{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl2"}}}, Source: &ast.VertexTableRef{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl3"}}}}, Dest: &ast.VertexTableRef{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl4"}}}}, Label: &ast.Ident{Name: "albl"}}}}},
		},
		{
			"createWithNoProperties",
			testToks(kw(CREATE), kw(PROPERTY), kw(GRAPH), id("mygraph"), kw(VERTEX), kw(TABLES), kw('('), id("atbl"), kw(')'), kw(EDGE), kw(TABLES), kw('('), id("atbl2"), kw(SOURCE), id("atbl3"), kw(DESTINATION), id("atbl4"), kw(NO), kw(PROPERTIES), kw(')'), kw(';')),
			[]ast.Stmt{&ast.CreateStmt{GraphName: &ast.QIdent{Names: []*ast.Ident{{Name: "mygraph"}}}, VertexTables: []*ast.VertexTableDecl{{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl"}}}}}, EdgeTables: []*ast.EdgeTableDecl{{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl2"}}}, Source: &ast.VertexTableRef{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl3"}}}}, Dest: &ast.VertexTableRef{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl4"}}}}, Props: &ast.PropsClause{None: true}}}}},
		},
		{
			"createWithEdgeSourceKey",
			testToks(kw(CREATE), kw(PROPERTY), kw(GRAPH), id("mygraph"), kw(VERTEX), kw(TABLES), kw('('), id("atbl"), kw(')'), kw(EDGE), kw(TABLES), kw('('), id("atbl2"), kw(SOURCE), kw(KEY), kw('('), id("acol"), kw(')'), kw(REFERENCES), id("atbl3"), kw('('), id("acol2"), kw(')'), kw(DESTINATION), id("atbl4"), kw(')'), kw(';')),
			[]ast.Stmt{&ast.CreateStmt{GraphName: &ast.QIdent{Names: []*ast.Ident{{Name: "mygraph"}}}, VertexTables: []*ast.VertexTableDecl{{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl"}}}}}, EdgeTables: []*ast.EdgeTableDecl{{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl2"}}}, Source: &ast.VertexTableRef{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl3"}}}, Keys: []*ast.Ident{{Name: "acol"}}, Columns: []*ast.Ident{{Name: "acol2"}}}, Dest: &ast.VertexTableRef{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl4"}}}}}}}},
		},
		{
			"createWithEdgeDestKey",
			testToks(kw(CREATE), kw(PROPERTY), kw(GRAPH), id("mygraph"), kw(VERTEX), kw(TABLES), kw('('), id("atbl"), kw(')'), kw(EDGE), kw(TABLES), kw('('), id("atbl2"), kw(SOURCE), id("atbl3"), kw(DESTINATION), kw(KEY), kw('('), id("acol"), kw(')'), kw(REFERENCES), id("atbl4"), kw('('), id("acol2"), kw(')'), kw(')'), kw(';')),
			[]ast.Stmt{&ast.CreateStmt{GraphName: &ast.QIdent{Names: []*ast.Ident{{Name: "mygraph"}}}, VertexTables: []*ast.VertexTableDecl{{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl"}}}}}, EdgeTables: []*ast.EdgeTableDecl{{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl2"}}}, Source: &ast.VertexTableRef{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl3"}}}}, Dest: &ast.VertexTableRef{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl4"}}}, Keys: []*ast.Ident{{Name: "acol"}}, Columns: []*ast.Ident{{Name: "acol2"}}}}}}},
		},
		{
			"createWithAliasWithoutAs",
			testToks(kw(CREATE), kw(PROPERTY), kw(GRAPH), id("mygraph"), kw(VERTEX), kw(TABLES), kw('('), id("atbl"), id("atbl2"), kw(')'), kw(';')),
			[]ast.Stmt{&ast.CreateStmt{GraphName: &ast.QIdent{Names: []*ast.Ident{{Name: "mygraph"}}}, VertexTables: []*ast.VertexTableDecl{{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl"}}}, TableAlias: &ast.Ident{Name: "atbl2"}}}}},
		},
		{
			"createWithKeys",
			testToks(kw(CREATE), kw(PROPERTY), kw(GRAPH), id("mygraph"), kw(VERTEX), kw(TABLES), kw('('), id("atbl"), kw(KEY), kw('('), id("acol"), kw(','), id("acol2"), kw(')'), kw(')'), kw(';')),
			[]ast.Stmt{&ast.CreateStmt{GraphName: &ast.QIdent{Names: []*ast.Ident{{Name: "mygraph"}}}, VertexTables: []*ast.VertexTableDecl{{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl"}}}, Keys: []*ast.Ident{{Name: "acol"}, {Name: "acol2"}}}}}},
		},
		{
			"createWithAllProperties",
			testToks(kw(CREATE), kw(PROPERTY), kw(GRAPH), id("mygraph"), kw(VERTEX), kw(TABLES), kw('('), id("atbl"), kw(PROPERTIES), kw(ALL), kw(COLUMNS), kw(')'), kw(';')),
			[]ast.Stmt{&ast.CreateStmt{GraphName: &ast.QIdent{Names: []*ast.Ident{{Name: "mygraph"}}}, VertexTables: []*ast.VertexTableDecl{{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl"}}}, Props: &ast.PropsClause{}}}}},
		},
		{
			"createWithAreAllProperties",
			testToks(kw(CREATE), kw(PROPERTY), kw(GRAPH), id("mygraph"), kw(VERTEX), kw(TABLES), kw('('), id("atbl"), kw(PROPERTIES), kw(ARE), kw(ALL), kw(COLUMNS), kw(')'), kw(';')),
			[]ast.Stmt{&ast.CreateStmt{GraphName: &ast.QIdent{Names: []*ast.Ident{{Name: "mygraph"}}}, VertexTables: []*ast.VertexTableDecl{{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl"}}}, Props: &ast.PropsClause{}}}}},
		},
		{
			"createWithAllPropertiesExcept",
			testToks(kw(CREATE), kw(PROPERTY), kw(GRAPH), id("mygraph"), kw(VERTEX), kw(TABLES), kw('('), id("atbl"), kw(PROPERTIES), kw(ALL), kw(COLUMNS), kw(EXCEPT), kw('('), id("acol"), kw(')'), kw(')'), kw(';')),
			[]ast.Stmt{&ast.CreateStmt{GraphName: &ast.QIdent{Names: []*ast.Ident{{Name: "mygraph"}}}, VertexTables: []*ast.VertexTableDecl{{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl"}}}, Props: &ast.PropsClause{Except: []*ast.Ident{{Name: "acol"}}}}}}},
		},
		{
			"createWithAllPropertiesExceptTwo",
			testToks(kw(CREATE), kw(PROPERTY), kw(GRAPH), id("mygraph"), kw(VERTEX), kw(TABLES), kw('('), id("atbl"), kw(PROPERTIES), kw(ALL), kw(COLUMNS), kw(EXCEPT), kw('('), id("acol"), kw(','), id("acol2"), kw(')'), kw(')'), kw(';')),
			[]ast.Stmt{&ast.CreateStmt{GraphName: &ast.QIdent{Names: []*ast.Ident{{Name: "mygraph"}}}, VertexTables: []*ast.VertexTableDecl{{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl"}}}, Props: &ast.PropsClause{Except: []*ast.Ident{{Name: "acol"}, {Name: "acol2"}}}}}}},
		},
		{
			"createWithProperty",
			testToks(kw(CREATE), kw(PROPERTY), kw(GRAPH), id("mygraph"), kw(VERTEX), kw(TABLES), kw('('), id("atbl"), kw(PROPERTIES), kw('('), id("acol"), kw(')'), kw(')'), kw(';')),
			[]ast.Stmt{&ast.CreateStmt{GraphName: &ast.QIdent{Names: []*ast.Ident{{Name: "mygraph"}}}, VertexTables: []*ast.VertexTableDecl{{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl"}}}, Props: &ast.PropsClause{Exprs: []*ast.PropExpr{{Column: &ast.Ident{Name: "acol"}}}}}}}},
		},
		{
			"createWithPropertyName",
			testToks(kw(CREATE), kw(PROPERTY), kw(GRAPH), id("mygraph"), kw(VERTEX), kw(TABLES), kw('('), id("atbl"), kw(PROPERTIES), kw('('), id("acol"), kw(AS), id("aprop"), kw(')'), kw(')'), kw(';')),
			[]ast.Stmt{&ast.CreateStmt{GraphName: &ast.QIdent{Names: []*ast.Ident{{Name: "mygraph"}}}, VertexTables: []*ast.VertexTableDecl{{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl"}}}, Props: &ast.PropsClause{Exprs: []*ast.PropExpr{{Name: &ast.Ident{Name: "aprop"}, Column: &ast.Ident{Name: "acol"}}}}}}}},
		},
		{
			"createWithTwoProperties",
			testToks(kw(CREATE), kw(PROPERTY), kw(GRAPH), id("mygraph"), kw(VERTEX), kw(TABLES), kw('('), id("atbl"), kw(PROPERTIES), kw('('), id("acol"), kw(','), id("acol2"), kw(')'), kw(')'), kw(';')),
			[]ast.Stmt{&ast.CreateStmt{GraphName: &ast.QIdent{Names: []*ast.Ident{{Name: "mygraph"}}}, VertexTables: []*ast.VertexTableDecl{{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl"}}}, Props: &ast.PropsClause{Exprs: []*ast.PropExpr{{Column: &ast.Ident{Name: "acol"}}, {Column: &ast.Ident{Name: "acol2"}}}}}}}},
		},
		// TODO: CastSpecification in PropsExpr.

		{
			"drop",
			testToks(kw(DROP), kw(PROPERTY), kw(GRAPH), id("mygraph"), kw(';')),
			[]ast.Stmt{&ast.DropStmt{GraphName: &ast.QIdent{Names: []*ast.Ident{{Name: "mygraph"}}}}},
		},

		{
			"quotedIdentifier",
			testToks(kw(DROP), kw(PROPERTY), kw(GRAPH), qid(`"my""graph"`), kw(';')),
			[]ast.Stmt{&ast.DropStmt{GraphName: &ast.QIdent{Names: []*ast.Ident{{Name: `my"graph`}}}}},
		},
	}
	for _, tst := range tsts {
		tst := tst
		t.Run(tst.Name, func(t *testing.T) {
			t.Parallel()

			l := &sliceLexer{
				toks: tst.Toks,
			}
			if ret := yyParse(l); ret != 0 || len(l.errs) > 0 {
				t.Fatalf("yyParse failed: %d, %s", ret, l.errs)
			}

			if diff := cmp.Diff(tst.Want, l.stmts); diff != "" {
				t.Errorf("yyParse: +got, -want:\n%s", diff)
			}
		})
	}
}

func testToks(toks ...testToken) []testToken {
	return toks
}

func kw(tok int) testToken {
	return testToken{Tok: tok}
}

func id(s string) testToken {
	return testToken{UNQUOTED_IDENTIFIER, yySymType{L: &lexValue{S: s}}}
}

func qid(qs string) testToken {
	return testToken{QUOTED_IDENTIFIER, yySymType{L: &lexValue{S: qs}}}
}

type sliceLexer struct {
	toks  []testToken
	errs  []string
	stmts []ast.Stmt
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

func (l *sliceLexer) Stmts(ss []ast.Stmt) {
	l.stmts = ss
}

type testToken struct {
	Tok  int
	LVal yySymType
}

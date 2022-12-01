package parser

import (
	"fmt"
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
		// Creating a Property Graph

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
		{
			"createWithPropertyName",
			testToks(kw(CREATE), kw(PROPERTY), kw(GRAPH), id("mygraph"), kw(VERTEX), kw(TABLES), kw('('), id("atbl"), kw(PROPERTIES), kw('('), id("acol"), kw(AS), id("aprop"), kw(')'), kw(')'), kw(';')),
			[]ast.Stmt{&ast.CreateStmt{GraphName: &ast.QIdent{Names: []*ast.Ident{{Name: "mygraph"}}}, VertexTables: []*ast.VertexTableDecl{{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl"}}}, Props: &ast.PropsClause{Exprs: []*ast.PropExpr{{Name: &ast.Ident{Name: "aprop"}, Column: &ast.Ident{Name: "acol"}}}}}}}},
		},
		{
			"createWithPropertyCast",
			testToks(kw(CREATE), kw(PROPERTY), kw(GRAPH), id("mygraph"), kw(VERTEX), kw(TABLES), kw('('), id("atbl"), kw(PROPERTIES), kw('('), kw(CAST), kw('('), id("acol"), kw(AS), kw(STRING), kw(')'), kw(')'), kw(')'), kw(';')),
			[]ast.Stmt{&ast.CreateStmt{GraphName: &ast.QIdent{Names: []*ast.Ident{{Name: "mygraph"}}}, VertexTables: []*ast.VertexTableDecl{{TableName: &ast.QIdent{Names: []*ast.Ident{{Name: "atbl"}}}, Props: &ast.PropsClause{Exprs: []*ast.PropExpr{{CastAs: &ast.CastExpr{Arg: &ast.Ident{Name: "acol"}, TypeKind: STRING}}}}}}}},
		},

		{
			"drop",
			testToks(kw(DROP), kw(PROPERTY), kw(GRAPH), id("mygraph"), kw(';')),
			[]ast.Stmt{&ast.DropStmt{GraphName: &ast.QIdent{Names: []*ast.Ident{{Name: "mygraph"}}}}},
		},

		// Graph Pattern Matching

		{
			"select",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}}},
		},
		{
			"selectElem",
			testToks(kw(SELECT), id("acolumn"), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{Sels: []*ast.SelectElem{{Named: &ast.NamedExpr{Expr: &ast.Ident{Name: "acolumn"}}}}, From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}}},
		},
		{
			"selectDistinct",
			testToks(kw(SELECT), kw(DISTINCT), id("acolumn"), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{Distinct: true, Sels: []*ast.SelectElem{{Named: &ast.NamedExpr{Expr: &ast.Ident{Name: "acolumn"}}}}, From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}}},
		},
		{
			"selectTwoElems",
			testToks(kw(SELECT), id("acolumn"), kw(','), id("acolumn2"), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{Sels: []*ast.SelectElem{{Named: &ast.NamedExpr{Expr: &ast.Ident{Name: "acolumn"}}}, {Named: &ast.NamedExpr{Expr: &ast.Ident{Name: "acolumn2"}}}}, From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}}},
		},
		{
			"selectElemAs",
			testToks(kw(SELECT), id("acolumn"), kw(AS), id("acolumn2"), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{Sels: []*ast.SelectElem{{Named: &ast.NamedExpr{Expr: &ast.Ident{Name: "acolumn"}, Name: &ast.Ident{Name: "acolumn2"}}}}, From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}}},
		},
		{
			"selectAllProps",
			testToks(kw(SELECT), id("atbl"), kw('.'), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{Sels: []*ast.SelectElem{{AllOf: &ast.Ident{Name: "atbl"}}}, From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}}},
		},
		{
			"selectAllPropsPrefix",
			testToks(kw(SELECT), id("atbl"), kw('.'), kw('*'), kw(PREFIX), str("aprefix"), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{Sels: []*ast.SelectElem{{AllOf: &ast.Ident{Name: "atbl"}, Prefix: strLit("aprefix")}}, From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}}},
		},
		{
			"selectTwoMatches",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(','), kw(MATCH), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}, {Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}}},
		},
		{
			"selectOn",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(ON), id("agraph"), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}, On: &ast.QIdent{Names: []*ast.Ident{{Name: "agraph"}}}}}}},
		},
		{
			"selectOneRowPerMatch",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(ONE), kw(ROW), kw(PER), kw(MATCH), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}, Rows: &ast.MatchRows{Kind: ast.OneRowPerMatch}}}}},
		},
		{
			"selectOneRowPerVertex",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(ONE), kw(ROW), kw(PER), kw(VERTEX), kw('('), id("avar"), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}, Rows: &ast.MatchRows{Kind: ast.OneRowPerVertex, Vars: []*ast.Ident{{Name: "avar"}}}}}}},
		},
		{
			"selectOneRowPerStep",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(ONE), kw(ROW), kw(PER), kw(STEP), kw('('), id("avar1"), kw(','), id("avar2"), kw(','), id("avar3"), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}, Rows: &ast.MatchRows{Kind: ast.OneRowPerStep, Vars: []*ast.Ident{{Name: "avar1"}, {Name: "avar2"}, {Name: "avar3"}}}}}}},
		},
		{
			"selectGraphPattern",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw('('), kw(')'), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}}},
		},
		{
			"selectTwoGraphPatterns",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw('('), kw(')'), kw(','), kw('('), kw(')'), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}, {Vs: []*ast.VertexPattern{{}}}}}}}},
		},
		{
			"selectPathPrimary",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(RARROW), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Es: []*ast.EdgePattern{{Dir: ast.Outgoing}}}}}}}}}},
		},
		{
			"selectTwoPathPrimary",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), id("v1"), kw(')'), kw(RARROW), kw('('), id("v2"), kw(')'), kw(LARROW), kw('('), id("v3"), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{Name: &ast.Ident{Name: "v1"}}, {Name: &ast.Ident{Name: "v2"}}, {Name: &ast.Ident{Name: "v3"}}}, Es: []*ast.PathPatternPrimary{{Es: []*ast.EdgePattern{{Dir: ast.Outgoing}}}, {Es: []*ast.EdgePattern{{Dir: ast.Incoming}}}}}}}}}},
		},
		{
			"selectReachability",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(LDASHSLASH), kw(':'), id("albl"), kw(RSLASHARROW), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Es: []*ast.EdgePattern{{LabelAlts: []*ast.Ident{{Name: "albl"}}, Dir: ast.Outgoing, Reachability: true}}}}}}}}}},
		},
		{
			"selectPathPrimaryIncoming",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(LARROW), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Es: []*ast.EdgePattern{{Dir: ast.Incoming}}}}}}}}}},
		},
		{
			"selectPathPrimaryAny",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw('-'), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Es: []*ast.EdgePattern{{Dir: ast.AnyDir}}}}}}}}}},
		},
		{
			"selectOutgoingVar",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(LDASHBRACKET), id("avar"), kw(RBRACKETARROW), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Es: []*ast.EdgePattern{{Name: &ast.Ident{Name: "avar"}, Dir: ast.Outgoing}}}}}}}}}},
		},
		{
			"selectIncomingVar",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(LARROWBRACKET), id("avar"), kw(RBRACKETDASH), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Es: []*ast.EdgePattern{{Name: &ast.Ident{Name: "avar"}, Dir: ast.Incoming}}}}}}}}}},
		},
		{
			"selectAnyDirVar",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(LDASHBRACKET), id("avar"), kw(RBRACKETDASH), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Es: []*ast.EdgePattern{{Name: &ast.Ident{Name: "avar"}, Dir: ast.AnyDir}}}}}}}}}},
		},
		{
			"selectPathPrimaryLabel",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(LDASHBRACKET), kw(':'), id("albl"), kw(RBRACKETARROW), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Es: []*ast.EdgePattern{{LabelAlts: []*ast.Ident{{Name: "albl"}}, Dir: ast.Outgoing}}}}}}}}}},
		},
		{
			"selectPathPrimaryLabelIs",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(LDASHBRACKET), kw(IS), id("albl"), kw(RBRACKETARROW), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Es: []*ast.EdgePattern{{LabelAlts: []*ast.Ident{{Name: "albl"}}, Dir: ast.Outgoing}}}}}}}}}},
		},
		{
			"selectWhere",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), off(kw(TRUE), 0), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: boolLit(true)}},
		},

		// Variable-Length Paths

		{
			"selectAnyPath",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw(ANY), kw('('), kw(')'), kw('-'), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Es: []*ast.EdgePattern{{Dir: ast.AnyDir}}}}, Cardinality: ast.AnyCardinality}}}}}},
		},
		{
			"selectAnyPathParen",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw(ANY), kw('('), kw('('), kw(')'), kw('-'), kw('('), kw(')'), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Es: []*ast.EdgePattern{{Dir: ast.AnyDir}}}}, Cardinality: ast.AnyCardinality}}}}}},
		},
		{
			"selectAnyPathQuantifier",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw(ANY), kw('('), kw(')'), kw('-'), kw('?'), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Es: []*ast.EdgePattern{{Dir: ast.AnyDir}}, Quantity: &ast.Quantifier{Max: uiLit(1)}}}, Cardinality: ast.AnyCardinality}}}}}},
		},
		{
			"selectAnyPathParenPattern",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw(ANY), kw('('), kw(')'), kw('('), kw('-'), kw(')'), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Vs: []*ast.VertexPattern{nil, nil}, Es: []*ast.EdgePattern{{Dir: ast.AnyDir}}}}, Cardinality: ast.AnyCardinality}}}}}},
		},
		{
			"selectAnyPathParenPatternVertices",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw(ANY), kw('('), kw(')'), kw('('), kw('('), kw(')'), kw('-'), kw('('), kw(')'), kw(')'), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.EdgePattern{{Dir: ast.AnyDir}}}}, Cardinality: ast.AnyCardinality}}}}}},
		},
		{
			"selectAnyPathParenPatternWhere",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw(ANY), kw('('), kw(')'), kw('('), kw('-'), kw(WHERE), off(kw(TRUE), 0), kw(')'), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Vs: []*ast.VertexPattern{nil, nil}, Es: []*ast.EdgePattern{{Dir: ast.AnyDir}}, Where: boolLit(true)}}, Cardinality: ast.AnyCardinality}}}}}},
		},
		{
			"selectAnyPathParenPatternCost",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw(ANY), kw('('), kw(')'), kw('('), kw('-'), kw(COST), ui(2), kw(')'), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Vs: []*ast.VertexPattern{nil, nil}, Es: []*ast.EdgePattern{{Dir: ast.AnyDir}}, Cost: uiLit(2)}}, Cardinality: ast.AnyCardinality}}}}}},
		},
		{
			"selectReachabilityPathPredicate",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(LDASHSLASH), kw(':'), id("albl"), kw('?'), kw(RSLASHARROW), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Es: []*ast.EdgePattern{{LabelAlts: []*ast.Ident{{Name: "albl"}}, Dir: ast.Outgoing, Reachability: true}}, Quantity: &ast.Quantifier{Max: uiLit(1)}}}}}}}}},
		},
		{
			"selectPathMacro",
			testToks(kw(PATH), id("apath"), kw(AS), kw('('), kw(')'), kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{PathMacros: []*ast.PathMacroClause{{Name: &ast.Ident{Name: "apath"}, Pattern: &ast.PathPattern{Vs: []*ast.VertexPattern{{}}}}}, From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}}},
		},
		{
			"selectTwoPathMacros",
			testToks(kw(PATH), id("apath"), kw(AS), kw('('), kw(')'), kw(PATH), id("apath2"), kw(AS), kw('('), kw(')'), kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{PathMacros: []*ast.PathMacroClause{{Name: &ast.Ident{Name: "apath"}, Pattern: &ast.PathPattern{Vs: []*ast.VertexPattern{{}}}}, {Name: &ast.Ident{Name: "apath2"}, Pattern: &ast.PathPattern{Vs: []*ast.VertexPattern{{}}}}}, From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}}},
		},
		{
			"selectPathMacroWhere",
			testToks(kw(PATH), id("apath"), kw(AS), kw('('), kw(')'), kw(WHERE), off(kw(TRUE), 0), kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{PathMacros: []*ast.PathMacroClause{{Name: &ast.Ident{Name: "apath"}, Pattern: &ast.PathPattern{Vs: []*ast.VertexPattern{{}}}, Where: boolLit(true)}}, From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}}},
		},
		{
			"selectAnyShortestPath",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw(ANY), kw(SHORTEST), kw('('), kw(')'), kw('-'), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Es: []*ast.EdgePattern{{Dir: ast.AnyDir}}}}, Cardinality: ast.AnyCardinality, Metric: ast.LengthMetric}}}}}},
		},
		{
			"selectAnyShortestPathParen",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw(ANY), kw(SHORTEST), kw('('), kw('('), kw(')'), kw('-'), kw('('), kw(')'), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Es: []*ast.EdgePattern{{Dir: ast.AnyDir}}}}, Cardinality: ast.AnyCardinality, Metric: ast.LengthMetric}}}}}},
		},
		{
			"selectAllShortestPath",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw(ALL), kw(SHORTEST), kw('('), kw(')'), kw('-'), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Es: []*ast.EdgePattern{{Dir: ast.AnyDir}}}}, Cardinality: ast.AllCardinality, Metric: ast.LengthMetric}}}}}},
		},
		{
			"selectAllShortestPathParen",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw(ALL), kw(SHORTEST), kw('('), kw('('), kw(')'), kw('-'), kw('('), kw(')'), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Es: []*ast.EdgePattern{{Dir: ast.AnyDir}}}}, Cardinality: ast.AllCardinality, Metric: ast.LengthMetric}}}}}},
		},
		{
			"selectTopShortestPath",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw(TOP), ui(2), kw(SHORTEST), kw('('), kw(')'), kw('-'), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Es: []*ast.EdgePattern{{Dir: ast.AnyDir}}}}, Cardinality: ast.TopCardinality, K: uiLit(2), Metric: ast.LengthMetric}}}}}},
		},
		{
			"selectTopShortestPathParen",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw(TOP), ui(2), kw(SHORTEST), kw('('), kw('('), kw(')'), kw('-'), kw('('), kw(')'), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Es: []*ast.EdgePattern{{Dir: ast.AnyDir}}}}, Cardinality: ast.TopCardinality, K: uiLit(2), Metric: ast.LengthMetric}}}}}},
		},
		{
			"selectAnyCheapestPath",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw(ANY), kw(CHEAPEST), kw('('), kw(')'), kw('-'), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Es: []*ast.EdgePattern{{Dir: ast.AnyDir}}}}, Cardinality: ast.AnyCardinality, Metric: ast.CostMetric}}}}}},
		},
		{
			"selectAnyCheapestPathParen",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw(ANY), kw(CHEAPEST), kw('('), kw('('), kw(')'), kw('-'), kw('('), kw(')'), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Es: []*ast.EdgePattern{{Dir: ast.AnyDir}}}}, Cardinality: ast.AnyCardinality, Metric: ast.CostMetric}}}}}},
		},
		{
			"selectTopCheapestPath",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw(TOP), ui(2), kw(CHEAPEST), kw('('), kw(')'), kw('-'), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Es: []*ast.EdgePattern{{Dir: ast.AnyDir}}}}, Cardinality: ast.TopCardinality, K: uiLit(2), Metric: ast.CostMetric}}}}}},
		},
		{
			"selectTopCheapestPathParen",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw(TOP), ui(2), kw(CHEAPEST), kw('('), kw('('), kw(')'), kw('-'), kw('('), kw(')'), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Es: []*ast.EdgePattern{{Dir: ast.AnyDir}}}}, Cardinality: ast.TopCardinality, K: uiLit(2), Metric: ast.CostMetric}}}}}},
		},
		{
			"selectAllPath",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw(ALL), kw('('), kw(')'), kw('-'), kw('?'), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Es: []*ast.EdgePattern{{Dir: ast.AnyDir}}, Quantity: &ast.Quantifier{Max: uiLit(1)}}}, Cardinality: ast.AllCardinality}}}}}},
		},
		{
			"selectAllPathParen",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw(ALL), kw('('), kw('('), kw(')'), kw('-'), kw('?'), kw('('), kw(')'), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Es: []*ast.EdgePattern{{Dir: ast.AnyDir}}, Quantity: &ast.Quantifier{Max: uiLit(1)}}}, Cardinality: ast.AllCardinality}}}}}},
		},

		// Number of Rows Per Match

		{
			"selectQuantifierZeroOrMore",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(LDASHSLASH), kw(':'), id("albl"), kw('*'), kw(RSLASHARROW), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Es: []*ast.EdgePattern{{LabelAlts: []*ast.Ident{{Name: "albl"}}, Dir: ast.Outgoing, Reachability: true}}, Quantity: &ast.Quantifier{Group: true}}}}}}}}},
		},
		{
			"selectQuantifierOneOrMore",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(LDASHSLASH), kw(':'), id("albl"), kw('+'), kw(RSLASHARROW), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Es: []*ast.EdgePattern{{LabelAlts: []*ast.Ident{{Name: "albl"}}, Dir: ast.Outgoing, Reachability: true}}, Quantity: &ast.Quantifier{Min: uiLit(1), Group: true}}}}}}}}},
		},
		{
			"selectQuantifierOptional",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(LDASHSLASH), kw(':'), id("albl"), kw('?'), kw(RSLASHARROW), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Es: []*ast.EdgePattern{{LabelAlts: []*ast.Ident{{Name: "albl"}}, Dir: ast.Outgoing, Reachability: true}}, Quantity: &ast.Quantifier{Max: uiLit(1)}}}}}}}}},
		},
		{
			"selectQuantifierExactlyN",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(LDASHSLASH), kw(':'), id("albl"), kw('{'), ui(2), kw('}'), kw(RSLASHARROW), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Es: []*ast.EdgePattern{{LabelAlts: []*ast.Ident{{Name: "albl"}}, Dir: ast.Outgoing, Reachability: true}}, Quantity: &ast.Quantifier{Min: uiLit(2), Max: uiLit(2), Group: true}}}}}}}}},
		},
		{
			"selectQuantifierNOrMore",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(LDASHSLASH), kw(':'), id("albl"), kw('{'), ui(2), kw(','), kw('}'), kw(RSLASHARROW), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Es: []*ast.EdgePattern{{LabelAlts: []*ast.Ident{{Name: "albl"}}, Dir: ast.Outgoing, Reachability: true}}, Quantity: &ast.Quantifier{Min: uiLit(2), Group: true}}}}}}}}},
		},
		{
			"selectQuantifierBetweenNAndM",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(LDASHSLASH), kw(':'), id("albl"), kw('{'), ui(2), kw(','), ui(3), kw('}'), kw(RSLASHARROW), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Es: []*ast.EdgePattern{{LabelAlts: []*ast.Ident{{Name: "albl"}}, Dir: ast.Outgoing, Reachability: true}}, Quantity: &ast.Quantifier{Min: uiLit(2), Max: uiLit(3), Group: true}}}}}}}}},
		},
		{
			"selectQuantifierBetweenZeroAndM",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(LDASHSLASH), kw(':'), id("albl"), kw('{'), kw(','), ui(2), kw('}'), kw(RSLASHARROW), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}, {}}, Es: []*ast.PathPatternPrimary{{Es: []*ast.EdgePattern{{LabelAlts: []*ast.Ident{{Name: "albl"}}, Dir: ast.Outgoing, Reachability: true}}, Quantity: &ast.Quantifier{Max: uiLit(2), Group: true}}}}}}}}},
		},

		// Grouping and Aggregation

		{
			"selectGroupBy",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(GROUP), kw(BY), id("acol"), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, GroupBy: []*ast.NamedExpr{{Expr: &ast.Ident{Name: "acol"}}}}},
		},
		{
			"selectAggCount",
			testToks(kw(SELECT), kw(COUNT), kw('('), kw('*'), kw(')'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{Sels: []*ast.SelectElem{{Named: &ast.NamedExpr{Expr: &ast.OpExpr{Op: COUNT}}}}, From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}}},
		},
		{
			"selectAggCountExpr",
			testToks(kw(SELECT), kw(COUNT), kw('('), ui(2), kw(')'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{Sels: []*ast.SelectElem{{Named: &ast.NamedExpr{Expr: &ast.OpExpr{Op: COUNT, Args: []ast.Expr{boolLit(false), uiLit(2)}}}}}, From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}}},
		},
		{
			"selectAggCountDistinct",
			testToks(kw(SELECT), kw(COUNT), kw('('), kw(DISTINCT), ui(2), kw(')'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{Sels: []*ast.SelectElem{{Named: &ast.NamedExpr{Expr: &ast.OpExpr{Op: COUNT, Args: []ast.Expr{boolLit(true), uiLit(2)}}}}}, From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}}},
		},
		{
			"selectAggMinExpr",
			testToks(kw(SELECT), kw(MIN), kw('('), ui(2), kw(')'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{Sels: []*ast.SelectElem{{Named: &ast.NamedExpr{Expr: &ast.OpExpr{Op: MIN, Args: []ast.Expr{boolLit(false), uiLit(2)}}}}}, From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}}},
		},
		{
			"selectAggMinDistinct",
			testToks(kw(SELECT), kw(MIN), kw('('), kw(DISTINCT), ui(2), kw(')'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{Sels: []*ast.SelectElem{{Named: &ast.NamedExpr{Expr: &ast.OpExpr{Op: MIN, Args: []ast.Expr{boolLit(true), uiLit(2)}}}}}, From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}}},
		},
		{
			"selectAggMaxExpr",
			testToks(kw(SELECT), kw(MAX), kw('('), ui(2), kw(')'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{Sels: []*ast.SelectElem{{Named: &ast.NamedExpr{Expr: &ast.OpExpr{Op: MAX, Args: []ast.Expr{boolLit(false), uiLit(2)}}}}}, From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}}},
		},
		{
			"selectAggMaxDistinct",
			testToks(kw(SELECT), kw(MAX), kw('('), kw(DISTINCT), ui(2), kw(')'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{Sels: []*ast.SelectElem{{Named: &ast.NamedExpr{Expr: &ast.OpExpr{Op: MAX, Args: []ast.Expr{boolLit(true), uiLit(2)}}}}}, From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}}},
		},
		{
			"selectAggAvgExpr",
			testToks(kw(SELECT), kw(AVG), kw('('), ui(2), kw(')'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{Sels: []*ast.SelectElem{{Named: &ast.NamedExpr{Expr: &ast.OpExpr{Op: AVG, Args: []ast.Expr{boolLit(false), uiLit(2)}}}}}, From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}}},
		},
		{
			"selectAggAvgDistinct",
			testToks(kw(SELECT), kw(AVG), kw('('), kw(DISTINCT), ui(2), kw(')'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{Sels: []*ast.SelectElem{{Named: &ast.NamedExpr{Expr: &ast.OpExpr{Op: AVG, Args: []ast.Expr{boolLit(true), uiLit(2)}}}}}, From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}}},
		},
		{
			"selectAggSumExpr",
			testToks(kw(SELECT), kw(SUM), kw('('), ui(2), kw(')'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{Sels: []*ast.SelectElem{{Named: &ast.NamedExpr{Expr: &ast.OpExpr{Op: SUM, Args: []ast.Expr{boolLit(false), uiLit(2)}}}}}, From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}}},
		},
		{
			"selectAggSumDistinct",
			testToks(kw(SELECT), kw(SUM), kw('('), kw(DISTINCT), ui(2), kw(')'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{Sels: []*ast.SelectElem{{Named: &ast.NamedExpr{Expr: &ast.OpExpr{Op: SUM, Args: []ast.Expr{boolLit(true), uiLit(2)}}}}}, From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}}},
		},
		{
			"selectAggArrayAggExpr",
			testToks(kw(SELECT), kw(ARRAY_AGG), kw('('), ui(2), kw(')'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{Sels: []*ast.SelectElem{{Named: &ast.NamedExpr{Expr: &ast.OpExpr{Op: ARRAY_AGG, Args: []ast.Expr{boolLit(false), uiLit(2)}}}}}, From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}}},
		},
		{
			"selectAggArrayAggDistinct",
			testToks(kw(SELECT), kw(ARRAY_AGG), kw('('), kw(DISTINCT), ui(2), kw(')'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{Sels: []*ast.SelectElem{{Named: &ast.NamedExpr{Expr: &ast.OpExpr{Op: ARRAY_AGG, Args: []ast.Expr{boolLit(true), uiLit(2)}}}}}, From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}}},
		},
		{
			"selectAggListaggExpr",
			testToks(kw(SELECT), kw(LISTAGG), kw('('), ui(2), kw(')'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{Sels: []*ast.SelectElem{{Named: &ast.NamedExpr{Expr: &ast.OpExpr{Op: LISTAGG, Args: []ast.Expr{boolLit(false), uiLit(2), nil}}}}}, From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}}},
		},
		{
			"selectAggListaggDistinct",
			testToks(kw(SELECT), kw(LISTAGG), kw('('), kw(DISTINCT), ui(2), kw(')'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{Sels: []*ast.SelectElem{{Named: &ast.NamedExpr{Expr: &ast.OpExpr{Op: LISTAGG, Args: []ast.Expr{boolLit(true), uiLit(2), nil}}}}}, From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}}},
		},
		{
			"selectAggListaggSep",
			testToks(kw(SELECT), kw(LISTAGG), kw('('), ui(2), kw(','), str("asep"), kw(')'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{Sels: []*ast.SelectElem{{Named: &ast.NamedExpr{Expr: &ast.OpExpr{Op: LISTAGG, Args: []ast.Expr{boolLit(false), uiLit(2), strLit("asep")}}}}}, From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}}},
		},
		{
			"selectHaving",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(HAVING), off(kw(TRUE), 0), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Having: boolLit(true)}},
		},

		// Sorting and Row Limiting

		{
			"selectOrderBy",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(ORDER), kw(BY), id("acol"), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, OrderBy: []*ast.OrderTerm{{Expr: &ast.Ident{Name: "acol"}, Order: ast.DefaultOrder}}}},
		},
		{
			"selectTwoOrderBy",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(ORDER), kw(BY), id("acol"), kw(','), id("acol2"), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, OrderBy: []*ast.OrderTerm{{Expr: &ast.Ident{Name: "acol"}, Order: ast.DefaultOrder}, {Expr: &ast.Ident{Name: "acol2"}, Order: ast.DefaultOrder}}}},
		},
		{
			"selectOrderByAsc",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(ORDER), kw(BY), id("acol"), kw(ASC), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, OrderBy: []*ast.OrderTerm{{Expr: &ast.Ident{Name: "acol"}, Order: ast.AscOrder}}}},
		},
		{
			"selectOrderByDesc",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(ORDER), kw(BY), id("acol"), kw(DESC), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, OrderBy: []*ast.OrderTerm{{Expr: &ast.Ident{Name: "acol"}, Order: ast.DescOrder}}}},
		},
		{
			"selectLimitOffset",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(LIMIT), ui(2), kw(OFFSET), ui(3), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Limit: uiLit(2), Offset: uiLit(3)}},
		},
		{
			"selectOffsetLimit",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(OFFSET), ui(3), kw(LIMIT), ui(2), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Limit: uiLit(2), Offset: uiLit(3)}},
		},
		{
			"selectLimit",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(LIMIT), ui(2), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Limit: uiLit(2)}},
		},
		{
			"selectOffset",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(OFFSET), ui(3), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Offset: uiLit(3)}},
		},
		{
			"selectLimitBindVar",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(LIMIT), kw('?'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Limit: &ast.BindVar{}}},
		},

		// Functions and Expressions

		{
			"exprVariableReference",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), id("avar"), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.Ident{Name: "avar"}}},
		},
		{
			"exprPropertyAccess",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), id("avar"), kw('.'), id("aprop"), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.QIdent{Names: []*ast.Ident{{Name: "avar"}, {Name: "aprop"}}}}},
		},
		{
			"exprBracketed",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), kw('('), ui(2), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: uiLit(2)}},
		},
		{
			"exprStringLiteral",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), str("astr"), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: strLit("astr")}},
		},
		{
			"exprIntLiteral",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), ui(2), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: uiLit(2)}},
		},
		{
			"exprDecLiteral",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), ud(2.5), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: udLit(2.5)}},
		},
		{
			"exprBoolLiteral",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), off(kw(FALSE), 0), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: boolLit(false)}},
		},
		{
			"exprDateLiteral",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), kw(DATE), str("adate"), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: lit(ast.DateKind, "'adate'")}},
		},
		{
			"exprTimeLiteral",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), kw(TIME), str("atime"), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: lit(ast.TimeKind, "'atime'")}},
		},
		{
			"exprTimestampLiteral",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), kw(TIMESTAMP), str("atimestamp"), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: lit(ast.TimestampKind, "'atimestamp'")}},
		},
		{
			"exprIntervalLiteral",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), kw(INTERVAL), str("aninterval"), kw(HOUR), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: lit(ast.IntervalKind, "'aninterval' HOUR")}},
		},
		{
			"exprBindVariable",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), kw('?'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.BindVar{}}},
		},
		{
			"exprUnaryMinus",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), kw('-'), ui(2), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.OpExpr{Op: '-', Args: []ast.Expr{uiLit(2)}}}},
		},
		{
			"exprStringConcat",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), ui(2), kw(DPIPE), ui(3), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.OpExpr{Op: DPIPE, Args: []ast.Expr{uiLit(2), uiLit(3)}}}},
		},
		{
			"exprMultiplication",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), ui(2), kw('*'), ui(3), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.OpExpr{Op: '*', Args: []ast.Expr{uiLit(2), uiLit(3)}}}},
		},
		{
			"exprDivision",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), ui(2), kw('/'), ui(3), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.OpExpr{Op: '/', Args: []ast.Expr{uiLit(2), uiLit(3)}}}},
		},
		{
			"exprModulo",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), ui(2), kw('%'), ui(3), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.OpExpr{Op: '%', Args: []ast.Expr{uiLit(2), uiLit(3)}}}},
		},
		{
			"exprAddition",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), ui(2), kw('+'), ui(3), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.OpExpr{Op: '+', Args: []ast.Expr{uiLit(2), uiLit(3)}}}},
		},
		{
			"exprSubtraction",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), ui(2), kw('-'), ui(3), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.OpExpr{Op: '-', Args: []ast.Expr{uiLit(2), uiLit(3)}}}},
		},
		{
			"exprEqual",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), ui(2), kw('-'), ui(3), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.OpExpr{Op: '-', Args: []ast.Expr{uiLit(2), uiLit(3)}}}},
		},
		{
			"exprNotEqual",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), ui(2), kw(LTGT), ui(3), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.OpExpr{Op: LTGT, Args: []ast.Expr{uiLit(2), uiLit(3)}}}},
		},
		{
			"exprGreater",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), ui(2), kw('>'), ui(3), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.OpExpr{Op: '>', Args: []ast.Expr{uiLit(2), uiLit(3)}}}},
		},
		{
			"exprLess",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), ui(2), kw('<'), ui(3), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.OpExpr{Op: '<', Args: []ast.Expr{uiLit(2), uiLit(3)}}}},
		},
		{
			"exprGreaterOrEqual",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), ui(2), kw(GTEQ), ui(3), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.OpExpr{Op: GTEQ, Args: []ast.Expr{uiLit(2), uiLit(3)}}}},
		},
		{
			"exprLessOrEqual",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), ui(2), kw(LTEQ), ui(3), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.OpExpr{Op: LTEQ, Args: []ast.Expr{uiLit(2), uiLit(3)}}}},
		},
		{
			"exprNot",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), kw(NOT), ui(2), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.OpExpr{Op: NOT, Args: []ast.Expr{uiLit(2)}}}},
		},
		{
			"exprAnd",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), ui(2), kw(AND), ui(3), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.OpExpr{Op: AND, Args: []ast.Expr{uiLit(2), uiLit(3)}}}},
		},
		{
			"exprOr",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), ui(2), kw(OR), ui(3), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.OpExpr{Op: OR, Args: []ast.Expr{uiLit(2), uiLit(3)}}}},
		},
		{
			"exprIsNull",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), ui(2), kw(IS), kw(NULL), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.OpExpr{Op: NULL, Args: []ast.Expr{uiLit(2)}}}},
		},
		{
			"exprIsNotNull",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), ui(2), kw(IS), kw(NOT), kw(NULL), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.OpExpr{Op: NOT_NULL, Args: []ast.Expr{uiLit(2)}}}},
		},
		{
			"exprSubstring",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), kw(SUBSTRING), kw('('), ui(2), kw(FROM), ui(3), kw(FOR), ui(4), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.OpExpr{Op: SUBSTRING, Args: []ast.Expr{uiLit(2), uiLit(3), uiLit(4)}}}},
		},
		{
			"exprSubstringNoFor",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), kw(SUBSTRING), kw('('), ui(2), kw(FROM), ui(3), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.OpExpr{Op: SUBSTRING, Args: []ast.Expr{uiLit(2), uiLit(3)}}}},
		},
		{
			"exprExtract",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), kw(EXTRACT), kw('('), kw(MINUTE), kw(FROM), ui(2), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.OpExpr{Op: EXTRACT, Args: []ast.Expr{&ast.Ident{Name: "MINUTE"}, uiLit(2)}}}},
		},
		{
			"exprFunctionInvocation",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), id("afunc"), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.CallExpr{Func: &ast.QIdent{Names: []*ast.Ident{{Name: "afunc"}}}}}},
		},
		{
			"exprFunctionInvocationQualified",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), id("anobj"), kw('.'), id("afunc"), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.CallExpr{Func: &ast.QIdent{Names: []*ast.Ident{{Name: "anobj"}, {Name: "afunc"}}}}}},
		},
		{
			"exprLabel",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), kw(LABEL), kw('('), ui(2), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.OpExpr{Op: LABEL, Args: []ast.Expr{uiLit(2)}}}},
		},
		{
			"exprLabels",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), kw(LABELS), kw('('), ui(2), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.OpExpr{Op: LABELS, Args: []ast.Expr{uiLit(2)}}}},
		},
		{
			"exprFunctionInvocationArg",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), id("afunc"), kw('('), ui(2), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.CallExpr{Func: &ast.QIdent{Names: []*ast.Ident{{Name: "afunc"}}}, Args: []ast.Expr{uiLit(2)}}}},
		},
		{
			"exprFunctionInvocationTwoArgs",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), id("afunc"), kw('('), ui(2), kw(','), ui(3), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.CallExpr{Func: &ast.QIdent{Names: []*ast.Ident{{Name: "afunc"}}}, Args: []ast.Expr{uiLit(2), uiLit(3)}}}},
		},
		{
			"exprCastExpression",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), kw(CAST), kw('('), ui(2), kw(AS), kw(STRING), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.CastExpr{Arg: uiLit(2), TypeKind: STRING}}},
		},
		{
			"exprSimpleCase",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), kw(CASE), ui(2), kw(WHEN), ui(3), kw(THEN), ui(4), kw(END), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.CaseExpr{Subject: uiLit(2), Whens: []*ast.WhenClause{{Cond: uiLit(3), Then: uiLit(4)}}}}},
		},
		{
			"exprSimpleCaseWithElse",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), kw(CASE), ui(2), kw(WHEN), ui(3), kw(THEN), ui(4), kw(ELSE), ui(5), kw(END), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.CaseExpr{Subject: uiLit(2), Whens: []*ast.WhenClause{{Cond: uiLit(3), Then: uiLit(4)}}, Else: uiLit(5)}}},
		},
		{
			"exprSearchedCase",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), kw(CASE), kw(WHEN), ui(3), kw(THEN), ui(4), kw(END), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.CaseExpr{Whens: []*ast.WhenClause{{Cond: uiLit(3), Then: uiLit(4)}}}}},
		},
		{
			"exprSearchedCaseWithElse",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), kw(CASE), kw(WHEN), ui(3), kw(THEN), ui(4), kw(ELSE), ui(5), kw(END), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.CaseExpr{Whens: []*ast.WhenClause{{Cond: uiLit(3), Then: uiLit(4)}}, Else: uiLit(5)}}},
		},
		{
			"exprSimpleCaseWithTwoWhen",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), kw(CASE), ui(2), kw(WHEN), ui(3), kw(THEN), ui(4), kw(WHEN), ui(5), kw(THEN), ui(6), kw(END), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.CaseExpr{Subject: uiLit(2), Whens: []*ast.WhenClause{{Cond: uiLit(3), Then: uiLit(4)}, {Cond: uiLit(5), Then: uiLit(6)}}}}},
		},
		{
			"exprIn",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), ui(2), kw(IN), kw('('), ui(3), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.InExpr{Subject: uiLit(2), Objects: []ast.Expr{uiLit(3)}}}},
		},
		{
			"exprNotIn",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), ui(2), kw(NOT), kw(IN), kw('('), ui(3), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.InExpr{Subject: uiLit(2), Objects: []ast.Expr{uiLit(3)}, Inv: true}}},
		},
		{
			"exprInTwoArgs",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), ui(2), kw(IN), kw('('), ui(3), kw(','), ui(4), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.InExpr{Subject: uiLit(2), Objects: []ast.Expr{uiLit(3), uiLit(4)}}}},
		},
		{
			"exprInBindVar",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), ui(2), kw(IN), kw('?'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.InExpr{Subject: uiLit(2)}}},
		},

		{
			"precMultiplicationAddition",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), ui(2), kw('*'), ui(3), kw('+'), ui(4), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.OpExpr{Op: '+', Args: []ast.Expr{&ast.OpExpr{Op: '*', Args: []ast.Expr{uiLit(2), uiLit(3)}}, uiLit(4)}}}},
		},
		{
			"precAdditionMultiplication",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), ui(2), kw('+'), ui(3), kw('*'), ui(4), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.OpExpr{Op: '+', Args: []ast.Expr{uiLit(2), &ast.OpExpr{Op: '*', Args: []ast.Expr{uiLit(3), uiLit(4)}}}}}},
		},

		// Subqueries

		{
			"exprExistsSubquery",
			testToks(kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), kw(EXISTS), kw('('), kw(SELECT), kw('*'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(')'), kw(';')),
			[]ast.Stmt{&ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: &ast.OpExpr{Op: EXISTS, Args: []ast.Expr{&ast.SubqueryExpr{Query: &ast.SelectStmt{From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}}}}}}},
		},

		// Graph Modification
		{
			"modifySimpleInsert",
			testToks(kw(INSERT), kw(VERTEX), id("avar"), kw(';')),
			[]ast.Stmt{&ast.ModifyStmt{Mods: []ast.ModClause{&ast.InsertClause{Vs: []*ast.VertexInsertion{{Var: &ast.Ident{Name: "avar"}}}}}}},
		},
		{
			"modifyFullDelete",
			testToks(kw(PATH), id("apath"), kw(AS), kw('('), kw(')'), kw(DELETE), id("avar"), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(WHERE), ui(2), kw(GROUP), kw(BY), ui(3), kw(HAVING), ui(4), kw(ORDER), kw(BY), ui(5), kw(LIMIT), ui(6), kw(OFFSET), ui(7), kw(';')),
			[]ast.Stmt{&ast.ModifyStmt{PathMacros: []*ast.PathMacroClause{{Name: &ast.Ident{Name: "apath"}, Pattern: &ast.PathPattern{Vs: []*ast.VertexPattern{{}}}}}, Mods: []ast.ModClause{&ast.DeleteClause{Vars: []*ast.Ident{{Name: "avar"}}}}, From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}, Where: uiLit(2), GroupBy: []*ast.NamedExpr{{Expr: uiLit(3)}}, Having: uiLit(4), OrderBy: []*ast.OrderTerm{{Expr: uiLit(5), Order: ast.DefaultOrder}}, Limit: uiLit(6), Offset: uiLit(7)}},
		},
		{
			"modifyTwoDeletes",
			testToks(kw(DELETE), id("avar"), kw(DELETE), id("avar2"), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.ModifyStmt{Mods: []ast.ModClause{&ast.DeleteClause{Vars: []*ast.Ident{{Name: "avar"}}}, &ast.DeleteClause{Vars: []*ast.Ident{{Name: "avar2"}}}}, From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}}},
		},
		{
			"modifySimpleInsertInto",
			testToks(kw(INSERT), kw(INTO), id("agraph"), kw(VERTEX), id("avar"), kw(';')),
			[]ast.Stmt{&ast.ModifyStmt{Mods: []ast.ModClause{&ast.InsertClause{Into: &ast.QIdent{Names: []*ast.Ident{{Name: "agraph"}}}, Vs: []*ast.VertexInsertion{{Var: &ast.Ident{Name: "avar"}}}}}}},
		},
		{
			"modifySimpleInsertVerticesAndEdges",
			testToks(kw(INSERT), kw(VERTEX), id("avar"), kw(','), kw(EDGE), id("avar2"), kw(BETWEEN), id("avar3"), kw(AND), id("avar4"), kw(','), kw(VERTEX), id("avar5"), kw(','), kw(EDGE), id("avar6"), kw(BETWEEN), id("avar7"), kw(AND), id("avar8"), kw(';')),
			[]ast.Stmt{&ast.ModifyStmt{Mods: []ast.ModClause{&ast.InsertClause{Vs: []*ast.VertexInsertion{{Var: &ast.Ident{Name: "avar"}}, {Var: &ast.Ident{Name: "avar5"}}}, Es: []*ast.EdgeInsertion{{Var: &ast.Ident{Name: "avar2"}, Source: &ast.Ident{Name: "avar3"}, Dest: &ast.Ident{Name: "avar4"}}, {Var: &ast.Ident{Name: "avar6"}, Source: &ast.Ident{Name: "avar7"}, Dest: &ast.Ident{Name: "avar8"}}}}}}},
		},
		{
			"modifySimpleInsertVertexWithLabel",
			testToks(kw(INSERT), kw(VERTEX), id("avar"), kw(LABELS), kw('('), id("albl"), kw(')'), kw(';')),
			[]ast.Stmt{&ast.ModifyStmt{Mods: []ast.ModClause{&ast.InsertClause{Vs: []*ast.VertexInsertion{{Var: &ast.Ident{Name: "avar"}, Labels: []*ast.Ident{{Name: "albl"}}}}}}}},
		},
		{
			"modifySimpleInsertVertexWithTwoLabels",
			testToks(kw(INSERT), kw(VERTEX), id("avar"), kw(LABELS), kw('('), id("albl"), kw(','), id("albl2"), kw(')'), kw(';')),
			[]ast.Stmt{&ast.ModifyStmt{Mods: []ast.ModClause{&ast.InsertClause{Vs: []*ast.VertexInsertion{{Var: &ast.Ident{Name: "avar"}, Labels: []*ast.Ident{{Name: "albl"}, {Name: "albl2"}}}}}}}},
		},
		{
			"modifySimpleInsertVertexWithProp",
			testToks(kw(INSERT), kw(VERTEX), id("avar"), kw(PROPERTIES), kw('('), id("avar2"), kw('.'), id("aprop"), kw('='), ui(2), kw(')'), kw(';')),
			[]ast.Stmt{&ast.ModifyStmt{Mods: []ast.ModClause{&ast.InsertClause{Vs: []*ast.VertexInsertion{{Var: &ast.Ident{Name: "avar"}, Props: []*ast.PropAssignment{{Prop: &ast.QIdent{Names: []*ast.Ident{{Name: "avar2"}, {Name: "aprop"}}}, Value: uiLit(2)}}}}}}}},
		},
		{
			"modifySimpleInsertVertexWithTwoProps",
			testToks(kw(INSERT), kw(VERTEX), id("avar"), kw(PROPERTIES), kw('('), id("avar2"), kw('.'), id("aprop"), kw('='), ui(2), kw(','), id("avar3"), kw('.'), id("aprop2"), kw('='), ui(3), kw(')'), kw(';')),
			[]ast.Stmt{&ast.ModifyStmt{Mods: []ast.ModClause{&ast.InsertClause{Vs: []*ast.VertexInsertion{{Var: &ast.Ident{Name: "avar"}, Props: []*ast.PropAssignment{{Prop: &ast.QIdent{Names: []*ast.Ident{{Name: "avar2"}, {Name: "aprop"}}}, Value: uiLit(2)}, {Prop: &ast.QIdent{Names: []*ast.Ident{{Name: "avar3"}, {Name: "aprop2"}}}, Value: uiLit(3)}}}}}}}},
		},
		{
			"modifySimpleInsertEdgeWithLabel",
			testToks(kw(INSERT), kw(EDGE), id("avar2"), kw(BETWEEN), id("avar3"), kw(AND), id("avar4"), kw(LABELS), kw('('), id("albl"), kw(')'), kw(';')),
			[]ast.Stmt{&ast.ModifyStmt{Mods: []ast.ModClause{&ast.InsertClause{Es: []*ast.EdgeInsertion{{Var: &ast.Ident{Name: "avar2"}, Source: &ast.Ident{Name: "avar3"}, Dest: &ast.Ident{Name: "avar4"}, Labels: []*ast.Ident{{Name: "albl"}}}}}}}},
		},
		{
			"modifyUpdate",
			testToks(kw(UPDATE), id("avar"), kw(SET), kw('('), id("avar2"), kw('.'), id("aprop"), kw('='), ui(2), kw(')'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.ModifyStmt{Mods: []ast.ModClause{&ast.UpdateClause{Updates: []*ast.Update{{Var: &ast.Ident{Name: "avar"}, Props: []*ast.PropAssignment{{Prop: &ast.QIdent{Names: []*ast.Ident{{Name: "avar2"}, {Name: "aprop"}}}, Value: uiLit(2)}}}}}}, From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}}},
		},
		{
			"modifyUpdateTwoProps",
			testToks(kw(UPDATE), id("avar"), kw(SET), kw('('), id("avar2"), kw('.'), id("aprop"), kw('='), ui(2), kw(','), id("avar3"), kw('.'), id("aprop2"), kw('='), ui(3), kw(')'), kw(FROM), kw(MATCH), kw('('), kw(')'), kw(';')),
			[]ast.Stmt{&ast.ModifyStmt{Mods: []ast.ModClause{&ast.UpdateClause{Updates: []*ast.Update{{Var: &ast.Ident{Name: "avar"}, Props: []*ast.PropAssignment{{Prop: &ast.QIdent{Names: []*ast.Ident{{Name: "avar2"}, {Name: "aprop"}}}, Value: uiLit(2)}, {Prop: &ast.QIdent{Names: []*ast.Ident{{Name: "avar3"}, {Name: "aprop2"}}}, Value: uiLit(3)}}}}}}, From: []*ast.MatchClause{{Patterns: []*ast.PathPattern{{Vs: []*ast.VertexPattern{{}}}}}}}},
		},

		// Other Syntactic Rules

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

func str(s string) testToken {
	return testToken{STRING_LITERAL, yySymType{L: &lexValue{S: "'" + s + "'"}}}
}

func ui(i uint) testToken {
	return testToken{UNSIGNED_INTEGER, yySymType{L: &lexValue{S: fmt.Sprint(i)}}}
}

func ud(v float64) testToken {
	return testToken{UNSIGNED_DECIMAL, yySymType{L: &lexValue{S: fmt.Sprint(v)}}}
}

func lit(kind ast.LitKind, s string) *ast.BasicLit {
	return &ast.BasicLit{S: s, Kind: kind}
}

func boolLit(b bool) *ast.BasicLit {
	return &ast.BasicLit{S: fmt.Sprint(b), Kind: ast.BoolKind}
}

func strLit(s string) *ast.BasicLit {
	return &ast.BasicLit{S: "'" + s + "'", Kind: ast.StringKind}
}

func uiLit(i uint) *ast.BasicLit {
	return &ast.BasicLit{S: fmt.Sprint(i), Kind: ast.UIntKind}
}

func udLit(v float64) *ast.BasicLit {
	return &ast.BasicLit{S: fmt.Sprint(v), Kind: ast.UDecKind}
}

func off(tok testToken, off int) testToken {
	tok.LVal.L = &lexValue{P: Position{Offset: off}}
	return tok
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

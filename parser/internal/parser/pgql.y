// Implements https://pgql-lang.org/spec/1.5/.
%{
package parser

import (
	"fmt"

	"github.com/itergia/pgql-go/ast"
	"github.com/itergia/pgql-go/token"
)

type yyParam interface {
	yyLexer
	Stmts([]ast.Stmt)
}

type lexValue struct {
	P Position
	PreWS []rune
	S string
}
%}

// Keywords.

%token ALL
%token ANY
%token ARE
%token AS
%token ASC DESC
%token BETWEEN
%token BY
%token CASE WHEN THEN ELSE END
%token CAST
%token COLUMNS
%token COST
%token COUNT MIN MAX AVG SUM ARRAY_AGG LISTAGG
%token CREATE
%token DATE
%token DESTINATION
%token DISTINCT
%token DROP
%token EDGE
%token EXCEPT
%token EXTRACT
%token FOR
%token FROM
%token GRAPH
%token GROUP
%token HAVING
%token INSERT UPDATE DELETE
%token INTERVAL
%token INTO
%token KEY
%token LABEL
%token LABELS
%token LIMIT OFFSET
%token MATCH
%token NO
%token NULL
%token ON
%token ONE
%token ORDER
%token PATH
%token PER
%token PREFIX
%token PROPERTIES
%token PROPERTY
%token REFERENCES
%token ROW
%token SELECT
%token SET
%token SHORTEST CHEAPEST
%token SOURCE
%token STEP
%token STRING BOOLEAN INTEGER INT LONG FLOAT DOUBLE
%token SUBSTRING
%token TABLES
%token TIME
%token TIMESTAMP
%token TIMEZONE_HOUR TIMEZONE_MINUTE
%token TOP
%token VERTEX
%token WHERE
%token WITH
%token YEAR MONTH DAY HOUR MINUTE SECOND
%token ZONE

// Operators.

%token LARROW RARROW
%token LARROWBRACKET RBRACKETARROW
%token LARROWSLASH RSLASHARROW
%token LDASHBRACKET RBRACKETDASH
%token LDASHSLASH RSLASHDASH
%token ':' '?'
%token '{' '}'
%token '(' ')'

%left     ';'
%left     ','
%left     '|'
%left     OR
%left     AND
%right    NOT
%nonassoc LTGT LTEQ GTEQ '=' '<' '>'
%left     IN
%left     '+' '-'
%left     '*' '/' '%'
%left     '.'
%left     DPIPE

%right EXISTS
%right IS
%right UMINUS

// Virtual tokens.

%token NOT_NULL
%token TIME_TZ TIMESTAMP_TZ

// Literals.

%token <L> TRUE FALSE UNSIGNED_INTEGER UNSIGNED_DECIMAL
%token <L> STRING_LITERAL UNQUOTED_IDENTIFIER QUOTED_IDENTIFIER

// Productions.

%type <Stmts> PgqlStatements PgqlStatement CreatePropertyGraph DropPropertyGraph Query SelectQuery ModifyQuery ModifyQueryFull
%type <VTables> VertexTables VertexTableList VertexTable LabelAndPropertiesClause
%type <ETables> EdgeTables OptEdgeTables EdgeTableList EdgeTable
%type <VTableRef> SourceVertexTable DestinationVertexTable
%type <Props> OptPropertiesClause PropertiesClause PropertiesAreAllColumns PropertyExpressions NoProperties
%type <PExprs> PropertyExpressionList PropertyExpression ColumnReferenceOrCastSpecification
%type <SelStmt> SelectClause
%type <PathMacros> OptPathPatternMacros PathPatternMacro
%type <SelElems> SelectElementList SelectElement AllProperties
%type <Matches> OptFromClause FromClause MatchClauseList MatchClause
%type <PPats> MatchPattern GraphPattern PathPatternList PathPattern SimplePathPattern AnyPathPattern AnyShortestPathPattern AllShortestPathPattern TopKShortestPathPattern AnyCheapestPathPattern TopKCheapestPathPattern AllPathPattern
%type <PPPats> PathPrimary QuantifiedPathPatternPrimary PathPatternPrimary ParenthesizedPathPatternExpression ReachabilityPathExpression OutgoingPathPattern IncomingPathPattern PathSpecification PathPredicate
%type <VPats> OptVertexPattern VertexPattern VariableSpecification SourceVertexPattern DestinationVertexPattern
%type <EPats> EdgePattern OutgoingEdgePattern IncomingEdgePattern AnyDirectedEdgePattern
%type <MatchRows> OptRowsPerMatch RowsPerMatch OneRowPerMatch OneRowPerVertex OneRowPerStep
%type <OTerms> OptOrderByClause OrderByClause OrderTermList OrderTerm
%type <Mods> ModificationList Modification InsertClause UpdateClause DeleteClause
%type <Insert> GraphElementInsertionList GraphElementInsertion LabelsAndProperties
%type <Updates> GraphElementUpdateList GraphElementUpdate
%type <PropAss> OptPropertiesSpecification PropertiesSpecification PropertyAssignmentList PropertyAssignment
%type <NExprs> ExpAsVarList ExpAsVar OptGroupByClause GroupByClause
%type <Exprs> OptLimitOffsetClauses LimitOffsetClauses OptArgumentList ArgumentList InValueList ValueExpressionList
%type <Expr> OptWhereClause WhereClause Aggregation CountAggregation MinAggregation MaxAggregation AvgAggregation SumAggregation ArrayAggregation ListaggAggregation OptListaggSeparator ListaggSeparator OptHavingClause HavingClause LimitClause OffsetClause LimitOffsetValue ValueExpression OptCostClause CostClause BracketedValueExpression BindVariable ArithmeticExpression UnaryMinus Multiplication Division Modulo Addition Subtraction RelationalExpression Equal NotEqual Greater Less GreaterOrEqual LessOrEqual LogicalExpression Not And Or StringConcat IsNullPredicate IsNotNullPredicate CharacterSubstring StartPosition StringLength ExtractFunction FunctionInvocation CastSpecification CaseExpression SimpleCase SearchedCase OptElseClause ElseClause InPredicate NotInPredicate ExistsPredicate Subquery ScalarSubquery
%type <Whens> WhenClauseList WhenClause
%type <Quant> OptGraphPatternQuantifier GraphPatternQuantifier ZeroOrMore OneOrMore Optional ExactlyN NOrMore BetweenNAndM BetweenZeroAndM
%type <QIdent> GraphName TableName SchemaQualifiedName SchemaIdentifierPart OptOnClause OnClause PropertyAccess OptIntoClause IntoClause
%type <Ident> Identifier TableAlias OptTableAlias ColumnName LabelClause OptLabelClause Label ColumnReference PropertyName ColumnReference OptVariableName VariableName VertexVariable VertexVariable1 VertexVariable2 EdgeVariable VariableReference ExtractField FunctionName VertexReference
%type <Idents> OptKeyClause KeyClause ColumnNameList OptExceptColumns ExceptColumns ColumnReferenceList OptLabelPredicate LabelPredicate LabelList LabelAlt OptLabelSpecification LabelSpecification VariableReferenceList
%type <BLit> AllPropertiesPrefix KValue Literal StringLiteral NumericLiteral BooleanLiteral DateLiteral TimeLiteral TimestampLiteral IntervalLiteral DateTimeField
%type <B> OptDistinct
%type <I> DataType

%union {
  L *lexValue

  Stmts []ast.Stmt
  VTables []*ast.VertexTableDecl
  ETables []*ast.EdgeTableDecl
  VTableRef *ast.VertexTableRef
  Props *ast.PropsClause
  PExprs []*ast.PropExpr
  SelStmt *ast.SelectStmt
  PathMacros []*ast.PathMacroClause
  SelElems []*ast.SelectElem
  Matches []*ast.MatchClause
  PPats []*ast.PathPattern
  PPPats []*ast.PathPatternPrimary
  VPats []*ast.VertexPattern
  EPats []*ast.EdgePattern
  MatchRows *ast.MatchRows
  OTerms []*ast.OrderTerm
  Mods []ast.ModClause
  Insert *ast.InsertClause
  Updates []*ast.Update
  PropAss []*ast.PropAssignment
  NExprs []*ast.NamedExpr
  Exprs []ast.Expr
  Expr ast.Expr
  Whens []*ast.WhenClause
  Quant *ast.Quantifier
  QIdent *ast.QIdent
  Ident *ast.Ident
  Idents []*ast.Ident
  BLit *ast.BasicLit
  B bool
  I int
}

%start start

%%

// Not part of the specification.

start: PgqlStatements  { yylex.(yyParam).Stmts($1) }
     ;

PgqlStatements: PgqlStatement ';'
              | PgqlStatements PgqlStatement ';'  { $$ = append($1, $2[0]) }
              ;

// Main Query Structure

PgqlStatement: CreatePropertyGraph
             | DropPropertyGraph
             | Query
             ;

// Creating a Property Graph

CreatePropertyGraph: CREATE PROPERTY GRAPH GraphName VertexTables OptEdgeTables  { $$ = []ast.Stmt{&ast.CreateStmt{GraphName: $4, VertexTables: $5, EdgeTables: $6}} }
                   ;

GraphName: SchemaQualifiedName
         ;

SchemaQualifiedName: Identifier                       { $$ = &ast.QIdent{Names: []*ast.Ident{$1}} }
                   | SchemaIdentifierPart Identifier  { $$ = &ast.QIdent{Names: append($1.Names, $2)} }
                   ;

SchemaIdentifierPart: Identifier '.'  { $$ = &ast.QIdent{Names: []*ast.Ident{$1}} }
                    ;

VertexTables: VERTEX TABLES '(' VertexTableList ')'  { $$ = $4 }
            ;

VertexTableList: VertexTable
               | VertexTableList ',' VertexTable  { $$ = append($1, $3[0]) }
               ;

OptEdgeTables: /* empty */  { $$ = nil }
             | EdgeTables
             ;

EdgeTables: EDGE TABLES '(' EdgeTableList ')'  { $$ = $4 }
          ;

EdgeTableList: EdgeTable
             | EdgeTableList ',' EdgeTable  { $$ = append($1, $3[0]) }
             ;

VertexTable: TableName OptTableAlias OptKeyClause LabelAndPropertiesClause  { $$ = []*ast.VertexTableDecl{{TableName: $1, TableAlias: $2, Label: $4[0].Label, Props: $4[0].Props, Keys: $3}} }
           ;

LabelAndPropertiesClause: OptLabelClause OptPropertiesClause  { $$ = []*ast.VertexTableDecl{{Label: $1, Props: $2}} }
                        ;

TableName: SchemaQualifiedName
         ;

EdgeTable: TableName OptTableAlias OptKeyClause SourceVertexTable DestinationVertexTable LabelAndPropertiesClause  { $$ = []*ast.EdgeTableDecl{{TableName: $1, TableAlias: $2, Source: $4, Dest: $5, Label: $6[0].Label, Props: $6[0].Props, Keys: $3}} }
         ;

// In the 1.5 spec, KEY and the referenced columns are missing.
SourceVertexTable: SOURCE KEY '(' ColumnNameList ')' REFERENCES TableName '(' ColumnNameList ')'  { $$ = &ast.VertexTableRef{Keys: $4, TableName: $7, Columns: $9} }
                 | SOURCE TableName                                                               { $$ = &ast.VertexTableRef{TableName: $2} }
                 ;

// In the 1.5 spec, KEY and the referenced columns are missing.
DestinationVertexTable: DESTINATION KEY '(' ColumnNameList ')' REFERENCES TableName '(' ColumnNameList ')'  { $$ = &ast.VertexTableRef{Keys: $4, TableName: $7, Columns: $9} }
                      | DESTINATION TableName                                                               { $$ = &ast.VertexTableRef{TableName: $2} }
                      ;

OptTableAlias: /* empty */  { $$ = nil }
             | TableAlias
             ;

TableAlias: AS Identifier  { $$ = $2 }
          | Identifier
          ;

OptKeyClause: /* empty */  { $$ = nil }
            | KeyClause
            ;

KeyClause: KEY '(' ColumnNameList ')'  { $$ = $3 }
         ;

ColumnNameList: ColumnName                     { $$ = []*ast.Ident{$1} }
              | ColumnNameList ',' ColumnName  { $$ = append($1, $3) }
              ;

ColumnName: Identifier
          ;

OptLabelClause: /* empty */  { $$ = nil }
              | LabelClause
              ;

LabelClause: LABEL Label  { $$ = $2 }
           ;

LabelList: Label                { $$ = []*ast.Ident{$1} }
         | LabelList ',' Label  { $$ = append($1, $3) }
         ;

Label: Identifier
     ;

OptPropertiesClause: /* empty */       { $$ = nil }
                   | PropertiesClause
                   ;

PropertiesClause: PropertiesAreAllColumns
                | PropertyExpressions
                | NoProperties
                ;

PropertiesAreAllColumns: PROPERTIES OptAre ALL COLUMNS OptExceptColumns  { $$ = &ast.PropsClause{Except: $5} }
                       ;

OptAre: /* empty */
      | ARE
      ;

OptExceptColumns: /* empty */    { $$ = nil }
                | ExceptColumns
                ;

ExceptColumns: EXCEPT '(' ColumnReferenceList ')'  { $$ = $3 }
             ;

ColumnReferenceList: ColumnReference                          { $$ = []*ast.Ident{$1} }
                   | ColumnReferenceList ',' ColumnReference  { $$ = append($1, $3) }
                   ;

PropertyExpressions: PROPERTIES '(' PropertyExpressionList ')'  { $$ = &ast.PropsClause{Exprs: $3} }
                   ;

PropertyExpressionList: PropertyExpression
                      | PropertyExpressionList ',' PropertyExpression  { $$ = append($1, $3[0]) }
                      ;

PropertyExpression: ColumnReferenceOrCastSpecification AS PropertyName  { $$ = $1; $$[0].Name = $3 }
                  | ColumnReferenceOrCastSpecification
                  ;

ColumnReferenceOrCastSpecification: ColumnReference    { $$ = []*ast.PropExpr{{Column: $1}} }
                                  | CastSpecification  { $$ = []*ast.PropExpr{{CastAs: $1.(*ast.CastExpr)}} }
                                  ;

PropertyName: Identifier
            ;

ColumnReference: Identifier
               ;

NoProperties: NO PROPERTIES  { $$ = &ast.PropsClause{None: true} }
            ;

DropPropertyGraph: DROP PROPERTY GRAPH GraphName  { $$ = []ast.Stmt{&ast.DropStmt{GraphName: $4}} }
                 ;

// Graph Pattern Matching

Query: SelectQuery
     | ModifyQuery
     ;

SelectQuery: OptPathPatternMacros SelectClause FromClause OptWhereClause OptGroupByClause OptHavingClause OptOrderByClause OptLimitOffsetClauses  { $$ = []ast.Stmt{&ast.SelectStmt{PathMacros: $1, Distinct: $2.Distinct, Sels: $2.Sels, From: $3, Where: $4, GroupBy: $5, Having: $6, OrderBy: $7, Limit: $8[0], Offset: $8[1]}} }
           ;

SelectClause: SELECT OptDistinct SelectElementList  { $$ = &ast.SelectStmt{Distinct: $2, Sels: $3} }
            | SELECT '*'                            { $$ = &ast.SelectStmt{} }
            ;

OptDistinct: /* empty */  { $$ = false }
           | DISTINCT     { $$ = true }
           ;

SelectElementList: SelectElement
                 | SelectElementList ',' SelectElement  { $$ = append($1, $3[0]) }
                 ;

SelectElement: ExpAsVar       { $$ = []*ast.SelectElem{{Named: $1[0]}} }
             | AllProperties
             ;

ExpAsVarList: ExpAsVar
            | ExpAsVarList ',' ExpAsVar  { $$ = append($1, $3[0]) }
            ;

ExpAsVar: ValueExpression AS VariableName  { $$ = []*ast.NamedExpr{{Expr: $1, Name: $3}} }
        | ValueExpression                  { $$ = []*ast.NamedExpr{{Expr: $1}} }
        ;

// The 1.5 spec says '.*' as one token, but we ignore space for symmetry.
AllProperties: Identifier '.' '*' AllPropertiesPrefix  { $$ = []*ast.SelectElem{{AllOf: $1, Prefix: $4}} }
             | Identifier '.' '*'                      { $$ = []*ast.SelectElem{{AllOf: $1}} }
             ;

AllPropertiesPrefix: PREFIX StringLiteral  { $$ = $2 }
                   ;

OptFromClause: /* empty */  { $$ = nil }
             | FromClause
             ;

FromClause: FROM MatchClauseList  { $$ = $2 }
          ;

MatchClauseList: MatchClause
               | MatchClauseList ',' MatchClause  { $$ = append($1, $3[0]) }
               ;

MatchClause: MATCH MatchPattern OptOnClause OptRowsPerMatch  { $$ = []*ast.MatchClause{{Patterns: $2, On: $3, Rows: $4}} }
           ;

MatchPattern: PathPattern
            | GraphPattern
            ;

GraphPattern: '(' PathPatternList ')'  { $$ = $2 }
            ;

PathPatternList: PathPattern
               | PathPatternList ',' PathPattern  { $$ = append($1, $3[0]) }
               ;

PathPattern: SimplePathPattern
           | AnyPathPattern
           | AnyShortestPathPattern
           | AllShortestPathPattern
           | TopKShortestPathPattern
           | AnyCheapestPathPattern
           | TopKCheapestPathPattern
           | AllPathPattern
           ;

SimplePathPattern: VertexPattern                                { $$ = []*ast.PathPattern{{Vs: $1}} }
                 | SimplePathPattern PathPrimary VertexPattern  { $$ = $1; p := $$[len($$)-1]; p.Vs = append(p.Vs, $3[0]); p.Es = append(p.Es, $2[0]) }
                 ;

OptVertexPattern: /* empty */    { $$ = nil }
                | VertexPattern
                ;

VertexPattern: '(' VariableSpecification ')'  { $$ = $2 }
             ;

PathPrimary: EdgePattern                 { $$ = []*ast.PathPatternPrimary{{Es: $1}} }
           | ReachabilityPathExpression
           ;

EdgePattern: OutgoingEdgePattern
           | IncomingEdgePattern
           | AnyDirectedEdgePattern
           ;

OutgoingEdgePattern: RARROW                                            { $$ = []*ast.EdgePattern{{Dir: ast.Outgoing}} }
                   | LDASHBRACKET VariableSpecification RBRACKETARROW  { $$ = []*ast.EdgePattern{{Name: $2[0].Name, LabelAlts: $2[0].LabelAlts, Dir: ast.Outgoing}} }
                   ;

IncomingEdgePattern: LARROW                                            { $$ = []*ast.EdgePattern{{Dir: ast.Incoming}} }
                   | LARROWBRACKET VariableSpecification RBRACKETDASH  { $$ = []*ast.EdgePattern{{Name: $2[0].Name, LabelAlts: $2[0].LabelAlts, Dir: ast.Incoming}} }
                   ;

AnyDirectedEdgePattern: '-'                                              { $$ = []*ast.EdgePattern{{Dir: ast.AnyDir}} }
                      | LDASHBRACKET VariableSpecification RBRACKETDASH  { $$ = []*ast.EdgePattern{{Name: $2[0].Name, LabelAlts: $2[0].LabelAlts, Dir: ast.AnyDir}} }
                      ;

VariableSpecification: OptVariableName OptLabelPredicate  { $$ = []*ast.VertexPattern{{Name: $1, LabelAlts: $2}} }
                     ;

OptVariableName: /* empty */   { $$ = nil }
               | VariableName
               ;

VariableName: Identifier
            ;

OptOnClause: /* empty */  { $$ = nil }
           | OnClause
           ;

OnClause: ON GraphName  { $$ = $2 }
        ;

OptLabelPredicate: /* empty */     { $$ = nil }
                 | LabelPredicate
                 ;

LabelPredicate: ':' LabelAlt  { $$ = $2 }
              | IS LabelAlt   { $$ = $2 }
              ;

LabelAlt: Label               { $$ = []*ast.Ident{$1} }
        | LabelAlt '|' Label  { $$ = append($1, $3) }
        ;

OptWhereClause: /* empty */  { $$ = nil }
              | WhereClause
              ;

WhereClause: WHERE ValueExpression  { $$ = $2 }
           ;

// Variable-Length Paths

OptGraphPatternQuantifier: /* empty */             { $$ = nil }
                         | GraphPatternQuantifier
                         ;

GraphPatternQuantifier: ZeroOrMore
                      | OneOrMore
                      | Optional
                      | ExactlyN
                      | NOrMore
                      | BetweenNAndM
                      | BetweenZeroAndM
                      ;

ZeroOrMore: '*'  { $$ = &ast.Quantifier{Min: nil, Max: nil, Group: true} }
          ;

OneOrMore: '+'  { $$ = &ast.Quantifier{Min: &ast.BasicLit{S: "1", Kind: ast.UIntKind}, Max: nil, Group: true} }
         ;

Optional: '?'  { $$ = &ast.Quantifier{Min: nil, Max: &ast.BasicLit{S: "1", Kind: ast.UIntKind}} }
        ;

ExactlyN: '{' UNSIGNED_INTEGER '}'  { $$ = &ast.Quantifier{Min: &ast.BasicLit{S: $2.S, Kind: ast.UIntKind}, Max: &ast.BasicLit{S: $2.S, Kind: ast.UIntKind}, Group: true} }
        ;

NOrMore: '{' UNSIGNED_INTEGER ',' '}'  { $$ = &ast.Quantifier{Min: &ast.BasicLit{S: $2.S, Kind: ast.UIntKind}, Max: nil, Group: true} }
       ;

BetweenNAndM: '{' UNSIGNED_INTEGER ',' UNSIGNED_INTEGER '}'  { $$ = &ast.Quantifier{Min: &ast.BasicLit{S: $2.S, Kind: ast.UIntKind}, Max: &ast.BasicLit{S: $4.S, Kind: ast.UIntKind}, Group: true} }
            ;

BetweenZeroAndM: '{' ',' UNSIGNED_INTEGER '}'  { $$ = &ast.Quantifier{Min: nil, Max: &ast.BasicLit{S: $3.S, Kind: ast.UIntKind}, Group: true} }
               ;

AnyPathPattern: ANY SourceVertexPattern QuantifiedPathPatternPrimary DestinationVertexPattern          { $$ = []*ast.PathPattern{{Vs: append($2, $4[0]), Es: $3, Cardinality: ast.AnyCardinality}} }
              | ANY '(' SourceVertexPattern QuantifiedPathPatternPrimary DestinationVertexPattern ')'  { $$ = []*ast.PathPattern{{Vs: append($3, $5[0]), Es: $4, Cardinality: ast.AnyCardinality}} }
              ;

SourceVertexPattern: VertexPattern
                   ;

DestinationVertexPattern: VertexPattern
                        ;

QuantifiedPathPatternPrimary: PathPatternPrimary OptGraphPatternQuantifier  { $$ = $1; $$[0].Quantity = $2 }
                            ;

PathPatternPrimary: EdgePattern                         { $$ = []*ast.PathPatternPrimary{{Es: $1}} }
                  | ParenthesizedPathPatternExpression
                  ;

ParenthesizedPathPatternExpression: '(' OptVertexPattern EdgePattern OptVertexPattern OptWhereClause OptCostClause ')'  { $$ = []*ast.PathPatternPrimary{{Vs: []*ast.VertexPattern{indexOr($2, 0, nil), indexOr($4, 0, nil)}, Es: $3, Where: $5, Cost: $6}} }
                                  ;

ReachabilityPathExpression: OutgoingPathPattern
                          | IncomingPathPattern
                          ;

OutgoingPathPattern: LDASHSLASH PathSpecification RSLASHARROW  { $$ = $2; $$[0].Es[0].Dir = ast.Outgoing }
                   ;

IncomingPathPattern: LARROWSLASH PathSpecification RSLASHDASH  { $$ = $2; $$[0].Es[0].Dir = ast.Incoming }
                   ;

PathSpecification: LabelPredicate  { $$ = []*ast.PathPatternPrimary{{Es: []*ast.EdgePattern{{LabelAlts: $1, Reachability: true}}}} }
                 | PathPredicate
                 ;

// The 1.5 spec says GraphPatternQuantifier is optional. LabelPredicate handles that case.
PathPredicate: ':' Label GraphPatternQuantifier  { $$ = []*ast.PathPatternPrimary{{Es: []*ast.EdgePattern{{LabelAlts: []*ast.Ident{$2}, Reachability: true}}, Quantity: $3}} }
             ;

OptPathPatternMacros: /* empty */                            { $$ = nil }
                    | OptPathPatternMacros PathPatternMacro  { $$ = append($1, $2[0]) }
                    ;

PathPatternMacro: PATH Identifier AS PathPattern OptWhereClause  { $$ = []*ast.PathMacroClause{{Name: $2, Pattern: $4[0], Where: $5}} }
                ;

AnyShortestPathPattern: ANY SHORTEST SourceVertexPattern QuantifiedPathPatternPrimary DestinationVertexPattern          { $$ = []*ast.PathPattern{{Vs: append($3, $5[0]), Es: $4, Cardinality: ast.AnyCardinality, Metric: ast.LengthMetric}} }
                      | ANY SHORTEST '(' SourceVertexPattern QuantifiedPathPatternPrimary DestinationVertexPattern ')'  { $$ = []*ast.PathPattern{{Vs: append($4, $6[0]), Es: $5, Cardinality: ast.AnyCardinality, Metric: ast.LengthMetric}} }
                      ;

AllShortestPathPattern: ALL SHORTEST SourceVertexPattern QuantifiedPathPatternPrimary DestinationVertexPattern          { $$ = []*ast.PathPattern{{Vs: append($3, $5[0]), Es: $4, Cardinality: ast.AllCardinality, Metric: ast.LengthMetric}} }
                      | ALL SHORTEST '(' SourceVertexPattern QuantifiedPathPatternPrimary DestinationVertexPattern ')'  { $$ = []*ast.PathPattern{{Vs: append($4, $6[0]), Es: $5, Cardinality: ast.AllCardinality, Metric: ast.LengthMetric}} }
                      ;

// The 1.5 spec is missing the SHORTEST keyword.
TopKShortestPathPattern: TOP KValue SHORTEST SourceVertexPattern QuantifiedPathPatternPrimary DestinationVertexPattern          { $$ = []*ast.PathPattern{{Vs: append($4, $6[0]), Es: $5, Cardinality: ast.TopCardinality, K: $2, Metric: ast.LengthMetric}} }
                       | TOP KValue SHORTEST '(' SourceVertexPattern QuantifiedPathPatternPrimary DestinationVertexPattern ')'  { $$ = []*ast.PathPattern{{Vs: append($5, $7[0]), Es: $6, Cardinality: ast.TopCardinality, K: $2, Metric: ast.LengthMetric}} }
                       ;

KValue: UNSIGNED_INTEGER  { $$ = &ast.BasicLit{S: $1.S, Kind: ast.UIntKind} }
      ;

AnyCheapestPathPattern: ANY CHEAPEST SourceVertexPattern QuantifiedPathPatternPrimary DestinationVertexPattern          { $$ = []*ast.PathPattern{{Vs: append($3, $5[0]), Es: $4, Cardinality: ast.AnyCardinality, Metric: ast.CostMetric}} }
                      | ANY CHEAPEST '(' SourceVertexPattern QuantifiedPathPatternPrimary DestinationVertexPattern ')'  { $$ = []*ast.PathPattern{{Vs: append($4, $6[0]), Es: $5, Cardinality: ast.AnyCardinality, Metric: ast.CostMetric}} }
                      ;

OptCostClause: /* empty */  { $$ = nil }
             | CostClause
             ;

CostClause: COST ValueExpression  { $$ = $2 }
          ;

// The 1.5 spec is missing the CHEAPEST keyword.
TopKCheapestPathPattern: TOP KValue CHEAPEST SourceVertexPattern QuantifiedPathPatternPrimary DestinationVertexPattern          { $$ = []*ast.PathPattern{{Vs: append($4, $6[0]), Es: $5, Cardinality: ast.TopCardinality, K: $2, Metric: ast.CostMetric}} }
                       | TOP KValue CHEAPEST '(' SourceVertexPattern QuantifiedPathPatternPrimary DestinationVertexPattern ')'  { $$ = []*ast.PathPattern{{Vs: append($5, $7[0]), Es: $6, Cardinality: ast.TopCardinality, K: $2, Metric: ast.CostMetric}} }
                       ;

// Quantifier must have an upper bound.
AllPathPattern: ALL SourceVertexPattern QuantifiedPathPatternPrimary DestinationVertexPattern          { $$ = []*ast.PathPattern{{Vs: append($2, $4[0]), Es: $3, Cardinality: ast.AllCardinality}} }
              | ALL '(' SourceVertexPattern QuantifiedPathPatternPrimary DestinationVertexPattern ')'  { $$ = []*ast.PathPattern{{Vs: append($3, $5[0]), Es: $4, Cardinality: ast.AllCardinality}} }
              ;

// Number of Rows Per Match

OptRowsPerMatch: /* empty */   { $$ = nil }
               | RowsPerMatch
               ;

RowsPerMatch: OneRowPerMatch
            | OneRowPerVertex
            | OneRowPerStep
            ;

OneRowPerMatch: ONE ROW PER MATCH  { $$ = &ast.MatchRows{Kind: ast.OneRowPerMatch} }
              ;

OneRowPerVertex: ONE ROW PER VERTEX '(' VertexVariable ')'  { $$ = &ast.MatchRows{Kind: ast.OneRowPerVertex, Vars: []*ast.Ident{$6}} }
               ;

VertexVariable: VariableName
              ;

OneRowPerStep: ONE ROW PER STEP '(' VertexVariable1 ',' EdgeVariable ',' VertexVariable2 ')'  { $$ = &ast.MatchRows{Kind: ast.OneRowPerStep, Vars: []*ast.Ident{$6, $8, $10}} }
             ;

VertexVariable1: VariableName
               ;

EdgeVariable: VariableName
            ;

VertexVariable2: VariableName
               ;

// Grouping and Aggregation

OptGroupByClause: /* empty */    { $$ = nil }
                | GroupByClause
                ;

GroupByClause: GROUP BY ExpAsVarList  { $$ = $3 }
             ;

Aggregation: CountAggregation
           | MinAggregation
           | MaxAggregation
           | AvgAggregation
           | SumAggregation
           | ArrayAggregation
           | ListaggAggregation
           ;

CountAggregation: COUNT '(' '*' ')'                          { $$ = &ast.OpExpr{Op: COUNT} }
                | COUNT '(' OptDistinct ValueExpression ')'  { $$ = &ast.OpExpr{Op: COUNT, Args: []ast.Expr{&ast.BasicLit{S: fmt.Sprint($3), Kind: ast.BoolKind}, $4}} }
                ;

MinAggregation: MIN '(' OptDistinct ValueExpression ')'  { $$ = &ast.OpExpr{Op: MIN, Args: []ast.Expr{&ast.BasicLit{S: fmt.Sprint($3), Kind: ast.BoolKind}, $4}} }
              ;

MaxAggregation: MAX '(' OptDistinct ValueExpression ')'  { $$ = &ast.OpExpr{Op: MAX, Args: []ast.Expr{&ast.BasicLit{S: fmt.Sprint($3), Kind: ast.BoolKind}, $4}} }
              ;

AvgAggregation: AVG '(' OptDistinct ValueExpression ')'  { $$ = &ast.OpExpr{Op: AVG, Args: []ast.Expr{&ast.BasicLit{S: fmt.Sprint($3), Kind: ast.BoolKind}, $4}} }
              ;

SumAggregation: SUM '(' OptDistinct ValueExpression ')'  { $$ = &ast.OpExpr{Op: SUM, Args: []ast.Expr{&ast.BasicLit{S: fmt.Sprint($3), Kind: ast.BoolKind}, $4}} }
              ;

ArrayAggregation: ARRAY_AGG '(' OptDistinct ValueExpression ')'  { $$ = &ast.OpExpr{Op: ARRAY_AGG, Args: []ast.Expr{&ast.BasicLit{S: fmt.Sprint($3), Kind: ast.BoolKind}, $4}} }
                ;

ListaggAggregation: LISTAGG '(' OptDistinct ValueExpression OptListaggSeparator ')'  { $$ = &ast.OpExpr{Op: LISTAGG, Args: append([]ast.Expr{&ast.BasicLit{S: fmt.Sprint($3), Kind: ast.BoolKind}}, $4, $5)} }
                  ;

OptListaggSeparator: /* empty */       { $$ = nil }
                   | ListaggSeparator
                   ;

ListaggSeparator: ',' StringLiteral  { $$ = $2 }
                ;

OptHavingClause: /* empty */   { $$ = nil }
               | HavingClause
               ;

HavingClause: HAVING ValueExpression  { $$ = $2 }
            ;

// Sorting and Row Limiting

OptOrderByClause: /* empty */    { $$ = nil }
                | OrderByClause
                ;

OrderByClause: ORDER BY OrderTermList  { $$ = $3 }
             ;

OrderTermList: OrderTerm
             | OrderTermList ',' OrderTerm  { $$ = append($1, $3[0]) }
             ;

OrderTerm: ValueExpression       { $$ = []*ast.OrderTerm{{Expr: $1, Order: ast.DefaultOrder}} }
         | ValueExpression ASC   { $$ = []*ast.OrderTerm{{Expr: $1, Order: ast.AscOrder}} }
         | ValueExpression DESC  { $$ = []*ast.OrderTerm{{Expr: $1, Order: ast.DescOrder}} }
         ;

OptLimitOffsetClauses: /* empty */         { $$ = []ast.Expr{nil, nil} }
                     | LimitOffsetClauses
                     ;

LimitOffsetClauses: LimitClause OffsetClause  { $$ = []ast.Expr{$1, $2} }
                  | OffsetClause LimitClause  { $$ = []ast.Expr{$2, $1} }
                  | LimitClause               { $$ = []ast.Expr{$1, nil} }
                  | OffsetClause              { $$ = []ast.Expr{nil, $1} }
                  ;

LimitClause: LIMIT LimitOffsetValue  { $$ = $2 }
           ;

OffsetClause: OFFSET LimitOffsetValue  { $$ = $2 }
            ;

LimitOffsetValue: UNSIGNED_INTEGER  { $$ = &ast.BasicLit{S: $1.S, Kind: ast.UIntKind, Pos: ast.Pos($1.P.Offset)} }
                | BindVariable
                ;

// Functions and Expressions

ValueExpression: VariableReference         { $$ = $1 }
               | PropertyAccess            { $$ = $1 }
               | Literal                   { $$ = $1 }
               | BindVariable              { $$ = $1 }
               | ArithmeticExpression      { $$ = $1 }
               | RelationalExpression      { $$ = $1 }
               | LogicalExpression         { $$ = $1 }
               | StringConcat              { $$ = $1 }
               | BracketedValueExpression
               | FunctionInvocation        { $$ = $1 }
               | CharacterSubstring        { $$ = $1 }
               | Aggregation               { $$ = $1 }
               | ExtractFunction           { $$ = $1 }
               | IsNullPredicate           { $$ = $1 }
               | IsNotNullPredicate        { $$ = $1 }
               | CastSpecification         { $$ = $1 }
               | CaseExpression            { $$ = $1 }
               | InPredicate               { $$ = $1 }
               | NotInPredicate            { $$ = $1 }
               | ExistsPredicate           { $$ = $1 }
               | ScalarSubquery            { $$ = $1 }
               ;

VariableReferenceList: VariableReference                            { $$ = []*ast.Ident{$1} }
                     | VariableReferenceList ',' VariableReference  { $$ = append($1, $3) }
                     ;

VariableReference: VariableName
                 ;

// The 1.5 spec uses VariableReference. We use Identifier to solve a conflict with FunctionInvocation.
PropertyAccess: Identifier '.' PropertyName  { $$ = &ast.QIdent{Names: []*ast.Ident{$1, $3}} }
              ;

BracketedValueExpression: '(' ValueExpression ')'  { $$ = $2 }
                        ;

// Time literals use STRING_LITERAL and must validate the string.
Literal: StringLiteral
       | NumericLiteral
       | BooleanLiteral
       | DateLiteral
       | TimeLiteral
       | TimestampLiteral
       | IntervalLiteral
       ;

StringLiteral: STRING_LITERAL  { $$ = &ast.BasicLit{S: $1.S, Kind: ast.StringKind, Pos: ast.Pos($1.P.Offset)} }
             ;

NumericLiteral: UNSIGNED_INTEGER  { $$ = &ast.BasicLit{S: $1.S, Kind: ast.UIntKind, Pos: ast.Pos($1.P.Offset)} }
              | UNSIGNED_DECIMAL  { $$ = &ast.BasicLit{S: $1.S, Kind: ast.UDecKind, Pos: ast.Pos($1.P.Offset)} }
              ;

// These are lower-case to match fmt.Sprint(true).
BooleanLiteral: TRUE   { $$ = &ast.BasicLit{S: "true", Kind: ast.BoolKind, Pos: ast.Pos($1.P.Offset)} }
              | FALSE  { $$ = &ast.BasicLit{S: "false", Kind: ast.BoolKind, Pos: ast.Pos($1.P.Offset)} }
              ;

DateLiteral: DATE STRING_LITERAL  { $$ = &ast.BasicLit{S: $2.S, Kind: ast.DateKind, Pos: ast.Pos($2.P.Offset)} }
           ;

TimeLiteral: TIME STRING_LITERAL  { $$ = &ast.BasicLit{S: $2.S, Kind: ast.TimeKind, Pos: ast.Pos($2.P.Offset)} }
           ;

TimestampLiteral: TIMESTAMP STRING_LITERAL  { $$ = &ast.BasicLit{S: $2.S, Kind: ast.TimestampKind, Pos: ast.Pos($2.P.Offset)} }
                ;

IntervalLiteral: INTERVAL StringLiteral DateTimeField  { $$ = $2; $$.S = $2.S + " " + $3.S; $$.Kind = ast.IntervalKind }
               ;

DateTimeField: YEAR    { $$ = &ast.BasicLit{S: "YEAR"} }
             | MONTH   { $$ = &ast.BasicLit{S: "MONTH"} }
             | DAY     { $$ = &ast.BasicLit{S: "DAY"} }
             | HOUR    { $$ = &ast.BasicLit{S: "HOUR"} }
             | MINUTE  { $$ = &ast.BasicLit{S: "MINUTE"} }
             | SECOND  { $$ = &ast.BasicLit{S: "SECOND"} }
             ;

BindVariable: '?'  { $$ = &ast.BindVar{} }
            ;

ArithmeticExpression: UnaryMinus
                    | Multiplication
                    | Division
                    | Modulo
                    | Addition
                    | Subtraction
                    ;

UnaryMinus: '-' ValueExpression  %prec UMINUS  { $$ = &ast.OpExpr{Op: '-', Args: []ast.Expr{$2}} }
          ;

StringConcat: ValueExpression DPIPE ValueExpression  { $$ = &ast.OpExpr{Op: DPIPE, Args: []ast.Expr{$1, $3}} }
            ;

Multiplication: ValueExpression '*' ValueExpression  { $$ = &ast.OpExpr{Op: '*', Args: []ast.Expr{$1, $3}} }
              ;

Division: ValueExpression '/' ValueExpression  { $$ = &ast.OpExpr{Op: '/', Args: []ast.Expr{$1, $3}} }
        ;

Modulo: ValueExpression '%' ValueExpression  { $$ = &ast.OpExpr{Op: '%', Args: []ast.Expr{$1, $3}} }
      ;

Addition: ValueExpression '+' ValueExpression  { $$ = &ast.OpExpr{Op: '+', Args: []ast.Expr{$1, $3}} }
        ;

Subtraction: ValueExpression '-' ValueExpression  { $$ = &ast.OpExpr{Op: '-', Args: []ast.Expr{$1, $3}} }
           ;

RelationalExpression: Equal
                    | NotEqual
                    | Greater
                    | Less
                    | GreaterOrEqual
                    | LessOrEqual
                    ;

Equal: ValueExpression '=' ValueExpression  { $$ = &ast.OpExpr{Op: '=', Args: []ast.Expr{$1, $3}} }
     ;

NotEqual: ValueExpression LTGT ValueExpression  { $$ = &ast.OpExpr{Op: LTGT, Args: []ast.Expr{$1, $3}} }
        ;

Greater: ValueExpression '>' ValueExpression  { $$ = &ast.OpExpr{Op: '>', Args: []ast.Expr{$1, $3}} }
       ;

Less: ValueExpression '<' ValueExpression  { $$ = &ast.OpExpr{Op: '<', Args: []ast.Expr{$1, $3}} }
    ;

GreaterOrEqual: ValueExpression GTEQ ValueExpression  { $$ = &ast.OpExpr{Op: GTEQ, Args: []ast.Expr{$1, $3}} }
              ;

LessOrEqual: ValueExpression LTEQ ValueExpression  { $$ = &ast.OpExpr{Op: LTEQ, Args: []ast.Expr{$1, $3}} }
           ;

LogicalExpression: Not
                 | And
                 | Or
                 ;

Not: NOT ValueExpression  { $$ = &ast.OpExpr{Op: NOT, Args: []ast.Expr{$2}} }
   ;

And: ValueExpression AND ValueExpression  { $$ = &ast.OpExpr{Op: AND, Args: []ast.Expr{$1, $3}} }
   ;

Or: ValueExpression OR ValueExpression  { $$ = &ast.OpExpr{Op: OR, Args: []ast.Expr{$1, $3}} }
  ;

IsNullPredicate: ValueExpression IS NULL  { $$ = &ast.OpExpr{Op: NULL, Args: []ast.Expr{$1}} }
               ;

IsNotNullPredicate: ValueExpression IS NOT NULL  { $$ = &ast.OpExpr{Op: NOT_NULL, Args: []ast.Expr{$1}} }
                  ;

CharacterSubstring: SUBSTRING '(' ValueExpression FROM StartPosition FOR StringLength ')'  { $$ = &ast.OpExpr{Op: SUBSTRING, Args: []ast.Expr{$3, $5, $7}} }
                  | SUBSTRING '(' ValueExpression FROM StartPosition ')'                   { $$ = &ast.OpExpr{Op: SUBSTRING, Args: []ast.Expr{$3, $5}} }
                  ;

StartPosition: ValueExpression
             ;

StringLength: ValueExpression
            ;

ExtractFunction: EXTRACT '(' ExtractField FROM ValueExpression ')'  { $$ = &ast.OpExpr{Op: EXTRACT, Args: []ast.Expr{$3, $5}} }
               ;

ExtractField: YEAR             { $$ = &ast.Ident{Name: "YEAR"} }
            | MONTH            { $$ = &ast.Ident{Name: "MONTH"} }
            | DAY              { $$ = &ast.Ident{Name: "DAY"} }
            | HOUR             { $$ = &ast.Ident{Name: "HOUR"} }
            | MINUTE           { $$ = &ast.Ident{Name: "MINUTE"} }
            | SECOND           { $$ = &ast.Ident{Name: "SECOND"} }
            | TIMEZONE_HOUR    { $$ = &ast.Ident{Name: "TIMEZONE_HOUR"} }
            | TIMEZONE_MINUTE  { $$ = &ast.Ident{Name: "TIMEZONE_MINUTE"} }
            ;

// The 1.5 spec uses PackageName. We use Identifier to solve a conflict with PropertyAccess.
FunctionInvocation: FunctionName '(' OptArgumentList ')'                 { $$ = &ast.CallExpr{Func: &ast.QIdent{Names: []*ast.Ident{$1}}, Args: $3} }
                  | Identifier '.' FunctionName '(' OptArgumentList ')'  { $$ = &ast.CallExpr{Func: &ast.QIdent{Names: []*ast.Ident{$1, $3}}, Args: $5} }
                  | LABEL '(' OptArgumentList ')'                        { $$ = &ast.OpExpr{Op: LABEL, Args: $3} }
                  | LABELS '(' OptArgumentList ')'                       { $$ = &ast.OpExpr{Op: LABELS, Args: $3} }
                  ;

FunctionName: Identifier
            ;

OptArgumentList: /* empty */   { $$ = nil }
               | ArgumentList
               ;

ArgumentList: ValueExpression                   { $$ = []ast.Expr{$1} }
            | ArgumentList ',' ValueExpression  { $$ = append($1, $3) }
            ;

CastSpecification: CAST '(' ValueExpression AS DataType ')'  { $$ = &ast.CastExpr{Arg: $3, TypeKind: $5} }
                 ;

DataType: STRING                    { $$ = STRING }
        | BOOLEAN                   { $$ = BOOLEAN }
        | INTEGER                   { $$ = INTEGER }
        | INT                       { $$ = INT }
        | LONG                      { $$ = LONG }
        | FLOAT                     { $$ = FLOAT }
        | DOUBLE                    { $$ = DOUBLE }
        | DATE                      { $$ = DATE }
        | TIME                      { $$ = TIME }
        | TIME WITH TIME ZONE       { $$ = TIME_TZ }
        | TIMESTAMP                 { $$ = TIMESTAMP }
        | TIMESTAMP WITH TIME ZONE  { $$ = TIMESTAMP_TZ }
        ;

CaseExpression: SimpleCase
              | SearchedCase
              ;

SimpleCase: CASE ValueExpression WhenClauseList OptElseClause END  { $$ = &ast.CaseExpr{Subject: $2, Whens: $3, Else: $4} }
          ;

SearchedCase: CASE WhenClauseList OptElseClause END  { $$ = &ast.CaseExpr{Whens: $2, Else: $3} }
            ;

WhenClauseList: WhenClause
              | WhenClauseList WhenClause  { $$ = append($1, $2[0]) }
              ;

WhenClause: WHEN ValueExpression THEN ValueExpression  { $$ = []*ast.WhenClause{{Cond: $2, Then: $4}} }
          ;

OptElseClause: /* empty */  { $$ = nil }
             | ElseClause
             ;

ElseClause: ELSE ValueExpression  { $$ = $2 }
          ;

InPredicate: ValueExpression IN InValueList  { $$ = &ast.InExpr{Subject: $1, Objects: $3} }
           ;

NotInPredicate: ValueExpression NOT IN InValueList  { $$ = &ast.InExpr{Subject: $1, Objects: $4, Inv: true} }
              ;

InValueList: '(' ValueExpressionList ')'  { $$ = $2 }
           | BindVariable                 { $$ = nil }
           ;

ValueExpressionList: ValueExpression                          { $$ = []ast.Expr{$1} }
                   | ValueExpressionList ',' ValueExpression  { $$ = append($1, $3) }
                   ;

// Subqueries

ExistsPredicate: EXISTS Subquery  { $$ = &ast.OpExpr{Op: EXISTS, Args: []ast.Expr{$2}} }
               ;

// The 1.5 spec uses Query, which would allow ModifyQuery.
Subquery: '(' SelectQuery ')'  { $$ = &ast.SubqueryExpr{Query: $2[0].(*ast.SelectStmt)} }
        ;

ScalarSubquery: Subquery
              ;

// Graph Modification

ModifyQuery: ModifyQueryFull
           ;

// The 1.5 spec requires FromClause and uses ModifyQuerySimple to
// allow an InsertClause alone. This creates a conflict. Parser code
// must validate that if FromClause is missing, then ModificationList
// is a single InsertClause and no other rules are present.
ModifyQueryFull: OptPathPatternMacros ModificationList OptFromClause OptWhereClause OptGroupByClause OptHavingClause OptOrderByClause OptLimitOffsetClauses  { $$ = []ast.Stmt{&ast.ModifyStmt{PathMacros: $1, Mods: $2, From: $3, Where: $4, GroupBy: $5, Having: $6, OrderBy: $7, Limit: $8[0], Offset: $8[1]}} }
               ;

ModificationList: Modification
                | ModificationList Modification  { $$ = append($1, $2[0]) }
                ;

Modification: InsertClause
            | UpdateClause
            | DeleteClause
            ;

InsertClause: INSERT OptIntoClause GraphElementInsertionList  { $$ = []ast.ModClause{&ast.InsertClause{Into: $2, Vs: $3.Vs, Es: $3.Es}} }
            ;

GraphElementInsertionList: GraphElementInsertion
                         | GraphElementInsertionList ',' GraphElementInsertion  { $$ = $1; $$.Vs = append($$.Vs, $3.Vs...); $$.Es = append($$.Es, $3.Es...) }
                         ;

OptIntoClause: /* empty */  { $$ = nil }
             | IntoClause
             ;

IntoClause: INTO GraphName  { $$ = $2 }
          ;

GraphElementInsertion: VERTEX OptVariableName LabelsAndProperties                                            { $$ = &ast.InsertClause{Vs: []*ast.VertexInsertion{{Var: $2, Labels: $3.Vs[0].Labels, Props: $3.Vs[0].Props}}} }
                     | EDGE OptVariableName BETWEEN VertexReference AND VertexReference LabelsAndProperties  { $$ = &ast.InsertClause{Es: []*ast.EdgeInsertion{{Var: $2, Source: $4, Dest: $6, Labels: $7.Vs[0].Labels, Props: $7.Vs[0].Props}}} }
                     ;

VertexReference: Identifier
               ;

// Reusing VertexInsertion to carry Labels and Props together.
LabelsAndProperties: OptLabelSpecification OptPropertiesSpecification  { $$ =  &ast.InsertClause{Vs: []*ast.VertexInsertion{{Labels: $1, Props: $2}}} }
                   ;

OptLabelSpecification: /* empty */         { $$ = nil }
                     | LabelSpecification
                     ;

LabelSpecification: LABELS '(' LabelList ')'  { $$ = $3 }
                  ;

OptPropertiesSpecification: /* empty */              { $$ = nil }
                          | PropertiesSpecification
                          ;

PropertiesSpecification: PROPERTIES '(' PropertyAssignmentList ')'  { $$ = $3 }
                       ;

PropertyAssignmentList: PropertyAssignment
                      | PropertyAssignmentList ',' PropertyAssignment  { $$ = append($1, $3[0]) }
                      ;

PropertyAssignment: PropertyAccess '=' ValueExpression  { $$ = []*ast.PropAssignment{{Prop: $1, Value: $3}} }
                  ;

UpdateClause: UPDATE GraphElementUpdateList  { $$ = []ast.ModClause{&ast.UpdateClause{Updates: $2}} }
            ;

GraphElementUpdateList: GraphElementUpdate                             { $$ = $1 }
                      | GraphElementUpdateList ',' GraphElementUpdate  { $$ = append($1, $3[0]) }
                      ;

GraphElementUpdate: VariableReference SET '(' PropertyAssignmentList ')'  { $$ = []*ast.Update{{Var: $1, Props: $4}} }
                  ;

DeleteClause: DELETE VariableReferenceList  { $$ = []ast.ModClause{&ast.DeleteClause{Vars: $2}} }
            ;

// Other Syntactic Rules

Identifier: UNQUOTED_IDENTIFIER  { $$ = &ast.Ident{Name: $1.S, Pos: ast.Pos($1.P.Offset)} }
          | QUOTED_IDENTIFIER    { $$ = &ast.Ident{Name: token.UnquoteIdentifier($1.S), Pos: ast.Pos($1.P.Offset)} }
          ;

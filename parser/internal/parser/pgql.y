// Implements https://pgql-lang.org/spec/1.5/.
%{
package parser

import (
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

// Literals.

%token <L> TRUE FALSE UNSIGNED_INTEGER UNSIGNED_DECIMAL
%token <L> STRING_LITERAL UNQUOTED_IDENTIFIER QUOTED_IDENTIFIER

// Productions.

%type <Stmts> PgqlStatements PgqlStatement CreatePropertyGraph DropPropertyGraph
%type <VTables> VertexTables VertexTableList VertexTable LabelAndPropertiesClause
%type <ETables> EdgeTables OptEdgeTables EdgeTableList EdgeTable
%type <VTableRef> SourceVertexTable DestinationVertexTable
%type <Props> OptPropertiesClause PropertiesClause PropertiesAreAllColumns PropertyExpressions NoProperties
%type <PExprs> PropertyExpressionList PropertyExpression ColumnReferenceOrCastSpecification
%type <QIdent> GraphName TableName SchemaQualifiedName SchemaIdentifierPart
%type <Ident> Identifier TableAlias OptTableAlias ColumnName LabelClause OptLabelClause Label ColumnReference PropertyName ColumnReference
%type <Idents> OptKeyClause KeyClause ColumnNameList OptExceptColumns ExceptColumns ColumnReferenceList

%union {
  L *lexValue

  Stmts []ast.Stmt
  VTables []*ast.VertexTableDecl
  ETables []*ast.EdgeTableDecl
  VTableRef *ast.VertexTableRef
  Props *ast.PropsClause
  PExprs []*ast.PropExpr
  QIdent *ast.QIdent
  Ident *ast.Ident
  Idents []*ast.Ident
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
             | Query                {/*TODO*/}
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

LabelList: Label
         | LabelList ',' Label
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
                                  | CastSpecification  { yylex.Error("TODO: CastSpecification not implemented"); return 1 }
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

SelectQuery: OptPathPatternMacros SelectClause FromClause OptWhereClause OptGroupByClause OptHavingClause OptOrderByClause OptLimitOffsetClauses
           ;

SelectClause: SELECT OptDistinct SelectElementList
            | SELECT '*'
            ;

OptDistinct: /* empty */
           | DISTINCT
           ;

SelectElementList: SelectElement
                 | SelectElementList ',' SelectElement
                 ;

SelectElement: ExpAsVar
             | AllProperties
             ;

ExpAsVarList: ExpAsVar
            | ExpAsVarList ',' ExpAsVar
            ;

ExpAsVar: ValueExpression AS VariableName
        | ValueExpression
        ;

// The 1.5 spec says '.*' as one token, but we ignore space for symmetry.
AllProperties: Identifier '.' '*' AllPropertiesPrefix
             | Identifier '.' '*'
             ;

AllPropertiesPrefix: PREFIX StringLiteral
                   ;

OptFromClause: /* empty */
             | FromClause
             ;

FromClause: FROM MatchClauseList
          ;

MatchClauseList: MatchClause
               | MatchClauseList ',' MatchClause
               ;

MatchClause: MATCH MatchPattern OptOnClause OptRowsPerMatch
           ;

MatchPattern: PathPattern
            | GraphPattern
            ;

GraphPattern: '(' PathPatternList ')'
            ;

PathPatternList: PathPattern
               | PathPatternList ',' PathPattern
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

SimplePathPattern: VertexPattern
                 | SimplePathPattern PathPrimary VertexPattern
                 ;

OptVertexPattern: /* empty */
                | VertexPattern
                ;

VertexPattern: '(' VariableSpecification ')'
             ;

PathPrimary: EdgePattern
           | ReachabilityPathExpression
           ;

EdgePattern: OutgoingEdgePattern
           | IncomingEdgePattern
           | AnyDirectedEdgePattern
           ;

OutgoingEdgePattern: RARROW
                   | LDASHBRACKET VariableSpecification RBRACKETARROW
                   ;

IncomingEdgePattern: LARROW
                   | LARROWBRACKET VariableSpecification RBRACKETDASH
                   ;

AnyDirectedEdgePattern: '-'
                      | LDASHBRACKET VariableSpecification RBRACKETDASH
                      ;

VariableSpecification: OptVariableName OptLabelPredicate
                     ;

OptVariableName: /* empty */
               | VariableName
               ;

VariableName: Identifier
            ;

OptOnClause: /* empty */
           | OnClause
           ;

OnClause: ON GraphName
        ;

OptLabelPredicate: /* empty */
                 | LabelPredicate
                 ;

LabelPredicate: ':' LabelAlt
              | IS LabelAlt
              ;

LabelAlt: Label
        | LabelAlt '|' Label
        ;

OptWhereClause: /* empty */
              | WhereClause
              ;

WhereClause: WHERE ValueExpression
           ;

// Variable-Length Paths

OptGraphPatternQuantifier: /* empty */
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

ZeroOrMore: '*'
          ;

OneOrMore: '+'
         ;

Optional: '?'
        ;

ExactlyN: '{' UNSIGNED_INTEGER '}'
        ;

NOrMore: '{' UNSIGNED_INTEGER ',' '}'
       ;

BetweenNAndM: '{' UNSIGNED_INTEGER ',' UNSIGNED_INTEGER '}'
            ;

BetweenZeroAndM: '{' ',' UNSIGNED_INTEGER '}'
               ;

AnyPathPattern: ANY SourceVertexPattern QuantifiedPathPatternPrimary DestinationVertexPattern
              | ANY '(' SourceVertexPattern QuantifiedPathPatternPrimary DestinationVertexPattern ')'
              ;

SourceVertexPattern: VertexPattern
                   ;

DestinationVertexPattern: VertexPattern
                        ;

QuantifiedPathPatternPrimary: PathPatternPrimary OptGraphPatternQuantifier
                            ;

PathPatternPrimary: EdgePattern
                  | ParenthesizedPathPatternExpression
                  ;

ParenthesizedPathPatternExpression: '(' OptVertexPattern EdgePattern OptVertexPattern OptWhereClause OptCostClause ')'
                                  ;

ReachabilityPathExpression: OutgoingPathPattern
                          | IncomingPathPattern
                          ;

OutgoingPathPattern: LDASHSLASH PathSpecification RSLASHARROW
                   ;

IncomingPathPattern: LARROWSLASH PathSpecification RSLASHDASH
                   ;

PathSpecification: LabelPredicate
                 | PathPredicate
                 ;

// The 1.5 spec says GraphPatternQuantifier is optional. LabelPredicate handles that case.
PathPredicate: ':' Label GraphPatternQuantifier
             ;

OptPathPatternMacros: /* empty */
                    | OptPathPatternMacros PathPatternMacro
                    ;

PathPatternMacro: PATH Identifier AS PathPattern OptWhereClause
                ;

AnyShortestPathPattern: ANY SHORTEST SourceVertexPattern QuantifiedPathPatternPrimary DestinationVertexPattern
                      | ANY SHORTEST '(' SourceVertexPattern QuantifiedPathPatternPrimary DestinationVertexPattern ')'
                      ;

AllShortestPathPattern: ALL SHORTEST SourceVertexPattern QuantifiedPathPatternPrimary DestinationVertexPattern
                      | ALL SHORTEST '(' SourceVertexPattern QuantifiedPathPatternPrimary DestinationVertexPattern ')'
                      ;

// The 1.5 spec is missing the SHORTEST keyword.
TopKShortestPathPattern: TOP KValue SHORTEST SourceVertexPattern QuantifiedPathPatternPrimary DestinationVertexPattern
                       | TOP KValue SHORTEST '(' SourceVertexPattern QuantifiedPathPatternPrimary DestinationVertexPattern ')'
                       ;

KValue: UNSIGNED_INTEGER
      ;

AnyCheapestPathPattern: ANY CHEAPEST SourceVertexPattern QuantifiedPathPatternPrimary DestinationVertexPattern
                      | ANY CHEAPEST '(' SourceVertexPattern QuantifiedPathPatternPrimary DestinationVertexPattern ')'
                      ;

OptCostClause: /* empty */
             | CostClause
             ;

CostClause: COST ValueExpression
          ;

// The 1.5 spec is missing the CHEAPEST keyword.
TopKCheapestPathPattern: TOP KValue CHEAPEST SourceVertexPattern QuantifiedPathPatternPrimary DestinationVertexPattern
                       | TOP KValue CHEAPEST '(' SourceVertexPattern QuantifiedPathPatternPrimary DestinationVertexPattern ')'
                       ;

AllPathPattern: ALL SourceVertexPattern QuantifiedPathPatternPrimary DestinationVertexPattern
              | ALL '(' SourceVertexPattern QuantifiedPathPatternPrimary DestinationVertexPattern ')'
              ;

// Number of Rows Per Match

OptRowsPerMatch: /* empty */
               | RowsPerMatch
               ;

RowsPerMatch: OneRowPerMatch
            | OneRowPerVertex
            | OneRowPerStep
            ;

OneRowPerMatch: ONE ROW PER MATCH
              ;

OneRowPerVertex: ONE ROW PER VERTEX '(' VertexVariable ')'
               ;

VertexVariable: VariableName
              ;

OneRowPerStep: ONE ROW PER STEP '(' VertexVariable1 ',' EdgeVariable ',' VertexVariable2 ')'
             ;

VertexVariable1: VariableName
               ;

EdgeVariable: VariableName
            ;

VertexVariable2: VariableName
               ;

// Grouping and Aggregation

OptGroupByClause: /* empty */
                | GroupByClause
                ;

GroupByClause: GROUP BY ExpAsVarList
             ;

Aggregation: CountAggregation
           | MinAggregation
           | MaxAggregation
           | AvgAggregation
           | SumAggregation
           | ArrayAggregation
           | ListaggAggregation
           ;

CountAggregation: COUNT '(' '*' ')'
                | COUNT '(' OptDistinct ValueExpression ')'
                ;

MinAggregation: MIN '(' OptDistinct ValueExpression ')'
              ;

MaxAggregation: MAX '(' OptDistinct ValueExpression ')'
              ;

AvgAggregation: AVG '(' OptDistinct ValueExpression ')'
              ;

SumAggregation: SUM '(' OptDistinct ValueExpression ')'
              ;

ArrayAggregation: ARRAY_AGG '(' OptDistinct ValueExpression ')'
                ;

ListaggAggregation: LISTAGG '(' OptDistinct ValueExpression OptListaggSeparator ')'
                  ;

OptListaggSeparator: /* empty */
                   | ListaggSeparator
                   ;

ListaggSeparator: ',' StringLiteral
                ;

OptHavingClause: /* empty */
               | HavingClause
               ;

HavingClause: HAVING ValueExpression
            ;

// Sorting and Row Limiting

OptOrderByClause: /* empty */
                | OrderByClause
                ;

OrderByClause: ORDER BY OrderTermList
             ;

OrderTermList: OrderTerm
             | OrderTermList ',' OrderTerm
             ;

OrderTerm: ValueExpression
         | ValueExpression ASC
         | ValueExpression DESC
         ;

OptLimitOffsetClauses: /* empty */
                     | LimitOffsetClauses
                     ;

LimitOffsetClauses: LimitClause OffsetClause
                  | OffsetClause LimitClause
                  | LimitClause
                  | OffsetClause
                  ;

LimitClause: LIMIT LimitOffsetValue
           ;

OffsetClause: OFFSET LimitOffsetValue
            ;

LimitOffsetValue: UNSIGNED_INTEGER
                | BindVariable
                ;

// Functions and Expressions

ValueExpression: VariableReference
               | PropertyAccess
               | Literal
               | BindVariable
               | ArithmeticExpression
               | RelationalExpression
               | LogicalExpression
               | StringConcat
               | BracketedValueExpression
               | FunctionInvocation
               | CharacterSubstring
               | Aggregation
               | ExtractFunction
               | IsNullPredicate
               | IsNotNullPredicate
               | CastSpecification
               | CaseExpression
               | InPredicate
               | NotInPredicate
               | ExistsPredicate
               | ScalarSubquery
               ;

VariableReferenceList: VariableReference
                     | VariableReferenceList ',' VariableReference
                     ;

VariableReference: VariableName
                 ;

// The 1.5 spec uses VariableReference. We use Identifier to solve a conflict with FunctionInvocation.
PropertyAccess: Identifier '.' PropertyName
              ;

BracketedValueExpression: '(' ValueExpression ')'
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

StringLiteral: STRING_LITERAL
             ;

NumericLiteral: UNSIGNED_INTEGER
              | UNSIGNED_DECIMAL
              ;

BooleanLiteral: TRUE
              | FALSE
              ;

DateLiteral: DATE STRING_LITERAL
           ;

TimeLiteral: TIME STRING_LITERAL
           ;

TimestampLiteral: TIMESTAMP STRING_LITERAL
                ;

IntervalLiteral: INTERVAL StringLiteral DateTimeField
               ;

DateTimeField: YEAR
             | MONTH
             | DAY
             | HOUR
             | MINUTE
             | SECOND
             ;

BindVariable: '?'
            ;

ArithmeticExpression: UnaryMinus
                    | Multiplication
                    | Division
                    | Modulo
                    | Addition
                    | Subtraction
                    ;

UnaryMinus: '-' ValueExpression  %prec UMINUS
          ;

StringConcat: ValueExpression DPIPE ValueExpression
            ;

Multiplication: ValueExpression '*' ValueExpression
              ;

Division: ValueExpression '/' ValueExpression
        ;

Modulo: ValueExpression '%' ValueExpression
      ;

Addition: ValueExpression '+' ValueExpression
        ;

Subtraction: ValueExpression '-' ValueExpression
           ;

RelationalExpression: Equal
                    | NotEqual
                    | Greater
                    | Less
                    | GreaterOrEqual
                    | LessOrEqual
                    ;

Equal: ValueExpression '=' ValueExpression
     ;

NotEqual: ValueExpression LTGT ValueExpression
        ;

Greater: ValueExpression '>' ValueExpression
       ;

Less: ValueExpression '<' ValueExpression
    ;

GreaterOrEqual: ValueExpression GTEQ ValueExpression
              ;

LessOrEqual: ValueExpression LTEQ ValueExpression
           ;

LogicalExpression: Not
                 | And
                 | Or
                 ;

Not: NOT ValueExpression
   ;

And: ValueExpression AND ValueExpression
   ;

Or: ValueExpression OR ValueExpression
  ;

IsNullPredicate: ValueExpression IS NULL
               ;

IsNotNullPredicate: ValueExpression IS NOT NULL
                  ;

CharacterSubstring: SUBSTRING '(' ValueExpression FROM StartPosition FOR StringLength ')'
                  | SUBSTRING '(' ValueExpression FROM StartPosition ')'
                  ;

StartPosition: ValueExpression
             ;

StringLength: ValueExpression
            ;

ExtractFunction: EXTRACT '(' ExtractField FROM ValueExpression ')'
               ;

ExtractField: YEAR
            | MONTH
            | DAY
            | HOUR
            | MINUTE
            | SECOND
            | TIMEZONE_HOUR
            | TIMEZONE_MINUTE
            ;

// The 1.5 spec uses PackageName. We use Identifier to solve a conflict with PropertyAccess.
FunctionInvocation: FunctionName '(' OptArgumentList ')'
                  | Identifier '.' FunctionName '(' OptArgumentList ')'
                  | LABEL '(' OptArgumentList ')'
                  | LABELS '(' OptArgumentList ')'
                  ;

FunctionName: Identifier
            ;

OptArgumentList: /* empty */
               | ArgumentList
               ;

ArgumentList: ValueExpression
            | ArgumentList ',' ValueExpression
            ;

CastSpecification: CAST '(' ValueExpression AS DataType ')'
                 ;

DataType: STRING
        | BOOLEAN
        | INTEGER
        | INT
        | LONG
        | FLOAT
        | DOUBLE
        | DATE
        | TIME
        | TIME WITH TIME ZONE
        | TIMESTAMP
        | TIMESTAMP WITH TIME ZONE
        ;

CaseExpression: SimpleCase
              | SearchedCase
              ;

SimpleCase: CASE ValueExpression WhenClauseList OptElseClause END
          ;

SearchedCase: CASE WhenClauseList OptElseClause END
            ;

WhenClauseList: WhenClause
              | WhenClauseList WhenClause
              ;

WhenClause: WHEN ValueExpression THEN ValueExpression
          ;

OptElseClause: /* empty */
             | ElseClause
             ;

ElseClause: ELSE ValueExpression
          ;

InPredicate: ValueExpression IN InValueList
           ;

NotInPredicate: ValueExpression NOT IN InValueList
              ;

InValueList: '(' ValueExpressionList ')'
           | BindVariable
           ;

ValueExpressionList: ValueExpression
                   | ValueExpressionList ',' ValueExpression
                   ;

// Subqueries

ExistsPredicate: EXISTS Subquery
               ;

// The 1.5 spec uses Query, which would allow ModifyQuery.
Subquery: '(' SelectQuery ')'
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
ModifyQueryFull: OptPathPatternMacros ModificationList OptFromClause OptWhereClause OptGroupByClause OptHavingClause OptOrderByClause OptLimitOffsetClauses
               ;

ModificationList: Modification
                | ModificationList Modification
                ;

Modification: InsertClause
            | UpdateClause
            | DeleteClause
            ;

InsertClause: INSERT OptIntoClause GraphElementInsertionList
            ;

GraphElementInsertionList: GraphElementInsertion
                         | GraphElementInsertionList ',' GraphElementInsertion
                         ;

OptIntoClause: /* empty */
             | IntoClause
             ;

IntoClause: INTO GraphName
          ;

GraphElementInsertion: VERTEX OptVariableName LabelsAndProperties
                     | EDGE OptVariableName BETWEEN VertexReference AND VertexReference LabelsAndProperties
                     ;

VertexReference: Identifier
               ;

LabelsAndProperties: OptLabelSpecification OptPropertiesSpecification
                   ;

OptLabelSpecification: /* empty */
                     | LabelSpecification
                     ;

LabelSpecification: LABELS '(' LabelList ')'
                  ;

OptPropertiesSpecification: /* empty */
                          | PropertiesSpecification
                          ;

PropertiesSpecification: PROPERTIES '(' PropertyAssignmentList ')'
                       ;

PropertyAssignmentList: PropertyAssignment
                      | PropertyAssignmentList ',' PropertyAssignment
                      ;

PropertyAssignment: PropertyAccess '=' ValueExpression
                  ;

UpdateClause: UPDATE GraphElementUpdateList
            ;

GraphElementUpdateList: GraphElementUpdate
                      | GraphElementUpdateList ',' GraphElementUpdate
                      ;

GraphElementUpdate: VariableReference SET '(' PropertyAssignmentList ')'
                  ;

DeleteClause: DELETE VariableReferenceList
            ;

// Other Syntactic Rules

Identifier: UNQUOTED_IDENTIFIER  { $$ = &ast.Ident{Name: $1.S, Pos: ast.Pos($1.P.Offset)} }
          | QUOTED_IDENTIFIER    { $$ = &ast.Ident{Name: token.UnquoteIdentifier($1.S), Pos: ast.Pos($1.P.Offset)} }
          ;

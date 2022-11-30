package ast

type SelectStmt struct {
	PathMacros []*PathMacroClause
	Sels       []*SelectElem
	From       []*MatchClause
	Where      Expr
	GroupBy    []*NamedExpr
	Having     Expr
	Limit      Expr
	Offset     Expr
	OrderBy    []*OrderTerm
	Distinct   bool
}

func (SelectStmt) stmtTag() {}

type PathMacroClause struct {
	Name    *Ident
	Pattern *PathPattern
	Where   Expr
}

type SelectElem struct {
	Named  *NamedExpr
	AllOf  *Ident
	Prefix *BasicLit
}

type MatchClause struct {
	On       *QIdent
	Rows     *MatchRows
	Patterns []*PathPattern
}

type PathPattern struct {
	K           *BasicLit
	Vs          []*VertexPattern
	Es          []*PathPatternPrimary
	Cardinality Cardinality
	Metric      Metric
}

type Cardinality int

const (
	NoCardinality Cardinality = iota
	AnyCardinality
	AllCardinality
	TopCardinality
)

type Metric int

const (
	NoMetric Metric = iota
	LengthMetric
	CostMetric
)

type PathPatternPrimary struct {
	Quantity *Quantifier
	Where    Expr
	Cost     Expr
	Vs       []*VertexPattern
	Es       []*EdgePattern
}

type VertexPattern struct {
	Name      *Ident
	LabelAlts []*Ident
}

type EdgePattern struct {
	Name         *Ident
	LabelAlts    []*Ident
	Dir          Dir
	Reachability bool
}

type Dir int

const (
	AnyDir Dir = iota
	Outgoing
	Incoming
)

type MatchRows struct {
	Vars []*Ident
	Kind MatchRowsKind
}

type MatchRowsKind int

const (
	DefaultMatchRows MatchRowsKind = iota
	OneRowPerMatch
	OneRowPerVertex
	OneRowPerStep
)

type OrderTerm struct {
	Expr  Expr
	Order Order
}

type Order int

const (
	DefaultOrder Order = iota
	AscOrder
	DescOrder
)

type Quantifier struct {
	Min, Max *BasicLit
	Group    bool
}

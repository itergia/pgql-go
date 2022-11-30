package ast

type Pos int

type Stmt interface {
	stmtTag()
}

type Expr interface {
	exprTag()
}

type OpExpr struct {
	Args []Expr
	Op   int
}

func (OpExpr) exprTag() {}

type CallExpr struct {
	Func *QIdent
	Args []Expr
}

func (CallExpr) exprTag() {}

type CastExpr struct {
	Arg      Expr
	TypeKind int
}

func (CastExpr) exprTag() {}

type CaseExpr struct {
	Subject Expr
	Else    Expr
	Whens   []*WhenClause
}

func (CaseExpr) exprTag() {}

type WhenClause struct {
	Cond Expr
	Then Expr
}

type InExpr struct {
	Subject Expr
	Objects []Expr // If empty, a bind variable exists.
	Inv     bool
}

func (InExpr) exprTag() {}

type SubqueryExpr struct {
	Query *SelectStmt
}

func (SubqueryExpr) exprTag() {}

type NamedExpr struct {
	Expr Expr
	Name *Ident
}

type QIdent struct {
	Names []*Ident
}

func (QIdent) exprTag() {}

type Ident struct {
	Name string
	Pos  Pos
}

func (Ident) exprTag() {}

type BasicLit struct {
	S    string
	Kind LitKind
	Pos  Pos
}

func (BasicLit) exprTag() {}

type LitKind int

const (
	UnknownLitKind LitKind = iota
	StringKind
	UIntKind
	UDecKind
	BoolKind
	DateKind
	TimeKind
	TimestampKind
	IntervalKind
)

type BindVar struct{}

func (BindVar) exprTag() {}

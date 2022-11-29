package ast

type Pos int

type Stmt interface {
	stmtTag()
}

type Expr interface {
	exprTag()
}

type QIdent struct {
	Names []*Ident
}

type Ident struct {
	Name string
	Pos  Pos
}

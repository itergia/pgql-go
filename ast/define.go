package ast

type CreateStmt struct {
	GraphName    *QIdent
	VertexTables []*VertexTableDecl
	EdgeTables   []*EdgeTableDecl
}

func (CreateStmt) stmtTag() {}

type DropStmt struct {
	GraphName *QIdent
}

func (DropStmt) stmtTag() {}

type VertexTableDecl struct {
	TableName  *QIdent
	TableAlias *Ident
	Label      *Ident
	Props      *PropsClause
	Keys       []*Ident
}

type EdgeTableDecl struct {
	TableName  *QIdent
	TableAlias *Ident
	Source     *VertexTableRef
	Dest       *VertexTableRef
	Label      *Ident
	Props      *PropsClause
	Keys       []*Ident
}

type PropsClause struct {
	Except []*Ident    // Valid iff None is false.
	Exprs  []*PropExpr // Valid iff None is false and Except is empty.
	None   bool
}

type VertexTableRef struct {
	Keys      []*Ident
	TableName *QIdent
	Columns   []*Ident
}

type PropExpr struct {
	Name   *Ident
	Column *Ident
	CastAs Expr // The 1.5 spec suggests that the ValueExpression in the cast can only be a simple column name.
}

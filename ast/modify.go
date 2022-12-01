package ast

type ModifyStmt struct {
	PathMacros []*PathMacroClause
	Mods       []ModClause
	From       []*MatchClause
	Where      Expr
	GroupBy    []*NamedExpr
	Having     Expr
	Limit      Expr
	Offset     Expr
	OrderBy    []*OrderTerm
}

func (ModifyStmt) stmtTag() {}

type ModClause interface {
	modTag()
}

type InsertClause struct {
	Into *QIdent
	Vs   []*VertexInsertion
	Es   []*EdgeInsertion
}

func (InsertClause) modTag() {}

type UpdateClause struct {
	Updates []*Update
}

func (UpdateClause) modTag() {}

type Update struct {
	Var   *Ident
	Props []*PropAssignment
}

type DeleteClause struct {
	Vars []*Ident
}

func (DeleteClause) modTag() {}

type VertexInsertion struct {
	Var    *Ident
	Labels []*Ident
	Props  []*PropAssignment
}

type PropAssignment struct {
	Prop  *QIdent
	Value Expr
}

type EdgeInsertion struct {
	Var    *Ident
	Source *Ident
	Dest   *Ident
	Labels []*Ident
	Props  []*PropAssignment
}

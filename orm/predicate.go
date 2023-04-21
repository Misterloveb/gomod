package orm

type op string

func (s op) ToString() string {
	return string(s)
}

const (
	opEq  op = "="
	opNot op = "NOT"
	opAnd op = "AND"
	opOr  op = "OR"
)

type Predicate struct {
	left  Expression
	op    op
	right Expression
}

func (Predicate) Expr() {}
func (p Predicate) And(pd Predicate) Predicate {
	return Predicate{
		left:  p,
		op:    opAnd,
		right: pd,
	}
}
func (p Predicate) Or(pd Predicate) Predicate {
	return Predicate{
		left:  p,
		op:    opOr,
		right: pd,
	}
}

func Not(pd Predicate) Predicate {
	return Predicate{
		op:    opNot,
		right: pd,
	}
}

type Column struct {
	name string
}

func C(column string) Column {
	return Column{
		name: column,
	}
}
func (c Column) Eq(arg any) Predicate {
	return Predicate{
		left: c,
		op:   opEq,
		right: Value{
			Arg: arg,
		},
	}
}
func (Column) Expr() {}

type Value struct {
	Arg any
}

func (Value) Expr() {}

//标记接口
type Expression interface {
	Expr()
}

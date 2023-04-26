package orm

import (
	"strings"

	"github.com/Misterloveb/gomod/orm/internel/err"
	"github.com/Misterloveb/gomod/orm/model"
)

type builder struct {
	model *model.Model
	where []Predicate
	args  []any
	str   strings.Builder
}

func (b *builder) buildExpression(p Expression) error {
	if p == nil {
		return nil
	}
	switch exp := p.(type) {
	case Predicate:
		_, ok := exp.left.(Predicate)
		if ok {
			b.str.WriteByte('(')
		}
		if err := b.buildExpression(exp.left); err != nil {
			return err
		}
		if ok {
			b.str.WriteByte(')')
		}
		b.str.WriteByte(' ')
		b.str.WriteString(exp.op.ToString())
		b.str.WriteByte(' ')
		_, ok = exp.right.(Predicate)
		if ok {
			b.str.WriteByte('(')
		}
		if err := b.buildExpression(exp.right); err != nil {
			return err
		}
		if ok {
			b.str.WriteByte(')')
		}
	case Column:
		col, ok := b.model.FieldMap[exp.name]
		if !ok {
			return err.ErrUnKnowColumn(exp.name)
		}
		b.str.WriteByte('`')
		b.str.WriteString(col.Column)
		b.str.WriteByte('`')
	case Value:
		b.str.WriteByte('?')
		b.args = append(b.args, exp.Arg)
	}
	return nil
}
func (b *builder) buildPredicate(pd []Predicate) error {
	p := pd[0]
	for i, wlen := 1, len(pd); i < wlen; i++ {
		p = p.And(pd[i])
	}
	return b.buildExpression(p)
}

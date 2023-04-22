package orm

import (
	"context"
	"strings"

	"github.com/Misterloveb/gomod/orm/internel/err"
)

type Selector[T any] struct {
	tablename string
	where     []Predicate
	model     *Model
	str       strings.Builder
	args      []any
}

func (s *Selector[T]) Build() (*Query, error) {
	s.str.WriteString("SELECT * FROM ")
	var err error
	s.model, err = ParseModel(new(T))
	if err != nil {
		return nil, err
	}
	if s.tablename == "" {
		s.str.WriteByte('`')
		s.str.WriteString(s.model.TableName)
		s.str.WriteByte('`')
	} else {
		s.str.WriteString(s.tablename)
	}
	if len(s.where) > 0 {
		s.str.WriteString(" WHERE ")
		p := s.where[0]
		for i, wlen := 1, len(s.where); i < wlen; i++ {
			p = p.And(s.where[i])
		}
		if err := s.buildExpression(p); err != nil {
			return nil, err
		}
	}

	s.str.WriteByte(';')
	return &Query{
		Sql:  s.str.String(),
		Args: s.args,
	}, nil
}
func (s *Selector[T]) buildExpression(p Expression) error {
	if p == nil {
		return nil
	}
	switch exp := p.(type) {
	case Predicate:
		_, ok := exp.left.(Predicate)
		if ok {
			s.str.WriteByte('(')
		}
		if err := s.buildExpression(exp.left); err != nil {
			return err
		}
		if ok {
			s.str.WriteByte(')')
		}
		s.str.WriteByte(' ')
		s.str.WriteString(exp.op.ToString())
		s.str.WriteByte(' ')
		_, ok = exp.right.(Predicate)
		if ok {
			s.str.WriteByte('(')
		}
		if err := s.buildExpression(exp.right); err != nil {
			return err
		}
		if ok {
			s.str.WriteByte(')')
		}
	case Column:
		col, ok := s.model.Field[exp.name]
		if !ok {
			return err.ErrUnKnowColumn(exp.name)
		}
		s.str.WriteByte('`')
		s.str.WriteString(col.Column)
		s.str.WriteByte('`')
	case Value:
		s.str.WriteByte('?')
		s.args = append(s.args, exp.Arg)
	}
	return nil
}
func (s *Selector[T]) From(name string) *Selector[T] {
	s.tablename = name
	return s
}
func (s *Selector[T]) Where(pd ...Predicate) *Selector[T] {
	s.where = pd
	return s
}
func (s *Selector[T]) Get(ctx context.Context) (T, error) {
	panic("not implemented") // TODO: Implement
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	panic("not implemented") // TODO: Implement
}

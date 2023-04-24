package orm

import (
	"context"
)

type Selector[T any] struct {
	builder
	tablename string
	where     []Predicate
	db        *DB
}

func NewSelector[T any](db *DB) *Selector[T] {
	return &Selector[T]{
		db: db,
	}
}
func (s *Selector[T]) Build() (*Query, error) {
	s.str.WriteString("SELECT * FROM ")
	var err error
	s.model, err = s.db.r.Registry(new(T))
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
		if err := s.buildPredicate(s.where); err != nil {
			return nil, err
		}
	}
	s.str.WriteByte(';')
	return &Query{
		Sql:  s.str.String(),
		Args: s.args,
	}, nil
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

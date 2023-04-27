package orm

import (
	"context"

	"github.com/Misterloveb/gomod/orm/internel/err"
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
func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	sqlobj, errs := s.Build()
	if errs != nil {
		return nil, errs
	}
	rows, errs := s.db.db.QueryContext(ctx, sqlobj.Sql, sqlobj.Args...)
	if errs != nil {
		return nil, errs
	}
	if !rows.Next() {
		return nil, err.ErrNoRows
	}
	ptr_T := new(T)
	s.db.valuer(ptr_T, s.model).SetColumnVal(rows)
	return ptr_T, nil
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	sqlobj, errs := s.Build()
	if errs != nil {
		return nil, errs
	}
	rows, errs := s.db.db.QueryContext(ctx, sqlobj.Sql, sqlobj.Args...)
	if errs != nil {
		return nil, errs
	}
	slice_T := make([]*T, 0, 10)
	for rows.Next() {
		p := new(T)
		s.db.valuer(p, s.model).SetColumnVal(rows)
		slice_T = append(slice_T, p)
	}
	if len(slice_T) == 0 {
		return nil, err.ErrNoRows
	}
	return slice_T, nil
}

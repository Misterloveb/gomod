package orm

import (
	"context"
	"database/sql"
)

type Deleteor[T any] struct {
	builder
	tablename string
	where     []Predicate
	db        *DB
}

func NewDeleteor[T any](db *DB) *Deleteor[T] {
	return &Deleteor[T]{
		db: db,
	}
}
func (d *Deleteor[T]) Build() (*Query, error) {
	d.str.WriteString("DELETE FROM ")
	var err error
	d.model, err = d.db.r.Get(new(T))
	if err != nil {
		return nil, err
	}
	if d.tablename == "" {
		d.str.WriteByte('`')
		d.str.WriteString(d.model.TableName)
		d.str.WriteByte('`')
	} else {
		d.str.WriteString(d.tablename)
	}
	if len(d.where) > 0 {
		d.str.WriteString(" WHERE ")
		if err := d.buildPredicate(d.where); err != nil {
			return nil, err
		}
	}

	d.str.WriteByte(';')
	return &Query{
		Sql:  d.str.String(),
		Args: d.args,
	}, nil
}

func (d *Deleteor[T]) Where(pd ...Predicate) *Deleteor[T] {
	d.where = pd
	return d
}
func (d *Deleteor[T]) Exec(ctx context.Context) (sql.Result, error) {
	panic("not implemented") // TODO: Implement
}
func (d *Deleteor[T]) From(table string) *Deleteor[T] {
	d.tablename = table
	return d
}

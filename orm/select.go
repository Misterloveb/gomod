package orm

import (
	"context"
	"reflect"
	"strings"
)

type Selector[T any] struct {
	tablename string
}

func (s *Selector[T]) Build() (*Query, error) {
	var str strings.Builder
	str.WriteString("SELECT * FROM ")
	if s.tablename == "" {
		var t T
		typ := reflect.TypeOf(t)
		str.WriteByte('`')
		str.WriteString(typ.Name())
		str.WriteByte('`')
	} else {
		str.WriteString(s.tablename)
	}
	str.WriteByte(';')
	return &Query{
		Sql: str.String(),
	}, nil
}
func (s *Selector[T]) From(name string) *Selector[T] {
	s.tablename = name
	return s
}
func (s *Selector[T]) Get(ctx context.Context) (T, error) {
	panic("not implemented") // TODO: Implement
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	panic("not implemented") // TODO: Implement
}

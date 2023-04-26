package valuer

import (
	"database/sql"
	"reflect"

	"github.com/Misterloveb/gomod/orm/internel/err"
	"github.com/Misterloveb/gomod/orm/model"
)

type reflectValuer struct {
	val   any //T的指针
	model *model.Model
}

func NewReflectValuer() Creator {
	return func(entity any, model *model.Model) ValuerFace {
		return &reflectValuer{
			val:   entity,
			model: model,
		}
	}
}
func (r *reflectValuer) SetColumnVal(rows *sql.Rows) error {
	columns, errs := rows.Columns()
	if errs != nil {
		return errs
	}
	cnum := len(columns)
	vals := make([]any, 0, cnum)
	colselem := make([]reflect.Value, 0, cnum)
	for _, col := range columns {
		cn, ok := r.model.ColumnMap[col]
		if !ok {
			return err.ErrUnKnowColumn(col)
		}
		newtyp := reflect.New(cn.Ctype)
		val := newtyp.Interface()
		vals = append(vals, val)
		colselem = append(colselem, newtyp.Elem())
	}
	if err := rows.Scan(vals...); err != nil {
		return err
	}
	ptr_T_value := reflect.ValueOf(r.val).Elem()
	for k, fdname := range columns {
		cn, ok := r.model.ColumnMap[fdname]
		if !ok {
			return err.ErrUnKnowColumn(fdname)
		}
		if ptr_T_value.FieldByName(cn.Goname).CanSet() {
			ptr_T_value.FieldByName(cn.Goname).Set(colselem[k])
		}
	}
	return nil
}

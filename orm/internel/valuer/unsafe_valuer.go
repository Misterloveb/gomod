package valuer

import (
	"database/sql"
	"reflect"
	"unsafe"

	"github.com/Misterloveb/gomod/orm/internel/err"
	"github.com/Misterloveb/gomod/orm/model"
)

type UnsafeValuer struct {
	val   any //T的指针类型
	model *model.Model
}

func NewUnsafeValuer() Creator {
	return func(entity any, model *model.Model) ValuerFace {
		return &UnsafeValuer{
			val:   entity,
			model: model,
		}
	}
}
func (u *UnsafeValuer) SetColumnVal(rows *sql.Rows) error {
	address_T := reflect.ValueOf(u.val).UnsafePointer() //T的起始地址
	columns, errs := rows.Columns()
	if errs != nil {
		return errs
	}
	vals := make([]any, 0, len(columns))
	for _, cl := range columns {
		fd, ok := u.model.ColumnMap[cl]
		if !ok {
			return err.ErrUnKnowColumn(cl)
		}
		fd_column_addr := reflect.NewAt(fd.Ctype, unsafe.Pointer(uintptr(address_T)+fd.Offset))
		vals = append(vals, fd_column_addr.Interface())
	}
	if err := rows.Scan(vals...); err != nil {
		return err
	}
	return nil
}

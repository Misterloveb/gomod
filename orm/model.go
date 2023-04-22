package orm

import (
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/Misterloveb/gomod/orm/internel/err"
)

type Model struct {
	TableName string            //表名
	Field     map[string]*Field //表字段
}

type Field struct {
	Column string //列名
}

//最多支持一级指针
func ParseModel(entity any) (*Model, error) {
	typ := reflect.TypeOf(entity)
	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		return nil, err.ErrPointerOnly
	}
	typ = typ.Elem()
	numFiled := typ.NumField()
	filemap := make(map[string]*Field, numFiled)
	for i := 0; i < numFiled; i++ {
		fd := typ.Field(i)
		filemap[fd.Name] = &Field{
			Column: UnderSourceName(fd.Name),
		}
	}
	return &Model{
		TableName: UnderSourceName(typ.Name()),
		Field:     filemap,
	}, nil
}

//大写字母转下划线
func UnderSourceName(str string) string {
	var restr strings.Builder
	if utf8.RuneCountInString(str) <= 2 {
		return strings.ToLower(str)
	}
	isup := 0
	for k, v := range str {
		if unicode.IsUpper(v) {
			isup++
		} else {
			isup = 0
		}
		if k != 0 && unicode.IsUpper(v) && isup == 1 {
			restr.WriteByte('_')
		}
		restr.WriteRune(unicode.ToLower(v))
	}
	return restr.String()
}

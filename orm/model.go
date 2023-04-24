package orm

import (
	"reflect"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"

	"github.com/Misterloveb/gomod/orm/internel/err"
)

const (
	tagName = "column" //tag自定义列名,column=XXX
)

type Registry interface {
	Get(val any) (*Model, error)
	Registry(val any, opt ...ModelOption) (*Model, error)
}
type ModelOption func(*Model) error
type Model struct {
	TableName string            //表名
	Field     map[string]*Field //表字段
}

type Field struct {
	Column string //列名
}
type registry struct {
	model sync.Map
}
type TestRegistry struct {
	registry
}

func NewRegistry() *registry {
	return &registry{}
}
func (r *registry) Get(entity any) (*Model, error) {
	typ := reflect.TypeOf(entity)
	m, ok := r.model.Load(typ)
	if !ok {
		res, err := r.Registry(entity)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
	return m.(*Model), nil
}

//最多支持一级指针
func (r *registry) Registry(entity any, opt ...ModelOption) (*Model, error) {
	typ := reflect.TypeOf(entity)
	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		return nil, err.ErrPointerOnly
	}
	typelem := typ.Elem()
	numFiled := typelem.NumField()
	filemap := make(map[string]*Field, numFiled)
	for i := 0; i < numFiled; i++ {
		fd := typelem.Field(i)
		usercol, err := r.parseTag(fd.Tag)
		if err != nil {
			return nil, err
		}
		coluser, ok := usercol[tagName]
		if ok {
			filemap[fd.Name] = &Field{
				Column: coluser,
			}
		} else {
			filemap[fd.Name] = &Field{
				Column: UnderSourceName(fd.Name),
			}
		}

	}
	var tabname string
	if tabobj, ok := entity.(TableName); ok {
		tabname = tabobj.SetTableName()
	} else {
		tabname = UnderSourceName(typelem.Name())
	}
	res := &Model{
		TableName: tabname,
		Field:     filemap,
	}
	for _, fn := range opt {
		if re := fn(res); re != nil {
			return nil, re
		}
	}
	r.model.Store(typ, res)
	return res, nil
}
func (r *registry) parseTag(tag reflect.StructTag) (map[string]string, error) {
	if len(tag) == 0 {
		return map[string]string{}, nil
	}
	tagstr, ok := tag.Lookup("orm")
	if !ok {
		return map[string]string{}, err.ErrTagNoOrm
	}
	tagslice := strings.Split(tagstr, ",")
	rsmap := make(map[string]string, len(tagslice))
	for _, tag := range tagslice {
		col := strings.Split(tag, "=")
		if len(col) != 2 {
			return map[string]string{}, err.ErrTagNoDeng
		}
		rsmap[col[0]] = col[1]
	}
	return rsmap, nil
}

//自定义表名
func ModelWithChangeTableName(name string) ModelOption {
	return func(m *Model) error {
		m.TableName = name
		return nil
	}
}

//自定义列名
func ModelWithChangeColunName(col string, newname string) ModelOption {
	return func(m *Model) error {
		res, ok := m.Field[col]
		if !ok {
			return err.ErrUnKnowColumn(col)
		}
		res.Column = newname
		return nil
	}
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

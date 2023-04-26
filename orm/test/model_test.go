package test

import (
	"reflect"
	"testing"

	"github.com/Misterloveb/gomod/orm/internel/err"
	"github.com/Misterloveb/gomod/orm/model"
	"github.com/stretchr/testify/assert"
)

func TestParseModel(t *testing.T) {
	tests := []struct {
		name    string
		entity  any
		want    *model.Model
		wantErr error
	}{
		{
			name:   "parsemodel error",
			entity: UserSetTable{},
			want: &model.Model{
				TableName: "user_table",
				FieldMap: map[string]*model.Field{
					"Name": {
						Column: "name",
						Goname: "Name",
						Ctype:  reflect.TypeOf(""),
					},
				},
				ColumnMap: map[string]*model.Field{
					"name": {
						Column: "name",
						Goname: "Name",
						Ctype:  reflect.TypeOf(""),
					},
				},
			},
			wantErr: err.ErrPointerOnly,
		},
		{
			name:   "parsemodel",
			entity: &UserSetTable{},
			want: &model.Model{
				TableName: "user_table",
				FieldMap: map[string]*model.Field{
					"Name": {
						Column: "name",
						Goname: "Name",
						Ctype:  reflect.TypeOf(""),
					},
				},
				ColumnMap: map[string]*model.Field{
					"name": {
						Column: "name",
						Goname: "Name",
						Ctype:  reflect.TypeOf(""),
					},
				},
			},
		},
		{
			name:   "use user tablename",
			entity: &UserSetTable{},
			want: &model.Model{
				TableName: "user_table",
				FieldMap: map[string]*model.Field{
					"Name": {
						Column: "name",
						Goname: "Name",
						Ctype:  reflect.TypeOf(""),
					},
				},
				ColumnMap: map[string]*model.Field{
					"name": {
						Column: "name",
						Goname: "Name",
						Ctype:  reflect.TypeOf(""),
					},
				},
			},
		},
		{
			name:   "not use user tablename",
			entity: &NotUserSetTable{},
			want: &model.Model{
				TableName: "not_user_set_table",
				FieldMap: map[string]*model.Field{
					"Name": {
						Column: "name",
						Goname: "Name",
						Ctype:  reflect.TypeOf(""),
					},
				},
				ColumnMap: map[string]*model.Field{
					"name": {
						Column: "name",
						Goname: "Name",
						Ctype:  reflect.TypeOf(""),
					},
				},
			},
		},
	}
	registry := model.NewRegistry()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := registry.Registry(tt.entity)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestChangeTableNameWithOpt(t *testing.T) {
	registry, err := model.NewRegistry().Registry(&UserModel{}, model.ModelWithChangeTableName("user_table_t"))
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	assert.Equal(t, "user_table_t", registry.TableName)
}
func TestChangeColumnNameWithOpt(t *testing.T) {
	testcase := []struct {
		name     string
		entity   any
		modelopt []model.ModelOption

		wantmodel *model.Model
		wanterr   error
	}{
		{
			name:   "change column name",
			entity: &NotUserSetTable{},
			modelopt: []model.ModelOption{
				model.ModelWithChangeColunName("Name", "first_name_t"),
			},
			wantmodel: &model.Model{
				TableName: "not_user_set_table",
				FieldMap: map[string]*model.Field{
					"Name": &model.Field{
						Column: "first_name_t",
						Goname: "Name",
						Ctype:  reflect.TypeOf(""),
					},
				},
				ColumnMap: map[string]*model.Field{
					"first_name_t": &model.Field{
						Column: "first_name_t",
						Goname: "Name",
						Ctype:  reflect.TypeOf(""),
					},
				},
			},
		},
		{
			name:     "change invoid column name",
			entity:   &NotUserSetTable{},
			modelopt: []model.ModelOption{model.ModelWithChangeColunName("age", "first_name_t")},
			wanterr:  err.ErrUnKnowColumn("age"),
		},
	}

	for _, tc := range testcase {
		t.Run(tc.name, func(t *testing.T) {
			registry, err := model.NewRegistry().Registry(tc.entity, tc.modelopt...)
			assert.Equal(t, tc.wanterr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantmodel, registry)
		})
	}
}

type UserSetTable struct {
	Name string
}
type NotUserSetTable struct {
	Name string
}

func (UserSetTable) SetTableName() string {
	return "user_table"
}

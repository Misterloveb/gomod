package test

import (
	"testing"

	"github.com/Misterloveb/gomod/orm"
	"github.com/Misterloveb/gomod/orm/internel/err"
	"github.com/stretchr/testify/assert"
)

func TestParseModel(t *testing.T) {
	tests := []struct {
		name    string
		entity  any
		want    *orm.Model
		wantErr error
	}{
		{
			name:   "parsemodel",
			entity: UserModel{},
			want: &orm.Model{
				TableName: "user_model",
				Field: map[string]*orm.Field{
					"Id": {
						Column: "id",
					},
					"FirstName": {
						Column: "first_name",
					},
					"Age": {
						Column: "age",
					},
					"LastName": {
						Column: "last_name",
					},
					"GetNAMEBYId": {
						Column: "get_namebyid",
					},
				},
			},
			wantErr: err.ErrPointerOnly,
		},
		{
			name:   "parsemodel",
			entity: &UserModel{},
			want: &orm.Model{
				TableName: "user_model",
				Field: map[string]*orm.Field{
					"Id": {
						Column: "id",
					},
					"FirstName": {
						Column: "first_name",
					},
					"Age": {
						Column: "age",
					},
					"LastName": {
						Column: "last_name",
					},
					"GetNAMEBYId": {
						Column: "get_namebyid",
					},
				},
			},
		},
		{
			name:   "use user tablename",
			entity: &UserSetTable{},
			want: &orm.Model{
				TableName: "user_table",
				Field: map[string]*orm.Field{
					"Name": {
						Column: "name",
					},
				},
			},
		},
		{
			name:   "not use user tablename",
			entity: &NotUserSetTable{},
			want: &orm.Model{
				TableName: "not_user_set_table",
				Field: map[string]*orm.Field{
					"Name": {
						Column: "name",
					},
				},
			},
		},
	}
	registry := orm.NewRegistry()
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
	registry, err := orm.NewRegistry().Registry(&UserModel{}, orm.ModelWithChangeTableName("user_table_t"))
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
		modelopt []orm.ModelOption

		wantmodel *orm.Model
		wanterr   error
	}{
		{
			name:   "change column name",
			entity: &NotUserSetTable{},
			modelopt: []orm.ModelOption{
				orm.ModelWithChangeColunName("Name", "first_name_t"),
			},
			wantmodel: &orm.Model{
				TableName: "not_user_set_table",
				Field: map[string]*orm.Field{
					"Name": &orm.Field{
						Column: "first_name_t",
					},
				},
			},
		},
		{
			name:     "change invoid column name",
			entity:   &NotUserSetTable{},
			modelopt: []orm.ModelOption{orm.ModelWithChangeColunName("age", "first_name_t")},
			wanterr:  err.ErrUnKnowColumn("age"),
		},
	}

	for _, tc := range testcase {
		t.Run(tc.name, func(t *testing.T) {
			registry, err := orm.NewRegistry().Registry(tc.entity, tc.modelopt...)
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

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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := orm.ParseModel(tt.entity)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

package test

import (
	"testing"

	"github.com/Misterloveb/gomod/orm"
	"github.com/stretchr/testify/assert"
)

func TestDeleteor(t *testing.T) {
	db, err := orm.NewDB()
	if err != nil {
		t.Fatal(err)
	}
	testcase := []struct {
		name         string
		querybuilder orm.QueryBuilder

		wantres *orm.Query
		wanterr error
	}{
		{
			name:         "delete",
			querybuilder: orm.NewDeleteor[UserModel](db).From("user_model"),
			wantres: &orm.Query{
				Sql: "DELETE FROM user_model;",
			},
		},
	}

	for _, tt := range testcase {
		t.Run(tt.name, func(t *testing.T) {
			res, err := tt.querybuilder.Build()
			assert.Equal(t, tt.wanterr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantres, res)
		})
	}
}

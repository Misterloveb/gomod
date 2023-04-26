package test

import (
	"testing"

	"github.com/Misterloveb/gomod/orm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func memoryDB(t *testing.T) *orm.DB {
	db, err := orm.Open("sqlite3", "file:test.db?cache=share&mode=memory")
	require.NoError(t, err)
	return db
}
func TestDeleteor(t *testing.T) {
	db := memoryDB(t)

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

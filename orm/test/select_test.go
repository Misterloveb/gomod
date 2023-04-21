package test

import (
	"database/sql"
	"testing"

	"github.com/Misterloveb/gomod/orm"

	"github.com/stretchr/testify/assert"
)

type UserModel struct {
	Id        int64
	FirstName string
	age       int8
	LastName  *sql.NullString
}

func TestSelector(t *testing.T) {
	testcases := []struct {
		name    string
		builder orm.QueryBuilder

		wantquery *orm.Query
		wanterror error
	}{
		{
			name:    "no from",
			builder: &orm.Selector[UserModel]{},
			wantquery: &orm.Query{
				Sql: "SELECT * FROM `UserModel`;",
			},
		},
		{
			name:    "from",
			builder: (&orm.Selector[UserModel]{}).From("user_model"),
			wantquery: &orm.Query{
				Sql: "SELECT * FROM user_model;",
			},
		},
		{
			name:    "where eq and not",
			builder: (&orm.Selector[UserModel]{}).Where(orm.C("name").Eq("lb"), orm.Not(orm.C("age").Eq(12))),
			wantquery: &orm.Query{
				Sql: "SELECT * FROM `UserModel` WHERE (`name` = ?) AND ( NOT (`age` = ?));",
				Args: []any{
					"lb", 12,
				},
			},
		},
		{
			name:    "where eq or eq",
			builder: (&orm.Selector[UserModel]{}).Where(orm.C("name").Eq("lb").Or(orm.C("age").Eq(15))),
			wantquery: &orm.Query{
				Sql: "SELECT * FROM `UserModel` WHERE (`name` = ?) OR (`age` = ?);",
				Args: []any{
					"lb", 15,
				},
			},
		},
		{
			name:    "where (eq or eq) and eq",
			builder: (&orm.Selector[UserModel]{}).Where(orm.C("name").Eq("lb").Or(orm.C("age").Eq(15)), orm.C("sex").Eq("男")),
			wantquery: &orm.Query{
				Sql: "SELECT * FROM `UserModel` WHERE ((`name` = ?) OR (`age` = ?)) AND (`sex` = ?);",
				Args: []any{
					"lb", 15, "男",
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.builder.Build()
			assert.Equal(t, err, tc.wanterror)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantquery, res)
		})
	}
}

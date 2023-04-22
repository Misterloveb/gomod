package test

import (
	"database/sql"
	"testing"

	"github.com/Misterloveb/gomod/orm"

	"github.com/stretchr/testify/assert"
)

type UserModel struct {
	Id          int64
	FirstName   string
	Age         int8
	LastName    *sql.NullString
	GetNAMEBYId string
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
				Sql: "SELECT * FROM `user_model`;",
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
			builder: (&orm.Selector[UserModel]{}).Where(orm.C("FirstName").Eq("lb"), orm.Not(orm.C("Age").Eq(12))),
			wantquery: &orm.Query{
				Sql: "SELECT * FROM `user_model` WHERE (`first_name` = ?) AND ( NOT (`age` = ?));",
				Args: []any{
					"lb", 12,
				},
			},
		},
		{
			name:    "where eq or eq",
			builder: (&orm.Selector[UserModel]{}).Where(orm.C("FirstName").Eq("lb").Or(orm.C("Age").Eq(15))),
			wantquery: &orm.Query{
				Sql: "SELECT * FROM `user_model` WHERE (`first_name` = ?) OR (`age` = ?);",
				Args: []any{
					"lb", 15,
				},
			},
		},
		{
			name:    "where (eq or eq) and eq",
			builder: (&orm.Selector[UserModel]{}).Where(orm.C("FirstName").Eq("lb").Or(orm.C("Age").Eq(15)), orm.C("GetNAMEBYId").Eq("男")),
			wantquery: &orm.Query{
				Sql: "SELECT * FROM `user_model` WHERE ((`first_name` = ?) OR (`age` = ?)) AND (`get_namebyid` = ?);",
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

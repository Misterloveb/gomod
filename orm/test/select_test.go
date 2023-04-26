package test

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Misterloveb/gomod/orm"
	"github.com/Misterloveb/gomod/orm/internel/err"
	_ "github.com/mattn/go-sqlite3"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type UserModel struct {
	Id          int64
	FirstName   string
	Age         int8
	LastName    *sql.NullString
	GetNAMEBYId string
	private     string
}

func TestSelector(t *testing.T) {
	db := memoryDB(t)
	testcases := []struct {
		name    string
		builder orm.QueryBuilder

		wantquery *orm.Query
		wanterror error
	}{
		{
			name:    "no from",
			builder: orm.NewSelector[UserModel](db),
			wantquery: &orm.Query{
				Sql: "SELECT * FROM `user_model`;",
			},
		},
		{
			name:    "from",
			builder: orm.NewSelector[UserModel](db).From("user_model"),
			wantquery: &orm.Query{
				Sql: "SELECT * FROM user_model;",
			},
		},
		{
			name:    "where eq and not",
			builder: orm.NewSelector[UserModel](db).Where(orm.C("FirstName").Eq("lb"), orm.Not(orm.C("Age").Eq(12))),
			wantquery: &orm.Query{
				Sql: "SELECT * FROM `user_model` WHERE (`first_name` = ?) AND ( NOT (`age` = ?));",
				Args: []any{
					"lb", 12,
				},
			},
		},
		{
			name:    "where eq or eq",
			builder: orm.NewSelector[UserModel](db).Where(orm.C("FirstName").Eq("lb").Or(orm.C("Age").Eq(15))),
			wantquery: &orm.Query{
				Sql: "SELECT * FROM `user_model` WHERE (`first_name` = ?) OR (`age` = ?);",
				Args: []any{
					"lb", 15,
				},
			},
		},
		{
			name:    "where (eq or eq) and eq",
			builder: orm.NewSelector[UserModel](db).Where(orm.C("FirstName").Eq("lb").Or(orm.C("Age").Eq(15)), orm.C("GetNAMEBYId").Eq("男")),
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

func TestTag(t *testing.T) {
	testcase := []struct {
		name   string
		entity any
		r      *orm.TestRegistry

		wantmodel *orm.Model
		wanter    error
	}{
		{
			name: "tag test",
			entity: func() any {
				type User struct {
					FirstName string `orm:"column=name"`
					FirstAge  string `orm:"column=firstage"`
					FirstSex  string ``
				}
				u := User{}
				return &u
			}(),
			r: &orm.TestRegistry{},
			wantmodel: &orm.Model{
				TableName: "user",
				FieldMap: map[string]*orm.Field{
					"FirstName": &orm.Field{
						Column: "name",
						Goname: "FirstName",
						Ctype:  reflect.TypeOf(""),
					},
					"FirstAge": &orm.Field{
						Column: "firstage",
						Goname: "FirstAge",
						Ctype:  reflect.TypeOf(""),
					},
					"FirstSex": &orm.Field{
						Column: "first_sex",
						Goname: "FirstSex",
						Ctype:  reflect.TypeOf(""),
					},
				},
				ColumnMap: map[string]*orm.Field{
					"name": &orm.Field{
						Column: "name",
						Goname: "FirstName",
						Ctype:  reflect.TypeOf(""),
					},
					"firstage": &orm.Field{
						Column: "firstage",
						Goname: "FirstAge",
						Ctype:  reflect.TypeOf(""),
					},
					"first_sex": &orm.Field{
						Column: "first_sex",
						Goname: "FirstSex",
						Ctype:  reflect.TypeOf(""),
					},
				},
			},
		},
	}
	for _, tc := range testcase {
		t.Run(tc.name, func(t *testing.T) {
			md, err := tc.r.Get(tc.entity)
			assert.Equal(t, tc.wanter, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantmodel, md)
		})
	}
}

func TestSelector_Get(t *testing.T) {
	mockdb, mock, errs := sqlmock.New()
	require.NoError(t, errs)
	db, errs := orm.OpenDB(mockdb)
	defer mockdb.Close()
	require.NoError(t, errs)
	//query error
	mock.ExpectQuery("SELECT .*").WillReturnError(errors.New("query error"))
	//no rows
	mock_rows := mock.NewRows([]string{"id", "first_name", "age"})
	mock.ExpectQuery("SELECT .*").WillReturnRows(mock_rows)
	//have row
	mock_rows = mock.NewRows([]string{"id", "first_name", "age", "last_name", "get_namebyid", "private"}).
		AddRow("1", "liubin", "18", "lb", "123", "priv")
	mock.ExpectQuery("SELECT .*").WillReturnRows(mock_rows)
	testcase := []struct {
		name     string
		selector *orm.Selector[UserModel]

		wantres *UserModel
		wanterr error
	}{
		{
			name:     "get invoild column",
			selector: (orm.NewSelector[UserModel](db)).Where(orm.C("sss").Eq(2)),
			wanterr:  err.ErrUnKnowColumn("sss"),
		},
		{
			name:     "query error",
			selector: (orm.NewSelector[UserModel](db).Where(orm.C("Id").Eq(2))),
			wanterr:  errors.New("query error"),
		},
		{
			name:     "query now rows",
			selector: (orm.NewSelector[UserModel](db).Where(orm.C("Id").Eq(2))),
			wanterr:  err.ErrNoRows,
		},
		{
			name:     "query now",
			selector: (orm.NewSelector[UserModel](db).Where(orm.C("Id").Eq(2))),
			wantres: &UserModel{
				Id:          1,
				FirstName:   "liubin",
				LastName:    &sql.NullString{Valid: true, String: "lb"},
				Age:         18,
				GetNAMEBYId: "123",
				private:     "",
			},
		},
	}

	for _, tc := range testcase {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.selector.Get(context.Background())
			assert.Equal(t, tc.wanterr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantres, res)
		})
	}
}

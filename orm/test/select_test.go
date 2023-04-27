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
	"github.com/Misterloveb/gomod/orm/internel/valuer"
	"github.com/Misterloveb/gomod/orm/model"
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
}
type UserModelPrivate struct {
	Id      int64
	private string
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
		r      *model.TestRegistry

		wantmodel *model.Model
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
			r: &model.TestRegistry{},
			wantmodel: &model.Model{
				TableName: "user",
				FieldMap: map[string]*model.Field{
					"FirstName": &model.Field{
						Column: "name",
						Goname: "FirstName",
						Ctype:  reflect.TypeOf(""),
					},
					"FirstAge": &model.Field{
						Column: "firstage",
						Goname: "FirstAge",
						Ctype:  reflect.TypeOf(""),
						Offset: 16,
					},
					"FirstSex": &model.Field{
						Column: "first_sex",
						Goname: "FirstSex",
						Ctype:  reflect.TypeOf(""),
						Offset: 32,
					},
				},
				ColumnMap: map[string]*model.Field{
					"name": &model.Field{
						Column: "name",
						Goname: "FirstName",
						Ctype:  reflect.TypeOf(""),
					},
					"firstage": &model.Field{
						Column: "firstage",
						Goname: "FirstAge",
						Ctype:  reflect.TypeOf(""),
						Offset: 16,
					},
					"first_sex": &model.Field{
						Column: "first_sex",
						Goname: "FirstSex",
						Ctype:  reflect.TypeOf(""),
						Offset: 32,
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
	mock_rows = mock.NewRows([]string{"id", "first_name", "age", "last_name", "get_namebyid"}).
		AddRow("1", "liubin", "18", "lb", "123")
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
			name:     "query no rows",
			selector: (orm.NewSelector[UserModel](db).Where(orm.C("Id").Eq(22))),
			wanterr:  err.ErrNoRows,
		},
		{
			name:     "query rows",
			selector: (orm.NewSelector[UserModel](db).Where(orm.C("Id").Eq(2))),
			wantres: &UserModel{
				Id:          1,
				FirstName:   "liubin",
				LastName:    &sql.NullString{Valid: true, String: "lb"},
				Age:         18,
				GetNAMEBYId: "123",
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

func TestSelector_GetMulity(t *testing.T) {
	mockdb, mock, errs := sqlmock.New()
	require.NoError(t, errs)
	db, errs := orm.OpenDB(mockdb, orm.WithColumnCreatorOpt(valuer.NewReflectValuer()))
	defer mockdb.Close()
	require.NoError(t, errs)
	//no rows
	mock_rows := mock.NewRows([]string{"id", "first_name", "age"})
	mock.ExpectQuery("SELECT .*").WillReturnRows(mock_rows)
	//have 1 row
	mock_rows = mock.NewRows([]string{"id", "first_name", "age", "last_name", "get_namebyid"}).
		AddRow("1", "liubin", "18", "lb", "123")
	mock.ExpectQuery("SELECT .*").WillReturnRows(mock_rows)
	//have 3 row
	mock_rows = mock.NewRows([]string{"id", "first_name", "age", "last_name", "get_namebyid"}).
		AddRow("1", "liubin", "18", "lb", "123").AddRow("2", "张三", "23", "zs", "ss").AddRow("3", "李四", "44", "ls", "dd")
	mock.ExpectQuery("SELECT .*").WillReturnRows(mock_rows)

	testcase := []struct {
		name     string
		selector *orm.Selector[UserModel]

		wantres []*UserModel
		wanterr error
	}{
		{
			name:     "query no rows",
			selector: (orm.NewSelector[UserModel](db).Where(orm.C("Id").Eq(22))),
			wanterr:  err.ErrNoRows,
		},
		{
			name:     "query 1 rows",
			selector: (orm.NewSelector[UserModel](db).Where(orm.C("Id").Eq(2))),
			wantres: []*UserModel{
				&UserModel{
					Id:          1,
					FirstName:   "liubin",
					LastName:    &sql.NullString{Valid: true, String: "lb"},
					Age:         18,
					GetNAMEBYId: "123",
				},
			},
		},
		{
			name:     "query 3 rows",
			selector: (orm.NewSelector[UserModel](db).Where(orm.C("Id").Eq(2))),
			wantres: []*UserModel{
				&UserModel{
					Id:          1,
					FirstName:   "liubin",
					LastName:    &sql.NullString{Valid: true, String: "lb"},
					Age:         18,
					GetNAMEBYId: "123",
				},
				&UserModel{
					Id:          2,
					FirstName:   "张三",
					LastName:    &sql.NullString{Valid: true, String: "zs"},
					Age:         23,
					GetNAMEBYId: "ss",
				},
				&UserModel{
					Id:          3,
					FirstName:   "李四",
					LastName:    &sql.NullString{Valid: true, String: "ls"},
					Age:         44,
					GetNAMEBYId: "dd",
				},
			},
		},
	}

	for _, tc := range testcase {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.selector.GetMulti(context.Background())
			assert.Equal(t, tc.wanterr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantres, res)
		})
	}
}

func TestPriveStruct_GetMulity(t *testing.T) {
	//指针方式写入结果集能返回私有成员值，反射方式不行
	mockdb, mock, errs := sqlmock.New()
	require.NoError(t, errs)
	db, errs := orm.OpenDB(mockdb, orm.WithColumnCreatorOpt(valuer.NewUnsafeValuer()))
	defer mockdb.Close()
	require.NoError(t, errs)
	//have 1 row
	mock_rows := mock.NewRows([]string{"id", "private"}).
		AddRow("1", "私有")
	mock.ExpectQuery("SELECT .*").WillReturnRows(mock_rows)
	testcase := []struct {
		name     string
		selector *orm.Selector[UserModelPrivate]

		wantres []*UserModelPrivate
		wanterr error
	}{
		{
			name:     "private struct",
			selector: (orm.NewSelector[UserModelPrivate](db).Where(orm.C("Id").Eq(2))),
			wantres: []*UserModelPrivate{
				&UserModelPrivate{
					Id:      1,
					private: "私有",
				},
			},
		},
	}

	for _, tc := range testcase {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.selector.GetMulti(context.Background())
			assert.Equal(t, tc.wanterr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantres, res)
		})
	}
}

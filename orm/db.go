package orm

import (
	"database/sql"

	"github.com/Misterloveb/gomod/orm/internel/valuer"
	"github.com/Misterloveb/gomod/orm/model"
)

type DB struct {
	r      *model.REgistry
	db     *sql.DB
	valuer valuer.Creator
}
type dbOption func(*DB)

func Open(driver, dsn string, opt ...dbOption) (*DB, error) {
	dbres, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return OpenDB(dbres, opt...)
}
func OpenDB(sqldb *sql.DB, opt ...dbOption) (*DB, error) {
	res := &DB{
		r:      model.NewRegistry(),
		db:     sqldb,
		valuer: valuer.NewUnsafeValuer(),
	}
	for _, fn := range opt {
		fn(res)
	}
	return res, nil
}

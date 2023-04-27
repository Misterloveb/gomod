package orm

import (
	"database/sql"

	"github.com/Misterloveb/gomod/orm/internel/valuer"
	"github.com/Misterloveb/gomod/orm/model"
)

type DB struct {
	r      model.Registry
	db     *sql.DB
	valuer valuer.Creator
}
type dbOption func(*DB)

func WithRegistryOpt(r model.Registry) dbOption {
	return func(d *DB) {
		d.r = r
	}
}
func WithColumnCreatorOpt(crea valuer.Creator) dbOption {
	return func(d *DB) {
		d.valuer = crea
	}
}
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

package orm

import "database/sql"

type DB struct {
	r  *registry
	db *sql.DB
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
		r:  NewRegistry(),
		db: sqldb,
	}
	for _, fn := range opt {
		fn(res)
	}
	return res, nil
}

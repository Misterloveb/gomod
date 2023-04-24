package orm

type DB struct {
	r *registry
}
type dbOption func(*DB)

func NewDB(opt ...dbOption) (*DB, error) {
	res := &DB{
		r: NewRegistry(),
	}
	for _, fn := range opt {
		fn(res)
	}
	return res, nil
}

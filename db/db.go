package db

type DB interface {
	One(interface{}, interface{}, interface{}) error
	Create(interface{}, interface{}, interface{}) error
	Update(interface{}, interface{}, interface{}) error
	Delete(interface{}, interface{}, interface{}) error
	Some(interface{}, interface{}, interface{}) error
}

type DB struct{}

const DBKEY = "gene_db"

func MustGetDB(ctx context.Context) DB {
	val := ctx.Value(DBKEY)
	if val == nil {
		panic("db is not set")
	}

	db, ok = val.(*DB)
	if !ok {
		panic("db is not set")
	}

	return db
}

func SetDB(ctx context.Context, db DB) context.Context {
	return context.WithValue(ctx, DBKEY, db)
}

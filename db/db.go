// Package db provides a simple interface for the db operations
package db

import "golang.org/x/net/context"

// DB holds the basic operation interfaces
type DB interface {
	One(interface{}, interface{}, interface{}) error
	Create(interface{}, interface{}, interface{}) error
	Update(interface{}, interface{}, interface{}) error
	Delete(interface{}, interface{}, interface{}) error
	Some(interface{}, interface{}, interface{}) error
}

// DBKEY holds the key for the db value in net.Context
const DBKEY = "gene_db"

// MustGetDB returns the DB from context, if db not found with it's key, panics
func MustGetDB(ctx context.Context) DB {
	val := ctx.Value(DBKEY)
	if val == nil {
		panic("db is not set")
	}

	d, ok := val.(DB)
	if !ok {
		panic("db is not set")
	}

	return d
}

// SetDB sets the db into context and returns the modified context
func SetDB(ctx context.Context, d DB) context.Context {
	return context.WithValue(ctx, DBKEY, d)
}

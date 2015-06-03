package twitterapi

import (
	"github.com/cihangir/gene/db"
	"github.com/fatih/invoker/models"
	"golang.org/x/net/context"
)

// New creates a new local Account handler
func NewAccount() *Account { return &Account{} }

// Account is for holding the api functions
type Account struct{}

func (a *Account) Create(ctx context.Context, req *models.Account, res *models.Account) error {
	return db.MustGetDB(ctx).Create(models.NewAccount(), req, res)
}

func (a *Account) Delete(ctx context.Context, req *models.Account, res *models.Account) error {
	return db.MustGetDB(ctx).Delete(models.NewAccount(), req, res)
}

func (a *Account) One(ctx context.Context, req *models.Account, res *models.Account) error {
	return db.MustGetDB(ctx).One(models.NewAccount(), req, res)
}

func (a *Account) Update(ctx context.Context, req *models.Account, res *models.Account) error {
	return db.MustGetDB(ctx).Update(models.NewAccount(), req, res)
}

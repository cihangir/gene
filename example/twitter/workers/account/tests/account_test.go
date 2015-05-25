package accounttests

import (
	"testing"

	"github.com/fatih/invoker/models"
	"github.com/fatih/invoker/tests"
	"github.com/fatih/invoker/workers/account/clients"
	"golang.org/x/net/context"
)

func TestAccountCreate(t *testing.T) {
	withAccountClient(t, func(c *accountclient.Account) {
		req := &models.Account{}
		res := &models.Account{}
		ctx := context.Background()
		err := c.Create(ctx, req, res)
		tests.Assert(t, err == nil, "Err should be nil while testing Account.Create")
	})
}

func TestAccountDelete(t *testing.T) {
	withAccountClient(t, func(c *accountclient.Account) {
		req := &models.Account{}
		res := &models.Account{}
		ctx := context.Background()
		err := c.Delete(ctx, req, res)
		tests.Assert(t, err == nil, "Err should be nil while testing Account.Delete")
	})
}

func TestAccountOne(t *testing.T) {
	withAccountClient(t, func(c *accountclient.Account) {
		req := &models.Account{}
		res := &models.Account{}
		ctx := context.Background()
		err := c.One(ctx, req, res)
		tests.Assert(t, err == nil, "Err should be nil while testing Account.One")
	})
}

func TestAccountSome(t *testing.T) {
	withAccountClient(t, func(c *accountclient.Account) {
		req := &models.Account{}
		res := &[]*models.Account{}
		ctx := context.Background()
		err := c.Some(ctx, req, res)
		tests.Assert(t, err == nil, "Err should be nil while testing Account.Some")
	})
}

func TestAccountUpdate(t *testing.T) {
	withAccountClient(t, func(c *accountclient.Account) {
		req := &models.Account{}
		res := &models.Account{}
		ctx := context.Background()
		err := c.Update(ctx, req, res)
		tests.Assert(t, err == nil, "Err should be nil while testing Account.Update")
	})
}

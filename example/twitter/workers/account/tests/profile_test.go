package accounttests

import (
	"testing"

	"github.com/fatih/invoker/models"
	"github.com/fatih/invoker/tests"
	"github.com/fatih/invoker/workers/account/clients"
	"golang.org/x/net/context"
)

func TestProfileCreate(t *testing.T) {
	withProfileClient(t, func(c *accountclient.Profile) {
		req := &models.Account{}
		res := &models.Account{}
		ctx := context.Background()
		err := c.Create(ctx, req, res)
		tests.Assert(t, err == nil, "Err should be nil while testing Profile.Create")
	})
}

func TestProfileDelete(t *testing.T) {
	withProfileClient(t, func(c *accountclient.Profile) {
		req := &models.Account{}
		res := &models.Account{}
		ctx := context.Background()
		err := c.Delete(ctx, req, res)
		tests.Assert(t, err == nil, "Err should be nil while testing Profile.Delete")
	})
}

func TestProfileOne(t *testing.T) {
	withProfileClient(t, func(c *accountclient.Profile) {
		req := &models.Account{}
		res := &models.Account{}
		ctx := context.Background()
		err := c.One(ctx, req, res)
		tests.Assert(t, err == nil, "Err should be nil while testing Profile.One")
	})
}

func TestProfileSome(t *testing.T) {
	withProfileClient(t, func(c *accountclient.Profile) {
		req := &models.Account{}
		res := &[]*models.Account{}
		ctx := context.Background()
		err := c.Some(ctx, req, res)
		tests.Assert(t, err == nil, "Err should be nil while testing Profile.Some")
	})
}

func TestProfileUpdate(t *testing.T) {
	withProfileClient(t, func(c *accountclient.Profile) {
		req := &models.Account{}
		res := &models.Account{}
		ctx := context.Background()
		err := c.Update(ctx, req, res)
		tests.Assert(t, err == nil, "Err should be nil while testing Profile.Update")
	})
}

package accountapi

import (
	"github.com/cihangir/gene/db"
	"github.com/cihangir/gene/example/twitter/models"
	"golang.org/x/net/context"
)

// New creates a new local Profile handler
func NewProfile() *Profile { return &Profile{} }

// Profile is for holding the api functions
type Profile struct{}

func (p *Profile) Create(ctx context.Context, req *models.Account, res *models.Account) error {
	return db.MustGetDB(ctx).Create(models.NewProfile(), req, res)
}

func (p *Profile) Delete(ctx context.Context, req *models.Account, res *models.Account) error {
	return db.MustGetDB(ctx).Delete(models.NewProfile(), req, res)
}

func (p *Profile) One(ctx context.Context, req *models.Account, res *models.Account) error {
	return db.MustGetDB(ctx).One(models.NewProfile(), req, res)
}

func (p *Profile) Update(ctx context.Context, req *models.Account, res *models.Account) error {
	return db.MustGetDB(ctx).Update(models.NewProfile(), req, res)
}

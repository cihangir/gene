package profile

import (
	"github.com/cihangir/gene/example/tinder/models"
	"golang.org/x/net/context"
)

// Profile represents a registered Account's Public Info
type ProfileService interface {
	Create(ctx context.Context, req *models.Account) (res *models.Account, err error)

	Delete(ctx context.Context, req *models.Account) (res *models.Account, err error)

	One(ctx context.Context, req *models.Account) (res *models.Account, err error)

	Update(ctx context.Context, req *models.Account) (res *models.Account, err error)
}

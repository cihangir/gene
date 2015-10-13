package profile

import (
	"github.com/cihangir/gene/example/tinder/models"
	"golang.org/x/net/context"
)

type ProfileService interface {
	Create(ctx context.Context, req *models.Profile) (res *models.Profile, err error)
	Delete(ctx context.Context, req *int64) (res *models.Profile, err error)
	MarkAs(ctx context.Context, req *models.MarkAsRequest) (res *models.Profile, err error)
	One(ctx context.Context, req *int64) (res *models.Profile, err error)
	Update(ctx context.Context, req *models.Profile) (res *models.Profile, err error)
}

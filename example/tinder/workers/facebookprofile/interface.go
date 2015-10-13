package facebookprofile

import (
	"github.com/cihangir/gene/example/tinder/models"
	"golang.org/x/net/context"
)

type FacebookProfileService interface {
	ByIDs(ctx context.Context, req *[]string) (res *[]*models.FacebookProfile, err error)
	Create(ctx context.Context, req *models.FacebookProfile) (res *models.FacebookProfile, err error)
	One(ctx context.Context, req *int64) (res *models.FacebookProfile, err error)
	Update(ctx context.Context, req *models.FacebookProfile) (res *models.FacebookProfile, err error)
}

package account

import (
	"github.com/cihangir/gene/example/tinder/models"
	"golang.org/x/net/context"
)

type AccountService interface {
	ByFacebookIDs(ctx context.Context, req *[]string) (res *[]*models.Account, err error)
	ByIDs(ctx context.Context, req *[]int64) (res *[]*models.Account, err error)
	Create(ctx context.Context, req *models.Account) (res *models.Account, err error)
	Delete(ctx context.Context, req *int64) (res *models.Account, err error)
	One(ctx context.Context, req *int64) (res *models.Account, err error)
	Update(ctx context.Context, req *models.Account) (res *models.Account, err error)
}

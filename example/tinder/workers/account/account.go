package account

import (
	"golang.org/x/net/context"

	"github.com/cihangir/gene/example/tinder/models"
)

type account struct{}

func NewAccount() AccountService {
	return account{}
}

func (mw account) ServiceName() string {
	return "account"
}

func (mw account) ByFacebookIDs(ctx context.Context, req *[]string) (res *[]*models.Account, err error) {
	return nil, nil
}

func (mw account) ByIDs(ctx context.Context, req *[]int64) (res *[]*models.Account, err error) {
	return nil, nil
}

func (mw account) Create(ctx context.Context, req *models.Account) (res *models.Account, err error) {
	return nil, nil
}

func (mw account) Delete(ctx context.Context, req *int64) (res *models.Account, err error) {
	return nil, nil
}

func (mw account) One(ctx context.Context, req *int64) (res *models.Account, err error) {
	return nil, nil
}

func (mw account) Update(ctx context.Context, req *models.Account) (res *models.Account, err error) {
	return nil, nil
}

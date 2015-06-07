package functions

import (
	"encoding/json"

	"testing"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/gene/testdata"
	"github.com/cihangir/schema"
)

func TestFunctions(t *testing.T) {
	s := &schema.Schema{}
	err := json.Unmarshal([]byte(testdata.TestDataFull), s)
	common.TestEquals(t, nil, err)

	s = s.Resolve(s)
	context := common.NewContext()

	a, err := generate(context, "test", s.Definitions["Account"])
	common.TestEquals(t, nil, err)
	common.TestEquals(t, expected, string(a))
}

const expected = `package testapi

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

func (a *Account) Some(ctx context.Context, req *models.Account, res *[]*models.Account) error {
	return db.MustGetDB(ctx).Some(models.NewAccount(), req, res)
}

func (a *Account) Update(ctx context.Context, req *models.Account, res *models.Account) error {
	return db.MustGetDB(ctx).Update(models.NewAccount(), req, res)
}
`

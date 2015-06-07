package clients

import (
	"encoding/json"

	"testing"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/gene/testdata"
	"github.com/cihangir/schema"
)

func TestClients(t *testing.T) {
	s := &schema.Schema{}
	if err := json.Unmarshal([]byte(testdata.JSON1), s); err != nil {
		t.Fatal(err.Error())
	}

	s = s.Resolve(s)

	a, err := New().Generate(common.NewContext(), s)
	common.TestEquals(t, nil, err)
	common.TestEquals(t, expected, string(a[0].Content))
}

const expected = `package accountclient

import (
	"github.com/youtube/vitess/go/rpcplus"
	"golang.org/x/net/context"
)

// New creates a new local Account rpc client
func NewAccount(client *rpcplus.Client) *Account {
	return &Account{
		client: client,
	}
}

// Account is for holding the api functions
type Account struct {
	client *rpcplus.Client
}

func (a *Account) Create(ctx context.Context, req *models.Account, res *models.Account) error {
	return a.client.Call(ctx, "Account.Create", req, res)
}

func (a *Account) Delete(ctx context.Context, req *models.Account, res *models.Account) error {
	return a.client.Call(ctx, "Account.Delete", req, res)
}

func (a *Account) One(ctx context.Context, req *models.Account, res *models.Account) error {
	return a.client.Call(ctx, "Account.One", req, res)
}

func (a *Account) Some(ctx context.Context, req *models.Account, res *[]*models.Account) error {
	return a.client.Call(ctx, "Account.Some", req, res)
}

func (a *Account) Update(ctx context.Context, req *models.Account, res *models.Account) error {
	return a.client.Call(ctx, "Account.Update", req, res)
}
`

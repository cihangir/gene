package clients

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"

	"testing"

	"github.com/cihangir/gene/testdata"
	"github.com/cihangir/schema"
)

const expected = `package testclient

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
	return m.client.Call(ctx, "Account.Create", req, res)
}

func (a *Account) Delete(ctx context.Context, req *models.Account, res *models.Account) error {
	return m.client.Call(ctx, "Account.Delete", req, res)
}

func (a *Account) One(ctx context.Context, req *models.Account, res *models.Account) error {
	return m.client.Call(ctx, "Account.One", req, res)
}

func (a *Account) Some(ctx context.Context, req *models.Account, res *[]*models.Account) error {
	return m.client.Call(ctx, "Account.Some", req, res)
}

func (a *Account) Update(ctx context.Context, req *models.Account, res *models.Account) error {
	return m.client.Call(ctx, "Account.Update", req, res)
}
`

func TestConstructors(t *testing.T) {
	s := &schema.Schema{}
	if err := json.Unmarshal([]byte(testdata.JSON1), s); err != nil {
		t.Fatal(err.Error())
	}

	s = s.Resolve(nil)

	a, err := generate("test", s.Definitions["Account"])
	equals(t, nil, err)
	equals(t, expected, string(a))
}

func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.Fail()
	}
}

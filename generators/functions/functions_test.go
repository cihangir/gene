package functions

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

func TestFunctions(t *testing.T) {
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

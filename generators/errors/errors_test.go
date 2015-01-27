package errors

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

const expected = `package errs

var (
	ErrAccountCreatedAtNotSet              = errors.New("Account.CreatedAt not set")
	ErrAccountEmailAddressNotSet           = errors.New("Account.EmailAddress not set")
	ErrAccountEmailStatusConstantNotSet    = errors.New("Account.EmailStatusConstant not set")
	ErrAccountIdNotSet                     = errors.New("Account.Id not set")
	ErrAccountPasswordNotSet               = errors.New("Account.Password not set")
	ErrAccountPasswordStatusConstantNotSet = errors.New("Account.PasswordStatusConstant not set")
	ErrAccountProfileIdNotSet              = errors.New("Account.ProfileId not set")
	ErrAccountSaltNotSet                   = errors.New("Account.Salt not set")
	ErrAccountStatusConstantNotSet         = errors.New("Account.StatusConstant not set")
	ErrAccountURLNotSet                    = errors.New("Account.URL not set")
	ErrAccountURLNameNotSet                = errors.New("Account.URLName not set")
	ErrConfigMongoNotSet                   = errors.New("Config.Mongo not set")
	ErrConfigPostgresNotSet                = errors.New("Config.Postgres not set")
	ErrProfileAvatarURLNotSet              = errors.New("Profile.AvatarURL not set")
	ErrProfileCreatedAtNotSet              = errors.New("Profile.CreatedAt not set")
	ErrProfileFirstNameNotSet              = errors.New("Profile.FirstName not set")
	ErrProfileIdNotSet                     = errors.New("Profile.Id not set")
	ErrProfileLastNameNotSet               = errors.New("Profile.LastName not set")
	ErrProfileNickNotSet                   = errors.New("Profile.Nick not set")
)
`

func TestConstructors(t *testing.T) {
	var s schema.Schema
	if err := json.Unmarshal([]byte(testdata.JSON1), &s); err != nil {
		t.Fatal(err.Error())
	}

	a, err := generate(&s)
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

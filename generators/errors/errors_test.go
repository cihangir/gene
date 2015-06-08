package errors

import (
	"encoding/json"

	"testing"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/gene/testdata"
	"github.com/cihangir/schema"
)

func TestErrors(t *testing.T) {
	s := &schema.Schema{}
	if err := json.Unmarshal([]byte(testdata.JSON1), s); err != nil {
		t.Fatal(err.Error())
	}

	s = s.Resolve(s)

	sts, err := (&Generator{}).Generate(common.NewContext(), s)
	common.TestEquals(t, nil, err)

	for i, s := range sts {
		common.TestEquals(t, expecteds[i], string(s.Content))
	}
}

var expecteds = []string{`package errs

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
)
`,
	`package errs

var (
	ErrProfileAvatarURLNotSet = errors.New("Profile.AvatarURL not set")
	ErrProfileCreatedAtNotSet = errors.New("Profile.CreatedAt not set")
	ErrProfileFirstNameNotSet = errors.New("Profile.FirstName not set")
	ErrProfileIdNotSet        = errors.New("Profile.Id not set")
	ErrProfileLastNameNotSet  = errors.New("Profile.LastName not set")
	ErrProfileNickNotSet      = errors.New("Profile.Nick not set")
)
`,
}

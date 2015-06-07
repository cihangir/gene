package errs

import "errors"

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

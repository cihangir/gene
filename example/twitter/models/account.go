// Generated struct for Account.
package models

import (
	"time"

	"github.com/cihangir/govalidator"
)

// AccountEmailStatusConstant holds the predefined enums
var AccountEmailStatusConstant = struct {
	Verified    string
	NotVerified string
}{
	Verified:    "verified",
	NotVerified: "notVerified",
}

// AccountPasswordStatusConstant holds the predefined enums
var AccountPasswordStatusConstant = struct {
	Valid      string
	NeedsReset string
	Generated  string
}{
	Valid:      "valid",
	NeedsReset: "needsReset",
	Generated:  "generated",
}

// AccountStatusConstant holds the predefined enums
var AccountStatusConstant = struct {
	Registered              string
	Unregistered            string
	NeedsManualVerification string
}{
	Registered:              "registered",
	Unregistered:            "unregistered",
	NeedsManualVerification: "needsManualVerification",
}

// Account represents a registered User
type Account struct {
	// The unique identifier for a Account
	ID int64 `json:"id,omitempty,string"`
	// The unique identifier for a Account's Profile
	ProfileID int64 `json:"profileId,omitempty,string"`
	// Salted Password of the Account
	Password string `json:"password"`
	// Unique URL of the Account
	URL string `json:"url,omitempty"`
	// Status of the Account's Password
	PasswordStatusConstant string `json:"passwordStatusConstant,omitempty"`
	// Salt used to hash Password of the Account
	Salt string `json:"salt,omitempty"`
	// Email Address of the Account
	EmailAddress string `json:"emailAddress"`
	// Status of the Account's Email
	EmailStatusConstant string `json:"emailStatusConstant,omitempty"`
	// Status of the Account
	StatusConstant string `json:"statusConstant,omitempty"`
	// Profile's creation time
	CreatedAt time.Time `json:"createdAt,omitempty"`
}

// NewAccount creates a new Account struct with default values
func NewAccount() *Account {
	return &Account{
		CreatedAt:              time.Now().UTC(),
		EmailStatusConstant:    AccountEmailStatusConstant.NotVerified,
		PasswordStatusConstant: AccountPasswordStatusConstant.Valid,
		StatusConstant:         AccountStatusConstant.Registered,
	}
}

// Validate validates the Account struct
func (a *Account) Validate() error {
	return govalidator.NewMulti(govalidator.Date(a.CreatedAt),
		govalidator.MaxLength(a.Salt, 255),
		govalidator.Min(float64(a.ID), 1.000000),
		govalidator.Min(float64(a.ProfileID), 1.000000),
		govalidator.MinLength(a.Password, 6),
		govalidator.MinLength(a.URL, 6),
		govalidator.OneOf(a.EmailStatusConstant, []string{
			AccountEmailStatusConstant.Verified,
			AccountEmailStatusConstant.NotVerified,
		}),
		govalidator.OneOf(a.PasswordStatusConstant, []string{
			AccountPasswordStatusConstant.Valid,
			AccountPasswordStatusConstant.NeedsReset,
			AccountPasswordStatusConstant.Generated,
		}),
		govalidator.OneOf(a.StatusConstant, []string{
			AccountStatusConstant.Registered,
			AccountStatusConstant.Unregistered,
			AccountStatusConstant.NeedsManualVerification,
		})).Validate()
}

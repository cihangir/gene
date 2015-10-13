// Generated struct for Profile.
package models

import (
	"time"

	"github.com/cihangir/govalidator"
)

// Profile represents a registered Account's Public Info
type Profile struct {
	// The unique identifier for a Account's Profile
	ID int64 `json:"id,omitempty,string"`
	// Full name associated with the profile. Maximum of 20 characters.
	ScreenName string `json:"screenName"`
	// URL associated with the profile. Will be prepended with “http://” if not
	// present. Maximum of 100 characters.
	URL string `json:"url,omitempty"`
	// The city or country describing where the user of the account is located. The
	// contents are not normalized or geocoded in any way. Maximum of 30 characters.
	Location string `json:"location,omitempty"`
	// A description of the user owning the account. Maximum of 160 characters.
	Description string `json:"description,omitempty"`
	// Sets a hex value that controls the color scheme of links used on the
	// authenticating user’s profile page on twitter.com. This must be a valid
	// hexadecimal value, and may be either three or six characters (ex: F00 or
	// FF0000).
	LinkColor string `json:"linkColor,omitempty"`
	// URL of the Account's Avatar
	AvatarURL string `json:"avatarUrl,omitempty"`
	// Profile's creation time
	CreatedAt time.Time `json:"createdAt,omitempty"`
}

// NewProfile creates a new Profile struct with default values
func NewProfile() *Profile {
	return &Profile{
		CreatedAt: time.Now().UTC(),
		LinkColor: "FF0000",
	}
}

// Validate validates the Profile struct
func (p *Profile) Validate() error {
	return govalidator.NewMulti(govalidator.Date(p.CreatedAt),
		govalidator.MaxLength(p.AvatarURL, 2000),
		govalidator.MaxLength(p.Description, 160),
		govalidator.MaxLength(p.LinkColor, 6),
		govalidator.MaxLength(p.Location, 30),
		govalidator.MaxLength(p.ScreenName, 20),
		govalidator.MaxLength(p.URL, 100),
		govalidator.Min(float64(p.ID), 1.000000),
		govalidator.MinLength(p.ScreenName, 4)).Validate()
}

// Generated struct for Tweet.
package models

import (
	"time"

	"github.com/cihangir/govalidator"
)

// Tweet represents a post a user created
type Tweet struct {
	// The unique identifier for a Tweet
	ID int64 `json:"id,omitempty,string"`
	// The unique identifier for a Account's Profile
	ProfileID int64 `json:"profileId,string"`
	// Text of the Tweet
	Body string `json:"body"`
	// Location of the Tweet's origin
	Location string `json:"location,omitempty"`
	// Aggregated Count for re-tweets of a tweet
	RetweetCount int32 `json:"retweetCount,omitempty"`
	// Aggregated Count for favourites of a tweet
	FavouritesCount int32 `json:"favouritesCount,omitempty"`
	// Mark tweet if it is possibly sensitive
	PossiblySensitive bool `json:"possiblySensitive,omitempty"`
	// Tweet's creation time
	CreatedAt time.Time `json:"createdAt,omitempty"`
}

// NewTweet creates a new Tweet struct with default values
func NewTweet() *Tweet {
	return &Tweet{
		CreatedAt: time.Now().UTC(),
	}
}

// Validate validates the Tweet struct
func (t *Tweet) Validate() error {
	return govalidator.NewMulti(govalidator.Date(t.CreatedAt),
		govalidator.Min(float64(t.ProfileID), 1.000000),
		govalidator.MinLength(t.Body, 1)).Validate()
}

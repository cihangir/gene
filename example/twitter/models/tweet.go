// Generated struct for models.
package models

import (
	"time"

	"github.com/cihangir/govalidator"
)

// Tweet represents a post a user created
type Tweet struct {
	ID                int64     `json:"id,omitempty,string"`         // The unique identifier for a Tweet
	ProfileID         int64     `json:"profileId,string"`            // The unique identifier for a Account's Profile
	Body              string    `json:"body"`                        // Text of the Tweet
	Location          string    `json:"location,omitempty"`          // Location of the Tweet's origin
	RetweetCount      int32     `json:"retweetCount,omitempty"`      // Aggregated Count for re-tweets of a tweet
	FavouritesCount   int32     `json:"favouritesCount,omitempty"`   // Aggregated Count for favourites of a tweet
	PossiblySensitive bool      `json:"possiblySensitive,omitempty"` // Mark tweet if it is possibly sensitive
	CreatedAt         time.Time `json:"createdAt,omitempty"`         // Tweet's creation time
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

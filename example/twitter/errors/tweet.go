package errs

import "errors"

var (
	ErrTweetBodyNotSet              = errors.New("Tweet.Body not set")
	ErrTweetCreatedAtNotSet         = errors.New("Tweet.CreatedAt not set")
	ErrTweetFavouritesCountNotSet   = errors.New("Tweet.FavouritesCount not set")
	ErrTweetIDNotSet                = errors.New("Tweet.ID not set")
	ErrTweetLocationNotSet          = errors.New("Tweet.Location not set")
	ErrTweetPossiblySensitiveNotSet = errors.New("Tweet.PossiblySensitive not set")
	ErrTweetProfileIDNotSet         = errors.New("Tweet.ProfileID not set")
	ErrTweetRetweetCountNotSet      = errors.New("Tweet.RetweetCount not set")
)

package errs

import "errors"

var (
	ErrTweetBodyNotSet              = errors.New("Tweet.Body not set")
	ErrTweetCreatedAtNotSet         = errors.New("Tweet.CreatedAt not set")
	ErrTweetFavouritesCountNotSet   = errors.New("Tweet.FavouritesCount not set")
	ErrTweetIdNotSet                = errors.New("Tweet.Id not set")
	ErrTweetLocationNotSet          = errors.New("Tweet.Location not set")
	ErrTweetPossiblySensitiveNotSet = errors.New("Tweet.PossiblySensitive not set")
	ErrTweetProfileIdNotSet         = errors.New("Tweet.ProfileId not set")
	ErrTweetRetweetCountNotSet      = errors.New("Tweet.RetweetCount not set")
)

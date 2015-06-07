package errs

import "errors"

var (
	ErrProfileAvatarURLNotSet   = errors.New("Profile.AvatarURL not set")
	ErrProfileCreatedAtNotSet   = errors.New("Profile.CreatedAt not set")
	ErrProfileDescriptionNotSet = errors.New("Profile.Description not set")
	ErrProfileIdNotSet          = errors.New("Profile.Id not set")
	ErrProfileLinkColorNotSet   = errors.New("Profile.LinkColor not set")
	ErrProfileLocationNotSet    = errors.New("Profile.Location not set")
	ErrProfileScreenNameNotSet  = errors.New("Profile.ScreenName not set")
	ErrProfileURLNotSet         = errors.New("Profile.URL not set")
)

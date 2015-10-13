package facebookfriends

import (
	"github.com/cihangir/gene/example/tinder/models"
	"golang.org/x/net/context"
)

type FacebookFriendsService interface {
	Create(ctx context.Context, req *models.FacebookFriends) (res *models.FacebookFriends, err error)
	Delete(ctx context.Context, req *models.FacebookFriends) (res *models.FacebookFriends, err error)
	Mutuals(ctx context.Context, req *[]*models.FacebookFriends) (res *[]string, err error)
	One(ctx context.Context, req *models.FacebookFriends) (res *models.FacebookFriends, err error)
}

// Package team6 contains code for team 6's client implementation
package team6

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const id = shared.Team6

const (
	hostile Friendship = iota
	neutral
	friendly
)

type GiftReceivedHistory map[shared.ClientID]shared.GiftOffer

type Friendship int

type FriendshipWith map[shared.ClientID]Friendship

func init() {
	baseclient.RegisterClient(
		id,
		&client{
			BaseClient: baseclient.NewClient(id),
			friendshipWith: FriendshipWith{
				shared.Team1: neutral,
				shared.Team2: neutral,
				shared.Team3: neutral,
				shared.Team4: neutral,
				shared.Team5: neutral,
			},
		},
	)
}

type client struct {
	*baseclient.BaseClient

	friendshipWith FriendshipWith
}

func (c *client) RaiseFriendshipLevel() {
	// increase the friendship level with some other islands
}

func (c *client) LowerFriendshipLevel() {
	// decrease the friendship level with some other islands
}

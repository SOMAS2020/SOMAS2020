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

// GiftReceivedHistory is what gifts that we have received from which islands
type GiftReceivedHistory map[shared.ClientID]shared.GiftOffer

// Friendship indicates friendship level
type Friendship int

// FriendshipWith indicates level with other islands
type FriendshipWith map[shared.ClientID]Friendship

// Config configures our island's initial state
type Config struct {
	friendshipWith FriendshipWith
}

func init() {
	baseclient.RegisterClient(
		id,
		&client{
			BaseClient: baseclient.NewClient(id),
			configInfo: clientConfigInfo,
		},
	)
}

type client struct {
	*baseclient.BaseClient

	configInfo Config
}

// increases the friendship level with some other islands
func (c *client) RaiseFriendshipLevel(clientID shared.ClientID) {
	currFriendship := c.configInfo.friendshipWith[clientID]

	switch {
	case currFriendship == hostile:
		c.configInfo.friendshipWith[clientID] = neutral
	case currFriendship == neutral:
		c.configInfo.friendshipWith[clientID] = friendly
	}
}

// decreases the friendship level with some other islands
func (c *client) LowerFriendshipLevel(clientID shared.ClientID) {
	currFriendship := c.configInfo.friendshipWith[clientID]

	switch {
	case currFriendship == friendly:
		c.configInfo.friendshipWith[clientID] = neutral
	case currFriendship == neutral:
		c.configInfo.friendshipWith[clientID] = hostile
	}
}

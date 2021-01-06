// Package team6 contains code for team 6's client implementation
package team6

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const (
	id = shared.Team6
)

type client struct {
	*baseclient.BaseClient

	friendship            Friendship
	giftsSentHistory      GiftsSentHistory
	giftsReceivedHistory  GiftsReceivedHistory
	giftsRequestedHistory GiftsRequestedHistory
	forageHistory         ForageHistory
	favourRules           FavourRules
	payingTax             shared.Resources

	clientConfig ClientConfig
}

func init() {
	baseclient.RegisterClient(id, NewClient(id))
}

// ########################
// ######  General  #######
// ########################

// NewClient creates a client objects for our island
func NewClient(clientID shared.ClientID) baseclient.Client {
	return &client{
		BaseClient:            baseclient.NewClient(clientID),
		friendship:            friendship,
		giftsSentHistory:      giftsSentHistory,
		giftsReceivedHistory:  giftsReceivedHistory,
		giftsRequestedHistory: giftsRequestedHistory,
		forageHistory:         forageHistory,
		favourRules:           favourRules,
		clientConfig:          clientConfig,
	}
}

// ------ TODO: OPTIONAL ------
func (c *client) Initialise(serverReadHandle baseclient.ServerReadHandle) {
	c.ServerReadHandle = serverReadHandle
}

func (c *client) StartOfTurn() {
	defer c.Logf("There are [%v] islands left in this game", c.getNumOfAliveIslands())

	c.updateConfig()
	c.updateFriendship()
}

// updateConfig will be called at the start of each turn to set the newest config
func (c *client) updateConfig() {
	defer c.Logf("Agent[%v] configuration has been updated", c.BaseClient.GetID())

	ourResources := c.ServerReadHandle.GetGameState().ClientInfo.Resources
	criticalCounter := c.ServerReadHandle.GetGameState().ClientInfo.CriticalConsecutiveTurnsCounter
	minThreshold := c.ServerReadHandle.GetGameConfig().MinimumResourceThreshold
	costOfLiving := c.ServerReadHandle.GetGameConfig().CostOfLiving
	maxCriticalCounter := c.ServerReadHandle.GetGameConfig().MaxCriticalConsecutiveTurns

	updatedConfig := ClientConfig{
		minFriendship:          0.0,
		maxFriendship:          100.0,
		friendshipChangingRate: 20.0,
		selfishThreshold:       minThreshold + 3.0*costOfLiving + ourResources/10.0,
		normalThreshold:        minThreshold + 6.0*costOfLiving + ourResources/10.0,
		payingTax:              shared.Resources(criticalCounter / maxCriticalCounter),
	}

	c.clientConfig = ClientConfig(updatedConfig)
}

// updateFriendship will be called at the start of each turn to update our friendships
func (c *client) updateFriendship() {
	defer c.Logf("Agent[%v] friendship information has been updated", c.BaseClient.GetID())

	for team, requested := range c.giftsRequestedHistory {
		if c.ServerReadHandle.GetGameState().ClientLifeStatuses[team] != shared.Alive {
			// doesn't judge if they are not able to survive themselves
			continue
		} else {
			received := c.giftsReceivedHistory[team]
			offered := c.giftsSentHistory[team]

			if received < offered && received < requested {
				// will be sad if the island give us very little
				c.lowerFriendshipLevel(team, c.clientConfig.friendshipChangingRate*FriendshipLevel(offered/(received+requested+shared.Resources(1.0))))
			} else if received >= offered && received >= requested {
				c.raiseFriendshipLevel(team, c.clientConfig.friendshipChangingRate*FriendshipLevel(received/(offered+requested+shared.Resources(1.0))))
			} else {
				// keeps the current friendship level
				continue
			}
		}
	}
}

// ------ TODO: OPTIONAL ------
func (c *client) DisasterNotification(dR disasters.DisasterReport, effects disasters.DisasterEffects) {
	c.BaseClient.DisasterNotification(dR, effects)
}

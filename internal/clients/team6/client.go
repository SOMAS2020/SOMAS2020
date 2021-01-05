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

	config Config
}

func init() {
	baseclient.RegisterClient(
		id,
		&client{
			BaseClient:            baseclient.NewClient(id),
			friendship:            friendship,
			giftsSentHistory:      giftsSentHistory,
			giftsReceivedHistory:  giftsReceivedHistory,
			giftsRequestedHistory: giftsRequestedHistory,
			forageHistory:         forageHistory,
			favourRules:           favourRules,
			config:                config,
		},
	)
}

// ########################
// ######  General  #######
// ########################

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

	updatedConfig := Config{
		minFriendship:          0.0,
		maxFriendship:          100.0,
		friendshipChangingRate: 20.0,
		selfishThreshold:       minThreshold + 3.0*costOfLiving + ourResources/10.0,
		normalThreshold:        minThreshold + 6.0*costOfLiving + ourResources/10.0,
		payingTax:              shared.Resources(criticalCounter / maxCriticalCounter),
	}

	c.config = Config(updatedConfig)
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
				c.lowerFriendshipLevel(team, c.config.friendshipChangingRate*FriendshipLevel(offered/(received+requested+shared.Resources(1.0))))
			} else if received >= offered && received >= requested {
				c.raiseFriendshipLevel(team, c.config.friendshipChangingRate*FriendshipLevel(received/(offered+requested+shared.Resources(1.0))))
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

// ########################
// ######    Utils   ######
// ########################

// increases the friendship level with some other islands
func (c *client) raiseFriendshipLevel(clientID shared.ClientID, increment FriendshipLevel) {
	currFriendship := c.friendship[clientID]
	logIncrement := (increment / (c.friendship[clientID] + increment)) * (c.config.maxFriendship / 5)
	raisedFriendship := currFriendship + logIncrement

	if raisedFriendship > c.config.maxFriendship {
		c.Logf("Friendship with island[%v] is at maximum!", clientID)
		c.friendship[clientID] = c.config.maxFriendship
	} else {
		c.friendship[clientID] = raisedFriendship
	}
}

// decreases the friendship level with some other islands
func (c *client) lowerFriendshipLevel(clientID shared.ClientID, deduction FriendshipLevel) {
	currFriendship := c.friendship[clientID]
	logDeduction := (deduction / (c.friendship[clientID] + deduction)) * (c.config.maxFriendship / 5)
	loweredFriendship := currFriendship - logDeduction

	if loweredFriendship > c.config.minFriendship {
		c.Logf("Friendship with island[%v] is at minimum!", clientID)
		c.friendship[clientID] = c.config.minFriendship
	} else {
		c.friendship[clientID] = loweredFriendship
	}
}

// gets friendship cofficients
func (c client) getFriendshipCoeffs() map[shared.ClientID]float64 {
	friendshipCoeffs := make(map[shared.ClientID]float64)

	for team, fs := range c.friendship {
		friendshipCoeffs[team] = float64(fs / c.config.maxFriendship)
	}

	return friendshipCoeffs
}

// gets our personality
func (c client) getPersonality() Personality {
	ourResources := c.ServerReadHandle.GetGameState().ClientInfo.Resources

	if ourResources <= c.config.selfishThreshold {
		return Selfish
	} else if ourResources <= c.config.normalThreshold {
		return Normal
	}

	// TODO: more cases to have?
	return Generous
}

// gets the number of alive islands
func (c client) getNumOfAliveIslands() uint {
	num := uint(0)

	for _, status := range c.ServerReadHandle.GetGameState().ClientLifeStatuses {
		if status == shared.Alive {
			num++
		}
	}

	return num
}

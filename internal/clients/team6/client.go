// Package team6 contains code for team 6's client implementation
package team6

import (
	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const (
	id = shared.Team6
)

type client struct {
	*baseclient.BaseClient

	friendship            Friendship
	trustRank             TrustRank
	giftsSentHistory      GiftsSentHistory
	giftsReceivedHistory  GiftsReceivedHistory
	giftsRequestedHistory GiftsRequestedHistory
	disastersHistory      DisastersHistory
	disasterPredictions   DisasterPredictions
	forageHistory         ForageHistory
	taxDemanded           shared.Resources
	sanctionDemanded      shared.Resources
	allocationAllowed     shared.Resources

	clientConfig ClientConfig
}

// DefaultClient creates the client that will be used for most simulations. All
// other personalities are considered alternatives. To give a different
// personality for your agent simply create another (exported) function with the
// same signature as "DefaultClient" that creates a different agent, and inform
// someone on the simulation team that you would like it to be included in
// testing
func DefaultClient(id shared.ClientID) baseclient.Client {
	return NewClient(id)
}

// NewClient creates a client objects for our island
func NewClient(clientID shared.ClientID) baseclient.Client {
	return &client{
		BaseClient:   baseclient.NewClient(clientID),
		clientConfig: getClientConfig(),
	}
}

func (c *client) Initialise(serverReadHandle baseclient.ServerReadHandle) {
	c.ServerReadHandle = serverReadHandle
	c.LocalVariableCache = rules.CopyVariableMap(serverReadHandle.GetGameState().RulesInfo.VariableMap)

	c.friendship = Friendship{}
	c.trustRank = TrustRank{}
	c.giftsSentHistory = GiftsSentHistory{}
	c.giftsReceivedHistory = GiftsReceivedHistory{}
	c.giftsRequestedHistory = GiftsRequestedHistory{}
	c.disastersHistory = DisastersHistory{}
	c.disasterPredictions = DisasterPredictions{}
	c.forageHistory = ForageHistory{}
	c.taxDemanded = 0.0
	c.sanctionDemanded = 0.0
	c.allocationAllowed = 0.0

	for _, team := range shared.TeamIDs {
		if team == c.GetID() {
			c.friendship[team] = c.clientConfig.maxFriendship
			c.trustRank[team] = 1

			continue
		}

		c.friendship[team] = 50
		c.trustRank[team] = 0.5
	}

	c.updateConfig()
}

func (c *client) StartOfTurn() {
	defer c.Logf("There are %v islands left in this game", c.getNumOfAliveIslands())

	for _, team := range shared.TeamIDs[:] {
		if team == c.GetID() {
			continue
		}
		if c.giftsReceivedHistory[team] == 0 && c.ServerReadHandle.GetGameState().Turn != 1 {
			c.lowerFriendshipLevel(team, 99*c.clientConfig.friendshipChangingRate)
		}
	}

	c.updateConfig()
}

// updateConfig will be called at the start of each turn to set the newest config
func (c *client) updateConfig() {
	defer c.Logf("Configuration status: %v", c.clientConfig)

	minThreshold := c.ServerReadHandle.GetGameConfig().MinimumResourceThreshold
	costOfLiving := c.ServerReadHandle.GetGameConfig().CostOfLiving

	updatedConfig := ClientConfig{
		minFriendship:          c.clientConfig.minFriendship,
		maxFriendship:          c.clientConfig.maxFriendship,
		friendshipChangingRate: c.clientConfig.friendshipChangingRate,
		selfishThreshold:       minThreshold + 3.0*costOfLiving + c.ServerReadHandle.GetGameState().CommonPool/12.0,
		normalThreshold:        minThreshold + 12.0*costOfLiving + c.ServerReadHandle.GetGameState().CommonPool/6.0,
		multiplier:             c.clientConfig.multiplier,
	}

	c.clientConfig = ClientConfig(updatedConfig)
}

// updateFriendship will be called every time a gift is exchanged; neg for sending, pos for receiving
func (c *client) updateFriendship(giftAmount shared.Resources, team shared.ClientID) {
	defer c.Logf("Friendship status: %v", c.friendship)

	friendshipChangingRate := c.clientConfig.friendshipChangingRate
	receivedSum := c.giftsReceivedHistory[team]
	sentSum := c.giftsSentHistory[team]
	requestedSum := c.giftsReceivedHistory[team]

	if sentSum+requestedSum+receivedSum == 0 {
		return
	}

	if receivedSum >= sentSum && giftAmount > 0 {
		c.raiseFriendshipLevel(team, friendshipChangingRate*FriendshipLevel(receivedSum-sentSum+giftAmount))
	} else if receivedSum < sentSum && giftAmount < 0 {
		c.lowerFriendshipLevel(team, friendshipChangingRate*FriendshipLevel(sentSum-receivedSum-giftAmount))
	}

}

func (c *client) DisasterNotification(dR disasters.DisasterReport, effects disasters.DisasterEffects) {
	defer c.Logf("Trust rank status: %v", c.trustRank)

	currTurn := c.ServerReadHandle.GetGameState().Turn

	disasterhappening := baseclient.DisasterInfo{
		CoordinateX: dR.X,
		CoordinateY: dR.Y,
		Magnitude:   dR.Magnitude,
		Turn:        currTurn,
	}

	ourDiff := math.Abs(c.disasterPredictions[c.GetID()].Magnitude - disasterhappening.Magnitude)
	theirDiff := float64(0)

	for team, prediction := range c.disasterPredictions {
		theirDiff = math.Abs(prediction.Magnitude - disasterhappening.Magnitude)

		if ourDiff != 0 {
			c.trustRank[team] += (ourDiff - theirDiff) / ourDiff
		} else {
			c.trustRank[team] = c.trustRank[team] / float64(2)
		}

		// sets the caps of trust ranks from 0 to 1
		if c.trustRank[team] < 0 {
			c.trustRank[team] = 0
		} else if c.trustRank[team] > 1 {
			c.trustRank[team] = 1
		}
	}

	c.disastersHistory = append(c.disastersHistory, disasterhappening)
}

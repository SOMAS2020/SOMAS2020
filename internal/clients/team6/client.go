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
	payingTax             shared.Resources
	payingSanction        shared.Resources

	clientConfig ClientConfig
}

func init() {
	baseclient.RegisterClientFactory(id, func() baseclient.Client { return NewClient(id) })
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
	c.payingTax = 0.0
	c.payingSanction = 0.0

	for _, team := range shared.TeamIDs {
		if team == c.GetID() {
			continue
		}

		c.friendship[team] = 50
		c.trustRank[team] = 0.5
	}
}

func (c *client) StartOfTurn() {
	defer c.Logf("There are %v islands left in this game", c.getNumOfAliveIslands())

	c.updateConfig()
	c.updateFriendship()
}

// updateConfig will be called at the start of each turn to set the newest config
func (c *client) updateConfig() {
	defer c.Logf("Configuration has been updated")

	ourResources := c.ServerReadHandle.GetGameState().ClientInfo.Resources
	minThreshold := c.ServerReadHandle.GetGameConfig().MinimumResourceThreshold
	costOfLiving := c.ServerReadHandle.GetGameConfig().CostOfLiving

	updatedConfig := ClientConfig{
		minFriendship:          c.clientConfig.minFriendship,
		maxFriendship:          c.clientConfig.maxFriendship,
		friendshipChangingRate: c.clientConfig.friendshipChangingRate,
		selfishThreshold:       minThreshold + 3.0*costOfLiving + ourResources/10.0,
		normalThreshold:        minThreshold + 6.0*costOfLiving + ourResources/10.0,
		multiplier:             c.clientConfig.multiplier,
	}

	c.clientConfig = ClientConfig(updatedConfig)
}

// updateFriendship will be called at the start of each turn to update our friendships
func (c *client) updateFriendship() {
	defer c.Logf("Friendship information has been updated")

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

func (c *client) DisasterNotification(dR disasters.DisasterReport, effects disasters.DisasterEffects) {
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

		// sets the cap of trust rank from 0 to 1
		if c.trustRank[team] < 0 {
			c.trustRank[team] = 0
		} else if c.turstRank[team] > 1 {
			c.trustRank[team] = 1
		}
	}

	c.disastersHistory = append(c.disastersHistory, disasterhappening)
}

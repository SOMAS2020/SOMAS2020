// Package team6 contains code for team 6's client implementation
package team6

import (
	"math/rand"
	"sort"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
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

func (c *client) StartOfTurn() {
	defer c.Logf("There are [%v] islands left in this game", c.getNumOfAliveIslands())

	c.updateConfig()
	c.updateFriendship()
}

// ########################
// ######  Foraging  ######
// ########################
func chooseForageType() int {
	return int(shared.DeerForageType)
}

func (c *client) DecideForage() (shared.ForageDecision, error) {
	ft := chooseForageType()
	return shared.ForageDecision{
		Type:         shared.ForageType(ft),
		Contribution: shared.Resources(rand.Float64() * 20),
	}, nil
}

func (c *client) ForageUpdate(shared.ForageDecision, shared.Resources) {}

/*
------ TODO: OPTIONAL ------
func (c *client) DisasterNotification(disasters.DisasterReport, map[shared.ClientID]shared.Magnitude)
*/

// ########################
// ######  Voting  ########
// ########################
// GetVoteForRule returns the client's vote in favour of or against a rule.
func (c *client) GetVoteForRule(ruleName string) bool {
	for _, val := range c.favourRules {
		if val == ruleName {
			return true
		}
	}
	return false
}

// GetVoteForElection returns the client's Borda vote for the role to be elected.
// COMPULSORY: use opinion formation to decide a rank for islands for the role
func (c *client) GetVoteForElection(roleToElect shared.Role) []shared.ClientID {
	// Done ;)
	// Get all alive islands
	aliveClients := rules.VariableMap[rules.IslandsAlive]
	// Convert to ClientID type and place into unordered map
	aliveClientIDs := map[int]shared.ClientID{}
	for i, v := range aliveClients.Values {
		aliveClientIDs[i] = shared.ClientID(int(v))
	}
	// Recombine map, in shuffled order
	var returnList []shared.ClientID
	for _, v := range aliveClientIDs {
		returnList = append(returnList, v)
	}
	return returnList
}

// ########################
// ######   Utils   ######
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

// Pair is used for sorting the friendship
type Pair struct {
	Key   shared.ClientID
	Value uint
}

// PairList is a slice of Pairs that implements sort.Interface to sort by Value.
type PairList []Pair

func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

// A function to turn a map into a PairList, then sort and return it.
func sortMapByValue(m map[shared.ClientID]uint) PairList {
	p := make(PairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = Pair{k, v}
	}
	sort.Sort(p)
	return p
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

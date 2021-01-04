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
	// MinFriendship is the minimum friendship level
	MinFriendship = 0.0
	// MaxFriendship is the maximum friendship level
	MaxFriendship = 100.0
)

type client struct {
	*baseclient.BaseClient

	friendship           map[shared.ClientID]float64
	giftsOfferedHistory  map[shared.ClientID]shared.Resources
	giftsReceivedHistory map[shared.ClientID]shared.Resources

	config        Config
	forageHistory ForageHistory
}

func init() {
	baseclient.RegisterClient(
		id,
		&client{
			BaseClient:           baseclient.NewClient(id),
			friendship:           friendship,
			giftsOfferedHistory:  giftsOfferedHistory,
			giftsReceivedHistory: giftsReceivedHistory,
			config:               config,
		},
	)
}

// ########################
// ######  General  #######
// ########################

//used for sorting the friendship
type Pair struct {
	Key   shared.ClientID
	Value uint
}

// A slice of Pairs that implements sort.Interface to sort by Value.
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

func (c *client) StartOfTurn() {}

// ########################
// ###### Friendship ######
// ########################

// increases the friendship level with some other islands
func (c *client) RaiseFriendshipLevel(clientID shared.ClientID) {
	currFriendship := c.friendship[clientID]

	if currFriendship == MaxFriendship {
		return
	}

	c.friendship[clientID]++
}

// decreases the friendship level with some other islands
func (c *client) LowerFriendshipLevel(clientID shared.ClientID) {
	currFriendship := c.friendship[clientID]

	if currFriendship == MinFriendship {
		return
	}

	c.friendship[clientID]--
}

// gets friendship cofficients
func (c client) getFriendshipCoeffs() map[shared.ClientID]float64 {
	friendshipCoeffs := make(map[shared.ClientID]float64)

	for team, fs := range c.friendship {
		friendshipCoeffs[team] = fs / MaxFriendship
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

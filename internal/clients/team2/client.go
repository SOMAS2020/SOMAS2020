// Package team2 contains code for team 2's client implementation
package team2

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const id = shared.Team2

type CommonPoolHistory map[uint]float64
type ResourcesLevelHistory map[uint]shared.Resources

// Type for Empathy level assigned to each other team
type EmpathyLevel int

const (
	Selfish EmpathyLevel = iota
	FairSharer
	Altruist
)

type IslandEmpathies map[shared.ClientID]EmpathyLevel

// Could be used to store our expectation on an island's behaviour (about anything)
// vs what they actually did
type ExpectationReality struct {
	exp  int
	real int
}

type Situation string

// us -> others: ReceivedRequests
// others -> us: GiftWeRequest

const (
	President     Situation = "President"
	Judge         Situation = "Judge"
	Speaker       Situation = "Speaker"
	Foraging      Situation = "Foraging"
	GiftWeRequest Situation = "Gifts"
)

type Opinion struct {
	Histories    map[Situation][]int
	Performances map[Situation]ExpectationReality
}

type ForageInfo struct {
	DecisionMade      shared.ForageDecision
	ResourcesObtained shared.Resources
}

type GiftInfo struct {
	requested shared.GiftRequest
	gifted    shared.GiftOffer
	reason    shared.AcceptReason
}

type GiftExchange struct {
	IslandRequest map[uint]GiftInfo
	OurRequest    map[uint]GiftInfo
}

type OpinionHist map[shared.ClientID]Opinion
type PredictionsHist map[shared.ClientID][]shared.DisasterPrediction
type ForagingReturnsHist map[shared.ClientID][]ForageInfo
type GiftHist map[shared.ClientID]GiftExchange

// we have to initialise our client somehow
type client struct {
	// should this have a * in front?
	*baseclient.BaseClient

	islandEmpathies      IslandEmpathies
	commonPoolHistory    CommonPoolHistory
	resourceLevelHistory ResourcesLevelHistory
	opinionHist          OpinionHist
	predictionHist       PredictionsHist
	foragingReturnsHist  ForagingReturnsHist
	giftHist             GiftHist
}

func init() {
	baseclient.RegisterClient(id, NewClient(id))
}

// TODO: are we using all of these or can they be removed
// TODO: we could experiment with how being more/less trustful affects agent performance
// i.e. start with assuming all islands selfish, normal, altruistic
func NewClient(clientID shared.ClientID) baseclient.Client {
	return &client{
		BaseClient:           baseclient.NewClient(clientID),
		commonPoolHistory:    CommonPoolHistory{},
		resourceLevelHistory: ResourcesLevelHistory{},
		opinionHist:          OpinionHist{},
		predictionHist:       PredictionsHist{},
		foragingReturnsHist:  ForagingReturnsHist{},
		giftHist:             GiftHist{},
		islandEmpathies:      IslandEmpathies{},

		//TODO: implement config to gather all changeable parameters in one place
	}
}

func (c *client) islandEmpathyLevel() EmpathyLevel {
	clientInfo := c.gameState().ClientInfo

	// switch statement to toggle between three levels
	// change our state based on these cases
	switch {
	case clientInfo.LifeStatus == shared.Critical:
		return Selfish
		// replace with some expression
	case (true):
		return Altruist
	default:
		return FairSharer
	}
}

func criticalStatus(c *client) bool {
	clientInfo := c.gameState().ClientInfo
	if clientInfo.LifeStatus == shared.Critical { //not sure about shared.Critical
		return true
	}
	return false
}

//TODO: how does this work?
// func (c *client) DisasterNotification(disasters.DisasterReport, map[shared.ClientID]shared.Magnitude)

//checkOthersCrit checks if anyone else is critical
func checkOthersCrit(c *client) bool {
	for clientID, status := range c.gameState().ClientLifeStatuses {
		if status == shared.Critical && clientID != c.GetID() {
			return true
		}
	}
	return false
}

func (c *client) gameState() gamestate.ClientGameState {
	return c.BaseClient.ServerReadHandle.GetGameState()
}

// Initialise initialises the base client.
// OPTIONAL: Overwrite, and make sure to keep the value of ServerReadHandle.
// You will need it to access the game state through its GetGameStateMethod.
func (c *client) Initialise(serverReadHandle baseclient.ServerReadHandle) {
	c.ServerReadHandle = serverReadHandle
	// loop through each island (there might not be 6)
	for clientID, _ := range c.gameState().ClientLifeStatuses {
		// set the confidence to 50 and initialise any other stuff
		Histories := make(map[Situation][]int)
		Histories["President"] = []int{50}
		Histories["Speaker"] = []int{50}
		Histories["Judge"] = []int{50}
		Histories["Foraging"] = []int{50}
		Histories["Gifts"] = []int{50}
		c.opinionHist[clientID] = Opinion{
			Histories:    Histories,
			Performances: map[Situation]ExpectationReality{},
		}
		c.giftHist[clientID] = GiftExchange{IslandRequest: map[uint]GiftInfo{},
			OurRequest: map[uint]GiftInfo{},
		}
	}
}

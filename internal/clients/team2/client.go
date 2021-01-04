// Package team2 contains code for team 2's client implementation
package team2

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const id = shared.Team2

func (c *client) StartOfTurn() {
	CommonPoolUpdate(c, c.commonPoolHistory) //add the common pool level
	ourResourcesHistoryUpdate(c, c.resourceLevelHistory)
}

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
	President        Situation = "President"
	Judge            Situation = "Judge"
	Speaker          Situation = "Speaker"
	Foraging         Situation = "Foraging"
	GiftWeRequest    Situation = "GiftWeRequest"
	ReceivedRequests Situation = "GiftGiven"
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

type DisasterOccurence struct {
	Turn   float64
	Report disasters.DisasterReport
}

type OpinionHist map[shared.ClientID]Opinion
type PredictionsHist map[shared.ClientID][]shared.DisasterPrediction
type ForagingReturnsHist map[shared.ClientID][]ForageInfo
type GiftHist map[shared.ClientID]GiftExchange
type DisasterHistory map[int]DisasterOccurence

func init() {
	baseclient.RegisterClient(
		id,
		&client{
			BaseClient: baseclient.NewClient(id),
			// add other things we want to remember (Histories)
			// commonpoolHistory: CommonPoolHistory{},
			// need to init to initially assume other islands are fair
			// IslandEmpathies: IslandEmpathies{},
		},
	)
}

// we have to initialise our client somehow
type client struct {
	// should this have a * in front?
	*baseclient.BaseClient

	islandEmpathies      IslandEmpathies
	commonPoolHistory    CommonPoolHistory
	resourceLevelHistory ResourcesLevelHistory
	opinionHist          OpinionHist
	predictionsHist      PredictionsHist
	foragingReturnsHist  ForagingReturnsHist
	giftHist             GiftHist
	disasterHistory      DisasterHistory
}

//NewClient After declaring the struct we have to actually make an object for the client
func NewClient(clientID shared.ClientID) baseclient.Client {
	// return a reference to the client struct variable's memory address
	return &client{
		BaseClient: baseclient.NewClient(clientID),
		// commonpoolHistory: CommonPoolHistory{},
		// we could experiment with how being more/less trustful affects agent performance
		// i.e. start with assuming all islands selfish, normal, altruistic
		islandEmpathies: IslandEmpathies{},
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
func (c *client) DisasterNotification(report disasters.DisasterReport, effects disasters.DisasterEffects) {
	c.disasterHistory[len(c.disasterHistory)] = DisasterOccurence{
		Turn:   float64(c.gameState().Turn),
		Report: report,
	}
}

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

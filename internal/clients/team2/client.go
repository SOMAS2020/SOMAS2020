// Package team2 contains code for team 2's client implementation
package team2

import (
	"github.com/SOMAS2020/SOMAS2020/internal/clients/team3/dynamics"
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
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

type DisasterOccurence struct {
	Turn   float64
	Report disasters.DisasterReport
}

// Currently what want to use to get archipelago geography but talking to Yannis to get this fixed
// Because it doesn't work atm
//archipelagoGeography := c.gamestate().Environment.Geography

// A set of constants that define tuning parameters
const (
	// Disasters
	TuningParamK             float64 = 1
	VarianceCapTimeRemaining float64 = 10000
	TuningParamG             float64 = 1
	VarianceCapMagnitude     float64 = 10000
)

type OpinionHist map[shared.ClientID]Opinion
type PredictionsHist map[shared.ClientID][]shared.DisasterPrediction
type ForagingReturnsHist map[shared.ClientID][]ForageInfo
type GiftHist map[shared.ClientID]GiftExchange
type DisasterHistory map[int]DisasterOccurence

// we have to initialise our client somehow
type client struct {
	// should this have a * in front?
	*baseclient.BaseClient

	ourSpeaker   speaker
	ourJudge     judge
	ourPresident president

	judgePerformance     map[shared.ClientID]int
	speakerPerformance   map[shared.ClientID]int
	presidentPerformance map[shared.ClientID]int

	// ## IIGO ##
	ruleVotedOn string

	// allVotes stores the votes of each island for/against each rule
	allVotes map[string]map[shared.ClientID]bool

	// iigoInfo caches information regarding iigo in the current turn
	iigoInfo iigoCommunicationInfo

	localVariableCache map[rules.VariableFieldName]rules.VariableValuePair

	localInputsCache map[rules.VariableFieldName]dynamics.Input
	// last sanction score cache to determine wheter or not we have been caugth in the last turn
	lastSanction roles.IIGOSanctionScore

	islandEmpathies      IslandEmpathies
	commonPoolHistory    CommonPoolHistory
	resourceLevelHistory ResourcesLevelHistory
	opinionHist          OpinionHist
	predictionHist       PredictionsHist
	foragingReturnsHist  ForagingReturnsHist
	giftHist             GiftHist
	disasterHistory      DisasterHistory
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

func (c *client) gameConfig() config.ClientConfig {
	return c.BaseClient.ServerReadHandle.GetGameConfig()
}

// Initialise initialises the base client.
// OPTIONAL: Overwrite, and make sure to keep the value of ServerReadHandle.
// You will need it to access the game state through its GetGameStateMethod.
func (c *client) Initialise(serverReadHandle baseclient.ServerReadHandle) {
	c.ServerReadHandle = serverReadHandle
	c.LocalVariableCache = rules.CopyVariableMap()
	// loop through each island (there might not be 6)
	for clientID := range c.gameState().ClientLifeStatuses {
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

type ruleVoteInfo struct {
	// ourVote needs to be updated accordingly
	ourVote         bool
	resultAnnounced bool
	// true -> yes, false -> no
	result bool
}

type iigoCommunicationInfo struct {
	// ourRole stores our current role in the IIGO
	ourRole *shared.Role
	// Retrieved fully from communications
	// commonPoolAllocation gives resources allocated by president from requests
	commonPoolAllocation shared.Resources
	// taxationAmount gives tax amount decided by president
	taxationAmount shared.Resources
	// monitoringOutcomes stores the outcome of the monitoring of an island.
	// key is the role being monitored.
	// true -> correct performance, false -> incorrect performance.
	monitoringOutcomes map[shared.Role]bool
	// monitoringDeclared stores as key the role being monitored and whether it was actually monitored.
	monitoringDeclared map[shared.Role]bool
	//Role IDs at the start of the turn
	startOfTurnPresidentID shared.ClientID
	startOfTurnJudgeID     shared.ClientID
	startOfTurnSpeakerID   shared.ClientID

	// Struct containing sanction information
	sanctions *sanctionInfo

	// Below need to be at least partially updated by our functions

	// ruleVotingResults is a map of rules and the corresponding info
	ruleVotingResults map[string]*ruleVoteInfo
	// ourRequest stores how much we requested from commonpool
	ourRequest shared.Resources
	// ourDeclaredResources stores how much we said we had to the president
	ourDeclaredResources shared.Resources
}

type sanctionInfo struct {
	// tierInfo provides tiers and sanction score required to get to that tier
	tierInfo map[roles.IIGOSanctionTier]roles.IIGOSanctionScore
	// rulePenalties provides sanction score given for breaking each rule
	rulePenalties map[string]roles.IIGOSanctionScore
	// islandSanctions stores sanction tier of each island (but not score)
	islandSanctions map[shared.ClientID]roles.IIGOSanctionTier
	// ourSanction is the sanction score for our island
	ourSanction roles.IIGOSanctionScore
}

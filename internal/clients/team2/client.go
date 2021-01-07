// Package team2 contains code for team 2's client implementation
package team2

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const id = shared.Team2

type CommonPoolHistory map[uint]float64
type ResourcesLevelHistory map[uint]shared.Resources

// IType for Empathy level assigned to each other team
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
	PresidentOp   Situation = "President"
	JudgeOp       Situation = "Judge"
	RoleOpinion   Situation = "RoleOpinion"
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

type IslandSanctionInfo struct {
	Turn   uint
	Tier   roles.IIGOSanctionTier
	Amount roles.IIGOSanctionScore
}

type CommonPoolInfo struct {
	tax             shared.Resources
	requestedToPres shared.Resources
	allocatedByPres shared.Resources
	takenFromCP     shared.Resources
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
type IslandSanctions map[shared.ClientID][]IslandSanctionInfo
type TierLevels map[roles.IIGOSanctionTier]roles.IIGOSanctionScore
type SanctionHist []IslandSanctionInfo
type CommonPoolHist map[uint]CommonPoolInfo

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
	disasterHistory      DisasterHistory

	currPresident President
	currJudge     Judge
	currSpeaker   Speaker

	taxAmount            shared.Resources
	commonPoolAllocation shared.Resources
	islandSanctions      IslandSanctions
	tierLevels           TierLevels
	sanctionHist         SanctionHist

	commonPoolHist CommonPoolHist

	//TODO: copied in from team3
	declaredResources   map[shared.ClientID]shared.Resources
	disasterPredictions []map[shared.ClientID]shared.DisasterPrediction
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
		sanctionHist:         SanctionHist{},
		tierLevels:           TierLevels{},
		islandSanctions:      IslandSanctions{},
		commonPoolHist:       CommonPoolHist{},

		//TODO: implement config to gather all changeable parameters in one place
	}
}

// Initialise initialises the base client.
// OPTIONAL: Overwrite, and make sure to keep the value of ServerReadHandle.
// You will need it to access the game state through its GetGameStateMethod.
func (c *client) Initialise(serverReadHandle baseclient.ServerReadHandle) {
	c.ServerReadHandle = serverReadHandle
	c.LocalVariableCache = rules.CopyVariableMap()
	// loop through each island (there might not be 6)
	for clientID, _ := range c.gameState().ClientLifeStatuses {
		// set the confidence to 50 and initialise any other stuff
		Histories := make(map[Situation][]int)
		Histories["President"] = []int{50}
		Histories["RoleOpinion"] = []int{50}
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
	c.currSpeaker = Speaker{c: c}
	c.currJudge = Judge{c: c}
	c.currPresident = President{c: c}
}

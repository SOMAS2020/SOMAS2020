// Package team2 contains code for team 2's client implementation
package team2

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const id = shared.Team2

type EmpathyLevel int

const (
	Selfish EmpathyLevel = iota
	FairSharer
	Altruist
)

type IslandEmpathies map[shared.ClientID]EmpathyLevel

type ExpectationReality struct {
	exp  int
	real int
}

type Situation string

const (
	PresidentOp Situation = "President"
	JudgeOp     Situation = "Judge"
	RoleOpinion Situation = "RoleOpinion"
	Foraging    Situation = "Foraging"
	Gifts       Situation = "Gifts"
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

type DisasterOccurrence struct {
	Turn   uint
	Report disasters.DisasterReport
}

type PredictionInfo struct {
	Prediction shared.DisasterPrediction
	Turn       uint
}

type IslandSanctionInfo struct {
	Turn   uint
	Tier   int
	Amount int
}

type CommonPoolInfo struct {
	turn            uint
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
	// Disasters (0, infinity]
	TuningParamK                     float64          = 1
	VarianceCapTimeRemaining         float64          = 10000
	TuningParamG                     float64          = 1
	VarianceCapMagnitude             float64          = 10000
	BaseResourcesToGiveDivisor       shared.Resources = 4
	BaseDisasterProtectionDivisor    shared.Resources = 4
	TimeLeftIncreaseDisProtection    float64          = 3
	DisasterSoonProtectionMultiplier float64          = 1.2
	DefaultFirstTurnContribution     shared.Resources = 20
	NoFreeRideAtStart                uint             = 3
	SwitchToFreeRideFactor           float64          = 0.5
	SwitchToAltruistFactor           float64          = 0.5
	FairShareFactorOfAvToGive        float64          = 1
	AltruistFactorOfAvToGive         float64          = 2
	ConfidenceRetrospectFactor       float64          = 0.5
	ForageDecisionThreshold          float64          = 0.6
	SlightRiskForageDivisor          shared.Resources = 2
	HelpCritOthersDivisor            shared.Resources = 2
	InitialDisasterTurnGuess         float64          = 7
)

type OpinionHist map[shared.ClientID]Opinion
type PredictionsHist map[shared.ClientID][]PredictionInfo
type ForagingReturnsHist map[shared.ClientID][]ForageInfo
type GiftHist map[shared.ClientID]GiftExchange
type CommonPoolHistory map[uint]shared.Resources
type ResourcesLevelHistory map[uint]shared.Resources
type DisasterHistory []DisasterOccurrence
type IslandSanctions map[shared.ClientID]IslandSanctionInfo
type TierLevels map[int]int
type SanctionHist map[shared.ClientID][]IslandSanctionInfo
type PresCommonPoolHist map[shared.ClientID][]CommonPoolInfo

type client struct {
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

	presCommonPoolHist PresCommonPoolHist

	declaredResources map[shared.ClientID]shared.Resources
}

func init() {
	baseclient.RegisterClientFactory(id, func() baseclient.Client { return NewClient(id) })
}

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
		presCommonPoolHist:   PresCommonPoolHist{},
		disasterHistory:      DisasterHistory{},

		//TODO: implement config to gather all changeable parameters in one place
	}
}

func (c *client) Initialise(serverReadHandle baseclient.ServerReadHandle) {
	c.ServerReadHandle = serverReadHandle
	c.LocalVariableCache = rules.CopyVariableMap(c.gameState().RulesInfo.VariableMap)

	for clientID := range c.gameState().ClientLifeStatuses {
		c.giftHist[clientID] = GiftExchange{
			IslandRequest: map[uint]GiftInfo{},
			OurRequest:    map[uint]GiftInfo{},
		}
	}

	c.currSpeaker = Speaker{c: c}
	c.currJudge = Judge{c: c}
	c.currPresident = President{c: c}
}

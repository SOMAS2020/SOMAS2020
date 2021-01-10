// Package team2 contains code for team 2's client implementation
package team2

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const id = shared.Team2

type AgentStrategy uint

// TODO: URGENT - THIS IS NOT BEING USED PROPERLY IN ANY OF THE SECTIONS
// TODO: THERE ARE FUNCTIONS THAT INCORRECTLY USE NUMBERS INSTEAD OF THE WORDS IN COMMON POOL AND OTHER PLACES
const (
	Selfish AgentStrategy = iota
	FairSharer
	Altruist
)

type IslandEmpathies map[shared.ClientID]AgentStrategy

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
	tax             shared.Resources
	requestedToPres shared.Resources
	allocatedByPres shared.Resources
	takenFromCP     shared.Resources
}

// Currently what want to use to get archipelago geography but talking to Yannis to get this fixed
// Because it doesn't work atm
//archipelagoGeography := c.gamestate().Environment.Geography

// A set of constants that define tuning parameters
type clientConfig struct {
	TuningParamK                     float64
	VarianceCapTimeRemaining         float64
	TuningParamG                     float64
	VarianceCapMagnitude             float64
	BaseResourcesToGiveDivisor       shared.Resources
	InitialThresholdProportionGuess  float64
	TimeLeftIncreaseDisProtection    uint
	DisasterSoonProtectionMultiplier float64
	DefaultFirstTurnContribution     shared.Resources
	SelfishStartTurns                uint
	SwitchToSelfishFactor            float64
	SwitchToAltruistFactor           float64
	FairShareFactorOfAvToGive        float64
	AltruistFactorOfAvToGive         float64
	ConfidenceRetrospectFactor       float64
	ForageDecisionThreshold          float64
	SlightRiskForageDivisor          shared.Resources
	HelpCritOthersDivisor            shared.Resources
	InitialDisasterTurnGuess         uint
	CombinedDisasterPred             shared.DisasterPrediction
	InitialCommonPoolThresholdGuess  shared.Resources
}

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
type PresCommonPoolHist map[shared.ClientID]map[uint]CommonPoolInfo

// DisasterVulnerabilityDict is a map from island ID to an islands DVP
type DisasterVulnerabilityDict map[shared.ClientID]float64

type client struct {
	*baseclient.BaseClient

	// TODO: naming convention on history objects is inconsistent
	islandEmpathies      IslandEmpathies
	commonPoolHistory    CommonPoolHistory
	resourceLevelHistory ResourcesLevelHistory
	opinionHist          OpinionHist
	predictionHist       PredictionsHist
	foragingReturnsHist  ForagingReturnsHist
	giftHist             GiftHist
	disasterHistory      DisasterHistory
	sanctionHist         SanctionHist
	presCommonPoolHist   PresCommonPoolHist

	currStrategy  AgentStrategy
	currPresident President
	currJudge     Judge
	currSpeaker   Speaker

	islandDVPs DisasterVulnerabilityDict

	// Define a global variable that holds the last prediction we shared
	AgentDisasterPred    shared.DisasterPredictionInfo
	CombinedDisasterPred shared.DisasterPrediction
	commonPoolThreshold  shared.Resources
	taxAmount            shared.Resources
	commonPoolAllocation shared.Resources
	islandSanctions      IslandSanctions
	tierLevels           TierLevels

	config clientConfig

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
		config: clientConfig{
			TuningParamK:                     1.0,
			VarianceCapTimeRemaining:         10000,
			TuningParamG:                     1.0,
			VarianceCapMagnitude:             10000,
			BaseResourcesToGiveDivisor:       4.0,
			InitialThresholdProportionGuess:  0.3,
			TimeLeftIncreaseDisProtection:    3.0,
			DisasterSoonProtectionMultiplier: 1.2,
			DefaultFirstTurnContribution:     20,
			SelfishStartTurns:                3,
			SwitchToSelfishFactor:            0.3,
			SwitchToAltruistFactor:           0.5,
			FairShareFactorOfAvToGive:        1.0,
			AltruistFactorOfAvToGive:         2.0,
			ConfidenceRetrospectFactor:       0.5,
			ForageDecisionThreshold:          0.6,
			SlightRiskForageDivisor:          2.0,
			HelpCritOthersDivisor:            2.0,
			InitialDisasterTurnGuess:         7.0,
			InitialCommonPoolThresholdGuess:  100, // this value is meaningless for now
		},
	}
}

func (c *client) Initialise(serverReadHandle baseclient.ServerReadHandle) {
	c.ServerReadHandle = serverReadHandle
	c.LocalVariableCache = rules.CopyVariableMap(c.gameState().RulesInfo.VariableMap)

	// We should probably initialise out empathy level
	for clientID := range c.gameState().ClientLifeStatuses {
		c.giftHist[clientID] = GiftExchange{
			IslandRequest: map[uint]GiftInfo{},
			OurRequest:    map[uint]GiftInfo{},
		}
	}

	// Set the initial strategy to selfish (can put anything here)
	c.currStrategy = Selfish

	// Initialise Disaster Prediction variables
	c.CombinedDisasterPred = shared.DisasterPrediction{}
	c.AgentDisasterPred = shared.DisasterPredictionInfo{}
	c.commonPoolThreshold = c.config.InitialCommonPoolThresholdGuess

	// Compute DVP for each Island based on Geography
	c.islandDVPs = DisasterVulnerabilityDict{}
	c.GetIslandDVPs(c.gameState().Geography)

	// Initialise Roles
	c.currSpeaker = Speaker{c: c}
	c.currJudge = Judge{c: c}
	c.currPresident = President{c: c}

}

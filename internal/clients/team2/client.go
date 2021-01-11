// Package team2 contains code for team 2's client implementation
package team2

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type AgentStrategy uint

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
	Disaster    Situation = "Disaster"
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
	DefaultContribution              shared.Resources
	SelfishStartTurns                uint
	SwitchToSelfishFactor            float64
	SwitchToAltruistFactor           float64
	FairShareFactorOfAvToGive        float64
	AltruistFactorOfAvToGive         float64
	ConfidenceRetrospectFactor       float64
	ForageDecisionThreshold          float64
	PatientTurns                     uint
	SlightRiskForageDivisor          shared.Resources
	HelpCritOthersDivisor            shared.Resources
	InitialDisasterTurnGuess         uint
	CombinedDisasterPred             shared.DisasterPrediction
	InitialCommonPoolThresholdGuess  shared.Resources
	TargetRequestGift                shared.Resources
	MaxGiftOffersMultiplier          shared.Resources
	AltruistMultiplier               shared.Resources
	FreeRiderMultipler               shared.Resources
	FairSharerMultipler              shared.Resources
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

type StrategyHistory map[uint]AgentStrategy

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
	strategyHistory      StrategyHistory

	currStrategy  AgentStrategy
	currPresident President
	currJudge     Judge
	currSpeaker   Speaker

	islandDVPs DisasterVulnerabilityDict

	// Define a global variable that holds the last prediction we shared
	AgentDisasterPred    shared.DisasterPredictionInfo
	CombinedDisasterPred shared.DisasterPrediction

	taxAmount            shared.Resources
	commonPoolAllocation shared.Resources
	islandSanctions      IslandSanctions
	tierLevels           TierLevels

	config clientConfig

	declaredResources map[shared.ClientID]shared.Resources

	lastForageType   shared.ForageType
	lastForageAmount shared.Resources
}

// DefaultClient creates the client that will be used for most simulations. All
// other personalities are considered alternatives. To give a different
// personality for your agent simply create another (exported) function with the
// same signature as "DefaultClient" that creates a different agent, and inform
// someone on the simulation team that you would like it to be included in
// testing
func DefaultClient(id shared.ClientID) baseclient.Client {
	return NewClient(id);
}

func NewClient(clientID shared.ClientID) baseclient.Client {
	return &client{
		BaseClient:           baseclient.NewClient(clientID),
		commonPoolHistory:    CommonPoolHistory{},
		resourceLevelHistory: ResourcesLevelHistory{},
		strategyHistory:      StrategyHistory{},
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
			PatientTurns:                     2,
			VarianceCapTimeRemaining:         10000,
			TuningParamG:                     1.0,
			VarianceCapMagnitude:             10000,
			BaseResourcesToGiveDivisor:       4.0,
			InitialThresholdProportionGuess:  0.3,
			TimeLeftIncreaseDisProtection:    3.0,
			DisasterSoonProtectionMultiplier: 1.2,
			DefaultContribution:              20,
			SelfishStartTurns:                3,
			SwitchToSelfishFactor:            0.3,
			SwitchToAltruistFactor:           0.9,
			FairShareFactorOfAvToGive:        1.0,
			AltruistFactorOfAvToGive:         2.0,
			ConfidenceRetrospectFactor:       0.5,
			ForageDecisionThreshold:          0.6,
			SlightRiskForageDivisor:          2.0,
			HelpCritOthersDivisor:            2.0,
			InitialDisasterTurnGuess:         7.0,
			InitialCommonPoolThresholdGuess:  100, // this value is meaningless for now
			TargetRequestGift:                1.5,
			MaxGiftOffersMultiplier:          0.5,
			AltruistMultiplier:               7.2,
			FreeRiderMultipler:               4.8,
			FairSharerMultipler:              6.4,
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
	c.strategyHistory = StrategyHistory{}

	// Compute DVP for each Island based on Geography
	c.islandDVPs = DisasterVulnerabilityDict{}
	c.getIslandDVPs(c.gameState().Geography)

	// Initialise Roles
	c.currSpeaker = Speaker{c: c}
	c.currJudge = Judge{c: c}
	c.currPresident = President{c: c}
}

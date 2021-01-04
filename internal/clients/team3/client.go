// Package team3 contains code for team 3's client implementation
package team3

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const id = shared.Team3

func init() {
	ourClient := &client{BaseClient: baseclient.NewClient(id)}

	baseclient.RegisterClient(id, ourClient)
}

type client struct {
	*baseclient.BaseClient
	// ## Gifting ##

	acceptedGifts        map[shared.ClientID]int
	requestedGiftAmounts map[shared.ClientID]int
	receivedResponses    []shared.GiftResponse

	// ## Trust ##

	theirTrustScore  map[shared.ClientID]float64
	trustScore       map[shared.ClientID]float64
	trustMapAgg      map[shared.ClientID][]float64
	theirTrustMapAgg map[shared.ClientID][]float64

	// ## Role performance ##

	judgePerformance     map[shared.ClientID]int
	speakerPerformance   map[shared.ClientID]int
	presidentPerformance map[shared.ClientID]int

	// ## IIGO ##
	ruleVotedOn    string
	ourVoteForRule bool

	// ## Game state & History ##
	criticalStatePrediction criticalStatePrediction

	// unused or replaced by getter functions
	// currentIteration iterationInfo
	// islandsAlive uint
	// localPool int

	// declaredResources is a map of all declared island resources
	declaredResources map[shared.ClientID]shared.Resources
	//disasterPredictions gives a list of predictions by island for each turn
	disasterPredictions []map[shared.ClientID]shared.DisasterPrediction

	// ## Compliance ##

	timeSinceCaught uint
	numTimeCaught   uint
	compliance      float64

	// allVotes stores the votes of each island for/against each rule
	allVotes map[string]map[shared.ClientID]bool

	// params is list of island wide function parameters
	params islandParams
	// iigoInfo caches information regarding iigo in the current turn
	iigoInfo iigoCommunicationInfo
}

type criticalStatePrediction struct {
	upperBound shared.Resources
	lowerBound shared.Resources
}

type islandParams struct {
	giftingThreshold            int
	equity                      float64
	complianceLevel             float64
	resourcesSkew               float64
	saveCriticalIsland          bool
	escapeCritcaIsland          bool
	selfishness                 float64
	minimumRequest              int
	disasterPredictionWeighting float64
	DesiredRuleSet              []string
	recidivism                  float64
	riskFactor                  float64
	friendliness                float64
	anger                       float64
	aggression                  float64
	sensitivity                 float64
}

type ruleVoteInfo struct {
	// ourVote needs to be updated accordingly
	ourVote         bool
	resultAnnounced bool
	// true -> yes, false -> no
	result bool
}

type iigoCommunicationInfo struct {
<<<<<<< HEAD
	// ourRole stores our current role in the IIGO
	ourRole *shared.Role
=======
	// Retrieved fully from communications

>>>>>>> origin/agent-strategy-biden
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
<<<<<<< HEAD
	//Role IDs at the start of the turn
	startOfTurnPresidentID shared.ClientID
	startOfTurnJudgeID     shared.ClientID
	startOfTurnSpeakerID   shared.ClientID

=======
>>>>>>> origin/agent-strategy-biden
	// Struct containing sanction information
	sanctions sanctionInfo

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

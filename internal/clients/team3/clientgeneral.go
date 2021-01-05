package team3

// General client functions

import (
	"math"
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func (c *client) DemoEvaluation() {
	evalResult, err := rules.BasicBooleanRuleEvaluator("Kinda Complicated Rule")
	if err != nil {
		panic(err.Error())
	}
	c.Logf("Rule Eval: %t", evalResult)
}

// NewClient initialises the island state
func NewClient(clientID shared.ClientID) baseclient.Client {
	ourClient := client{
		// Initialise variables here
		BaseClient: baseclient.NewClient(clientID),
		params: islandParams{
			// Define parameter values here
			selfishness: 0.5,
		},
	}

	// Set trust scores
	for _, islandID := range shared.TeamIDs {

		// Initialise trust scores for all islands except our own
		if islandID == ourClient.GetID() {
			continue
		}
		ourClient.trustScore[islandID] = 50
		ourClient.theirTrustScore[islandID] = 50
	}

	return &ourClient
}

func (c *client) StartOfTurn() {
	// c.Logf("Start of turn!")
	// TODO add any functions and vairable changes here

	// Initialise trustMap and theirtrustMap local cache to empty maps
	c.inittrustMapAgg()
	c.inittheirtrustMapAgg()
  
	c.updateCompliance()
	c.resetIIGOInfo()
}

func (c *client) Initialise(serverReadHandle baseclient.ServerReadHandle) {
	c.ServerReadHandle = serverReadHandle
	c.ourSpeaker = speaker{c: c}
	c.ourJudge = judge{c: c}
	c.ourPresident = president{c: c}
	// Initialise variables
}

// updatetrustMapAgg adds the amount to the aggregate trust map list for given client
func (c *client) updatetrustMapAgg(ClientID shared.ClientID, amount float64) {
	c.trustMapAgg[ClientID] = append(c.trustMapAgg[ClientID], amount)
}

// updatetheirtrustMapAgg adds the amount to the their aggregate trust map list for given client
func (c *client) updatetheirtrustMapAgg(ClientID shared.ClientID, amount float64) {
	c.theirTrustMapAgg[ClientID] = append(c.theirTrustMapAgg[ClientID], amount)
}

// inittrustMapAgg initialises the trustMapAgg to empty list values ready for each turn
func (c *client) inittrustMapAgg() {

	for _, islandID := range shared.TeamIDs {
		if islandID == c.BaseClient.GetID() {
			continue
		}
		c.trustMapAgg[islandID] = []float64{}
	}
}

// inittheirtrustMapAgg initialises the theirTrustMapAgg to empty list values ready for each turn
func (c *client) inittheirtrustMapAgg() {
	for _, islandID := range shared.TeamIDs {
		if islandID == c.BaseClient.GetID() {
			continue
		}
		c.theirTrustMapAgg[islandID] = []float64{}
	}
}

// updateTrustScore obtains average of all accumulated trust changes
// and updates the trustScore global map with new values
// ensuring that the values do not drop below 0 or exceed 100
func (c *client) updateTrustScore(trustMapAgg map[shared.ClientID][]float64) {
	for client, val := range trustMapAgg {
		avgScore := getAverage(val)
		if c.trustScore[client]+avgScore > 100.0 {
			avgScore = 100.0 - c.trustScore[client]
		}
		if c.trustScore[client]+avgScore < 0.0 {
			avgScore = 0.0 - c.trustScore[client]
		}
		c.trustScore[client] += avgScore
	}
}

// updateTheirTrustScore obtains average of all accumulated trust changes
// and updates the trustScore global map with new values
// ensuring that the values do not drop below 0 or exceed 100
func (c *client) updateTheirTrustScore(theirTrustMapAgg map[shared.ClientID][]float64) {
	for client, val := range theirTrustMapAgg {
		avgScore := getAverage(val)
		if c.theirTrustScore[client]+avgScore > 100.0 {
			avgScore = 100.0 - c.theirTrustScore[client]
		}
		if c.theirTrustScore[client]+avgScore < 0.0 {
			avgScore = 0.0 - c.theirTrustScore[client]
		}
		c.theirTrustScore[client] += avgScore
	}
}

//updateCriticalThreshold updates our predicted value of what is the resources threshold of critical state
// it uses estimated resources to find these bound. isIncriticalState is a boolean to indicate if the island
// is in the critical state and the estimated resources is our estimated resources of the island i.e.
// trust-adjusted resources.
func (c *client) updateCriticalThreshold(isInCriticalState bool, estimatedResource shared.Resources) {
	if !isInCriticalState {
		if estimatedResource < c.criticalStatePrediction.upperBound {
			c.criticalStatePrediction.upperBound = estimatedResource
		}
	} else {
		if estimatedResource > c.criticalStatePrediction.lowerBound {
			c.criticalStatePrediction.lowerBound = estimatedResource
		}
	}
}

// updateCompliance updates the compliance variable at the beginning of each turn.
// In the case that our island has been caught cheating in the previous turn, it is
// reset to 1 (aka. we fully comply and do not cheat)
func (c *client) updateCompliance() {
	if c.timeSinceCaught == 0 {
		c.compliance = 1
		c.numTimeCaught += 1
	} else {
		c.compliance = c.params.complianceLevel + (1.0-c.params.complianceLevel)*
			math.Exp(-float64(c.timeSinceCaught)/math.Pow((float64(c.numTimeCaught)+1.0), c.params.recidivism))
		c.timeSinceCaught += 1
	}
}

// shouldICheat returns whether or not our agent should cheat based
// the compliance at a specific time in the game. If the compliance is
// 1, we expect this method to always return False.
func (c *client) shouldICheat() bool {
	var should_i_cheat = rand.Float64() > c.compliance
	return should_i_cheat
}

/*
	DisasterNotification(disasters.DisasterReport, map[shared.ClientID]shared.Magnitude)
	updateCompliance
	shouldICheat
	updateCriticalThreshold
	evalPresidentPerformance
	evalSpeakerPerformance
	evalJudgePerformance
*/

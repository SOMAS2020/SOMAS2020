package team3

// General client functions

import (
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
		ourClient.trustScore[islandID] = 50
		ourClient.theirTrustScore[islandID] = 50
	}
	// Set our trust in ourselves to 100
	ourClient.theirTrustScore[id] = 100

	return &ourClient
}

func (c *client) StartOfTurn() {
	// c.Logf("Start of turn!")
	// TODO add any functions and vairable changes here
}

func (c *client) Initialise(serverReadHandle baseclient.ServerReadHandle) {
	c.ServerReadHandle = serverReadHandle
	// Initialise variables
}

// updatetrustMapAgg adds the amount to the aggregate trust map list for given client
func (c *client) updatetrustMapAgg(ClientID shared.ClientID, amount float64) {
	c.trustMapAgg[ClientID] = append(c.trustMapAgg[ClientID], amount)
}

// inittrustMapAgg initialises the trustMapAgg to empty list values ready for each turn
func (c *client) inittrustMapAgg() {
	c.trustMapAgg = map[shared.ClientID][]float64{
		0: []float64{},
		1: []float64{},
		3: []float64{},
		4: []float64{},
		5: []float64{},
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

/*
	ReceiveCommunication(sender shared.ClientID, data map[shared.CommunicationFieldName]shared.CommunicationContent)
	GetCommunications() *map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent
	DisasterNotification(disasters.DisasterReport, map[shared.ClientID]shared.Magnitude)

	updateCompliance
	shouldICheat
	getCompliance

	updateCriticalThreshold

	updateTrustScore
	updateTheirTrustScore
	decayTrust (may not be needed)

	evalPresidentPerformance
	evalSpeakerPerformance
	evalJudgePerformance
*/

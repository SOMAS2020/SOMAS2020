package team3

// General client functions

import (
	"math"
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// NewClient initialises the island state
func NewClient(clientID shared.ClientID) baseclient.Client {
	ourClient := client{
		// Initialise variables here
		BaseClient: baseclient.NewClient(clientID),
		params:     getislandParams(),
	}

	return &ourClient
}

func (c *client) StartOfTurn() {
	c.clientPrint("Start of turn!")

	// update Trust Scores at the start of every turn
	c.updateTrustScore(c.trustMapAgg)
	c.updateTheirTrustScore(c.theirTrustMapAgg)

	// Initialise trustMap and theirtrustMap local cache to empty maps
	c.inittrustMapAgg()
	c.inittheirtrustMapAgg()

	if c.checkIfCaught() {
		c.clientPrint("We've been caught")
		c.timeSinceCaught = 0
	}

	c.updateCompliance()
	c.lastSanction = c.iigoInfo.sanctions.ourSanction

	c.updateCriticalThreshold(c.ServerReadHandle.GetGameState().ClientInfo.LifeStatus, c.ServerReadHandle.GetGameState().ClientInfo.Resources)

	c.resetIIGOInfo()
	c.Logf("Our Status: %+v\n", c.ServerReadHandle.GetGameState().ClientInfo)
}

func (c *client) Initialise(serverReadHandle baseclient.ServerReadHandle) {
	c.ServerReadHandle = serverReadHandle
	c.LocalVariableCache = rules.CopyVariableMap()
	c.ourSpeaker = speaker{c: c}
	c.ourJudge = judge{c: c}
	c.ourPresident = president{c: c}

	c.initgiftOpinions()

	// Set trust scores
	c.trustScore = make(map[shared.ClientID]float64)
	c.theirTrustScore = make(map[shared.ClientID]float64)
	//c.localVariableCache = rules.CopyVariableMap()
	for _, islandID := range shared.TeamIDs {
		// Initialise trust scores for all islands except our own
		if islandID == c.BaseClient.GetID() {
			continue
		}
		c.trustScore[islandID] = 50
		c.theirTrustScore[islandID] = 50
	}

	// Set our trust in ourselves to 100
	c.theirTrustScore[id] = 100

	c.iigoInfo = iigoCommunicationInfo{
		sanctions: &sanctionInfo{
			tierInfo:        make(map[roles.IIGOSanctionTier]roles.IIGOSanctionScore),
			rulePenalties:   make(map[string]roles.IIGOSanctionScore),
			islandSanctions: make(map[shared.ClientID]roles.IIGOSanctionTier),
			ourSanction:     roles.IIGOSanctionScore(0),
		},
	}
	c.criticalStatePrediction.upperBound = serverReadHandle.GetGameState().ClientInfo.Resources
	c.criticalStatePrediction.lowerBound = serverReadHandle.GetGameState().ClientInfo.Resources
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
	c.trustMapAgg = map[shared.ClientID][]float64{}

	for _, islandID := range shared.TeamIDs {
		if islandID+1 == c.BaseClient.GetID() {
			continue
		}
		c.trustMapAgg[islandID] = []float64{}
	}
}

// inittheirtrustMapAgg initialises the theirTrustMapAgg to empty list values ready for each turn
func (c *client) inittheirtrustMapAgg() {
	c.theirTrustMapAgg = map[shared.ClientID][]float64{}

	for _, islandID := range shared.TeamIDs {
		if islandID+1 == c.BaseClient.GetID() {
			continue
		}
		c.theirTrustMapAgg[islandID] = []float64{}
	}
}

// inittheirtrustMapAgg initialises the theirTrustMapAgg to empty list values ready for each turn
func (c *client) initgiftOpinions() {
	c.giftOpinions = map[shared.ClientID]int{}

	for _, islandID := range shared.TeamIDs {
		if islandID+1 == c.BaseClient.GetID() {
			continue
		}
		c.giftOpinions[islandID] = 10
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
func (c *client) updateCriticalThreshold(state shared.ClientLifeStatus, estimatedResource shared.Resources) {
	isInCriticalState := state == shared.Critical
	if !isInCriticalState {
		if estimatedResource < c.criticalStatePrediction.upperBound {
			c.criticalStatePrediction.upperBound = estimatedResource
			if c.criticalStatePrediction.upperBound < c.criticalStatePrediction.lowerBound {
				c.criticalStatePrediction.lowerBound = estimatedResource
			}
		}
	} else {
		if estimatedResource > c.criticalStatePrediction.lowerBound {
			c.criticalStatePrediction.lowerBound = estimatedResource
			if c.criticalStatePrediction.upperBound < c.criticalStatePrediction.lowerBound {
				c.criticalStatePrediction.upperBound = estimatedResource
			}
		}
	}
}

// updateCompliance updates the compliance variable at the beginning of each turn.
// In the case that our island has been caught cheating in the previous turn, it is
// reset to 1 (aka. we fully comply and do not cheat)
func (c *client) updateCompliance() {
	if c.timeSinceCaught == 0 {
		c.compliance = 1
		c.numTimeCaught++
	} else {
		c.compliance = c.params.complianceLevel + (1.0-c.params.complianceLevel)*
			math.Exp(-float64(c.timeSinceCaught)/math.Pow((float64(c.numTimeCaught)+1.0), c.params.recidivism))
		c.timeSinceCaught++
	}
}

// shouldICheat returns whether or not our agent should cheat based
// the compliance at a specific time in the game. If the compliance is
// 1, we expect this method to always return False.
func (c *client) shouldICheat() bool {
	return rand.Float64() > c.compliance
}

// checkIfCaught, checks if the island has been caught during the last turn
// If it has been caught, it returns True, otherwise False.
func (c *client) checkIfCaught() bool {
	return c.iigoInfo.sanctions.ourSanction > c.lastSanction
}

// ResourceReport overides the basic method to mis-report when we have a low compliance score
func (c *client) ResourceReport() shared.ResourcesReport {
	resource := c.BaseClient.ServerReadHandle.GetGameState().ClientInfo.Resources
	if c.areWeCritical() || !c.shouldICheat() {
		return shared.ResourcesReport{ReportedAmount: resource, Reported: true}
	}
	skewedResource := resource / shared.Resources(c.params.resourcesSkew)
	return shared.ResourcesReport{ReportedAmount: skewedResource, Reported: true}
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

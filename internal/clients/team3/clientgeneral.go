package team3

// General client functions

import (
	"math"
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
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
	c.LocalVariableCache = rules.CopyVariableMap(c.ServerReadHandle.GetGameState().RulesInfo.VariableMap)
	c.ourSpeaker = speaker{c: c, BaseSpeaker: &baseclient.BaseSpeaker{GameState: c.ServerReadHandle.GetGameState()}}
	c.ourJudge = judge{c: c, BaseJudge: &baseclient.BaseJudge{GameState: c.ServerReadHandle.GetGameState()}}
	c.ourPresident = president{c: c, BasePresident: &baseclient.BasePresident{GameState: c.ServerReadHandle.GetGameState()}}

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
			tierInfo:        make(map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore),
			rulePenalties:   make(map[string]shared.IIGOSanctionsScore),
			islandSanctions: make(map[shared.ClientID]shared.IIGOSanctionsTier),
			ourSanction:     shared.IIGOSanctionsScore(0),
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

// Internal function that evaluates the performance of the judge for the purposes of opinion formation.
// This is called AFTER IIGO FINISHES.
func (c *client) evalJudgePerformance() {
	previousJudgeID := c.iigoInfo.startOfTurnJudgeID
	previousPresidentID := c.iigoInfo.startOfTurnPresidentID
	evalOfJudge := float64(c.judgePerformance[previousJudgeID])

	// If the judge didn't evaluate the speaker, the judge didn't do a good job
	if c.iigoInfo.monitoringDeclared[shared.Speaker] == false {
		evalOfJudge -= c.trustScore[previousJudgeID] * c.params.sensitivity
	}

	// Use the president's evaluation of the judge to determine how well the judge performed
	var presidentEvalOfJudge bool
	if c.iigoInfo.monitoringDeclared[shared.Judge] == true {
		presidentEvalOfJudge = c.iigoInfo.monitoringOutcomes[shared.Judge]
	}
	if presidentEvalOfJudge == true {
		evalOfJudge += c.trustScore[previousPresidentID] * c.params.sensitivity
	} else {
		evalOfJudge -= c.trustScore[previousPresidentID] * c.params.sensitivity
	}

	// Did the judge support our vote for president?
	ourVoteForPresident := c.VoteForElection(shared.President, c.getAliveIslands())
	electedPresident := c.ServerReadHandle.GetGameState().PresidentID
	var ourRankingChosen int
	for index, islandID := range ourVoteForPresident {
		if islandID == electedPresident {
			ourRankingChosen = index
		}
	}
	// If our third choice was voted in (ourRankingChosen == 2), no effect on Judge Performance.
	// Anything better/worse than third is rewarded/penalized proportionally.
	evalOfJudge += c.params.sensitivity * float64((2 - ourRankingChosen))

	// Did the judge sanction us?
	sanctionAmount := c.iigoInfo.sanctions.ourSanction
	evalOfJudge -= float64(sanctionAmount) * c.params.sensitivity

}

// Internal function that evaluates the performance of the president for the purposes of opinion formation.
// This is called AFTER IIGO FINISHES.
func (c *client) evalPresidentPerformance() {
	previousSpeakerID := c.iigoInfo.startOfTurnSpeakerID
	previousPresidentID := c.iigoInfo.startOfTurnPresidentID
	evalOfPresident := float64(c.presidentPerformance[previousPresidentID])

	// If the president didn't evaluate the judge, the president didn't do a good job
	if c.iigoInfo.monitoringDeclared[shared.Judge] == false {
		evalOfPresident -= c.trustScore[previousPresidentID] * c.params.sensitivity
	}

	// Use the speaker's evaluation of the president to determine how well the president performed
	var speakerEvalofPresident bool
	if c.iigoInfo.monitoringDeclared[shared.President] == true {
		speakerEvalofPresident = c.iigoInfo.monitoringOutcomes[shared.President]
	}
	if speakerEvalofPresident == true {
		evalOfPresident += c.trustScore[previousSpeakerID] * c.params.sensitivity
	} else {
		evalOfPresident -= c.trustScore[previousSpeakerID] * c.params.sensitivity
	}

	evalOfPresident += float64(c.iigoInfo.commonPoolAllocation-c.CommonPoolResourceRequest()) * c.params.sensitivity

	// if c.ourPresident.PickRuleToVote() == c.ruleVotedOn {
	// 	evalOfPresident += c.params.sensitivity
	// } else {
	// 	evalOfPresident -= c.params.sensitivity
	// }

	// evalOfPresident += (SetTaxationAmount() - c.iigoInfo.taxationAmount) * c.params.sensitivity

	// Did the president support our vote for speaker?
	ourVoteForSpeaker := c.VoteForElection(shared.Speaker, c.getAliveIslands())
	electedSpeaker := c.ServerReadHandle.GetGameState().SpeakerID
	var ourRankingChosen int
	for index, islandID := range ourVoteForSpeaker {
		if islandID == electedSpeaker {
			ourRankingChosen = index
		}
	}
	// If our third choice was voted in (ourRankingChosen == 2), no effect on President Performance.
	// Anything better/worse than third is rewarded/penalized proportionally.
	evalOfPresident += c.params.sensitivity * float64((2 - ourRankingChosen))
}

// Internal function that evaluates the performance of the speaker for the purposes of opinion formation.
// This is called AFTER IIGO FINISHES.
func (c *client) evalSpeakerPerformance() {
	previousJudgeID := c.iigoInfo.startOfTurnJudgeID
	previousSpeakerID := c.iigoInfo.startOfTurnJudgeID
	evalOfSpeaker := float64(c.speakerPerformance[previousSpeakerID])

	// If the speaker didn't evaluate the president, the speaker didn't do a good job
	if c.iigoInfo.monitoringDeclared[shared.President] == false {
		evalOfSpeaker -= c.trustScore[previousSpeakerID] * c.params.sensitivity
	}

	// Use the judge's evaluation of the speaker to determine how well the speaker performed
	var judgeEvalofSpeaker bool

	if c.iigoInfo.monitoringDeclared[shared.Speaker] == true {
		judgeEvalofSpeaker = c.iigoInfo.monitoringOutcomes[shared.Speaker]
	}

	if judgeEvalofSpeaker == true {
		evalOfSpeaker += c.trustScore[previousJudgeID] * c.params.sensitivity
	} else {
		evalOfSpeaker -= c.trustScore[previousJudgeID] * c.params.sensitivity
	}

	ruleVoteInfo := *c.iigoInfo.ruleVotingResults[c.ruleVotedOn]
	if ruleVoteInfo.ourVote != ruleVoteInfo.result {
		evalOfSpeaker += c.params.sensitivity
	} else {
		evalOfSpeaker -= c.params.sensitivity
	}

	// Did the speaker support our vote for judge?
	ourVoteForJudge := c.VoteForElection(shared.Judge, c.getAliveIslands())
	electedJudge := c.ServerReadHandle.GetGameState().JudgeID
	var ourRankingChosen int
	for index, islandID := range ourVoteForJudge {
		if islandID == electedJudge {
			ourRankingChosen = index
		}
	}
	// If our third choice was voted in (ourRankingChosen == 2), no effect on President Performance.
	// Anything better/worse than third is rewarded/penalized proportionally.
	evalOfSpeaker += c.params.sensitivity * float64((2 - ourRankingChosen))
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

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
	c.resetIIGOInfo()
}

func (c *client) Initialise(serverReadHandle baseclient.ServerReadHandle) {
	c.ServerReadHandle = serverReadHandle
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
	c.trustMapAgg = map[shared.ClientID][]float64{
		0: []float64{},
		1: []float64{},
		3: []float64{},
		4: []float64{},
		5: []float64{},
	}
}

// inittheirtrustMapAgg initialises the theirTrustMapAgg to empty list values ready for each turn
func (c *client) inittheirtrustMapAgg() {
	c.theirTrustMapAgg = map[shared.ClientID][]float64{
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
	ourVoteForPresident := c.GetVoteForElection(shared.President)
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
	sanctionAmount = c.iigoInfo.sanctions.ourSanction
	evalOfJudge -= float64(sanctionAmount) * c.params.sensitivity

}

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

	evalOfPresident += (c.iigoInfo.commonPoolAllocation - c.CommonPoolResourceRequest()) * c.params.sensitivity

	if PickRuleToVote() == c.ruleVotedOn {
		evalOfPresident += c.params.sensitivity
	} else {
		evalOfPresident -= c.params.sensitivity
	}

	evalOfPresident += (SetTaxationAmount() - c.iigoInfo.taxationAmount) * c.params.sensitivity

	// Did the president support our vote for speaker?
	ourVoteForSpeaker := c.GetVoteForElection(shared.Speaker)
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

	if c.ourVoteForRule != c.iigoInfo.ruleVotingResults[c.ruleVotedOn] {
		evalOfSpeaker += c.params.sensitivity
	} else {
		evalOfSpeaker -= c.params.sensitivity
	}

	// Did the speaker support our vote for judge?
	ourVoteForJudge := c.GetVoteForElection(shared.Judge)
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

/*
	DisasterNotification(disasters.DisasterReport, map[shared.ClientID]shared.Magnitude)

	updateCompliance
	shouldICheat
	getCompliance

	updateCriticalThreshold

	evalPresidentPerformance
	evalSpeakerPerformance
	evalJudgePerformance
*/

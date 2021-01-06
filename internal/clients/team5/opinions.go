package team5

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type opinionScore float64

type trustBasis int

// basis on which we trust another team. May may only care about their forecasting reputation, for example.
const (
	generalBasis trustBasis = iota
	forecastingBasis
)

// an opinion of another island is our general perception of it and is characterised by a score
// between 1.0 (highest/best perception) and -1.0 (worst/lowest perception). Potential events that may
// affect our opinion of other teams:
// - allocating/distributing less than they initially offer for gifting
// - hunting deer when the population is below critical population
// - corrupt IIGO roles
// - others?
type opinion struct {
	score              opinionScore // broad opinion score
	forecastReputation opinionScore
}

// opinions of each team. Need opinion as a pointer so we can modify it
type opinionMap map[shared.ClientID]*opinion

// history of opinionMaps (opinions per team) across turns
type opinionHistory map[uint]opinionMap // key is turn, value is opinion

// creates initial opinions of clients and creates
func (c *client) initOpinions() {
	c.opinionHistory = opinionHistory{}
	c.opinions = opinionMap{}
	for _, team := range c.getAliveTeams(true) { // true to include our team if alive
		c.opinions[team] = &opinion{score: 0} // start with neutral opinion score
	}
	c.opinionHistory[startTurn] = c.opinions // 0th turn is how we start before the game starts - our initial bias
	c.Logf("Opinions at first turn (turn %v): %v", startTurn, c.opinionHistory)
}

func (o opinion) String() string {
	return fmt.Sprintf("opinion{score: %.2f}", o.score)
}

// "Normalize" opinion to keep it in range -1 and 1
func (c *client) normalizeOpinion() {
	for clientID := range c.opinions {
		if c.opinions[clientID].score < -1.0 {
			c.opinions[clientID].score = -1
		} else if c.opinions[clientID].score > 1.0 {
			c.opinions[clientID].score = 1
		}
	}
}

//Evaluate if the roles are corrupted or not based on their budget spending versus total tax paid to common pool
//Either everyone is corrupted or not
func (c *client) evaluateRoles() {
	speakerID := c.ServerReadHandle.GetGameState().SpeakerID
	judgeID := c.ServerReadHandle.GetGameState().JudgeID
	presidentID := c.ServerReadHandle.GetGameState().PresidentID
	//compute total budget
	budget := c.ServerReadHandle.GetGameState().IIGORolesBudget
	var totalBudget shared.Resources = 0
	for role := range budget {
		totalBudget += budget[role]
	}
	// compute total maximum tax to cp
	var totalTax shared.Resources
	numberAliveTeams := len(c.getAliveTeams(true)) //include us
	for i := 0; i < numberAliveTeams; i++ {
		totalTax += c.taxAmount
	}
	// Not corrupt
	if totalBudget <= totalTax {
		c.opinions[speakerID].score += 0.1 //arbitrary number
		c.opinions[judgeID].score += 0.1
		c.opinions[presidentID].score += 0.1
	} else {
		c.opinions[speakerID].score -= 0.1 //arbitrary number
		c.opinions[judgeID].score -= 0.1
		c.opinions[presidentID].score -= 0.1
	}
	c.normalizeOpinion()
}

// getTrustedTeams finds teams whose opinion scores (our opinion of them) exceed a threshold `trustThresh`. Furthermore,
// if `proportional` is true, the scores of the trusted teams will be relative (such that sum of scores = 1). If not,
// the absolute opinion scores for each client are returned. // TODO: decide if our team should be included here or not
func (c client) getTrustedTeams(trustThresh opinionScore, proportional bool, basis trustBasis) (trustedTeams map[shared.ClientID]float64) {
	totalTrustedOpScore := 0.0
	trustedTeams = map[shared.ClientID]float64{}
	for team, opinion := range c.opinions {
		opValue := opinionScore(0.0)

		switch basis { // get opinion value based on the trust basis specified
		case generalBasis:
			opValue = opinion.score
		case forecastingBasis:
			opValue = opinion.forecastReputation
		}

		if opValue >= trustThresh && c.isClientAlive(team) { // trust team and they're alive
			trustedTeams[team] = float64(opinion.score)
			totalTrustedOpScore += float64(opinion.score) // store this in case proportional score scaling is req.
		}
	}
	if !proportional {
		return trustedTeams // return here if we only want absolute op. scores
	}
	for team, score := range trustedTeams {
		trustedTeams[team] = score / totalTrustedOpScore // scale proportionally to other trusted teams for weighting
	}
	return trustedTeams // weighted proportional scores
}

func (c *client) giftOpinions() {
	for team := range c.gameState().ClientLifeStatuses { // for each ID
		// ======================= Bad =======================
		// If we get OFFERED LESS than we Requested
		if shared.Resources(c.giftHistory[team].ourRequest[c.getTurn()].offered) <
			shared.Resources(c.giftHistory[team].ourRequest[c.getTurn()].requested) {
			c.opinions[team].score -= 0.05
		}

		// If we ACTUALLY get LESS than they OFFERED us
		if shared.Resources(c.giftHistory[team].ourRequest[c.getTurn()].actualRecieved) <
			shared.Resources(c.giftHistory[team].ourRequest[c.getTurn()].offered) {
			c.opinions[team].score -= 0.10
		}

		// If they REQUEST the MOST compared to other islands
		highestRequest := shared.Team1
		if c.giftHistory[highestRequest].theirRequest[c.getTurn()].requested <
			c.giftHistory[team].theirRequest[c.getTurn()].requested {
			highestRequest = team
		}
		c.opinions[highestRequest].score -= 0.05

		// ======================= Good =======================
		// If they GIVE MORE than OFFERED then increase it a bit (can be abused)
		if shared.Resources(c.giftHistory[team].ourRequest[c.getTurn()].actualRecieved) >
			shared.Resources(c.giftHistory[team].ourRequest[c.getTurn()].offered) {
			c.opinions[team].score += 0.025
		}

		// If we RECIEVE MORE than WE REQUESTED and they OFFERED
		if shared.Resources(c.giftHistory[team].ourRequest[c.getTurn()].actualRecieved) >
			shared.Resources(c.giftHistory[team].ourRequest[c.getTurn()].offered) &&
			shared.Resources(c.giftHistory[team].ourRequest[c.getTurn()].actualRecieved) >
				shared.Resources(c.giftHistory[team].ourRequest[c.getTurn()].requested) {
			c.opinions[team].score += 0.2
		}

		// If they REQUEST the LEAST compared to other islands
		lowestRequest := shared.Team1
		if c.giftHistory[lowestRequest].theirRequest[c.getTurn()].requested >
			c.giftHistory[team].theirRequest[c.getTurn()].requested {
			lowestRequest = team
		}
		c.opinions[lowestRequest].score += 0.05
	}
}

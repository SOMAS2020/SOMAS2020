package team5

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type opinionScore float64

type opinionBasis int

// basis on which we trust another team. May may only care about their forecasting reputation, for example.
const (
	generalBasis opinionBasis = iota
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

type wrappedOpininon struct {
	opinion opinion
}

// opinions of each team. Need opinion as a pointer so we can modify it
type opinionMap map[shared.ClientID]*wrappedOpininon

// history of opinionMaps (opinions per team) across turns
type opinionHistory map[uint]opinionMap // key is turn, value is opinion

// String implements Stringer
func (wo wrappedOpininon) String() string {
	return fmt.Sprintf("wrappedOpinion{opinion: %v}", wo.opinion)
}

func (wo wrappedOpininon) getScore() opinionScore {
	return wo.opinion.score
}

func (wo wrappedOpininon) getForecastingRep() opinionScore {
	return wo.opinion.forecastReputation
}

func (wo *wrappedOpininon) updateOpinion(basis opinionBasis, increment float64) {
	op := wo.opinion
	switch basis {
	case generalBasis:
		newScore := float64(op.score) + increment
		op.score = opinionScore(absoluteCap(newScore, 1.0))
	case forecastingBasis:
		newScore := float64(op.forecastReputation) + increment
		op.forecastReputation = opinionScore(absoluteCap(newScore, 1.0))
	}
	wo.opinion = op // update opinion
}

// creates initial opinions of clients and creates
func (c *client) initOpinions() {
	c.opinionHistory = opinionHistory{}
	c.opinions = opinionMap{}
	for _, team := range c.getAliveTeams(true) { // true to include our team if alive
		c.opinions[team] = &wrappedOpininon{opinion: opinion{score: 0, forecastReputation: 0}} // start with neutral opinion score
	}
	c.opinionHistory[startTurn] = c.opinions // 0th turn is how we start before the game starts - our initial bias
	// c.Logf("Opinions at first turn (turn %v): %v", startTurn, c.opinionHistory)
}

func (o opinion) String() string {
	return fmt.Sprintf("opinion{score: %.2f}", o.score)
}

// getTrustedTeams finds teams whose opinion scores (our opinion of them) exceed a threshold `trustThresh`. Furthermore,
// if `proportional` is true, the scores of the trusted teams will be relative (such that sum of scores = 1). If not,
// the absolute opinion scores for each client are returned. // TODO: decide if our team should be included here or not
func (c client) getTrustedTeams(trustThresh opinionScore, proportional bool, basis opinionBasis) (trustedTeams map[shared.ClientID]float64) {
	totalTrustedOpScore := 0.0
	trustedTeams = map[shared.ClientID]float64{}
	for team, opinion := range c.opinions {
		opValue := opinionScore(0.0)

		switch basis { // get opinion value based on the trust basis specified
		case generalBasis:
			opValue = opinion.getScore()
		case forecastingBasis:
			opValue = opinion.getForecastingRep()
		}

		if opValue >= trustThresh && c.isClientAlive(team) { // trust team and they're alive
			trustedTeams[team] = float64(opValue)
			totalTrustedOpScore += float64(opValue) // store this in case proportional score scaling is req.
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

package team5

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type opinionScore float64

// an opinion of another island is our general perception of it and is characterised by a score
// between 1.0 (highest/best perception) and -1.0 (worst/lowest perception). Potential events that may
// affect our opinion of other teams:
// - allocating/distributing less than they initially offer for gifting
// - hunting deer when the population is below critical population
// - corrupt IIGO roles
// - others?
type opinion struct {
	score opinionScore // may want to define other atrributes that form an opinion
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

// getTrustedTeams finds teams whose opinion scores (our opinion of them) exceed a threshold `trustThresh`. Furthermore,
// if `proportional` is true, the scores of the trusted teams will be relative (such that sum of scores = 1). If not,
// the absolute opinion scores for each client are returned.
func (c client) getTrustedTeams(trustThresh opinionScore, proportional bool) (trustedTeams map[shared.ClientID]float64) {
	totalTrustedOpScore := 0.0
	for team, opinion := range c.opinions {
		if opinion.score >= trustThresh && c.isClientAlive(team) { // trust team and they're alive
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

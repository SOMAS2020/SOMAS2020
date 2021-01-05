package team5

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

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

type opinionMap map[shared.ClientID]opinion // opinions of each team

// history of opinionMaps (opinions per team) across turns
type opinionHistory map[uint]opinionMap // key is turn, value is opinion

// creates initial opinions of clients and creates
func (c *client) initOpinions() {
	c.opinionHistory = opinionHistory{}
	c.opinions = opinionMap{}
	for _, team := range c.getAliveTeams(true) { // true to include our team if alive
		c.opinions[team] = opinion{score: 0} // start with neutral opinion score
	}
	c.opinionHistory[startTurn] = c.opinions // 0th turn is how we start before the game starts - our initial bias
}

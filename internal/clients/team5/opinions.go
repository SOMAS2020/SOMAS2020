package team5

import (
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

type opinionMap map[shared.ClientID]*opinion // opinions of each team

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
}

func (c *client) giftOpinions() {
	// ======================= Bad =======================

	// If we get OFFERED LESS than we Requested
	for team := range c.gameState().ClientLifeStatuses { // for each ID
		if shared.Resources(c.giftHistory[team].OurRequest[c.getTurn()].offered) <
			shared.Resources(c.giftHistory[team].OurRequest[c.getTurn()].requested) {
			c.opinions[team].score -= 0.05
		}

		// If we ACTUALLY get LESS than they OFFERED us
		if shared.Resources(c.giftHistory[team].OurRequest[c.getTurn()].actualRecieved) <
			shared.Resources(c.giftHistory[team].OurRequest[c.getTurn()].offered) {
			c.opinions[team].score -= 0.10
		}

		// If they REQUEST the MOST compared to other islands
		highestRequest := shared.Team1
		if c.giftHistory[highestRequest].TheirRequest[c.getTurn()].requested <
			c.giftHistory[team].TheirRequest[c.getTurn()].requested {
			highestRequest = team
		}
		c.opinions[highestRequest].score -= 0.05

		// ======================= Good =======================
		// If they GIVE MORE than OFFERED then increase it a bit
		if shared.Resources(c.giftHistory[team].OurRequest[c.getTurn()].actualRecieved) >
			shared.Resources(c.giftHistory[team].OurRequest[c.getTurn()].offered) {
			c.opinions[team].score += 0.025
		}

		// If they OFFER MORE than we REQUESTED + we RECIEVE MORE than they offered
		if shared.Resources(c.giftHistory[team].OurRequest[c.getTurn()].actualRecieved) >
			shared.Resources(c.giftHistory[team].OurRequest[c.getTurn()].offered) &&
			shared.Resources(c.giftHistory[team].OurRequest[c.getTurn()].actualRecieved) >
				shared.Resources(c.giftHistory[team].OurRequest[c.getTurn()].requested) {
			c.opinions[team].score += 0.10
		}

		// If they REQUEST the LEAST compared to other islands
		lowestRequest := shared.Team1
		if c.giftHistory[lowestRequest].TheirRequest[c.getTurn()].requested >
			c.giftHistory[team].TheirRequest[c.getTurn()].requested {
			lowestRequest = team
		}
		c.opinions[lowestRequest].score += 0.05
	}
}

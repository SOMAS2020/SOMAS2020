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

// "Normalize" opinion to keep it in range -1 and 1
func (c *client) normalizeOpinion() {
	for clientID, _ := range c.opinions {
		if c.opinions[clientID].score < -1.0 {
			c.opinions[clientID].score = -1
		} else if c.opinions[clientID].score > 1.0 {
			c.opinions[clientID].score = 1
		}
	}
}

//Evaluate if the roles are corrupted or not based on their budget spending
//	versus total tax paid to common pool
//Either everyone is corrupted or not
func (c *client) evaluateRoles() {
	speakerID := c.ServerReadHandle.GetGameState().SpeakerID
	judgeID := c.ServerReadHandle.GetGameState().JudgeID
	presidentID := c.ServerReadHandle.GetGameState().PresidentID
	//compute total budget
	budget := c.ServerReadHandle.GetGameState().IIGORolesBudget
	var totalBudget shared.Resources = 0
	for role, _ := range budget {
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

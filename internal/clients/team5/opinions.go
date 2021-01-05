package team5

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

// an opinion of another island is our general perception of it and is characterised by a score
// between 1.0 (highest/best perception) and -1.0 (worst/lowest perception). Potential events that may
// affect our opinion of other teams:
// - allocating/distributing less than they initially offer for gifting
// - hunting deer when the population is below critical population
// - corrupt IIGO roles
// - others?
type opinion map[shared.ClientID]struct {
	score float64
}

// history of opinions for each team across turns
type opinionHistory map[uint]opinion // key is turn, value is opinion

func (c client) initOpinionHistory() {
}

package team5

import (
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func TestGetVoteForElection(t *testing.T) {
	c.opinions = opinionMap{
		shared.Team1: &wrappedOpininon{opinion{score: 0.2}},
		shared.Team2: &wrappedOpininon{opinion{score: 0.3}},
		shared.Team3: &wrappedOpininon{opinion{score: 0.4}},
		shared.Team4: &wrappedOpininon{opinion{score: 0.5}},
		shared.Team5: &wrappedOpininon{opinion{score: 0.6}},
		shared.Team6: &wrappedOpininon{opinion{score: 0.7}},
	}
	candidateList := []shared.ClientID{shared.Team1, shared.Team2, shared.Team3,
		shared.Team4, shared.Team5, shared.Team6}
	ballot := c.VoteForElection(shared.President, candidateList)

	var w []shared.ClientID
	w = append(w, shared.Team5, shared.Team6, shared.Team4,
		shared.Team3, shared.Team2, shared.Team1)
	for i := 0; i < 6; i++ {
		if w[i] != ballot[i] {
			t.Errorf("Ballot not generating properly. Want %v, got %v", w, ballot)
		}
	}
}

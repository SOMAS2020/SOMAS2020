package team5

import (
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func TestGetVoteForElection(t *testing.T) {
	c.opinions = opinionMap{
		shared.Teams["Team1"]: &wrappedOpininon{opinion{score: 0.2}},
		shared.Teams["Team2"]: &wrappedOpininon{opinion{score: 0.3}},
		shared.Teams["Team3"]: &wrappedOpininon{opinion{score: 0.4}},
		shared.Teams["Team4"]: &wrappedOpininon{opinion{score: 0.5}},
		shared.Teams["Team5"]: &wrappedOpininon{opinion{score: 0.6}},
		shared.Teams["Team6"]: &wrappedOpininon{opinion{score: 0.7}},
	}
	candidateList := []shared.ClientID{shared.Teams["Team1"], shared.Teams["Team2"], shared.Teams["Team3"],
		shared.Teams["Team4"], shared.Teams["Team5"], shared.Teams["Team6"]}
	ballot := c.VoteForElection(shared.President, candidateList)

	var w []shared.ClientID
	w = append(w, shared.Teams["Team5"], shared.Teams["Team6"], shared.Teams["Team4"],
		shared.Teams["Team3"], shared.Teams["Team2"], shared.Teams["Team1"])
	for i := 0; i < 6; i++ {
		if w[i] != ballot[i] {
			t.Errorf("Ballot not generating properly. Want %v, got %v", w, ballot)
		}
	}
}

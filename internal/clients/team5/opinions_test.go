package team5

import (
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

var opinions = opinionMap{
	shared.Team1: &wrappedOpininon{opinion: opinion{score: 0, forecastReputation: 0}},
	shared.Team2: &wrappedOpininon{opinion: opinion{score: 0, forecastReputation: 0}},
}
var opHistory = opinionHistory{}

func TestUpdateOpinion(t *testing.T) {
	nTurns := 5

	op1 := opinions[shared.Team1]
	op2 := opinions[shared.Team2]

	for i := 1; i <= nTurns; i++ {
		// update history
		opHistory[uint(i)] = opinions
		op1.updateOpinion(generalBasis, 0.1)
		op2.updateOpinion(forecastingBasis, -0.1)

		t.Logf("Log history: %v", opHistory)

		w := 0.1 * float64(i)
		if float64(op1.getScore()) != w {
			t.Errorf("Score not updating properly. Want %v, got %v", w, op1.getScore())
		}

		w = 0.1 * float64(i) * -1
		if float64(op2.getForecastingRep()) != w {
			t.Errorf("Score not updating properly. Want %v, got %v", w, op2.getScore())
		}
	}
}

func TestOpinionCapping(t *testing.T) {
	op1 := opinions[shared.Team1]

	op1.updateOpinion(generalBasis, -0.9)
	op1.updateOpinion(generalBasis, -0.9)

	op1.updateOpinion(forecastingBasis, 0.9)
	op1.updateOpinion(forecastingBasis, 0.9)

	if op1.getScore() != -1 {
		t.Errorf("Got %v, want %v", op1.getScore(), -1)
	}

	if op1.getForecastingRep() != 1 {
		t.Errorf("Got %v, want %v", op1.getScore(), 1)
	}

}

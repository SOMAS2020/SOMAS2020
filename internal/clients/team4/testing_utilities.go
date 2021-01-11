package team4

import (
	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"gonum.org/v1/gonum/mat"
)

type fakeServerHandle struct {
	PresidentID        shared.ClientID
	JudgeID            shared.ClientID
	SpeakerID          shared.ClientID
	TermLengths        map[shared.Role]uint
	clients            []shared.ClientID
	ElectionRuleInPlay bool
}

func (s fakeServerHandle) GetGameState() gamestate.ClientGameState {
	lifeStatuses := map[shared.ClientID]shared.ClientLifeStatus{}
	for _, clientID := range s.clients {
		lifeStatuses[clientID] = shared.Alive
	}

	fakeGS := gamestate.ClientGameState{
		SpeakerID:          s.SpeakerID,
		JudgeID:            s.JudgeID,
		PresidentID:        s.PresidentID,
		ClientLifeStatuses: lifeStatuses,
	}

	if s.ElectionRuleInPlay {
		fakeGS.RulesInfo.CurrentRulesInPlay = registerTestElectionRule()
	}

	return fakeGS
}

func (s fakeServerHandle) GetGameConfig() config.ClientConfig {
	return config.ClientConfig{
		IIGOClientConfig: config.IIGOConfig{
			IIGOTermLengths: s.TermLengths,
		},
	}
}

func registerTestElectionRule() map[string]rules.RuleMatrix {
	rulesStore := map[string]rules.RuleMatrix{}

	name := "roles_must_hold_election"
	reqVar := []rules.VariableFieldName{
		rules.TermEnded,
		rules.ElectionHeld,
	}
	v := []float64{1, -1, 0}
	CoreMatrix := mat.NewDense(1, 3, v)

	aux := []float64{0}
	AuxiliaryVector := mat.NewVecDense(1, aux)

	rm := rules.RuleMatrix{RuleName: name, RequiredVariables: reqVar, ApplicableMatrix: *CoreMatrix, AuxiliaryVector: *AuxiliaryVector, Mutable: false}
	rulesStore[name] = rm
	return rulesStore
}

func floatEqual(x float64, y float64) bool {
	tolerance := 0.0001
	return math.Abs(x-y) < tolerance // if x - y < tolerance then x == y
}

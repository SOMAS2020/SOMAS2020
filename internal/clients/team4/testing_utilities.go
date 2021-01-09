package team4

import (
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
	ElectionRuleInPlay bool
}

func (s fakeServerHandle) GetGameState() gamestate.ClientGameState {
	fakeGS := gamestate.ClientGameState{
		SpeakerID:   s.SpeakerID,
		JudgeID:     s.JudgeID,
		PresidentID: s.PresidentID,
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

package iigointernal

import (
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func TestPutSalaryBack(t *testing.T) {
	fakeGameState := gamestate.GameState{
		CommonPool: 149,
		IIGORolesBudget: map[string]shared.Resources{
			"president": 0,
			"speaker":   0,
			"judge":     0,
		},
		ClientInfos: map[shared.ClientID]gamestate.ClientInfo{
			shared.Team1: {Resources: 0, LifeStatus: shared.Alive},
			shared.Team2: {Resources: 0, LifeStatus: shared.Alive},
			shared.Team3: {Resources: 0, LifeStatus: shared.Alive},
		},
		SpeakerID:   shared.Team1,
		JudgeID:     shared.Team2,
		PresidentID: shared.Team3,
	}
	fakeClientMap := map[shared.ClientID]baseclient.Client{
		shared.Team1: baseclient.NewClient(shared.Team1),
		shared.Team2: baseclient.NewClient(shared.Team2),
		shared.Team3: baseclient.NewClient(shared.Team3),
	}
	goodRun, _ := RunIIGO(&fakeGameState, &fakeClientMap)

	if goodRun {
		t.Errorf("IIGO didn't throw error when salaries couldn't be paid")
	} else {
		if fakeGameState.CommonPool != 149 {
			t.Errorf("Common pool contained '%v' expected '%v'", fakeGameState.CommonPool, 149)
		}
	}
}

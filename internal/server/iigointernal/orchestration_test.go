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
	}
	goodRun, _ := RunIIGO(&fakeGameState, &map[shared.ClientID]baseclient.Client{
		shared.Team1: &baseclient.BaseClient{},
		shared.Team2: &baseclient.BaseClient{},
		shared.Team3: &baseclient.BaseClient{},
	})
	if goodRun {
		t.Errorf("IIGO didn't throw error when salaries couldn't be paid")
	} else {
		if fakeGameState.CommonPool != 149 {
			t.Errorf("Common pool contained '%v' expected '%v'", fakeGameState.CommonPool, 149)
		}
	}
}

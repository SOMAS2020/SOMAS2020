package iigointernal

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"testing"
)

func TestPutSalaryBack(t *testing.T) {
	fakeGameState := gamestate.GameState{
		CommonPool: 149,
	}
	err := RunIIGO(&fakeGameState, &map[shared.ClientID]baseclient.Client{
		shared.Team1: &baseclient.BaseClient{},
		shared.Team2: &baseclient.BaseClient{},
		shared.Team3: &baseclient.BaseClient{},
	})
	if err == nil {
		t.Errorf("IIGO didn't throw error when salaries couldn't be paid")
	} else {
		if fakeGameState.CommonPool != 149 {
			t.Errorf("Common pool contained '%v' expected '%v'", fakeGameState.CommonPool, 149)
		}
	}
}



// we need more tests:
// below are the test proposed:
//1. test salary paid are moved into island private resources pools
//2. test that the first round will works well
//3. add experimentatl

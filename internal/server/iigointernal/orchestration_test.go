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
	err := RunIIGO(&fakeGameState, &map[shared.ClientID]baseclient.Client{})
	if err == nil {
		t.Errorf("IIGO didn't throw error when salaries couldn't be paid")
	} else {
		if fakeGameState.CommonPool != 149 {
			t.Errorf("Common pool contained '%v' expected '%v'", fakeGameState.CommonPool, 149)
		}
	}
}

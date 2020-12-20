package roles

import (
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common"
)

func TestWithdrawFromCommonPoolThrowsError(t *testing.T) {
	fakeGameState := common.GameState{CommonPool: 100}
	// Withdraw more than we have in it
	valueToWithdraw := 120
	err := WithdrawFromCommonPool(valueToWithdraw, &fakeGameState)
	if err == nil {
		t.Errorf("We can withdraw more from the common pool than it actually has.")
	}
}

func TestWithdrawFromCommonPoolDeductsValue(t *testing.T) {
	fakeGameState := common.GameState{CommonPool: 100}
	valueToWithdraw := 60
	_ = WithdrawFromCommonPool(valueToWithdraw, &fakeGameState)
	unexpectedAmountRemaining := fakeGameState.CommonPool != 40
	if unexpectedAmountRemaining == true {
		t.Errorf("Not withdrawing resources from CommonPool correctly.")
	}
}

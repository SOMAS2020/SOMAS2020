package roles

import (
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common"
)

// TestBasicRuleEvaluatorPositive Checks whether rule we expect to evaluate as true actually evaluates as such
func TestBasicRuleEvaluatorPositive(t *testing.T) {
	result, err := BasicBooleanRuleEvaluator("Kinda Complicated Rule")
	if !result {
		t.Errorf("Rule evaluation came as false, when it was expected to be true, potential error with value '%v'", err)
	}
}

func TestEnoughInCommonPoolNegative(t *testing.T) {
fakeGameState:
	common.GameState{CommonPool: 100}
	// Withdraw more than we have in it
	valueToWithdraw := 120
	status = CheckEnoughInCommonPool(valueToWithdraw, fakeGameState)
	if status {
		t.Errorf("We were able to withdraw 120 from the common pool, when there was only 100 in it. Error in internal.server.roles.utilities.CheckEnoughInCommonPool")
	}
}

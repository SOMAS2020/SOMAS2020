package roles

import (
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
)

//base President Object
type basePresident struct {
	id                 int
	budget             int
	speakerSalary      int
	resourceRequests   map[int]int
	resourceAllocation map[int]int
	ruleToVote         int
	taxAmount          int
}

func (p *basePresident) withdrawSpeakerSalary(gameState *common.GameState) error {
	var speakerSalary = int(rules.VariableMap["speakerSalary"].Values[0])
	var withdrawError = WithdrawFromCommonPool(speakerSalary, gameState)
	if withdrawError != nil {
		Base_President.speakerSalary = speakerSalary
	}
	return withdrawError
}

// Pay the speaker
func (p *basePresident) paySpeaker(gameState *common.GameState) {
	Base_speaker.budget = Base_President.speakerSalary
	Base_President.speakerSalary = 0
}

func (p *basePresident) signalAllocationRequests(int) map[int]int {
	return nil
}
func (p *basePresident) replyAllocationRequests(int) error {
	return nil
}
func (p *basePresident) sendRuleToSpeaker(int) error {
	return nil
}
func (p *basePresident) appointNextSpeaker() int {
	return rand.Intn(5)
}

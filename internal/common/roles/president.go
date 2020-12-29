package roles

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

//President Object
type President interface {
	PaySpeaker(salary shared.Resources) shared.PresidentReturnContent
	SetTaxationAmount(map[shared.ClientID]shared.ResourcesReport) shared.PresidentReturnContent
	EvaluateAllocationRequests(map[shared.ClientID]shared.Resources, shared.Resources) shared.PresidentReturnContent
	PickRuleToVote([]string) shared.PresidentReturnContent
}

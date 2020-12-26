package roles

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

//President Object
type President interface {
	PaySpeaker(salary shared.Resources) (shared.Resources, error)
	SetTaxationAmount(map[shared.ClientID]shared.Resources) (map[shared.ClientID]shared.Resources, error)
	EvaluateAllocationRequests(map[shared.ClientID]shared.Resources, shared.Resources) (map[shared.ClientID]shared.Resources, error)
	PickRuleToVote([]string) (string, bool)
}

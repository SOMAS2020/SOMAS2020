package roles

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

//President Object
type President interface {
	PaySpeaker() (shared.Resources, bool)
	SetTaxationAmount(map[shared.ClientID]shared.Resources) (map[shared.ClientID]shared.Resources, bool)
	EvaluateAllocationRequests(map[shared.ClientID]shared.Resources, shared.Resources) (map[shared.ClientID]shared.Resources, bool)
	PickRuleToVote([]string) (string, bool)
}

package roles

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

// President is an interface that is implemented by basePresident but can also be
// optionally implemented by individual islands.
type President interface {
	PaySpeaker(salary shared.Resources) (shared.Resources, bool)
	SetTaxationAmount(map[shared.ClientID]shared.Resources) (map[shared.ClientID]shared.Resources, bool)
	EvaluateAllocationRequests(map[shared.ClientID]shared.Resources, shared.Resources) (map[shared.ClientID]shared.Resources, bool)
	PickRuleToVote([]string) (string, bool)
	CallSpeakerElection(int, []shared.ClientID) shared.ElectionSettings
	DecideNextSpeaker(shared.ClientID) shared.ClientID
}

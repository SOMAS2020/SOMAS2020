package roles

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

// President is an interface that is implemented by basePresident but can also be
// optionally implemented by individual islands.
type President interface {
	PaySpeaker() shared.PresidentReturnContent
	GetSpeakerSalary() shared.Resources
	SetTaxationAmount(map[shared.ClientID]shared.ResourcesReport) shared.PresidentReturnContent
	EvaluateAllocationRequests(map[shared.ClientID]shared.Resources, shared.Resources) shared.PresidentReturnContent
	PickRuleToVote([]string) shared.PresidentReturnContent
	CallSpeakerElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) shared.ElectionSettings
	DecideNextSpeaker(shared.ClientID) shared.ClientID
}

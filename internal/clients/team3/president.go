package team3

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type president struct {
	*baseclient.BasePresident
}

func (p *president) PaySpeaker(salary shared.Resources) (shared.Resources, bool) {
	return shared.Resources(0), false
}
func (p *president) SetTaxationAmount(map[shared.ClientID]shared.Resources) (map[shared.ClientID]shared.Resources, bool) {
	return nil, false
}
func (p *president) EvaluateAllocationRequests(map[shared.ClientID]shared.Resources, shared.Resources) (map[shared.ClientID]shared.Resources, bool) {
	return nil, false
}
func (p *president) PickRuleToVote([]string) (string, bool) {
	return "", false
}
func (p *president) CallSpeakerElection(int, []shared.ClientID) shared.ElectionSettings {
	return shared.ElectionSettings{}
}
func (p *president) DecideNextSpeaker(shared.ClientID) shared.ClientID {
	return shared.ClientID(0)
}

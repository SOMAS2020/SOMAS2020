package team3

import (
	// "github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	// "github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	// "github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

/*
	//IIGO: COMPULSORY
	MonitorIIGORole(shared.Role) bool
	DecideIIGOMonitoringAnnouncement(bool) (bool, bool)

	GetVoteForRule(ruleName string) bool
	GetVoteForElection(roleToElect shared.Role) []shared.ClientID

	CommonPoolResourceRequest() shared.Resources
	ResourceReport() shared.Resources
	RuleProposal() string
	GetClientPresidentPointer() roles.President
	GetClientJudgePointer() roles.Judge
	GetClientSpeakerPointer() roles.Speaker
	TaxTaken(shared.Resources)
	GetTaxContribution() shared.Resources
	RequestAllocation() shared.Resources
*/

func (c *client) GetClientSpeakerPointer() roles.Speaker {
	// c.Logf("became speaker")
	return &speaker{c: c}
}

func (c *client) GetClientJudgePointer() roles.Judge {
	// c.Logf("became judge")
	return &judge{c: c}
}

func (c *client) GetClientPresidentPointer() roles.President {
	// c.Logf("became president")
	return &president{c: c}
}

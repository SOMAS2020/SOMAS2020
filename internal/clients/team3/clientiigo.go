package team3

import (
	"github.com/SOMAS2020/SOMAS2020/internal/clients/team3/dynamics"
	// "github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
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

func (c *client) RuleProposal() string {
	c.locationService.syncGameState(c.ServerReadHandle.GetGameState())
	c.locationService.syncTrustScore(c.trustScore)
	// Magically will be available
	coolMap := make(map[string]rules.RuleMatrix)
	coolmap2 := make(map[rules.VariableFieldName]dynamics.Input)

	// Will fix properly later
	shortestSoFar := 999999.0
	longestSoFar := 0.0
	selectedRule := ""
	for key, rule := range rules.AvailableRules {
		if _, ok := rules.RulesInPlay[key]; !ok {
			idealLoc, valid := c.locationService.checkIfIdealLocationAvailable(rule)
			if valid {
				ruleDynamics := dynamics.BuildAllDynamics(rule, rule.AuxiliaryVector)
				distance := dynamics.GetDistanceToSubspace(ruleDynamics, idealLoc)
				if distance != -1 {
					if shortestSoFar > distance {
						if _, ok := rules.RulesInPlay[rule.RuleName]; !ok {
							shortestSoFar = distance
							selectedRule = rule.RuleName
						}
					}
				}
			}
		} else {
			lstRules := dynamics.RemoveFromMap(coolMap, key)
			dist := dynamics.CalculateDistanceFromRuleSpace(lstRules, coolmap2)
			if dist > longestSoFar {
				selectedRule = rule.RuleName
			}
		}
	}
	if selectedRule == "" {
		return "inspect_ballot_rule"
	}
	return selectedRule
}

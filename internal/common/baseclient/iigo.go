package baseclient

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// CommonPoolResourceRequest is called by the President in IIGO to
// request an allocation of resources from the common pool.
func (c *BaseClient) CommonPoolResourceRequest() shared.Resources {
	// TODO: Implement needs based resource request.
	return 20
}

// ResourceReport is an island's self-report of its own resources. This is called by
// the President to help work out how many resources to allocate each island.
// OPTIONAL : as is, this function will always report island resources accurately
func (c *BaseClient) ResourceReport() shared.ResourcesReport {
	return shared.ResourcesReport{
		ReportedAmount: c.ServerReadHandle.GetGameState().ClientInfo.Resources,
		Reported:       true,
	}
}

// RuleProposal is for an island to propose a rule to be voted on. It is called
// by the President in IIGO. If the returned ruleMatrix is one of the rules
// in AvailableRules cache with unchanged content, then the proposal is for
// putting the rule in/out of play. However, if the returned ruleMatrix is
// one of the rules in AvailableRules cache with changed content, the proposal
// is then for modifying the rule's content only, it won't put the rule in/out of
// play. Only a mutable rule's content can be modified.
func (c *BaseClient) RuleProposal() rules.RuleMatrix {
	allRules := rules.AvailableRules
	for _, ruleMatrix := range allRules {
		return ruleMatrix
	}
	return rules.RuleMatrix{}
}

// GetClientPresidentPointer is called by IIGO to get the client's implementation of the President Role
// COMPULSORY: ovverride to return a pointer to your own President object
func (c *BaseClient) GetClientPresidentPointer() roles.President {
	return &BasePresident{}
}

// GetClientJudgePointer is called by IIGO to get the client's implementation of the Judge Role
// COMPULSORY: ovverride to return a pointer to your own Judge object
func (c *BaseClient) GetClientJudgePointer() roles.Judge {
	return &BaseJudge{}
}

// GetClientSpeakerPointer is called by IIGO to get the client's implementation of the Speaker Role
// COMPULSORY: ovverride to return a pointer to your own Speaker object
func (c *BaseClient) GetClientSpeakerPointer() roles.Speaker {
	return &BaseSpeaker{}
}

// GetTaxContribution gives value of how much the island wants to pay in taxes
// The tax is the minimum contribution, you can pay more if you want to
// COMPULSORY
func (c *BaseClient) GetTaxContribution() shared.Resources {
	// TODO: Implement common pool contribution greater than or equal to tax.
	return 0
}

// GetSanctionPayment is called at the end of turn to pay any sanctions that have been
// imposed on the island, it is up to the island if they choose to pay the sanction or not.
// COMPULSORY
func (c *BaseClient) GetSanctionPayment() shared.Resources {
	return 0
}

// RequestAllocation is called at the end of the turn to request resources from the
// common pool. If there are enough resources in the common pool, the server will
// update your island's pool with the resources you requested.
// COMPULSORY
func (c *BaseClient) RequestAllocation() shared.Resources {
	// TODO: Implement request equal to the allocation permitted by President.
	return 0
}

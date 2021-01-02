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

// ResourceReport is an island's self-report of its own resources.
func (c *BaseClient) ResourceReport() shared.Resources {
	return c.ServerReadHandle.GetGameState().ClientInfo.Resources
}

// RuleProposal is called by the President in IIGO to propose a
// rule to be voted on.
func (c *BaseClient) RuleProposal() string {
	allRules := rules.AvailableRules
	for k := range allRules {
		return k
	}
	return ""
}

// GetClientPresidentPointer is called by IIGO to get the client's implementation of the President Role
func (c *BaseClient) GetClientPresidentPointer() roles.President {
	return &BasePresident{}
}

// GetClientJudgePointer is called by IIGO to get the client's implementation of the Judge Role
func (c *BaseClient) GetClientJudgePointer() roles.Judge {
	return &BaseJudge{}
}

// GetClientSpeakerPointer is called by IIGO to get the client's implementation of the Speaker Role
func (c *BaseClient) GetClientSpeakerPointer() roles.Speaker {
	return &BaseSpeaker{}
}

func (c *BaseClient) TaxTaken(shared.Resources) {
	// Just an update. Ignore
}

// GetTaxContribution FIXME: Add documentation. What does this function do?
func (c *BaseClient) GetTaxContribution() shared.Resources {
	// TODO: Implement common pool contribution greater than or equal to tax.
	return 0
}

// RequestAllocation FIXME: Add documentation. What does this function do?
func (c *BaseClient) RequestAllocation() shared.Resources {
	// TODO: Implement request equal to the allocation permitted by President.
	return 0
}

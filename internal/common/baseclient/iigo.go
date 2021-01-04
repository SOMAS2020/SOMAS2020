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
func (c *BaseClient) ResourceReport() shared.ResourcesReport {
	return shared.ResourcesReport{
		ReportedAmount: c.ServerReadHandle.GetGameState().ClientInfo.Resources,
		Reported:       true,
	}
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

// GetTaxContribution gives value of how much the island wants to pay in taxes
func (c *BaseClient) GetTaxContribution() shared.Resources {
	// If no new communication was received and there is no rule governing taxation, it will still recieve last valid tax
	clientGameState := c.ServerReadHandle.GetGameState()
	presidentCommunications := c.Communications[clientGameState.PresidentID]
	newTaxDecision := shared.TaxDecision{TaxDecided: false}

	for i := range presidentCommunications {
		msg := presidentCommunications[len(presidentCommunications)-1-i] // in reverse order
		if msg[shared.Tax].T == shared.CommunicationTax {
			newTaxDecision = msg[shared.Tax].TaxDecision
			break
		}
	}

	if newTaxDecision.TaxDecided {
		// gotVariable := newTaxDecision.ExpectedTax
		// gotRule := newTaxDecision.TaxRule
		return newTaxDecision.TaxAmount
	}

	return 0 //default if no tax message was found
}

// GetSanctionPayment gives the value of how much the island is paying in sanctions
func (c *BaseClient) GetSanctionPayment() shared.Resources {
	return 0
}

// RequestAllocation FIXME: Add documentation. What does this function do?
func (c *BaseClient) RequestAllocation() shared.Resources {
	// TODO: Implement request equal to the allocation permitted by President.
	clientGameState := c.ServerReadHandle.GetGameState()
	presidentCommunications := c.Communications[clientGameState.PresidentID]
	newAllocationDecision := shared.AllocationDecision{AllocationDecided: false}
	for i := range presidentCommunications {
		msg := presidentCommunications[len(presidentCommunications)-1-i] // in reverse order
		if msg[shared.Allocation].T == shared.CommunicationAllocation {
			newAllocationDecision = msg[shared.Allocation].AllocationDecision
			break
		}
	}

	if newAllocationDecision.AllocationDecided {
		// gotVariable := newAllocationDecision.ExpectedAllocation
		// gotRule := newAllocationDecision.AllocationRule
		return newAllocationDecision.AllocationAmount
	}

	return clientGameState.CommonPool //default if no allocation was made
}

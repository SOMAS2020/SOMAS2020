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

func (c *BaseClient) TaxTaken(shared.Resources) {
	// Just an update. Ignore
}

// GetTaxContribution gives value of how much the island wants to pay in taxes
// The tax is the minimum contribution, you can pay more if you want to
// COMPULSORY
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
		clientContribution := float64(newTaxDecision.TaxAmount)

		expectedVariable := newTaxDecision.ExpectedTax
		contribuitionVariable := rules.VariableValuePair{
			VariableName: rules.IslandTaxContribution,
			Values:       []float64{clientContribution},
		}
		variableMap := map[rules.VariableFieldName]rules.VariableValuePair{}
		variableMap[expectedVariable.VariableName] = expectedVariable
		variableMap[contribuitionVariable.VariableName] = contribuitionVariable

		gotRule := newTaxDecision.TaxRule

		eval, err := rules.BasicLocalBooleanRuleEvaluator(gotRule, variableMap)
		if err == nil && eval {
			return newTaxDecision.TaxAmount
		}
	}

	return 0 //default if no tax message was found
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
		clientAllocation := float64(newAllocationDecision.AllocationAmount)

		expectedVariable := newAllocationDecision.ExpectedAllocation
		contribuitionVariable := rules.VariableValuePair{
			VariableName: rules.IslandAllocation,
			Values:       []float64{clientAllocation},
		}
		variableMap := map[rules.VariableFieldName]rules.VariableValuePair{}
		variableMap[expectedVariable.VariableName] = expectedVariable
		variableMap[contribuitionVariable.VariableName] = contribuitionVariable

		gotRule := newAllocationDecision.AllocationRule

		eval, err := rules.BasicLocalBooleanRuleEvaluator(gotRule, variableMap)
		if err == nil && eval {
			return newAllocationDecision.AllocationAmount
		}
	}

	return clientGameState.CommonPool //default if no allocation was made
}

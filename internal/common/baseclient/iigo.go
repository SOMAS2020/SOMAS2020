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
	amountReported := c.ServerReadHandle.GetGameState().ClientInfo.Resources
	c.LocalVariableCache[rules.IslandReportedResources] = rules.MakeVariableValuePair(rules.IslandReportedResources, []float64{float64(amountReported)})
	return shared.ResourcesReport{
		ReportedAmount: amountReported,
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

// GetTaxContribution gives value of how much the island wants to pay in taxes
// The tax is the minimum contribution, you can pay more if you want to
// COMPULSORY
func (c *BaseClient) GetTaxContribution() shared.Resources {
	// TODO: Implement common pool contribution greater than or equal to tax.
	valToBeReturned := shared.Resources(0)
	c.LocalVariableCache[rules.IslandTaxContribution] = rules.VariableValuePair{
		VariableName: rules.IslandTaxContribution,
		Values:       []float64{float64(valToBeReturned)},
	}
	isCompliant := c.CheckCompliance(rules.IslandTaxContribution)
	if isCompliant {
		// TODO: with this compliance check, agents can see whether they'd like to continue returning this value
		valToBeReturned = 0
	}
	// Use the toolkit to recommend a value
	newVal, success := c.GetRecommendation(rules.IslandTaxContribution)
	if success {
		// TODO: Choose whether to use this compliant value
		valToBeReturned = shared.Resources(newVal.Values[rules.SingleValueVariableEntry])
	}
	return valToBeReturned
}

// GetSanctionPayment is called at the end of turn to pay any sanctions that have been
// imposed on the island, it is up to the island if they choose to pay the sanction or not.
// COMPULSORY
func (c *BaseClient) GetSanctionPayment() shared.Resources {
	valToBeReturned := shared.Resources(0)
	c.LocalVariableCache[rules.SanctionPaid] = rules.VariableValuePair{
		VariableName: rules.SanctionPaid,
		Values:       []float64{float64(valToBeReturned)},
	}
	isCompliant := c.CheckCompliance(rules.SanctionPaid)
	if isCompliant {
		// TODO: with this compliance check, agents can see whether they'd like to continue returning this value
		valToBeReturned = 0
	}
	// Use the toolkit to recommend a value
	newVal, success := c.GetRecommendation(rules.SanctionPaid)
	if success {
		// TODO: Choose whether to use this compliant value
		valToBeReturned = shared.Resources(newVal.Values[rules.SingleValueVariableEntry])
	}
	return valToBeReturned
}

// RequestAllocation is called at the end of the turn to request resources from the
// common pool. If there are enough resources in the common pool, the server will
// update your island's pool with the resources you requested.
// COMPULSORY
func (c *BaseClient) RequestAllocation() shared.Resources {
	// TODO: Implement request equal to the allocation permitted by President.
	valToBeReturned := shared.Resources(0)
	c.LocalVariableCache[rules.IslandAllocation] = rules.VariableValuePair{
		VariableName: rules.IslandAllocation,
		Values:       []float64{float64(valToBeReturned)},
	}
	isCompliant := c.CheckCompliance(rules.IslandAllocation)
	if isCompliant {
		// TODO: with this compliance check, agents can see whether they'd like to continue returning this value
		return valToBeReturned
	}
	// Use the toolkit to recommend a value
	newVal, success := c.GetRecommendation(rules.IslandAllocation)
	if success {
		// TODO: Choose whether to use this compliant value
		valToBeReturned = shared.Resources(newVal.Values[rules.SingleValueVariableEntry])
	}
	return valToBeReturned
}

// CheckCompliance provides clients with an easy interface to feed a variable and check whether it is compliant
// with all the rules that are affected by it
// OPTIONAL
func (c *BaseClient) CheckCompliance(variable rules.VariableFieldName) bool {
	rulesAffected, found := rules.PickUpRulesByVariable(variable, rules.RulesInPlay, c.LocalVariableCache)
	complianceCheck := true
	if !found {
		return false
	}
	for _, ruleName := range rulesAffected {
		compliant, err := rules.ComplianceCheck(rules.RulesInPlay[ruleName], c.LocalVariableCache)
		if err != nil {
			c.Logf("Attempted to evaluate rule, failed with error %v", err)
			return false
		}
		complianceCheck = complianceCheck && compliant
	}
	return complianceCheck
}

// GetRecommendation provides clients with a way of working out (in reasonably simple cases) what value
// a given variable must be to ensure compliance
// OPTIONAL
func (c *BaseClient) GetRecommendation(variable rules.VariableFieldName) (compliantValue rules.VariableValuePair, success bool) {
	rulesAffected, found := rules.PickUpRulesByVariable(variable, rules.RulesInPlay, c.LocalVariableCache)
	if !found {
		return c.LocalVariableCache[variable], false
	}
	for _, ruleName := range rulesAffected {
		newMap, ok := rules.ComplianceRecommendation(rules.RulesInPlay[ruleName], c.LocalVariableCache)
		if ok {
			return newMap[variable], ok
		}
	}
	return c.LocalVariableCache[variable], false
}

package rules

import "github.com/pkg/errors"

func ComplianceCheck(rule RuleMatrix, variables map[VariableFieldName]VariableValuePair) (compliant bool, ruleError RuleError) {
	if checkAllVariablesAvailable(rule.RequiredVariables, variables) {

	}
	return false, RuleError{
		ErrorType: VariableCacheDidNotHaveAllRequiredVariables,
		Err:       errors.Errorf("Variable cache provided didn't have all the errors for this rule"),
	}
}

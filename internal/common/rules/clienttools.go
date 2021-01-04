package rules

import "github.com/pkg/errors"

func ComplianceCheck(rule RuleMatrix, variables map[VariableFieldName]VariableValuePair) (compliant bool, ruleError error) {
	if checkAllVariablesAvailable(rule.RequiredVariables, variables) {
		evalResult := EvaluateRuleFromCaches(rule.RuleName, map[string]RuleMatrix{
			rule.RuleName: rule,
		}, variables)
		return evalResult.RulePasses, evalResult.EvalError
	}
	return false, &RuleError {
		ErrorType: VariableCacheDidNotHaveAllRequiredVariables,
		Err:       errors.Errorf("Variable cache provided didn't have all the errors for this rule"),
	}
}

//func ComplianceRecommendation(rule RuleMatrix, variables map[VariableFieldName]VariableValuePair) map[VariableFieldName]VariableValuePair {
//
//}
//
//
//func fetchRequiredVariables(reqVariables []VariableFieldName, variables map[VariableFieldName]VariableValuePair) []float64 {
//
//}
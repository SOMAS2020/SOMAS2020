package rules

import (
	"github.com/pkg/errors"
	"gonum.org/v1/gonum/mat"
)

// AvailableRules is a global cache of all rules that are available to agents
var AvailableRules = map[string]RuleMatrix{}

// RulesInPlay is a global cache of all rules currently in effect
var RulesInPlay = map[string]RuleMatrix{}

// RegisterNewRule Creates and registers new rule based on inputs
func RegisterNewRule(ruleName string, requiredVariables []string, applicableMatrix mat.Dense, auxiliaryVector mat.VecDense) (*RuleMatrix, error) {
	return registerNewRuleInternal(ruleName, requiredVariables, applicableMatrix, auxiliaryVector, AvailableRules)
}

// registerNewRuleInternal provides primal register logic for any rule cache
func registerNewRuleInternal(ruleName string, requiredVariables []string, applicableMatrix mat.Dense, auxiliaryVector mat.VecDense, ruleStore map[string]RuleMatrix) (*RuleMatrix, error) {
	if _, ok := ruleStore[ruleName]; ok {
		return nil, errors.Errorf("Rule '%v' already registered", ruleName)
	}

	rm := RuleMatrix{RuleName: ruleName, RequiredVariables: requiredVariables, ApplicableMatrix: applicableMatrix, AuxiliaryVector: auxiliaryVector}
	ruleStore[ruleName] = rm
	return &rm, nil
}

// PullRuleIntoPlay provides engagement logic for global rules in play cache
func PullRuleIntoPlay(rulename string) error {
	return pullRuleIntoPlayInternal(rulename, AvailableRules, RulesInPlay)
}

// pullRuleIntoPlayInternal provides primal rule engagement logic for any pair of caches
func pullRuleIntoPlayInternal(rulename string, allRules map[string]RuleMatrix, playRules map[string]RuleMatrix) error {
	if _, ok := allRules[rulename]; ok {
		if _, ok := playRules[rulename]; ok {
			return errors.Errorf("Rule '%v' is already in play", rulename)
		}
		playRules[rulename] = allRules[rulename]
		return nil
	} else {
		return errors.Errorf("Rule '%v' is not available in rules cache", rulename)
	}
}

// PullRuleOutOfPlay provides disengagement logic for global rules in play cache
func PullRuleOutOfPlay(rulename string) error {
	return pullRuleOutOfPlayInternal(rulename, AvailableRules, RulesInPlay)
}

// pullRuleOutOfPlayInternal provides primal rule disengagement logic for any pair of caches
func pullRuleOutOfPlayInternal(rulename string, allRules map[string]RuleMatrix, playRules map[string]RuleMatrix) error {
	if _, ok := allRules[rulename]; ok {
		if _, ok := playRules[rulename]; ok {
			delete(playRules, rulename)
			return nil
		} else {
			return errors.Errorf("Rule '%v' is not in play", rulename)
		}
	} else {
		return errors.Errorf("Rule '%v' is not available in rules cache", rulename)
	}
}

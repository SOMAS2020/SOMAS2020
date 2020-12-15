package rules

import (
	"github.com/pkg/errors"
	"gonum.org/v1/gonum/mat"
)

var AvailableRules = map[string]RuleMatrix{}

var RulesInPlay = map[string]RuleMatrix{}

// RegisterNewRule Creates and registers new rule based on inputs
func RegisterNewRule(ruleName string, requiredVariables []string, applicableMatrix mat.Dense, auxiliaryVector mat.VecDense) (*RuleMatrix, error) {
	if _, ok := AvailableRules[ruleName]; ok {
		return nil, errors.Errorf("Rule '%v' already registered", ruleName)
	}

	rm := RuleMatrix{ruleName: ruleName, RequiredVariables: requiredVariables, ApplicableMatrix: applicableMatrix, AuxiliaryVector: auxiliaryVector}
	AvailableRules[ruleName] = rm
	return &rm, nil
}

func PullRuleIntoPlay(rulename string) error {
	if _, ok := AvailableRules[rulename]; ok {
		if _, ok := RulesInPlay[rulename]; ok {
			return errors.Errorf("Rule '%v' is already in play", rulename)
		}
		RulesInPlay[rulename] = AvailableRules[rulename]
		return nil
	} else {
		return errors.Errorf("Rule '%v' is not available in rules cache", rulename)
	}
}

func PullRuleOutOfPlay(rulename string) error {
	if _, ok := RulesInPlay[rulename]; ok {
		delete(RulesInPlay, rulename)
		return nil
	} else {
		return errors.Errorf("Rule '%v' is already in play", rulename)
	}
}

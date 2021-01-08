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
func RegisterNewRule(ruleName string, requiredVariables []VariableFieldName, applicableMatrix mat.Dense, auxiliaryVector mat.VecDense, mutable bool, link RuleLink) (constructedMatrix *RuleMatrix, Error error) {
	return RegisterNewRuleInternal(ruleName, requiredVariables, applicableMatrix, auxiliaryVector, AvailableRules, mutable, link)
}

// RegisterNewRuleInternal provides primal register logic for any rule cache
func RegisterNewRuleInternal(ruleName string, requiredVariables []VariableFieldName, applicableMatrix mat.Dense, auxiliaryVector mat.VecDense, ruleStore map[string]RuleMatrix, mutable bool, link RuleLink) (constructedMatrix *RuleMatrix, Error error) {
	if _, ok := ruleStore[ruleName]; ok {
		return nil, &RuleError{Err: errors.Errorf("Rule '%v' already in rule cache", ruleName), ErrorType: TriedToReRegisterRule}
	}

	rm := RuleMatrix{RuleName: ruleName, RequiredVariables: requiredVariables, ApplicableMatrix: applicableMatrix, AuxiliaryVector: auxiliaryVector, Mutable: mutable, Link: link}
	ruleStore[ruleName] = rm
	return &rm, nil
}

// PullRuleIntoPlay provides engagement logic for global rules in play cache
func PullRuleIntoPlay(rulename string) error {
	return PullRuleIntoPlayInternal(rulename, AvailableRules, RulesInPlay)
}

// PullRuleIntoPlayInternal provides primal rule engagement logic for any pair of caches
func PullRuleIntoPlayInternal(rulename string, allRules map[string]RuleMatrix, playRules map[string]RuleMatrix) error {
	if _, ok := allRules[rulename]; ok {
		if _, ok := playRules[rulename]; ok {
			return &RuleError{Err: errors.Errorf("Rule '%v' is already in play", rulename), ErrorType: RuleIsAlreadyInPlay}
		}
		linkRule, linked := checkLinking(rulename, allRules)
		if linked {
			playRules[linkRule] = allRules[linkRule]
		}
		playRules[rulename] = allRules[rulename]
		return nil
	}
	return &RuleError{Err: errors.Errorf("Rule '%v' does not exist in available rules", rulename), ErrorType: RuleNotInAvailableRulesCache}
}

// PullRuleOutOfPlay provides disengagement logic for global rules in play cache
func PullRuleOutOfPlay(rulename string) error {
	return PullRuleOutOfPlayInternal(rulename, AvailableRules, RulesInPlay)
}

// PullRuleOutOfPlayInternal provides primal rule disengagement logic for any pair of caches
func PullRuleOutOfPlayInternal(rulename string, allRules map[string]RuleMatrix, playRules map[string]RuleMatrix) error {
	if _, ok := allRules[rulename]; ok {
		if _, ok := playRules[rulename]; ok {
			linkRule, linked := checkLinking(rulename, allRules)
			if linked {
				delete(playRules, linkRule)
			}
			delete(playRules, rulename)
			return nil
		}
		return &RuleError{Err: errors.Errorf("Rule '%v' is not in play", rulename), ErrorType: RuleIsNotInPlay}
	}
	return &RuleError{Err: errors.Errorf("Rule '%v' does not exist in available rules cache", rulename), ErrorType: RuleNotInAvailableRulesCache}
}

// ModifyRule allows for rules that are flagged as mutable to be modified
func ModifyRule(rulename string, newMatrix mat.Dense, newAuxiliary mat.VecDense) error {
	return ModifyRuleInternal(rulename, newMatrix, newAuxiliary, AvailableRules, RulesInPlay)
}

func ModifyRuleInternal(rulename string, newMatrix mat.Dense, newAuxiliary mat.VecDense, rulesCache map[string]RuleMatrix, inPlayCache map[string]RuleMatrix) error {
	if _, ok := rulesCache[rulename]; ok {
		oldRuleMatrix := rulesCache[rulename]
		if !oldRuleMatrix.Mutable {
			return &RuleError{Err: errors.Errorf("Rule '%v' is not mutable", rulename), ErrorType: RuleRequestedForModificationWasImmutable}
		}
		oldMatrix := oldRuleMatrix.ApplicableMatrix
		_, ocols := oldMatrix.Dims()
		nrows, ncols := newMatrix.Dims()
		if ocols != ncols {
			return &RuleError{Err: errors.Errorf("Provided Rule matrix '%v' has a dimension mismatch with old matrix '%v'", newMatrix, oldMatrix), ErrorType: ModifiedRuleMatrixDimensionMismatch}
		}
		auxRows, _ := newAuxiliary.Dims()
		if nrows != auxRows {
			return &RuleError{Err: errors.Errorf("Aux vector '%v' has dimension mismatch with Rule Matrix '%v'", newAuxiliary, newMatrix), ErrorType: AuxVectorDimensionDontMatchRuleMatrix}
		}
		newRuleMatrix := RuleMatrix{
			RuleName:          oldRuleMatrix.RuleName,
			RequiredVariables: oldRuleMatrix.RequiredVariables,
			ApplicableMatrix:  newMatrix,
			AuxiliaryVector:   newAuxiliary,
			Mutable:           true,
		}
		rulesCache[rulename] = newRuleMatrix
		if _, ok := inPlayCache[rulename]; ok {
			inPlayCache[rulename] = rulesCache[rulename]
		}
		return nil
	}
	return &RuleError{Err: errors.Errorf("Rule '%v' does not exist in available rules cache", rulename), ErrorType: RuleNotInAvailableRulesCache}
}

func checkLinking(ruleName string, availableRules map[string]RuleMatrix) (string, bool) {
	if rule, ok := availableRules[ruleName]; ok {
		if rule.Link.Linked {
			return rule.Link.LinkedRule, true
		}
		for _, ruleVal := range availableRules {
			if ruleVal.Link.Linked {
				if ruleVal.Link.LinkedRule == ruleName {
					return ruleVal.RuleName, true
				}
			}
		}
	}
	return "", false
}

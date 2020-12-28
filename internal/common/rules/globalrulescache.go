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
func RegisterNewRule(ruleName string, requiredVariables []VariableFieldName, applicableMatrix mat.Dense, auxiliaryVector mat.VecDense, mutable bool) (constructedMatrix *RuleMatrix, Error error) {
	return registerNewRuleInternal(ruleName, requiredVariables, applicableMatrix, auxiliaryVector, AvailableRules, mutable)
}

// registerNewRuleInternal provides primal register logic for any rule cache
func registerNewRuleInternal(ruleName string, requiredVariables []VariableFieldName, applicableMatrix mat.Dense, auxiliaryVector mat.VecDense, ruleStore map[string]RuleMatrix, mutable bool) (constructedMatrix *RuleMatrix, Error error) {
	if _, ok := ruleStore[ruleName]; ok {
		return nil, &RuleError{err: errors.Errorf("Rule '%v' already in rule cache", ruleName), errorType: TriedToReRegisterRule}
	}

	rm := RuleMatrix{RuleName: ruleName, RequiredVariables: requiredVariables, ApplicableMatrix: applicableMatrix, AuxiliaryVector: auxiliaryVector, Mutable: mutable}
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
			return &RuleError{err: errors.Errorf("Rule '%v' is already in play", rulename), errorType: RuleIsAlreadyInPlay}
		}
		playRules[rulename] = allRules[rulename]
		return nil
	}
	return &RuleError{err: errors.Errorf("Rule '%v' does not exist in available rules", rulename), errorType: RuleNotInAvailableRulesCache}
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
		}
		return &RuleError{err: errors.Errorf("Rule '%v' is not in play", rulename), errorType: RuleIsNotInPlay}
	}
	return &RuleError{err: errors.Errorf("Rule '%v' does not exist in available rules cache", rulename), errorType: RuleNotInAvailableRulesCache}
}

func ModifyRule(rulename string, newMatrix mat.Dense, newAuxiliary mat.VecDense) error {
	return modifyRuleInternal(rulename, newMatrix, newAuxiliary, AvailableRules, RulesInPlay)
}

func modifyRuleInternal(rulename string, newMatrix mat.Dense, newAuxiliary mat.VecDense, rulesCache map[string]RuleMatrix, inPlayCache map[string]RuleMatrix) error {
	if _, ok := rulesCache[rulename]; ok {
		oldRuleMatrix := rulesCache[rulename]
		if !oldRuleMatrix.Mutable {
			return &RuleError{err: errors.Errorf("Rule '%v' is not mutable", rulename), errorType: RuleRequestedForModificationWasImmutable}
		}
		oldMatrix := oldRuleMatrix.ApplicableMatrix
		_, ocols := oldMatrix.Dims()
		nrows, ncols := newMatrix.Dims()
		if ocols != ncols {
			return &RuleError{err: errors.Errorf("Provided Rule matrix '%v' has a dimension mismatch with old matrix '%v'", newMatrix, oldMatrix), errorType: ModifiedRuleMatrixDimensionMismatch}
		}
		auxRows, _ := newAuxiliary.Dims()
		if nrows != auxRows {
			return &RuleError{err: errors.Errorf("Aux vector '%v' has dimension mismatch with Rule Matrix '%v'", newAuxiliary, newMatrix), errorType: AuxVectorDimensionDontMatchRuleMatrix}
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
	} else {
		return &RuleError{err: errors.Errorf("Rule '%v' does not exist in available rules cache", rulename), errorType: RuleNotInAvailableRulesCache}
	}
}

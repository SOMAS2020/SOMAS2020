package rules

import (
	"github.com/pkg/errors"
	"gonum.org/v1/gonum/mat"
)

// AvailableRules is a global cache of all rules that are available to agents
var AvailableRules = map[string]RuleMatrix{}

// RulesInPlay is a global cache of all rules currently in effect
var RulesInPlay = map[string]RuleMatrix{}

// A RuleError is a non-critical issue which can be caused by an island trying to modify, register or pick rules which isn't mechanically feasible
type RuleError int

const (
	RuleNotInAvailableRulesCache RuleError = iota
	ModifiedRuleMatrixDimensionMismatch
	AuxVectorDimensionDontMatchRuleMatrix
	RuleRequestedForModificationWasImmutable
	None
)

// RegisterNewRule Creates and registers new rule based on inputs
func RegisterNewRule(ruleName string, requiredVariables []VariableFieldName, applicableMatrix mat.Dense, auxiliaryVector mat.VecDense, mutable bool) (*RuleMatrix, error) {
	return registerNewRuleInternal(ruleName, requiredVariables, applicableMatrix, auxiliaryVector, AvailableRules, mutable)
}

// registerNewRuleInternal provides primal register logic for any rule cache
func registerNewRuleInternal(ruleName string, requiredVariables []VariableFieldName, applicableMatrix mat.Dense, auxiliaryVector mat.VecDense, ruleStore map[string]RuleMatrix, mutable bool) (*RuleMatrix, error) {
	if _, ok := ruleStore[ruleName]; ok {
		return nil, errors.Errorf("Rule '%v' already registered", ruleName)
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
			return errors.Errorf("Rule '%v' is already in play", rulename)
		}
		playRules[rulename] = allRules[rulename]
		return nil
	}
	return errors.Errorf("Rule '%v' is not available in rules cache", rulename)
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
		return errors.Errorf("Rule '%v' is not in play", rulename)
	}
	return errors.Errorf("Rule '%v' is not available in rules cache", rulename)
}

func ModifyRule(rulename string, newMatrix mat.Dense, newAuxiliary mat.VecDense) (successfulModification bool, status RuleError) {
	return modifyRuleInternal(rulename, newMatrix, newAuxiliary, AvailableRules, RulesInPlay)
}

func modifyRuleInternal(rulename string, newMatrix mat.Dense, newAuxiliary mat.VecDense, rulesCache map[string]RuleMatrix, inPlayCache map[string]RuleMatrix) (success bool, message RuleError) {
	if _, ok := rulesCache[rulename]; ok {
		oldRuleMatrix := rulesCache[rulename]
		if !oldRuleMatrix.Mutable {
			return false, RuleRequestedForModificationWasImmutable
		}
		oldMatrix := oldRuleMatrix.ApplicableMatrix
		_, ocols := oldMatrix.Dims()
		nrows, ncols := newMatrix.Dims()
		if ocols != ncols {
			return false, ModifiedRuleMatrixDimensionMismatch
		}
		auxRows, _ := newAuxiliary.Dims()
		if nrows != auxRows {
			return false, AuxVectorDimensionDontMatchRuleMatrix
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
		return true, None
	} else {
		return false, RuleNotInAvailableRulesCache
	}
}

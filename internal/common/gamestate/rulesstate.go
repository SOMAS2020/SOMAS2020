package gamestate

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"gonum.org/v1/gonum/mat"
)

// Rules cache manipulation functions

func (g *GameState) RegisterNewRule(ruleName string, requiredVariables []rules.VariableFieldName, applicableMatrix mat.Dense, auxiliaryVector mat.VecDense, mutable bool, link rules.RuleLink) (constructedMatrix *rules.RuleMatrix, Error error) {
	return rules.RegisterNewRuleInternal(ruleName, requiredVariables, applicableMatrix, auxiliaryVector, g.RulesInfo.AvailableRules, mutable, link)
}

func (g *GameState) PullRuleIntoPlay(rulename string) error {
	return rules.PullRuleIntoPlayInternal(rulename, g.RulesInfo.AvailableRules, g.RulesInfo.CurrentRulesInPlay)
}

func (g *GameState) PullRuleOutOfPlay(rulename string) error {
	return rules.PullRuleOutOfPlayInternal(rulename, g.RulesInfo.AvailableRules, g.RulesInfo.CurrentRulesInPlay)
}

func (g *GameState) ModifyRule(rulename string, newMatrix mat.Dense, newAuxiliary mat.VecDense) error {
	return rules.ModifyRuleInternal(rulename, newMatrix, newAuxiliary, g.RulesInfo.AvailableRules, g.RulesInfo.CurrentRulesInPlay)
}

// Variables cache manipulation functions

func (g *GameState) RegisterNewVariable(pair rules.VariableValuePair) error {
	return rules.RegisterNewVariableInternal(pair, g.RulesInfo.VariableMap)
}

func (g *GameState) UpdateVariable(variableName rules.VariableFieldName, newValue rules.VariableValuePair) bool {
	return rules.UpdateVariableInternal(variableName, newValue, g.RulesInfo.VariableMap)
}

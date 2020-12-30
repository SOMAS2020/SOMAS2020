package iigointernal

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type monitor struct {
	judgeID           shared.ClientID
	speakerID         shared.ClientID
	presidentID       shared.ClientID
	internalIIGOCache []shared.Accountability
}

func (m *monitor) addToCache(roleToMonitorID shared.ClientID, variables []rules.VariableFieldName, values [][]float64) {
	pairs := []rules.VariableValuePair{}
	for index, variable := range variables {
		pairs = append(pairs, rules.VariableValuePair{
			VariableName: variable,
			Values:       values[index],
		})
	}
	m.internalIIGOCache = append(m.internalIIGOCache, shared.Accountability{
		ClientID: roleToMonitorID,
		Pairs:    pairs,
	})
}

func (m *monitor) monitorRole(roleToMonitorID shared.ClientID) bool {
	performedRoleCorrectly := true
	for _, entry := range m.internalIIGOCache {
		if entry.ClientID == roleToMonitorID {
			variablePairs := entry.Pairs
			var rulesAffected []string
			for _, variable := range variablePairs {
				valuesToBeAdded, foundRules := rules.PickUpRulesByVariable(variable.VariableName, rules.RulesInPlay)
				if foundRules {
					rulesAffected = append(rulesAffected, valuesToBeAdded...)
				}
				rules.UpdateVariable(variable.VariableName, variable)
			}
			for _, rule := range rulesAffected {
				evaluation, err := rules.BasicBooleanRuleEvaluator(rule)
				if err != nil {
					continue
				}
				performedRoleCorrectly = evaluation && performedRoleCorrectly
			}
		}
	}
	return performedRoleCorrectly
}

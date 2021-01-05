package iigointernal

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/pkg/errors"
)

type monitor struct {
	judgeID           shared.ClientID
	speakerID         shared.ClientID
	presidentID       shared.ClientID
	internalIIGOCache []shared.Accountability
	TermLengths       map[shared.Role]uint
}

func (m *monitor) addToCache(roleToMonitorID shared.ClientID, variables []rules.VariableFieldName, values [][]float64) {
	pairs := []rules.VariableValuePair{}
	if len(variables) == len(values) {
		for index, variable := range variables {
			pairs = append(pairs, rules.MakeVariableValuePair(variable, values[index]))
		}
		m.internalIIGOCache = append(m.internalIIGOCache, shared.Accountability{
			ClientID: roleToMonitorID,
			Pairs:    pairs,
		})
	}
}

func (m *monitor) monitorRole(g *gamestate.GameState, roleAccountable baseclient.Client) shared.MonitorResult {
	roleToMonitor, roleName, err := m.findRoleToMonitor(roleAccountable.GetID())
	if err == nil {
		decideToMonitor := roleAccountable.MonitorIIGORole(roleName)
		evaluationResult := false
		if decideToMonitor {
			evaluationResult = m.evaluateCache(roleToMonitor, rules.RulesInPlay)
		}

		evaluationResultAnnounce, announce := roleAccountable.DecideIIGOMonitoringAnnouncement(evaluationResult)

		//announce == decideToMonitor
		variablesToCache := []rules.VariableFieldName{rules.MonitorRoleAnnounce, rules.MonitorRoleDecideToMonitor}
		valuesToCache := [][]float64{{boolToFloat(decideToMonitor)}, {boolToFloat(announce)}}
		m.addToCache(roleAccountable.GetID(), variablesToCache, valuesToCache)

		if announce {
			//check if evalResult = o.g. evalResult
			variablesToCache := []rules.VariableFieldName{rules.MonitorRoleEvalResult, rules.MonitorRoleEvalResultDecide}
			valuesToCache := [][]float64{{boolToFloat(evaluationResult)}, {boolToFloat(evaluationResultAnnounce)}}
			m.addToCache(roleAccountable.GetID(), variablesToCache, valuesToCache)

			message := generateMonitoringMessage(roleName, evaluationResult)
			broadcastToAllIslands(roleAccountable.GetID(), message)

			g.IIGOTurnsInPower[roleName] = m.TermLengths[roleName]
		}

		result := shared.MonitorResult{Performed: decideToMonitor, Result: evaluationResult}
		return result
	}
	result := shared.MonitorResult{Performed: false, Result: false}
	return result
}

func (m *monitor) evaluateCache(roleToMonitorID shared.ClientID, ruleStore map[string]rules.RuleMatrix) bool {
	performedRoleCorrectly := true
	for _, entry := range m.internalIIGOCache {
		if entry.ClientID == roleToMonitorID {
			variablePairs := entry.Pairs
			var rulesAffected []string
			for _, variable := range variablePairs {
				valuesToBeAdded, foundRules := rules.PickUpRulesByVariable(variable.VariableName, ruleStore, rules.VariableMap)
				if foundRules {
					rulesAffected = append(rulesAffected, valuesToBeAdded...)
				}
				rules.UpdateVariable(variable.VariableName, variable)
			}
			for _, rule := range rulesAffected {
				evaluation, err := rules.BasicBooleanRuleEvaluator(rule)
				if err == nil {
					performedRoleCorrectly = evaluation && performedRoleCorrectly
				}
			}
		}
	}
	return performedRoleCorrectly
}

func (m *monitor) findRoleToMonitor(roleAccountable shared.ClientID) (shared.ClientID, shared.Role, error) {
	switch roleAccountable {
	case m.speakerID:
		return m.presidentID, shared.President, nil
	case m.presidentID:
		return m.judgeID, shared.Judge, nil
	case m.judgeID:
		return m.speakerID, shared.Speaker, nil
	default:
		return shared.ClientID(-1), shared.Speaker, errors.Errorf("Monitoring by island that is not an IIGO Role")
	}
}

func generateMonitoringMessage(role shared.Role, result bool) map[shared.CommunicationFieldName]shared.CommunicationContent {
	returnMap := map[shared.CommunicationFieldName]shared.CommunicationContent{}

	returnMap[shared.RoleMonitored] = shared.CommunicationContent{
		T:            shared.CommunicationIIGORole,
		IIGORoleData: role,
	}
	returnMap[shared.MonitoringResult] = shared.CommunicationContent{
		T:           shared.CommunicationBool,
		BooleanData: result,
	}

	return returnMap
}

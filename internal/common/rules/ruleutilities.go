package rules

// PickUpRulesByVariable returns a list of rule_id's which are affected by certain variables.
func PickUpRulesByVariable(variableName VariableFieldName, ruleStore map[string]RuleMatrix) ([]string, bool) {
	var Rules []string
	if _, ok := VariableMap[variableName]; ok {
		for k, v := range ruleStore {
			_, found := searchForVariableInArray(variableName, v.RequiredVariables)
			if found {
				Rules = append(Rules, k)
			}
		}
		return Rules, true
	}
	return []string{}, false
}

func searchForVariableInArray(val VariableFieldName, array []VariableFieldName) (int, bool) {
	for i, v := range array {
		if v == val {
			return i, true
		}
	}
	return -1, false
}

// MakeVariableValuePair creates a VariableValuePair from the variable name and its value
func MakeVariableValuePair(variable VariableFieldName, value []float64) VariableValuePair {
	return VariableValuePair{
		VariableName: variable,
		Values:       value,
	}
}

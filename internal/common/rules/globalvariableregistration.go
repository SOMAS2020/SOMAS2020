package rules

func init() {
	for _, v := range StaticVariables {
		e := RegisterNewVariable(v)
		if e != nil {
			panic("variable registration gone wrong, variable: " + v.VariableName + " has been registered multiple times")
		}
	}
}

var StaticVariables = []VariableValuePair{
	{
		VariableName: "number_of_islands_contributing_to_common_pool",
		Multivalued:  false,
		SingleValue:  5,
		MultiValue:   nil,
	},
	{
		VariableName: "number_of_failed_forages",
		Multivalued:  false,
		SingleValue:  0.5,
		MultiValue:   nil,
	},
	{
		VariableName: "number_of_broken_agreements",
		Multivalued:  false,
		SingleValue:  1,
		MultiValue:   nil,
	},
	{
		VariableName: "max_severity_of_sanctions",
		Multivalued:  false,
		SingleValue:  2,
		MultiValue:   nil,
	},
}

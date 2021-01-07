package rules

import "testing"

func TestIIGOMonitorRulePermission1(t *testing.T) {
	cases := []struct {
		decideToMonitor float64
		announce        float64
		expected        bool
	}{
		{
			decideToMonitor: 0,
			announce:        0,
			expected:        true,
		},
		{
			decideToMonitor: 0,
			announce:        1,
			expected:        false,
		},
		{
			decideToMonitor: 1,
			announce:        0,
			expected:        false,
		},
		{
			decideToMonitor: 1,
			announce:        1,
			expected:        true,
		},
	}

	for _, tc := range cases {
		UpdateVariable(MonitorRoleDecideToMonitor, MakeVariableValuePair(MonitorRoleDecideToMonitor, []float64{tc.decideToMonitor}))
		UpdateVariable(MonitorRoleAnnounce, MakeVariableValuePair(MonitorRoleAnnounce, []float64{tc.announce}))
		eval := EvaluateRule("iigo_monitor_rule_permission_1")

		if eval.RulePasses != tc.expected {
			t.Errorf("MonitorRulePermission1 - Failed. Input (%v,%v). Rule evaluated to %v, but expected %v.",
				tc.decideToMonitor, tc.announce, eval, tc.expected)
		}
	}
}

func TestIIGOMonitorRulePermission2(t *testing.T) {
	cases := []struct {
		evalResult       float64
		evalResultDecide float64
		expected         bool
	}{
		{
			evalResult:       0,
			evalResultDecide: 0,
			expected:         true,
		},
		{
			evalResult:       0,
			evalResultDecide: 1,
			expected:         false,
		},
		{
			evalResult:       1,
			evalResultDecide: 0,
			expected:         false,
		},
		{
			evalResult:       1,
			evalResultDecide: 1,
			expected:         true,
		},
	}

	for _, tc := range cases {
		UpdateVariable(MonitorRoleEvalResult, MakeVariableValuePair(MonitorRoleEvalResult, []float64{tc.evalResult}))
		UpdateVariable(MonitorRoleEvalResultDecide, MakeVariableValuePair(MonitorRoleEvalResultDecide, []float64{tc.evalResultDecide}))
		eval := EvaluateRule("iigo_monitor_rule_permission_2")

		if eval.RulePasses != tc.expected {
			t.Errorf("MonitorRulePermission2 - Failed. Input (%v,%v). Rule evaluated to %v, but expected %v.",
				tc.evalResult, tc.evalResultDecide, eval, tc.expected)
		}
	}
}

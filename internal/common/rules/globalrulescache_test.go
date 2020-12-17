package rules

import (
	"github.com/SOMAS2020/SOMAS2020/pkg/testutils"
	"github.com/pkg/errors"
	"gonum.org/v1/gonum/mat"
	"testing"
)

// TestRegisterNewRule Tests whether the global rule cache is able to register new rules
func TestRegisterNewRule(t *testing.T) {
	registerTestRule()
	if _, ok := AvailableRules["Kinda Test Rule"]; !ok {
		t.Errorf("Global rule register unable to register new rules")
	}
}

func TestPullRuleIntoPlay(t *testing.T) {
	registerTestRule()
	_ = PullRuleIntoPlay("Kinda Test Rule 2")
	cases := []struct {
		name string
		rule string
		want error
	}{
		{
			name: "normal working",
			rule: "Kinda Test Rule",
			want: nil,
		},
		{
			name: "unidentified rule name",
			rule: "Unknown Rule",
			want: errors.Errorf("Rule '%v' is not available in rules cache", "Unknown Rule"),
		},
		{
			name: "Rule already in play",
			rule: "Kinda Test Rule 2",
			want: errors.Errorf("Rule '%v' is already in play", "Kinda Test Rule 2"),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := PullRuleIntoPlay(tc.rule)
			testutils.CompareTestErrors(tc.want, got, t)
		})
	}
}

func TestPullRuleOutOfPlay(t *testing.T) {
	registerTestRule()
	_ = PullRuleOutOfPlay("Kinda Test Rule 2")
	_ = PullRuleIntoPlay("Kinda Test Rule")
	cases := []struct {
		name string
		rule string
		want error
	}{
		{
			name: "normal working",
			rule: "Kinda Test Rule",
			want: nil,
		},
		{
			name: "Rule already in play",
			rule: "Kinda Test Rule 2",
			want: errors.Errorf("Rule '%v' is not in play", "Kinda Test Rule 2"),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := PullRuleOutOfPlay(tc.rule)
			testutils.CompareTestErrors(tc.want, got, t)
		})
	}
}

func registerTestRule() {

	//A very contrived rule//
	name := "Kinda Test Rule"
	reqVar := []string{
		"number_of_islands_contributing_to_common_pool",
		"number_of_failed_forages",
		"number_of_broken_agreements",
		"max_severity_of_sanctions",
	}

	v := []float64{1, 0, 0, 0, -4, 0, -1, -1, 0, 2, 0, 0, 0, 1, -2, 0, 0, 1, 0, -1}
	CoreMatrix := mat.NewDense(4, 5, v)
	aux := []float64{2, 3, 3, 2}
	AuxiliaryVector := mat.NewVecDense(4, aux)

	_, _ = RegisterNewRule(name, reqVar, *CoreMatrix, *AuxiliaryVector)
	// Check internal/clients/team3/client.go for an implementation of a basic evaluator for this rule
	//A very contrived rule//
	name = "Kinda Test Rule 2"

	_, _ = RegisterNewRule(name, reqVar, *CoreMatrix, *AuxiliaryVector)
}

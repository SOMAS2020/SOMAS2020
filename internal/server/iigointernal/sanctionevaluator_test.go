package iigointernal

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"reflect"
	"testing"
)

func TestEvaluateSanction(t *testing.T) {
	cases := []struct {
		name           string
		sanction       rules.RuleMatrix
		localVarCache  map[rules.VariableFieldName]rules.VariableValuePair
		expectedResult shared.Resources
	}{
		{
			name:           "Basic Sanction Test case",
			sanction:       generateArbitrarySanctionMatrix("Basic", 0, 0),
			localVarCache:  generateAndLoadDummyLocalVarCache(5, 5, 1),
			expectedResult: 0,
		},
		{
			name:           "Normal Sanction Test case",
			sanction:       generateArbitrarySanctionMatrix("Normal", 0.5, 50),
			localVarCache:  generateAndLoadDummyLocalVarCache(500, 1, 1),
			expectedResult: 300,
		},
		{
			name:           "More complicated Sanction Test case",
			sanction:       generateArbitrarySanctionMatrix("Normal", 0.7, 50),
			localVarCache:  generateAndLoadDummyLocalVarCache(500, 12, 3),
			expectedResult: 950,
		},
		{
			name:           "0 returned because sanction times out",
			sanction:       generateArbitrarySanctionMatrix("Normal", 0.7, 50),
			localVarCache:  generateAndLoadDummyLocalVarCache(500, 12, 0),
			expectedResult: 0,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := evaluateSanction(tc.sanction, tc.localVarCache)
			if !reflect.DeepEqual(res, tc.expectedResult) {
				t.Errorf("Expected %v got %v", tc.expectedResult, res)
			}
		})
	}
}

func generateAndLoadDummyLocalVarCache(islandReportedResources float64, constSanctionAmount float64, turnsLeft float64) map[rules.VariableFieldName]rules.VariableValuePair {
	emptyCache := generateDummyVariableCache()
	emptyCache[rules.IslandReportedResources] = rules.MakeVariableValuePair(rules.IslandReportedResources, []float64{islandReportedResources})
	emptyCache[rules.ConstSanctionAmount] = rules.MakeVariableValuePair(rules.ConstSanctionAmount, []float64{constSanctionAmount})
	emptyCache[rules.TurnsLeftOnSanction] = rules.MakeVariableValuePair(rules.TurnsLeftOnSanction, []float64{turnsLeft})
	return emptyCache
}

func generateArbitrarySanctionMatrix(ruleName string, currentResourcesFactor float64, constant float64) rules.RuleMatrix {
	coreVect := []float64{0, 0, 1, 0, currentResourcesFactor, constant, 0, 0}
	auxVect := []float64{1, 4}
	return rules.CompileRuleCase(rules.RawRuleSpecification{
		Name: ruleName,
		ReqVar: []rules.VariableFieldName{
			rules.IslandReportedResources,
			rules.ConstSanctionAmount,
			rules.TurnsLeftOnSanction,
		},
		V:       coreVect,
		Aux:     auxVect,
		Mutable: false,
	})
}

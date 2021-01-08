package rules

import (
	"reflect"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/pkg/testutils"
	"gonum.org/v1/gonum/mat"
)

func TestGenerateBoolList(t *testing.T) {
	cases := []struct {
		name        string
		fillVal     bool
		length      int
		expectedArr []bool
	}{
		{
			name:        "Basic fill test false",
			fillVal:     false,
			length:      5,
			expectedArr: []bool{false, false, false, false, false},
		},
		{
			name:        "Basic fill test true",
			fillVal:     true,
			length:      3,
			expectedArr: []bool{true, true, true},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := generateBoolList(tc.fillVal, tc.length)
			if !reflect.DeepEqual(res, tc.expectedArr) {
				t.Errorf("Expected %v got %v", tc.expectedArr, res)
			}
		})
	}
}

func TestUnfetchRequiredVariables(t *testing.T) {
	cases := []struct {
		name        string
		requiredVar []VariableFieldName
		variables   map[VariableFieldName]VariableValuePair
		newData     []float64
		expectedRes map[VariableFieldName]VariableValuePair
	}{
		{
			name: "Basic unfetch",
			requiredVar: []VariableFieldName{
				VoteCalled,
				NumberOfBallotsCast,
			},
			variables: map[VariableFieldName]VariableValuePair{
				VoteCalled: {
					VariableName: VoteCalled,
					Values:       []float64{5},
				},
				NumberOfBallotsCast: {
					VariableName: NumberOfBallotsCast,
					Values:       []float64{6},
				},
			},
			newData: []float64{7, 8},
			expectedRes: map[VariableFieldName]VariableValuePair{
				VoteCalled: {
					VariableName: VoteCalled,
					Values:       []float64{7},
				},
				NumberOfBallotsCast: {
					VariableName: NumberOfBallotsCast,
					Values:       []float64{8},
				},
			},
		},
		{
			name: "Complex unfetch",
			requiredVar: []VariableFieldName{
				VoteCalled,
				NumberOfBallotsCast,
			},
			variables: map[VariableFieldName]VariableValuePair{
				VoteCalled: {
					VariableName: VoteCalled,
					Values:       []float64{5, 7, 8},
				},
				NumberOfBallotsCast: {
					VariableName: NumberOfBallotsCast,
					Values:       []float64{6, 13, 5},
				},
			},
			newData: []float64{1, 2, 3, 4, 5, 6},
			expectedRes: map[VariableFieldName]VariableValuePair{
				VoteCalled: {
					VariableName: VoteCalled,
					Values:       []float64{1, 2, 3},
				},
				NumberOfBallotsCast: {
					VariableName: NumberOfBallotsCast,
					Values:       []float64{4, 5, 6},
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := unfetchRequiredVariables(tc.requiredVar, tc.variables, tc.newData)
			if !reflect.DeepEqual(res, tc.expectedRes) {
				t.Errorf("Expected %v got %v", tc.expectedRes, res)
			}
		})
	}
}

func TestFetchRequiredVariables(t *testing.T) {
	cases := []struct {
		name          string
		requiredVar   []VariableFieldName
		variables     map[VariableFieldName]VariableValuePair
		expectedData  []float64
		expectedFixed []bool
		expectedRes   bool
	}{
		{
			name: "Basic fetch",
			requiredVar: []VariableFieldName{
				VoteCalled,
				NumberOfBallotsCast,
			},
			variables: map[VariableFieldName]VariableValuePair{
				VoteCalled: {
					VariableName: VoteCalled,
					Values:       []float64{5},
				},
				NumberOfBallotsCast: {
					VariableName: NumberOfBallotsCast,
					Values:       []float64{6},
				},
			},
			expectedData:  []float64{5, 6},
			expectedFixed: []bool{false, false},
			expectedRes:   true,
		},
		{
			name: "Complex fetch",
			requiredVar: []VariableFieldName{
				VoteCalled,
				NumberOfBallotsCast,
				IslandAllocation,
			},
			variables: map[VariableFieldName]VariableValuePair{
				VoteCalled: {
					VariableName: VoteCalled,
					Values:       []float64{1, 2, 3, 4},
				},
				NumberOfBallotsCast: {
					VariableName: NumberOfBallotsCast,
					Values:       []float64{5},
				},
				IslandAllocation: {
					VariableName: NumberOfBallotsCast,
					Values:       []float64{6, 7, 8},
				},
			},
			expectedData:  []float64{1, 2, 3, 4, 5, 6, 7, 8},
			expectedFixed: []bool{false, false, false, false, false, true, true, true},
			expectedRes:   true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dat, fix, res := fetchRequiredVariables(tc.requiredVar, tc.variables)
			if !reflect.DeepEqual(res, tc.expectedRes) {
				t.Errorf("Expected %v got %v", tc.expectedRes, res)
			}
			if !reflect.DeepEqual(dat, tc.expectedData) {
				t.Errorf("Expected %v got %v", tc.expectedData, dat)
			}
			if !reflect.DeepEqual(fix, tc.expectedFixed) {
				t.Errorf("Expected %v got %v", tc.expectedFixed, fix)
			}
		})
	}
}

func TestRecommendCorrection(t *testing.T) {
	cases := []struct {
		name        string
		prev        float64
		res         float64
		aux         float64
		multiplier  float64
		expectedRes float64
	}{
		{
			name:        "Basic correction",
			prev:        1,
			res:         4,
			aux:         0,
			multiplier:  1,
			expectedRes: -3,
		},
		{
			name:        "Type 1 correction",
			prev:        1,
			res:         -3,
			aux:         1,
			multiplier:  2,
			expectedRes: 5.5,
		},
		{
			name:        "Type 2 correction",
			prev:        0,
			res:         -3,
			aux:         2,
			multiplier:  2,
			expectedRes: 1.5,
		},
		{
			name:        "Type 3 correction",
			prev:        1,
			res:         -3,
			aux:         3,
			multiplier:  2,
			expectedRes: -2,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			newVal := recommendCorrection(tc.prev, tc.res, tc.aux, tc.multiplier)
			if !reflect.DeepEqual(tc.expectedRes, newVal) {
				t.Errorf("Expected %v got %v", tc.expectedRes, newVal)
			}
		})
	}
}

func TestSatisfy(t *testing.T) {
	cases := []struct {
		name        string
		res         float64
		aux         float64
		expectedRes bool
	}{
		{
			name:        "Basic correction",
			res:         4,
			aux:         0,
			expectedRes: false,
		},
		{
			name:        "Type 1 correction",
			res:         3,
			aux:         1,
			expectedRes: true,
		},
		{
			name:        "Type 2 correction",
			res:         0,
			aux:         2,
			expectedRes: true,
		},
		{
			name:        "Type 3 correction",
			res:         0,
			aux:         3,
			expectedRes: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			newVal := satisfy(tc.res, tc.aux)
			if !reflect.DeepEqual(tc.expectedRes, newVal) {
				t.Errorf("Expected %v got %v", tc.expectedRes, newVal)
			}
		})
	}
}

func TestFindFirst(t *testing.T) {
	cases := []struct {
		name          string
		search        bool
		values        []bool
		line          mat.Vector
		expectedIndex int
	}{
		{
			name:          "Can't be found",
			search:        false,
			values:        []bool{true, true},
			line:          mat.NewVecDense(2, []float64{1, 1}),
			expectedIndex: -1,
		},
		{
			name:          "Can be found",
			search:        false,
			values:        []bool{true, true, false},
			line:          mat.NewVecDense(3, []float64{1, 1, 1}),
			expectedIndex: 2,
		},
		{
			name:          "CCan't be found due to multiplier",
			search:        false,
			values:        []bool{true, true, false},
			line:          mat.NewVecDense(3, []float64{1, 1, 0}),
			expectedIndex: -1,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			newVal := findFirst(tc.search, tc.values, tc.line)
			if !reflect.DeepEqual(tc.expectedIndex, newVal) {
				t.Errorf("Expected %v got %v", tc.expectedIndex, newVal)
			}
		})
	}
}

func TestAnalysisEngine(t *testing.T) {
	cases := []struct {
		name        string
		line        mat.Vector
		currentVal  []float64
		fixed       []bool
		auxCode     float64
		expectedRet []float64
		success     bool
	}{
		{
			name:        "Simple adjust",
			line:        mat.NewVecDense(3, []float64{1, 1, 0}),
			currentVal:  []float64{1, 1},
			fixed:       []bool{true, false},
			auxCode:     0,
			expectedRet: []float64{-1, 1},
			success:     true,
		},
		{
			name:        "Complex adjust",
			line:        mat.NewVecDense(4, []float64{0, 1, 1, 0}),
			currentVal:  []float64{1, 1, 1},
			fixed:       []bool{true, false, true},
			auxCode:     0,
			expectedRet: []float64{1, 1, -1},
			success:     true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			newVals, res := analysisEngine(tc.line, tc.currentVal, tc.fixed, tc.auxCode)
			if !reflect.DeepEqual(tc.expectedRet, newVals) {
				t.Errorf("Expected %v got %v", tc.expectedRet, newVals)
			}
			if res != tc.success {
				t.Errorf("Expected %v got %v", tc.success, res)
			}
		})
	}
}

func TestComplianceRecommendation(t *testing.T) {
	cases := []struct {
		name         string
		rule         RuleMatrix
		variables    map[VariableFieldName]VariableValuePair
		newVariables map[VariableFieldName]VariableValuePair
		success      bool
	}{
		{
			name: "Simple adjust",
			rule: RuleMatrix{
				RuleName: "Agrims rule",
				RequiredVariables: []VariableFieldName{
					VoteCalled,
					IslandAllocation,
				},
				ApplicableMatrix: *mat.NewDense(1, 3, []float64{1, 1, 0}),
				AuxiliaryVector:  *mat.NewVecDense(1, []float64{0}),
				Mutable:          false,
				Link: RuleLink{
					Linked: false,
				},
			},
			variables: map[VariableFieldName]VariableValuePair{
				VoteCalled: {
					VariableName: VoteCalled,
					Values:       []float64{1},
				},
				IslandAllocation: {
					VariableName: IslandAllocation,
					Values:       []float64{1},
				},
			},
			newVariables: map[VariableFieldName]VariableValuePair{
				VoteCalled: {
					VariableName: VoteCalled,
					Values:       []float64{1},
				},
				IslandAllocation: {
					VariableName: IslandAllocation,
					Values:       []float64{-1},
				},
			},
			success: true,
		},
		{
			name: "Complex adjust",
			rule: RuleMatrix{
				RuleName: "Some rule name",
				RequiredVariables: []VariableFieldName{
					VoteCalled,
					IslandAllocation,
				},
				ApplicableMatrix: *mat.NewDense(2, 3, []float64{1, 1, 0, 2, 2, 1}),
				AuxiliaryVector:  *mat.NewVecDense(2, []float64{0, 1}),
				Mutable:          false,
				Link: RuleLink{
					Linked: false,
				},
			},
			variables: map[VariableFieldName]VariableValuePair{
				VoteCalled: {
					VariableName: VoteCalled,
					Values:       []float64{1},
				},
				IslandAllocation: {
					VariableName: IslandAllocation,
					Values:       []float64{1},
				},
			},
			newVariables: map[VariableFieldName]VariableValuePair{
				VoteCalled: {
					VariableName: VoteCalled,
					Values:       []float64{1},
				},
				IslandAllocation: {
					VariableName: IslandAllocation,
					Values:       []float64{-1},
				},
			},
			success: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			newVals, res := ComplianceRecommendation(tc.rule, tc.variables)
			if !reflect.DeepEqual(tc.newVariables, newVals) {
				t.Errorf("Expected %v got %v", tc.newVariables, newVals)
			}
			if res != tc.success {
				t.Errorf("Expected %v got %v", tc.success, res)
			}
		})
	}
}

func TestComplianceCheck(t *testing.T) {
	cases := []struct {
		name      string
		rule      RuleMatrix
		variables map[VariableFieldName]VariableValuePair
		success   bool
		err       error
	}{
		{
			name: "Simple Compliance check",
			rule: RuleMatrix{
				RuleName: "Some rule name",
				RequiredVariables: []VariableFieldName{
					VoteCalled,
					IslandAllocation,
				},
				ApplicableMatrix: *mat.NewDense(1, 3, []float64{1, 1, 0}),
				AuxiliaryVector:  *mat.NewVecDense(1, []float64{0}),
				Mutable:          false,
				Link: RuleLink{
					Linked: false,
				},
			},
			variables: map[VariableFieldName]VariableValuePair{
				VoteCalled: {
					VariableName: VoteCalled,
					Values:       []float64{1},
				},
				IslandAllocation: {
					VariableName: IslandAllocation,
					Values:       []float64{1},
				},
			},
			success: false,
			err:     nil,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, inPlay := InitialRuleRegistration(true)
			res, err := ComplianceCheck(tc.rule, tc.variables, inPlay)
			testutils.CompareTestErrors(tc.err, err, t)
			if res != tc.success {
				t.Errorf("Expected %v got %v", tc.success, res)
			}
		})
	}
}

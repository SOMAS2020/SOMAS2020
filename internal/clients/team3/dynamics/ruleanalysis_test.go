package dynamics

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/pkg/testutils"
	"gonum.org/v1/gonum/mat"
	"reflect"
	"testing"
)

func TestDropAllInputStructs(t *testing.T) {
	base, med, hard := generateInputMaps()
	cases := []struct {
		name     string
		input    map[rules.VariableFieldName]Input
		expected map[rules.VariableFieldName][]float64
	}{
		{
			name:  "Basic input unpack",
			input: base,
			expected: map[rules.VariableFieldName][]float64{
				rules.NumberOfBallotsCast: {5},
			},
		},
		{
			name:  "Medium input unpack",
			input: med,
			expected: map[rules.VariableFieldName][]float64{
				rules.NumberOfBallotsCast:  {5},
				rules.NumberOfIslandsAlive: {7},
				rules.IslandsAlive:         {1, 2, 3, 4, 5, 6},
			},
		},
		{
			name:  "Hard input unpack",
			input: hard,
			expected: map[rules.VariableFieldName][]float64{
				rules.NumberOfBallotsCast:    {5},
				rules.NumberOfIslandsAlive:   {7},
				rules.IslandsAlive:           {1, 2, 3, 4, 5, 6},
				rules.AllocationRequestsMade: {3, 5},
				rules.AllocationMade:         {1},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := DropAllInputStructs(tc.input)
			if !reflect.DeepEqual(res, tc.expected) {
				t.Errorf("Expected %v got %v", tc.expected, res)
			}
		})
	}
}

func TestEvaluateSingle(t *testing.T) {
	cases := []struct {
		name       string
		rm         rules.RuleMatrix
		values     map[rules.VariableFieldName][]float64
		outputVect mat.VecDense
		success    bool
	}{
		{
			name: "Basic Input evaluation",
			rm: rules.RuleMatrix{
				RuleName: "Some rule name",
				RequiredVariables: []rules.VariableFieldName{
					rules.AllocationMade,
					rules.AllocationRequestsMade,
				},
				ApplicableMatrix: *mat.NewDense(1, 3, []float64{1, -1, 0}),
				AuxiliaryVector:  *mat.NewVecDense(1, []float64{0}),
				Mutable:          false,
			},
			values: map[rules.VariableFieldName][]float64{
				rules.AllocationMade:         {5},
				rules.AllocationRequestsMade: {5},
			},
			outputVect: *mat.NewVecDense(1, []float64{0}),
			success:    true,
		},
		{
			name: "More complex evaluation",
			rm: rules.RuleMatrix{
				RuleName: "Some rule name",
				RequiredVariables: []rules.VariableFieldName{
					rules.IslandsAlive,
				},
				ApplicableMatrix: *mat.NewDense(1, 5, []float64{2, -1, 0, 0, 0}),
				AuxiliaryVector:  *mat.NewVecDense(1, []float64{0}),
				Mutable:          false,
			},
			values: map[rules.VariableFieldName][]float64{
				rules.IslandsAlive: {1, 2, 3, 4},
			},
			outputVect: *mat.NewVecDense(1, []float64{0}),
			success:    true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, succ := evaluateSingle(tc.rm, tc.values)
			if !reflect.DeepEqual(res, tc.outputVect) {
				t.Errorf("Expected %v got %v", tc.outputVect, res)
			}
			if succ != tc.success {
				t.Errorf("Expected %v got %v", succ, tc.success)
			}
		})
	}
}

func TestIdentifyDeficiencies(t *testing.T) {
	cases := []struct {
		name        string
		inputVec    mat.VecDense
		aux         mat.VecDense
		expectedDef []int
		expectedErr error
	}{
		{
			name:        "Basic deficiencies",
			inputVec:    *mat.NewVecDense(3, []float64{1, 2, 3}),
			aux:         *mat.NewVecDense(3, []float64{1, 1, 1}),
			expectedDef: []int{},
			expectedErr: nil,
		},
		{
			name:        "Complex deficiencies",
			inputVec:    *mat.NewVecDense(5, []float64{0, 1, 1, 1, 1}),
			aux:         *mat.NewVecDense(5, []float64{1, 1, 1, 1, 1}),
			expectedDef: []int{0},
			expectedErr: nil,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := IdentifyDeficiencies(tc.inputVec, tc.aux)
			if !reflect.DeepEqual(res, tc.expectedDef) {
				t.Errorf("Expected %v got %v", tc.expectedDef, res)
			}
			testutils.CompareTestErrors(tc.expectedErr, err, t)
		})
	}
}

func TestBuildSelectedDynamics(t *testing.T) {
	cases := []struct {
		name        string
		rm          rules.RuleMatrix
		aux         mat.VecDense
		deficients  []int
		expectedDyn []dynamic
	}{
		{
			name: "Basic deficiencies",
			rm: rules.RuleMatrix{
				RuleName: "Some rule",
				RequiredVariables: []rules.VariableFieldName{
					rules.NumberOfBallotsCast,
					rules.AllocationMade,
				},
				ApplicableMatrix: *mat.NewDense(1, 3, []float64{1, -1, 0}),
				AuxiliaryVector:  *mat.NewVecDense(1, []float64{0}),
				Mutable:          false,
			},
			aux:         *mat.NewVecDense(1, []float64{0}),
			deficients:  []int{},
			expectedDyn: []dynamic{},
		},
		{
			name: "Complex deficiencies",
			rm: rules.RuleMatrix{
				RuleName: "Some rule",
				RequiredVariables: []rules.VariableFieldName{
					rules.NumberOfBallotsCast,
					rules.AllocationMade,
				},
				ApplicableMatrix: *mat.NewDense(1, 3, []float64{1, -1, 0}),
				AuxiliaryVector:  *mat.NewVecDense(1, []float64{0}),
				Mutable:          false,
			},
			aux:        *mat.NewVecDense(1, []float64{0}),
			deficients: []int{0},
			expectedDyn: []dynamic{
				{
					w:   *mat.NewVecDense(2, []float64{1, -1}),
					b:   0.0,
					aux: 0,
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := BuildSelectedDynamics(tc.rm, tc.aux, tc.deficients)
			if !reflect.DeepEqual(res, tc.expectedDyn) {
				t.Errorf("Expected %v got %v", tc.expectedDyn, res)
			}
		})
	}
}

func TestFetchImmutableInputs(t *testing.T) {
	cases := []struct {
		name         string
		namedInputs  map[rules.VariableFieldName]Input
		expectedVars []rules.VariableFieldName
	}{
		{
			name: "Basic Immutable check",
			namedInputs: map[rules.VariableFieldName]Input{
				rules.AllocationMade: {
					ClientAdjustable: false,
					Name:             rules.AllocationMade,
					Value:            []float64{5},
				},
			},
			expectedVars: []rules.VariableFieldName{rules.AllocationMade},
		},
		{
			name: "Complex Immutable check",
			namedInputs: map[rules.VariableFieldName]Input{
				rules.AllocationMade: {
					ClientAdjustable: false,
					Name:             rules.AllocationMade,
					Value:            []float64{5},
				},
				rules.IslandsAlive: {
					ClientAdjustable: true,
					Name:             rules.IslandsAlive,
					Value:            []float64{3, 4},
				},
				rules.NumberOfIslandsContributingToCommonPool: {
					ClientAdjustable: false,
					Name:             rules.NumberOfIslandsContributingToCommonPool,
					Value:            []float64{4},
				},
			},
			expectedVars: []rules.VariableFieldName{rules.AllocationMade, rules.NumberOfIslandsContributingToCommonPool},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := fetchImmutableInputs(tc.namedInputs)
			if !checkVariableFieldNameListsAreEqual(res, tc.expectedVars) {
				t.Errorf("Expected %v got %v", tc.expectedVars, res)
			}
		})
	}
}

func TestCollapseInputMap(t *testing.T) {
	cases := []struct {
		name           string
		namedInputs    map[rules.VariableFieldName]Input
		expectedInputs []Input
	}{
		{
			name: "Basic collapse",
			namedInputs: map[rules.VariableFieldName]Input{
				rules.AllocationMade: {
					ClientAdjustable: false,
					Name:             rules.AllocationMade,
					Value:            []float64{5},
				},
			},
			expectedInputs: []Input{
				{
					ClientAdjustable: false,
					Name:             rules.AllocationMade,
					Value:            []float64{5},
				},
			},
		},
		{
			name: "Complex collapse",
			namedInputs: map[rules.VariableFieldName]Input{
				rules.AllocationMade: {
					ClientAdjustable: false,
					Name:             rules.AllocationMade,
					Value:            []float64{5},
				},
				rules.IslandsAlive: {
					ClientAdjustable: true,
					Name:             rules.IslandsAlive,
					Value:            []float64{3, 4},
				},
				rules.NumberOfIslandsContributingToCommonPool: {
					ClientAdjustable: false,
					Name:             rules.NumberOfIslandsContributingToCommonPool,
					Value:            []float64{4},
				},
			},
			expectedInputs: []Input{
				{
					ClientAdjustable: false,
					Name:             rules.AllocationMade,
					Value:            []float64{5},
				},
				{
					ClientAdjustable: true,
					Name:             rules.IslandsAlive,
					Value:            []float64{3, 4},
				},
				{
					ClientAdjustable: false,
					Name:             rules.NumberOfIslandsContributingToCommonPool,
					Value:            []float64{4},
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := collapseInputsMap(tc.namedInputs)
			if !checkInputListsAreEqual(res, tc.expectedInputs) {
				t.Errorf("Expected %v got %v", tc.expectedInputs, res)
			}
		})
	}
}

func TestCalculateFullSpaceSize(t *testing.T) {
	cases := []struct {
		name         string
		inputs       []Input
		expectedSize int
	}{
		{
			name: "Basic space sizing",
			inputs: []Input{
				{
					ClientAdjustable: true,
					Name:             rules.IslandsAlive,
					Value:            []float64{1, 2, 3, 4, 5, 6},
				},
			},
			expectedSize: 6,
		},
		{
			name: "Complex space sizing",
			inputs: []Input{
				{
					ClientAdjustable: true,
					Name:             rules.IslandsAlive,
					Value:            []float64{1, 2, 3, 4, 5, 6},
				},
				{
					ClientAdjustable: true,
					Name:             rules.NumberOfIslandsContributingToCommonPool,
					Value:            []float64{1, 2, 3, 4, 5, 6},
				},
				{
					ClientAdjustable: true,
					Name:             rules.AllocationMade,
					Value:            []float64{1, 2, 3, 4, 5, 6},
				},
			},
			expectedSize: 18,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := calculateFullSpaceSize(tc.inputs)
			if !reflect.DeepEqual(res, tc.expectedSize) {
				t.Errorf("Expected %v got %v", tc.expectedSize, res)
			}
		})
	}
}

func TestSelectInputs(t *testing.T) {
	cases := []struct {
		name           string
		immutables     []rules.VariableFieldName
		namedInputs    map[rules.VariableFieldName]Input
		expectedInputs []Input
	}{
		{
			name: "Basic select inputs",
			immutables: []rules.VariableFieldName{
				rules.AllocationMade,
			},
			namedInputs: map[rules.VariableFieldName]Input{
				rules.AllocationMade: {
					ClientAdjustable: false,
					Name:             rules.AllocationMade,
					Value:            []float64{5},
				},
			},
			expectedInputs: []Input{
				{
					ClientAdjustable: false,
					Name:             rules.AllocationMade,
					Value:            []float64{5},
				},
			},
		},
		{
			name: "Complex select inputs",
			immutables: []rules.VariableFieldName{
				rules.AllocationMade,
			},
			namedInputs: map[rules.VariableFieldName]Input{
				rules.AllocationMade: {
					ClientAdjustable: false,
					Name:             rules.AllocationMade,
					Value:            []float64{5},
				},
				rules.IslandsAlive: {
					ClientAdjustable: false,
					Name:             rules.IslandsAlive,
					Value:            []float64{5},
				},
			},
			expectedInputs: []Input{
				{
					ClientAdjustable: false,
					Name:             rules.AllocationMade,
					Value:            []float64{5},
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := selectInputs(tc.namedInputs, tc.immutables)
			if !reflect.DeepEqual(res, tc.expectedInputs) {
				t.Errorf("Expected %v got %v", tc.expectedInputs, res)
			}
		})
	}
}

func TestConstructAllImmutableDynamics(t *testing.T) {
	cases := []struct {
		name             string
		immutables       []Input
		ruleMat          rules.RuleMatrix
		fullSize         int
		expectedDynamics []dynamic
		expectedSuccess  bool
	}{
		{
			name: "Basic construct",
			immutables: []Input{
				{
					Name:             rules.IslandsAlive,
					ClientAdjustable: false,
					Value:            []float64{5},
				},
			},
			ruleMat: rules.RuleMatrix{
				RuleName: "Some Rule",
				RequiredVariables: []rules.VariableFieldName{
					rules.IslandsAlive,
				},
				ApplicableMatrix: mat.Dense{},
				AuxiliaryVector:  mat.VecDense{},
				Mutable:          false,
			},
			fullSize: 1,
			expectedDynamics: []dynamic{
				{
					w:   *mat.NewVecDense(1, []float64{1}),
					b:   -5.0,
					aux: 0,
				},
			},
			expectedSuccess: true,
		},
		{
			name: "Complex construct",
			immutables: []Input{
				{
					Name:             rules.IslandsAlive,
					ClientAdjustable: false,
					Value:            []float64{5},
				},
				{
					Name:             rules.NumberOfBallotsCast,
					ClientAdjustable: false,
					Value:            []float64{1, 2, 3, 4},
				},
			},
			ruleMat: rules.RuleMatrix{
				RuleName: "Some Rule",
				RequiredVariables: []rules.VariableFieldName{
					rules.IslandsAlive,
					rules.NumberOfBallotsCast,
				},
				ApplicableMatrix: mat.Dense{},
				AuxiliaryVector:  mat.VecDense{},
				Mutable:          false,
			},
			fullSize: 5,
			expectedDynamics: []dynamic{
				{
					w:   *mat.NewVecDense(5, []float64{1, 0, 0, 0, 0}),
					b:   -5.0,
					aux: 0,
				},
				{
					w:   *mat.NewVecDense(5, []float64{0, 1, 1, 1, 1}),
					b:   -10.0,
					aux: 0,
				},
			},
			expectedSuccess: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, succ := constructAllImmutableDynamics(tc.immutables, tc.ruleMat, tc.fullSize)
			if !checkDynamicListsAreEqual(res, tc.expectedDynamics) {
				t.Errorf("Expected %v got %v", tc.expectedDynamics, res)
			}
			if succ != tc.expectedSuccess {
				t.Errorf("Expected %v got %v", tc.expectedSuccess, succ)
			}
		})
	}
}

func TestFindClosestApproachInSubspace(t *testing.T) {
	cases := []struct {
		name     string
		ruleMat  rules.RuleMatrix
		dynamics []dynamic
		location mat.VecDense
		result   mat.VecDense
	}{
		{
			name: "Basic no dynamic test",
			ruleMat: rules.RuleMatrix{
				RuleName: "Some rule",
				RequiredVariables: []rules.VariableFieldName{
					rules.NumberOfBallotsCast,
				},
				ApplicableMatrix: *mat.NewDense(1, 2, []float64{1, -6}),
				AuxiliaryVector:  *mat.NewVecDense(1, []float64{0}),
				Mutable:          false,
			},
			dynamics: []dynamic{},
			location: *mat.NewVecDense(1, []float64{6}),
			result:   *mat.NewVecDense(1, []float64{6}),
		},
		{
			name: "Basic single dynamic test",
			ruleMat: rules.RuleMatrix{
				RuleName: "Some rule",
				RequiredVariables: []rules.VariableFieldName{
					rules.NumberOfBallotsCast,
				},
				ApplicableMatrix: *mat.NewDense(1, 2, []float64{1, -6}),
				AuxiliaryVector:  *mat.NewVecDense(1, []float64{0}),
				Mutable:          false,
			},
			dynamics: []dynamic{
				{
					w:   *mat.NewVecDense(1, []float64{1}),
					b:   -6.0,
					aux: 0,
				},
			},
			location: *mat.NewVecDense(1, []float64{5}),
			result:   *mat.NewVecDense(1, []float64{6}),
		},
		{
			name: "Basic multi dynamic test",
			ruleMat: rules.RuleMatrix{
				RuleName: "Some rule",
				RequiredVariables: []rules.VariableFieldName{
					rules.NumberOfBallotsCast,
					rules.NumberOfIslandsAlive,
				},
				ApplicableMatrix: *mat.NewDense(2, 3, []float64{1, 0, -6, 0, 1, -2}),
				AuxiliaryVector:  *mat.NewVecDense(2, []float64{0, 0}),
				Mutable:          false,
			},
			dynamics: []dynamic{
				{
					w:   *mat.NewVecDense(2, []float64{1, 0}),
					b:   -6.0,
					aux: 0,
				},
				{
					w:   *mat.NewVecDense(2, []float64{0, 1}),
					b:   -2.0,
					aux: 0,
				},
			},
			location: *mat.NewVecDense(2, []float64{5, 1}),
			result:   *mat.NewVecDense(2, []float64{6, 2}),
		},
		{
			name: "Complex multi dynamic test",
			ruleMat: rules.RuleMatrix{
				RuleName: "Some rule",
				RequiredVariables: []rules.VariableFieldName{
					rules.NumberOfBallotsCast,
					rules.NumberOfIslandsAlive,
					rules.IslandsAlive,
				},
				ApplicableMatrix: *mat.NewDense(3, 4, []float64{1, 0, 0, -6, 0, 1, 0, -2, 0, 0, 1, -3}),
				AuxiliaryVector:  *mat.NewVecDense(3, []float64{0, 0, 0}),
				Mutable:          false,
			},
			dynamics: []dynamic{
				{
					w:   *mat.NewVecDense(3, []float64{1, 0, 0}),
					b:   -6.0,
					aux: 0,
				},
				{
					w:   *mat.NewVecDense(3, []float64{0, 1, 0}),
					b:   -2.0,
					aux: 0,
				},
				{
					w:   *mat.NewVecDense(3, []float64{0, 0, 1}),
					b:   -3.0,
					aux: 0,
				},
			},
			location: *mat.NewVecDense(3, []float64{5, 1, 1}),
			result:   *mat.NewVecDense(3, []float64{6, 2, 3}),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := findClosestApproachInSubspace(tc.ruleMat, tc.dynamics, tc.location)
			if !reflect.DeepEqual(res, tc.result) {
				t.Errorf("Expected %v got %v", tc.result, res)
			}
		})
	}
}

func TestShift(t *testing.T) {
	_, medMap, _ := generateInputMaps()
	cases := []struct {
		name        string
		inputMatrix rules.RuleMatrix
		inputMap    map[rules.VariableFieldName]Input
		edited      bool
		outputMat   rules.RuleMatrix
	}{
		{
			name: "Basic non edit test",
			inputMatrix: rules.RuleMatrix{
				RuleName: "Basic rule",
				RequiredVariables: []rules.VariableFieldName{
					rules.NumberOfBallotsCast,
					rules.NumberOfIslandsAlive,
				},
				ApplicableMatrix: *mat.NewDense(2, 3, []float64{1, 1, 1, 1, 1, 1}),
				AuxiliaryVector:  *mat.NewVecDense(2, []float64{1, 1}),
				Mutable:          false,
				Link:             rules.RuleLink{},
			},
			inputMap: medMap,
			outputMat: rules.RuleMatrix{
				RuleName: "Basic rule",
				RequiredVariables: []rules.VariableFieldName{
					rules.NumberOfBallotsCast,
					rules.NumberOfIslandsAlive,
				},
				ApplicableMatrix: *mat.NewDense(2, 3, []float64{1, 1, 1, 1, 1, 1}),
				AuxiliaryVector:  *mat.NewVecDense(2, []float64{1, 1}),
				Mutable:          false,
				Link:             rules.RuleLink{},
			},
			edited: false,
		},
		{
			name: "Basic edit test",
			inputMatrix: rules.RuleMatrix{
				RuleName: "Basic rule",
				RequiredVariables: []rules.VariableFieldName{
					rules.NumberOfBallotsCast,
					rules.NumberOfIslandsAlive,
				},
				ApplicableMatrix: *mat.NewDense(2, 3, []float64{1, 1, 1, 1, 1, 1}),
				AuxiliaryVector:  *mat.NewVecDense(2, []float64{0, 1}),
				Mutable:          false,
				Link:             rules.RuleLink{},
			},
			inputMap: medMap,
			outputMat: rules.RuleMatrix{
				RuleName: "Basic rule",
				RequiredVariables: []rules.VariableFieldName{
					rules.NumberOfBallotsCast,
					rules.NumberOfIslandsAlive,
				},
				ApplicableMatrix: *mat.NewDense(2, 3, []float64{1, 1, -12, 1, 1, 1}),
				AuxiliaryVector:  *mat.NewVecDense(2, []float64{0, 1}),
				Mutable:          false,
				Link:             rules.RuleLink{},
			},
			edited: true,
		},
		{
			name: "Complex edit test",
			inputMatrix: rules.RuleMatrix{
				RuleName: "Basic rule",
				RequiredVariables: []rules.VariableFieldName{
					rules.NumberOfBallotsCast,
					rules.NumberOfIslandsAlive,
				},
				ApplicableMatrix: *mat.NewDense(2, 3, []float64{1, 1, 1, 1, 0, 1}),
				AuxiliaryVector:  *mat.NewVecDense(2, []float64{0, 0}),
				Mutable:          false,
				Link:             rules.RuleLink{},
			},
			inputMap: medMap,
			outputMat: rules.RuleMatrix{
				RuleName: "Basic rule",
				RequiredVariables: []rules.VariableFieldName{
					rules.NumberOfBallotsCast,
					rules.NumberOfIslandsAlive,
				},
				ApplicableMatrix: *mat.NewDense(2, 3, []float64{1, 1, -12, 1, 0, -5}),
				AuxiliaryVector:  *mat.NewVecDense(2, []float64{0, 0}),
				Mutable:          false,
				Link:             rules.RuleLink{},
			},
			edited: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, succ := Shift(tc.inputMatrix, tc.inputMap)
			if !reflect.DeepEqual(res, tc.outputMat) {
				t.Errorf("Expected %v got %v", tc.outputMat, res)
			}
			if succ != tc.edited {
				t.Errorf("Expected %v got %v", succ, tc.edited)
			}
		})
	}
}

// Oh boy
func TestFindClosestApproach(t *testing.T) {
	cases := []struct {
		name         string
		ruleMatrix   rules.RuleMatrix
		namedInputs  map[rules.VariableFieldName]Input
		namedOutputs map[rules.VariableFieldName]Input
	}{
		{
			name: "Basic compliant inputs test",
			ruleMatrix: rules.RuleMatrix{
				RuleName: "Some rule",
				RequiredVariables: []rules.VariableFieldName{
					rules.NumberOfBallotsCast,
					rules.IslandsAlive,
				},
				ApplicableMatrix: *mat.NewDense(1, 3, []float64{1, 1, -2}),
				AuxiliaryVector:  *mat.NewVecDense(1, []float64{0}),
				Mutable:          false,
			},
			namedInputs: map[rules.VariableFieldName]Input{
				rules.NumberOfBallotsCast: {
					ClientAdjustable: false,
					Name:             rules.NumberOfBallotsCast,
					Value:            []float64{1},
				},
				rules.IslandsAlive: {
					ClientAdjustable: true,
					Name:             rules.IslandsAlive,
					Value:            []float64{0},
				},
			},
			namedOutputs: map[rules.VariableFieldName]Input{
				rules.NumberOfBallotsCast: {
					ClientAdjustable: false,
					Name:             rules.NumberOfBallotsCast,
					Value:            []float64{1},
				},
				rules.IslandsAlive: {
					ClientAdjustable: true,
					Name:             rules.IslandsAlive,
					Value:            []float64{1},
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := FindClosestApproach(tc.ruleMatrix, tc.namedInputs)
			if !reflect.DeepEqual(res, tc.namedOutputs) {
				t.Errorf("Expected %v got %v", tc.namedOutputs, res)
			}
		})
	}
}

func generateInputMaps() (map[rules.VariableFieldName]Input, map[rules.VariableFieldName]Input, map[rules.VariableFieldName]Input) {
	baseMap := map[rules.VariableFieldName]Input{}
	baseMap[rules.NumberOfBallotsCast] = Input{
		Name:             rules.NumberOfBallotsCast,
		ClientAdjustable: false,
		Value:            []float64{5},
	}
	mediumMap := copyMap(baseMap)
	mediumMap[rules.NumberOfIslandsAlive] = Input{
		Name:             rules.NumberOfIslandsAlive,
		ClientAdjustable: true,
		Value:            []float64{7},
	}
	mediumMap[rules.IslandsAlive] = Input{
		Name:             rules.IslandsAlive,
		ClientAdjustable: false,
		Value:            []float64{1, 2, 3, 4, 5, 6},
	}
	hardMap := copyMap(mediumMap)
	hardMap[rules.AllocationRequestsMade] = Input{
		Name:             rules.AllocationRequestsMade,
		ClientAdjustable: true,
		Value:            []float64{3, 5},
	}
	hardMap[rules.AllocationMade] = Input{
		Name:             rules.AllocationMade,
		ClientAdjustable: false,
		Value:            []float64{1},
	}
	return baseMap, mediumMap, hardMap
}

func copyMap(inp map[rules.VariableFieldName]Input) map[rules.VariableFieldName]Input {
	ret := make(map[rules.VariableFieldName]Input)
	for key, val := range inp {
		ret[key] = val
	}
	return ret
}

func checkInputListsAreEqual(lst1 []Input, lst2 []Input) bool {
	map1 := make(map[rules.VariableFieldName]int)
	map2 := make(map[rules.VariableFieldName]int)
	for _, inp := range lst1 {
		map1[inp.Name] = 0
	}
	for _, inp := range lst2 {
		map2[inp.Name] = 0
	}
	return reflect.DeepEqual(map1, map2)
}

func checkDynamicListsAreEqual(lst1 []dynamic, lst2 []dynamic) bool {
	map1 := make(map[float64]int)
	map2 := make(map[float64]int)
	for _, inp := range lst1 {
		map1[inp.b] = 0
	}
	for _, inp := range lst2 {
		map2[inp.b] = 0
	}
	return reflect.DeepEqual(map1, map2)
}

func checkVariableFieldNameListsAreEqual(lst1 []rules.VariableFieldName, lst2 []rules.VariableFieldName) bool {
	map1 := make(map[rules.VariableFieldName]int)
	map2 := make(map[rules.VariableFieldName]int)
	for _, inp := range lst1 {
		map1[inp] = 0
	}
	for _, inp := range lst2 {
		map2[inp] = 0
	}
	return reflect.DeepEqual(map1, map2)
}

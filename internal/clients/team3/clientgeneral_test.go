package team3

// General client functions testing

import (
	"math"
	"reflect"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func TestUpdateTrustMapAgg(t *testing.T) {
	cases := []struct {
		name        string
		ourClient   client
		clientID    shared.ClientID
		amount      float64
		expectedVal map[shared.ClientID][]float64
	}{
		{
			name: "Basic test",
			ourClient: client{
				trustMapAgg: map[shared.ClientID][]float64{
					0: []float64{},
					1: []float64{},
					3: []float64{},
					4: []float64{},
					5: []float64{},
				},
			},
			clientID: 1,
			amount:   10.34,
			expectedVal: map[shared.ClientID][]float64{
				0: []float64{},
				1: []float64{10.34},
				3: []float64{},
				4: []float64{},
				5: []float64{},
			},
		},
		{
			name: "Basic Test 1",
			ourClient: client{
				trustMapAgg: map[shared.ClientID][]float64{
					0: []float64{5.92},
					1: []float64{62.78},
					3: []float64{17.62},
					4: []float64{-10.3},
					5: []float64{6.42},
				},
			},
			clientID: 4,
			amount:   -9.56,
			expectedVal: map[shared.ClientID][]float64{
				0: []float64{5.92},
				1: []float64{62.78},
				3: []float64{17.62},
				4: []float64{-10.3, -9.56},
				5: []float64{6.42},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.ourClient.updatetrustMapAgg(tc.clientID, tc.amount)
			if !reflect.DeepEqual(tc.ourClient.trustMapAgg, tc.expectedVal) {
				t.Errorf("Expected final transgressions to be %v got %v", tc.expectedVal, tc.ourClient.trustMapAgg)
			}
		})
	}
}

func TestInitTrustMapAgg(t *testing.T) {
	cases := []struct {
		name        string
		ourClient   client
		expectedVal map[shared.ClientID][]float64
	}{
		{
			name: "Basic test",
			ourClient: client{
				trustMapAgg: map[shared.ClientID][]float64{
					0: []float64{5.92},
					1: []float64{62.78},
					3: []float64{17.62},
					4: []float64{-10.3},
					5: []float64{6.42},
				},
			},
			expectedVal: map[shared.ClientID][]float64{
				0: []float64{},
				1: []float64{},
				3: []float64{},
				4: []float64{},
				5: []float64{},
			},
		},
		{
			name: "Complex test",
			ourClient: client{
				trustMapAgg: map[shared.ClientID][]float64{
					0: []float64{5.92, 8.97, 19.23},
					1: []float64{62.78, 55.89, -10.65, -20.76},
					3: []float64{17.62, 5.64, -15.67, 45.86, -99.80},
					4: []float64{-10.3, 6.58, 3.74, -65.78, -78.98, 34.56},
					5: []float64{6.42, 69.69, 98.87, -60.7857, 99.9999, 0.00001, 0.05},
				},
			},
			expectedVal: map[shared.ClientID][]float64{
				0: []float64{},
				1: []float64{},
				3: []float64{},
				4: []float64{},
				5: []float64{},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.ourClient.inittrustMapAgg()
			if !reflect.DeepEqual(tc.ourClient.trustMapAgg, tc.expectedVal) {
				t.Errorf("Expected final transgressions to be %v got %v", tc.expectedVal, tc.ourClient.trustMapAgg)
			}
		})
	}
}

func TestUpdateTrustScore(t *testing.T) {
	cases := []struct {
		name        string
		ourClient   client
		trustMapAgg map[shared.ClientID][]float64
		expectedVal map[shared.ClientID]float64
	}{
		{
			name: "Basic test",
			ourClient: client{
				trustScore: map[shared.ClientID]float64{
					0: 50.0,
					1: 50.0,
					3: 50.0,
					4: 50.0,
					5: 50.0,
				},
			},
			trustMapAgg: map[shared.ClientID][]float64{
				0: []float64{5.92},
				1: []float64{62.78},
				3: []float64{17.62},
				4: []float64{-50.3},
				5: []float64{6.42},
			},
			expectedVal: map[shared.ClientID]float64{
				0: 55.92,
				1: 100,
				3: 67.62,
				4: 0,
				5: 56.42,
			},
		},
		{
			name: "Complex test",
			ourClient: client{
				trustScore: map[shared.ClientID]float64{
					0: 55.87,
					1: 88.98,
					3: 23.45,
					4: 5.05,
					5: 69.69,
				},
			},
			trustMapAgg: map[shared.ClientID][]float64{
				0: []float64{5.92, 8.97, 10.75},
				1: []float64{12.78, -13.45, 23.45},
				3: []float64{17.62, 25.62},
				4: []float64{-86.56, 43.43, 48.99},
				5: []float64{6.42, 0.001, -5.96, -45.45},
			},
			expectedVal: map[shared.ClientID]float64{
				0: 64.41666666666666,
				1: 96.57333333333334,
				3: 45.07,
				4: 7.003333333333333,
				5: 58.44275,
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.ourClient.updateTrustScore(tc.trustMapAgg)
			if !reflect.DeepEqual(tc.ourClient.trustScore, tc.expectedVal) {
				t.Errorf("Expected final transgressions to be %v got %v", tc.expectedVal, tc.ourClient.trustScore)
			}
		})
	}
}

func TestUpdateCriticalThreshold(t *testing.T) {
	cases := []struct {
		name              string
		ourClient         client
		isInCriticalState bool
		estimatedResource shared.Resources
		expected          criticalStatePrediction
	}{
		{
			name: "in Critical Test",
			ourClient: client{
				criticalStatePrediction: criticalStatePrediction{upperBound: 70, lowerBound: 30}},
			isInCriticalState: true,
			estimatedResource: shared.Resources(40),
			expected:          criticalStatePrediction{upperBound: 70, lowerBound: 40},
		},
		{
			name: "Not in Critical Test",
			ourClient: client{
				criticalStatePrediction: criticalStatePrediction{upperBound: 70, lowerBound: 30}},
			isInCriticalState: false,
			estimatedResource: shared.Resources(60),
			expected:          criticalStatePrediction{upperBound: 60, lowerBound: 30},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.ourClient.updateCriticalThreshold(tc.isInCriticalState, tc.estimatedResource)
			ans := tc.ourClient.criticalStatePrediction
			if ans != tc.expected {
				t.Errorf("got %f-%f, want %f-%f", ans.lowerBound, ans.upperBound, tc.expected.lowerBound, tc.expected.upperBound)
			}
		})
	}
}

func TestUpdateCompliance(t *testing.T) {
	cases := []struct {
		name        string
		ourClient   client
		expectedVal float64
	}{
		{
			name: "Just caught!",
			ourClient: client{
				timeSinceCaught: 0,
				numTimeCaught:   100,
				compliance:      0.2,
				params: islandParams{
					recidivism:      1.0,
					complianceLevel: 0.1,
				},
			},
			expectedVal: 1.0,
		},
		{
			name: "Compliance decay - non-compliant agent",
			ourClient: client{
				timeSinceCaught: 1,
				numTimeCaught:   1,
				compliance:      1.0,
				params: islandParams{
					recidivism:      1.0,
					complianceLevel: 0.0,
				},
			},
			expectedVal: math.Exp(-0.5),
		},
		{
			name: "Compliance decay - fully-compliant agent",
			ourClient: client{
				timeSinceCaught: 1,
				numTimeCaught:   10,
				compliance:      1.0,
				params: islandParams{
					recidivism:      1.0,
					complianceLevel: 1.0,
				},
			},
			expectedVal: 1.0,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.ourClient.updateCompliance()
			if tc.ourClient.compliance != tc.expectedVal {
				t.Errorf("Expected final transgressions to be %v got %v", tc.expectedVal, tc.ourClient.compliance)
			}
		})
	}
}

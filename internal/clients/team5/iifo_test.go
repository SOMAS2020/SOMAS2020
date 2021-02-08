package team5

import (
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

var c = initClient()

func TestGenerateForecast(t *testing.T) {
	// can use same spatial and mag info because we're only assessing period
	dInfo := disasterInfo{report: disasters.DisasterReport{X: 0, Y: 0, Magnitude: 1}}
	dh2 := disasterHistory{8: dInfo}
	dh3 := disasterHistory{3: dInfo, 5: dInfo, 7: dInfo, 9: dInfo}

	var tests = []struct {
		name       string
		dh         disasterHistory
		wantPeriod uint
		wantConf   float64
	}{
		{"1 past disaster", dh2, 7, 50},
		{"many periodic disasters", dh3, 2, 100},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			n := len(tc.dh)
			obs := []float64{}
			for _, p := range tc.dh.getPastDisasterPeriods() {
				obs = append(obs, p)
			}
			dModel := disasterModel{period: newKdeModel(obs)}

			t.Logf("%v samples: obs: %v, period pdf %v", n, obs, 2)

			ansPeriod, _ := dModel.period.getStatistics(uint(n))

			if uint(ansPeriod.mean) != tc.wantPeriod {
				t.Errorf("period ans %.4f", ansPeriod)
				t.Errorf("period: got %d, want %d", uint(ansPeriod.mean), tc.wantPeriod)
			}
		})
	}
}

func TestUpdateForecastingReputations(t *testing.T) {
	receivedPreds := shared.ReceivedDisasterPredictionsDict{
		shared.Team1: shared.ReceivedDisasterPredictionInfo{
			PredictionMade: shared.DisasterPrediction{
				Confidence: 60,
			},
			SharedFrom: shared.Team1,
		},
		shared.Team2: shared.ReceivedDisasterPredictionInfo{
			PredictionMade: shared.DisasterPrediction{
				Confidence: 20,
			},
			SharedFrom: shared.Team2,
		},
		shared.Team3: shared.ReceivedDisasterPredictionInfo{
			PredictionMade: shared.DisasterPrediction{
				Confidence: 100,
			},
			SharedFrom: shared.Team3,
		},
	}
	c.opinions = opinionMap{
		shared.Team1: &wrappedOpininon{opinion{forecastReputation: 0.0}},
		shared.Team2: &wrappedOpininon{opinion{forecastReputation: 0.0}},
		shared.Team3: &wrappedOpininon{opinion{forecastReputation: 0.0}},
	}
	c.disasterHistory = disasterHistory{} // no disasters recorded
	c.updateForecastingReputations(receivedPreds)
	if c.opinions[shared.Team1].getForecastingRep() >= 0 {
		t.Error("Received prediction with confidence > 50 percent with no disasters")
	}

	if c.opinions[shared.Team2].getForecastingRep() < 0 {
		t.Error("Expected no negative change to reputation after sensible prediction")
	}
	c.disasterHistory = disasterHistory{1: disasterInfo{}, 5: disasterInfo{}} // no disasters recorded
	c.updateForecastingReputations(receivedPreds)

	if c.opinions[shared.Team3].getForecastingRep() >= 0 {
		t.Error("Received perfectly confident prediction. Expected rep. to decrease.")
	}

}

func TestComputeForecastPerformance(t *testing.T) {
	d1 := disasterInfo{report: disasters.DisasterReport{X: 1.0}}
	d2 := disasterInfo{report: disasters.DisasterReport{X: 0.0}}
	dh := disasterHistory{3: d1, 5: d2}
	fh := forecastHistory{
		1: forecastInfo{epiX: 0, period: 0},
		2: forecastInfo{epiX: 0.5, period: 1},
		3: forecastInfo{epiX: 0.75, period: 2},
		4: forecastInfo{epiX: 0.5, period: 2},
		5: forecastInfo{epiX: 0.5, period: 2},
	}
	periodSamples := [][]float64{{0, 1, 2}, {2, 2}}
	periodTargets := []float64{2, 2}

	xSamples := [][]float64{{0, 0.5, 0.75}, {0.5, 0.5}}
	xTargets := []float64{1, 0}

	conf := createClient(shared.Team5).config
	decay := conf.forecastTemporalDecay

	forecastErrors, err := computeForecastingPerformance(dh, fh, conf)

	if err != nil {
		t.Logf("Error analysing disaster history: %v", err)
	}

	for i := range dh.sortKeys() {
		errorMap := forecastErrors[i]

		expectedXError := univariateWeightedMSE(xSamples[i], xTargets[i], decay)
		expectedPeriodError := univariateWeightedMSE(periodSamples[i], periodTargets[i], decay)

		if errorMap[x] != expectedXError {
			t.Logf("computed x epicentre error was incorrect. Expected %v, got %v", expectedXError, errorMap[x])
		}
		if errorMap[period] != expectedPeriodError {
			t.Logf("computed period error was incorrect. Expected %v, got %v", expectedPeriodError, errorMap[period])
		}
	}
}

func initClient() *client {
	c := createClient(shared.Team5)
	return c
}

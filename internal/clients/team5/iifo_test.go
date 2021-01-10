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
			// if ansConf != tc.wantConf {
			// 	t.Errorf("conf: got %.3f, want %.3f", ansConf, tc.wantConf)
			// }
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

func TestAnalyseDisasterHistory(t *testing.T) {
	di := disasterInfo{}
	fi := forecastInfo{}
	dh := disasterHistory{4: di, 8: di}
	fh := forecastHistory{1: fi, 2: fi, 3: fi, 4: fi, 5: fi, 6: fi, 7: fi, 8: fi}

	analyseDisasterHistory(dh, fh)
	t.Error("dummy error")
}

func initClient() *client {
	c := createClient()
	return c
}

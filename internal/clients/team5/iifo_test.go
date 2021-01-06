package team5

import (
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
)

var c = initClient()

func TestAnalyseDisasterPeriod(t *testing.T) {
	// can use same spatial and mag info because we're only assessing period
	dInfo := disasterInfo{report: disasters.DisasterReport{X: 0, Y: 0, Magnitude: 1}}
	dh1 := disasterHistory{}
	dh2 := disasterHistory{8: dInfo}
	dh3 := disasterHistory{3: dInfo, 5: dInfo, 7: dInfo, 9: dInfo}

	var tests = []struct {
		name       string
		dh         disasterHistory
		wantPeriod uint // output tier
		wantConf   float64
	}{
		{"no past disasters", dh1, 0, 0},
		{"1 past disaster", dh2, 7, 50},
		{"many periodic disasters", dh3, 2, 100},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c.disasterHistory = tc.dh
			ansPeriod, ansConf := c.analyseDisasterPeriod()
			if ansPeriod != tc.wantPeriod {
				t.Errorf("period: got %d, want %d", ansPeriod, tc.wantPeriod)
			}
			if ansConf != tc.wantConf {
				t.Errorf("conf: got %.3f, want %.3f", ansConf, tc.wantConf)
			}
		})
	}
}

func initClient() *client {
	c := createClient()
	return c
}

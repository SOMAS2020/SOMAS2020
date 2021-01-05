package config

import (
	"reflect"
	"testing"
)

// Tests GetClientDisasterConfig too
func TestGetClientConfig(t *testing.T) {
	cases := []struct {
		name   string
		config Config
		want   ClientConfig
	}{
		{
			name: "all visible",
			config: Config{
				CostOfLiving:                1,
				MinimumResourceThreshold:    2,
				MaxCriticalConsecutiveTurns: 3,
				MaxSeasons:                  4, // not visible
				DisasterConfig: DisasterConfig{
					CommonpoolThreshold:        6,
					CommonpoolThresholdVisible: true,
					Period:                     4, // not visible
					PeriodVisible:              true,
					StochasticPeriod:           true,
				},
			},
			want: ClientConfig{
				CostOfLiving:                1,
				MinimumResourceThreshold:    2,
				MaxCriticalConsecutiveTurns: 3,
				DisasterConfig: ClientDisasterConfig{
					CommonpoolThreshold: SelectivelyVisibleResources{
						Value: 6,
						Valid: true,
					},
					DisasterPeriod: SelectivelyVisibleDisasterPeriod{
						Period: 4,
						Valid:  true,
					},
					StochasticDisasters: true,
				},
			},
		},
		{
			name: "all selectively visible invisible",
			config: Config{
				CostOfLiving:                1,
				MinimumResourceThreshold:    2,
				MaxCriticalConsecutiveTurns: 3,
				MaxSeasons:                  4, // not visible
				DisasterConfig: DisasterConfig{
					CommonpoolThreshold:        6,
					CommonpoolThresholdVisible: false,
					Period:                     4, // not visible
				},
			},
			want: ClientConfig{
				CostOfLiving:                1,
				MinimumResourceThreshold:    2,
				MaxCriticalConsecutiveTurns: 3,
				DisasterConfig: ClientDisasterConfig{
					CommonpoolThreshold: SelectivelyVisibleResources{
						Valid: false,
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.config.GetClientConfig()
			if !reflect.DeepEqual(tc.want, got) {
				t.Errorf("want '%v' got '%v'", tc.want, got)
			}
		})
	}
}

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
					StochasticPeriodVisible:    true,
				},
				IIGOConfig: IIGOConfig{
					GetRuleForSpeakerActionCost:        50,
					BroadcastTaxationActionCost:        50,
					ReplyAllocationRequestsActionCost:  50,
					RequestAllocationRequestActionCost: 50,
					RequestRuleProposalActionCost:      50,
					AppointNextSpeakerActionCost:       50,
					InspectHistoryActionCost:           50,
					InspectBallotActionCost:            50,
					InspectAllocationActionCost:        50,
					AppointNextPresidentActionCost:     50,
					SanctionCacheDepth:                 3,
					HistoryCacheDepth:                  3,
					AssumedResourcesNoReport:           500,
					SanctionLength:                     3,
					SetVotingResultActionCost:          50,
					SetRuleToVoteActionCost:            50,
					AnnounceVotingResultActionCost:     50,
					UpdateRulesActionCost:              50,
					AppointNextJudgeActionCost:         50,
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
					DisasterPeriod: SelectivelyVisibleUint{
						Value: 4,
						Valid: true,
					},
					StochasticDisasters: SelectivelyVisibleBool{
						Value: true,
						Valid: true,
					},
				},
				IIGOClientConfig: IIGOConfig{
					GetRuleForSpeakerActionCost:        50,
					BroadcastTaxationActionCost:        50,
					ReplyAllocationRequestsActionCost:  50,
					RequestAllocationRequestActionCost: 50,
					RequestRuleProposalActionCost:      50,
					AppointNextSpeakerActionCost:       50,
					InspectHistoryActionCost:           50,
					InspectBallotActionCost:            50,
					InspectAllocationActionCost:        50,
					AppointNextPresidentActionCost:     50,
					SanctionCacheDepth:                 3,
					HistoryCacheDepth:                  3,
					AssumedResourcesNoReport:           500,
					SanctionLength:                     3,
					SetVotingResultActionCost:          50,
					SetRuleToVoteActionCost:            50,
					AnnounceVotingResultActionCost:     50,
					UpdateRulesActionCost:              50,
					AppointNextJudgeActionCost:         50,
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
				IIGOConfig: IIGOConfig{
					GetRuleForSpeakerActionCost:        50,
					BroadcastTaxationActionCost:        50,
					ReplyAllocationRequestsActionCost:  50,
					RequestAllocationRequestActionCost: 50,
					RequestRuleProposalActionCost:      50,
					AppointNextSpeakerActionCost:       50,
					InspectHistoryActionCost:           50,
					InspectBallotActionCost:            50,
					InspectAllocationActionCost:        50,
					AppointNextPresidentActionCost:     50,
					SanctionCacheDepth:                 3,
					HistoryCacheDepth:                  3,
					AssumedResourcesNoReport:           500,
					SanctionLength:                     3,
					SetVotingResultActionCost:          50,
					SetRuleToVoteActionCost:            50,
					AnnounceVotingResultActionCost:     50,
					UpdateRulesActionCost:              50,
					AppointNextJudgeActionCost:         50,
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
				IIGOClientConfig: IIGOConfig{
					GetRuleForSpeakerActionCost:        50,
					BroadcastTaxationActionCost:        50,
					ReplyAllocationRequestsActionCost:  50,
					RequestAllocationRequestActionCost: 50,
					RequestRuleProposalActionCost:      50,
					AppointNextSpeakerActionCost:       50,
					InspectHistoryActionCost:           50,
					InspectBallotActionCost:            50,
					InspectAllocationActionCost:        50,
					AppointNextPresidentActionCost:     50,
					SanctionCacheDepth:                 3,
					HistoryCacheDepth:                  3,
					AssumedResourcesNoReport:           500,
					SanctionLength:                     3,
					SetVotingResultActionCost:          50,
					SetRuleToVoteActionCost:            50,
					AnnounceVotingResultActionCost:     50,
					UpdateRulesActionCost:              50,
					AppointNextJudgeActionCost:         50,
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

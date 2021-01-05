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
				},
				IIGOClientConfig: ClientIIGOConfig{
					GetRuleForSpeakerActionCost: SelectivelyVisibleResources{
						Value: 50,
						Valid: true,
					},
					BroadcastTaxationActionCost: SelectivelyVisibleResources{
						Value: 50,
						Valid: true,
					},
					ReplyAllocationRequestsActionCost: SelectivelyVisibleResources{
						Value: 50,
						Valid: true,
					},
					RequestAllocationRequestActionCost: SelectivelyVisibleResources{
						Value: 50,
						Valid: true,
					},
					RequestRuleProposalActionCost: SelectivelyVisibleResources{
						Value: 50,
						Valid: true,
					},
					AppointNextSpeakerActionCost: SelectivelyVisibleResources{
						Value: 50,
						Valid: true,
					},
					InspectHistoryActionCost: SelectivelyVisibleResources{
						Value: 50,
						Valid: true,
					},
					InspectBallotActionCost: SelectivelyVisibleResources{
						Value: 50,
						Valid: true,
					},
					InspectAllocationActionCost: SelectivelyVisibleResources{
						Value: 50,
						Valid: true,
					},
					AppointNextPresidentActionCost: SelectivelyVisibleResources{
						Value: 50,
						Valid: true,
					},
					SanctionCacheDepth: SelectivelyVisibleInteger{
						Value: 3,
						Valid: true,
					},
					HistoryCacheDepth: SelectivelyVisibleInteger{
						Value: 3,
						Valid: true,
					},
					AssumedResourcesNoReport: SelectivelyVisibleResources{
						Value: 500,
						Valid: true,
					},
					SanctionLength: SelectivelyVisibleInteger{
						Value: 3,
						Valid: true,
					},
					SetVotingResultActionCost: SelectivelyVisibleResources{
						Value: 50,
						Valid: true,
					},
					SetRuleToVoteActionCost: SelectivelyVisibleResources{
						Value: 50,
						Valid: true,
					},
					AnnounceVotingResultActionCost: SelectivelyVisibleResources{
						Value: 50,
						Valid: true,
					},
					UpdateRulesActionCost: SelectivelyVisibleResources{
						Value: 50,
						Valid: true,
					},
					AppointNextJudgeActionCost: SelectivelyVisibleResources{
						Value: 50,
						Valid: true,
					},
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
				IIGOClientConfig: ClientIIGOConfig{
					GetRuleForSpeakerActionCost: SelectivelyVisibleResources{
						Value: 50,
						Valid: true,
					},
					BroadcastTaxationActionCost: SelectivelyVisibleResources{
						Value: 50,
						Valid: true,
					},
					ReplyAllocationRequestsActionCost: SelectivelyVisibleResources{
						Value: 50,
						Valid: true,
					},
					RequestAllocationRequestActionCost: SelectivelyVisibleResources{
						Value: 50,
						Valid: true,
					},
					RequestRuleProposalActionCost: SelectivelyVisibleResources{
						Value: 50,
						Valid: true,
					},
					AppointNextSpeakerActionCost: SelectivelyVisibleResources{
						Value: 50,
						Valid: true,
					},
					InspectHistoryActionCost: SelectivelyVisibleResources{
						Value: 50,
						Valid: true,
					},
					InspectBallotActionCost: SelectivelyVisibleResources{
						Value: 50,
						Valid: true,
					},
					InspectAllocationActionCost: SelectivelyVisibleResources{
						Value: 50,
						Valid: true,
					},
					AppointNextPresidentActionCost: SelectivelyVisibleResources{
						Value: 50,
						Valid: true,
					},
					SanctionCacheDepth: SelectivelyVisibleInteger{
						Value: 3,
						Valid: true,
					},
					HistoryCacheDepth: SelectivelyVisibleInteger{
						Value: 3,
						Valid: true,
					},
					AssumedResourcesNoReport: SelectivelyVisibleResources{
						Value: 500,
						Valid: true,
					},
					SanctionLength: SelectivelyVisibleInteger{
						Value: 3,
						Valid: true,
					},
					SetVotingResultActionCost: SelectivelyVisibleResources{
						Value: 50,
						Valid: true,
					},
					SetRuleToVoteActionCost: SelectivelyVisibleResources{
						Value: 50,
						Valid: true,
					},
					AnnounceVotingResultActionCost: SelectivelyVisibleResources{
						Value: 50,
						Valid: true,
					},
					UpdateRulesActionCost: SelectivelyVisibleResources{
						Value: 50,
						Valid: true,
					},
					AppointNextJudgeActionCost: SelectivelyVisibleResources{
						Value: 50,
						Valid: true,
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

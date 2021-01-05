package config

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

// ClientConfig contains config information visible to clients.
type ClientConfig struct {
	CostOfLiving                shared.Resources
	MinimumResourceThreshold    shared.Resources
	MaxCriticalConsecutiveTurns uint
	DisasterConfig              ClientDisasterConfig
	IIGOClientConfig            ClientIIGOConfig
}

// ClientIIGOConfig contains iigo config fields that is visible to clients
type ClientIIGOConfig struct {
	// Executive branch
	GetRuleForSpeakerActionCost        shared.Resources
	BroadcastTaxationActionCost        shared.Resources
	ReplyAllocationRequestsActionCost  shared.Resources
	RequestAllocationRequestActionCost shared.Resources
	RequestRuleProposalActionCost      shared.Resources
	AppointNextSpeakerActionCost       shared.Resources
	// Judiciary branch
	InspectHistoryActionCost       shared.Resources
	InspectBallotActionCost        shared.Resources
	InspectAllocationActionCost    shared.Resources
	AppointNextPresidentActionCost shared.Resources
	SanctionCacheDepth             int
	HistoryCacheDepth              int
	AssumedResourcesNoReport       shared.Resources
	SanctionLength                 int
	// Legislative branch
	SetVotingResultActionCost      shared.Resources
	SetRuleToVoteActionCost        shared.Resources
	AnnounceVotingResultActionCost shared.Resources
	UpdateRulesActionCost          shared.Resources
	AppointNextJudgeActionCost     shared.Resources
}

// ClientDisasterConfig contains disaster config information visible to clients.
type ClientDisasterConfig struct {
	CommonpoolThreshold SelectivelyVisibleResources
}

// GetClientConfig gets ClientConfig.
func (c Config) GetClientConfig() ClientConfig {
	return ClientConfig{
		CostOfLiving:                c.CostOfLiving,
		MinimumResourceThreshold:    c.MinimumResourceThreshold,
		MaxCriticalConsecutiveTurns: c.MaxCriticalConsecutiveTurns,
		DisasterConfig:              c.DisasterConfig.GetClientDisasterConfig(),
		IIGOClientConfig:            c.IIGOConfig.GetClientIIGOConfig(),
	}
}

// GetClientDisasterConfig gets ClientDisasterConfig
func (c DisasterConfig) GetClientDisasterConfig() ClientDisasterConfig {
	return ClientDisasterConfig{
		CommonpoolThreshold: getSelectivelyVisibleResources(
			c.CommonpoolThreshold,
			c.CommonpoolThresholdVisible,
		),
	}
}

// GetClientIIGOConfig gets ClientIIGOConfig
func (c IIGOConfig) GetClientIIGOConfig() ClientIIGOConfig {
	return ClientIIGOConfig{
		// Executive branch
		GetRuleForSpeakerActionCost:        c.GetRuleForSpeakerActionCost,
		BroadcastTaxationActionCost:        c.BroadcastTaxationActionCost,
		ReplyAllocationRequestsActionCost:  c.ReplyAllocationRequestsActionCost,
		RequestAllocationRequestActionCost: c.RequestAllocationRequestActionCost,
		RequestRuleProposalActionCost:      c.RequestRuleProposalActionCost,
		AppointNextSpeakerActionCost:       c.AppointNextSpeakerActionCost,
		// Judiciary branch
		InspectHistoryActionCost:       c.InspectHistoryActionCost,
		InspectBallotActionCost:        c.InspectBallotActionCost,
		InspectAllocationActionCost:    c.InspectAllocationActionCost,
		AppointNextPresidentActionCost: c.AppointNextPresidentActionCost,
		SanctionCacheDepth:             c.SanctionCacheDepth,
		HistoryCacheDepth:              c.HistoryCacheDepth,
		AssumedResourcesNoReport:       c.AssumedResourcesNoReport,
		SanctionLength:                 c.SanctionLength,
		// Legislative branch
		SetVotingResultActionCost:      c.SetVotingResultActionCost,
		SetRuleToVoteActionCost:        c.SetRuleToVoteActionCost,
		AnnounceVotingResultActionCost: c.AnnounceVotingResultActionCost,
		UpdateRulesActionCost:          c.UpdateRulesActionCost,
		AppointNextJudgeActionCost:     c.AppointNextJudgeActionCost,
	}
}

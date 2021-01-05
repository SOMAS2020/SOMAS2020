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
	GetRuleForSpeakerActionCost        SelectivelyVisibleResources
	BroadcastTaxationActionCost        SelectivelyVisibleResources
	ReplyAllocationRequestsActionCost  SelectivelyVisibleResources
	RequestAllocationRequestActionCost SelectivelyVisibleResources
	RequestRuleProposalActionCost      SelectivelyVisibleResources
	AppointNextSpeakerActionCost       SelectivelyVisibleResources
	// Judiciary branch
	InspectHistoryActionCost       SelectivelyVisibleResources
	InspectBallotActionCost        SelectivelyVisibleResources
	InspectAllocationActionCost    SelectivelyVisibleResources
	AppointNextPresidentActionCost SelectivelyVisibleResources
	SanctionCacheDepth             SelectivelyVisibleInteger
	HistoryCacheDepth              SelectivelyVisibleInteger
	AssumedResourcesNoReport       SelectivelyVisibleResources
	SanctionLength                 SelectivelyVisibleInteger
	// Legislative branch
	SetVotingResultActionCost      SelectivelyVisibleResources
	SetRuleToVoteActionCost        SelectivelyVisibleResources
	AnnounceVotingResultActionCost SelectivelyVisibleResources
	UpdateRulesActionCost          SelectivelyVisibleResources
	AppointNextJudgeActionCost     SelectivelyVisibleResources
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
		GetRuleForSpeakerActionCost: getSelectivelyVisibleResources(
			c.GetRuleForSpeakerActionCost,
			true,
		),
		BroadcastTaxationActionCost: getSelectivelyVisibleResources(
			c.BroadcastTaxationActionCost,
			true,
		),
		ReplyAllocationRequestsActionCost: getSelectivelyVisibleResources(
			c.ReplyAllocationRequestsActionCost,
			true,
		),
		RequestAllocationRequestActionCost: getSelectivelyVisibleResources(
			c.RequestAllocationRequestActionCost,
			true,
		),
		RequestRuleProposalActionCost: getSelectivelyVisibleResources(
			c.RequestRuleProposalActionCost,
			true,
		),
		AppointNextSpeakerActionCost: getSelectivelyVisibleResources(
			c.AppointNextSpeakerActionCost,
			true,
		),
		// Judiciary branch
		InspectHistoryActionCost: getSelectivelyVisibleResources(
			c.InspectHistoryActionCost,
			true,
		),
		InspectBallotActionCost: getSelectivelyVisibleResources(
			c.InspectBallotActionCost,
			true,
		),
		InspectAllocationActionCost: getSelectivelyVisibleResources(
			c.InspectAllocationActionCost,
			true,
		),
		AppointNextPresidentActionCost: getSelectivelyVisibleResources(
			c.AppointNextPresidentActionCost,
			true,
		),
		SanctionCacheDepth: getSelectivelyVisibleInteger(
			c.SanctionCacheDepth,
			true,
		),
		HistoryCacheDepth: getSelectivelyVisibleInteger(
			c.HistoryCacheDepth,
			true,
		),
		AssumedResourcesNoReport: getSelectivelyVisibleResources(
			c.AssumedResourcesNoReport,
			true,
		),
		SanctionLength: getSelectivelyVisibleInteger(
			c.SanctionLength,
			true,
		),
		// Legislative branch
		SetVotingResultActionCost: getSelectivelyVisibleResources(
			c.SetVotingResultActionCost,
			true,
		),
		SetRuleToVoteActionCost: getSelectivelyVisibleResources(
			c.SetRuleToVoteActionCost,
			true,
		),
		AnnounceVotingResultActionCost: getSelectivelyVisibleResources(
			c.AnnounceVotingResultActionCost,
			true,
		),
		UpdateRulesActionCost: getSelectivelyVisibleResources(
			c.UpdateRulesActionCost,
			true,
		),
		AppointNextJudgeActionCost: getSelectivelyVisibleResources(
			c.AppointNextJudgeActionCost,
			true,
		),
	}
}

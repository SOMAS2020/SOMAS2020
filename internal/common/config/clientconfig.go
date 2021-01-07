package config

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

// ClientConfig contains config information visible to clients.
type ClientConfig struct {
	CostOfLiving                shared.Resources
	MinimumResourceThreshold    shared.Resources
	MaxCriticalConsecutiveTurns uint
	DisasterConfig              ClientDisasterConfig
	IIGOClientConfig            IIGOConfig
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
	SanctionCacheDepth             uint
	HistoryCacheDepth              uint
	AssumedResourcesNoReport       shared.Resources
	SanctionLength                 uint
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
	DisasterPeriod      SelectivelyVisibleUint
	StochasticDisasters SelectivelyVisibleBool // if true, period between disasters becomes random. If false, it will be consistent (deterministic)
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
		DisasterPeriod: getSelectivelyVisibleUint(
			c.Period,
			c.PeriodVisible,
		),
		StochasticDisasters: getSelectivelyVisibleBool(
			c.StochasticPeriod,
			c.StochasticPeriodVisible,
		),
	}
}

// GetClientIIGOConfig gets ClientIIGOConfig
func (c IIGOConfig) GetClientIIGOConfig() IIGOConfig {
	return c
}

package config

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

// ClientConfig contains config information visible to clients.
type ClientConfig struct {
	CostOfLiving                shared.Resources
	MinimumResourceThreshold    shared.Resources
	MaxCriticalConsecutiveTurns uint
	DisasterConfig              ClientDisasterConfig
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

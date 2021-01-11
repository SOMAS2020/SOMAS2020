package team4

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

// ClientConfig is the data type for the different client configs available
type ClientConfig int

// Team 4 client configurations
const (
	honest ClientConfig = iota
	dishonest
	moderate
	forTesting
)

func configureClient(config ClientConfig) (retConfig internalParameters) {
	switch config {
	case honest: // honest
		retConfig = internalParameters{
			greediness:                   0.2,
			selfishness:                  0.4,
			fairness:                     0.8,
			collaboration:                0.7,
			riskTaking:                   0.3,
			maxPardonTime:                2,
			maxTierToPardon:              shared.SanctionTier3,
			minTrustToPardon:             0.65,
			historyWeight:                1,   // in range [0,1]
			historyFullTruthfulnessBonus: 0.1, // in range [0,1]
			monitoringWeight:             1,   // in range [0,1]
			monitoringResultChange:       0.1, // in range [0,1]
		}
	case dishonest: // dishonest
		retConfig = internalParameters{
			greediness:                   0.8,
			selfishness:                  0.7,
			fairness:                     0.3,
			collaboration:                0.5,
			riskTaking:                   0.8,
			maxPardonTime:                10,
			maxTierToPardon:              shared.SanctionTier5,
			minTrustToPardon:             0.3,
			historyWeight:                1,    // in range [0,1]
			historyFullTruthfulnessBonus: 0.05, // in range [0,1]
			monitoringWeight:             0.7,  // in range [0,1]
			monitoringResultChange:       0.05, // in range [0,1]
		}
	case moderate: // middle
		retConfig = internalParameters{
			greediness:                   0.5,
			selfishness:                  0.5,
			fairness:                     0.5,
			collaboration:                0.5,
			riskTaking:                   0.5,
			maxPardonTime:                5,
			maxTierToPardon:              shared.SanctionTier3,
			minTrustToPardon:             0.6,
			historyWeight:                1,    // in range [0,1]
			historyFullTruthfulnessBonus: 0.15, // in range [0,1]
			monitoringWeight:             1,    // in range [0,1]
			monitoringResultChange:       0.07, // in range [0,1]
		}
	case forTesting: // testing only
		retConfig = internalParameters{
			greediness:                   0.5,
			selfishness:                  0.5,
			fairness:                     0.5,
			collaboration:                0.5,
			riskTaking:                   0.5,
			maxPardonTime:                5,
			maxTierToPardon:              shared.SanctionTier3,
			minTrustToPardon:             0.6,
			historyWeight:                1,   // in range [0,1]
			historyFullTruthfulnessBonus: 0.1, // in range [0,1]
			monitoringWeight:             1,   // in range [0,1]
			monitoringResultChange:       0.1, // in range [0,1]
		}
	default:
		retConfig = internalParameters{
			greediness:                   0.2,
			selfishness:                  0.4,
			fairness:                     0.8,
			collaboration:                0.7,
			riskTaking:                   0.3,
			maxPardonTime:                2,
			maxTierToPardon:              shared.SanctionTier3,
			minTrustToPardon:             0.65,
			historyWeight:                1,   // in range [0,1]
			historyFullTruthfulnessBonus: 0.1, // in range [0,1]
			monitoringWeight:             1,   // in range [0,1]
			monitoringResultChange:       0.1, // in range [0,1]
		}
	}
	return retConfig
}

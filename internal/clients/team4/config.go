package team4

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

// ClientConfig is the data type for the different client configs available
type ClientConfig int

// TeamIDs
const (
	C1 ClientConfig = iota
	C2
)

func configureClient(config ClientConfig) internalParameters {
	switch config {
	case C1:
		return internalParameters{
			greediness:                   0,
			selfishness:                  0,
			fairness:                     0,
			collaboration:                0,
			riskTaking:                   0,
			maxPardonTime:                3,
			maxTierToPardon:              shared.SanctionTier3,
			minTrustToPardon:             0.6,
			historyWeight:                1,   // in range [0,1]
			historyFullTruthfulnessBonus: 0.1, // in range [0,1]
			monitoringWeight:             1,   // in range [0,1]
			monitoringResultChange:       0.1, // in range [0,1]
		}
	case C2:
		return internalParameters{
			greediness:                   0,
			selfishness:                  0,
			fairness:                     0,
			collaboration:                0,
			riskTaking:                   0,
			maxPardonTime:                3,
			maxTierToPardon:              shared.SanctionTier3,
			minTrustToPardon:             0.6,
			historyWeight:                1,   // in range [0,1]
			historyFullTruthfulnessBonus: 0.1, // in range [0,1]
			monitoringWeight:             1,   // in range [0,1]
			monitoringResultChange:       0.1, // in range [0,1]
		}
	}

}

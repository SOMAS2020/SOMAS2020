package baseclient

import (
	"math"
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// DecideForage makes a foraging decision, currently you can only forage deer, but fishing will be available later
// the forageContribution can not be larger than the total resources available
func (c *BaseClient) DecideForage() (shared.ForageDecision, error) {
	forageType := int(math.Round(rand.Float64())) // 0 or 1 with equal prob.
	return shared.ForageDecision{Type: forageType, Contribution: rand.Float64()}, nil
}

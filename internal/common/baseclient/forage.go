package baseclient

import (
	"math"
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// DecideForage makes a foraging decision
// the forageContribution can not be larger than the total resources available
func (c *BaseClient) DecideForage() (shared.ForageDecision, error) {
	ft := int(math.Round(rand.Float64())) // 0 or 1 with equal prob.
	return shared.ForageDecision{
		Type:         shared.ForageType(ft),
		Contribution: shared.Resources(rand.Float64() * 20),
	}, nil
}

// ForageUpdate is called by the server upon completion of a foraging session. This handler can be used by clients to
// analyse their returns - resources returned to them, as well as number of fish/deer caught.
func (c *BaseClient) ForageUpdate(initialDecision shared.ForageDecision, resourceReturn shared.Resources, numberCaught uint) {
}

package baseclient

import (
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// DecideForage makes a foraging decision, currently you can only forage deer, but fishing will be available later
// the forageContribution can not be larger than the total resources available
func (c *BaseClient) DecideForage() (shared.ForageDecision, error) {
	return shared.ForageDecision{
		Type: shared.DeerForageType,
		Contribution: shared.Resources(rand.Float64() * 5),
	}, nil
}

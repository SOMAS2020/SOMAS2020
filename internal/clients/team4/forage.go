package team4

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// DecideForage makes a foraging decision, currently you can only forage deer, but fishing will be available later
// the forageContribution can not be larger than the total resources available
type forageStorage struct {
	preferedForageMethod int
}

func (c *client) DecideForage() (shared.ForageDecision, error) {
	ft := c.forage.preferedForageMethod
	scale := 5 * c.getSafeResourceLevel()
	resources := c.getResources() - c.getSafeResourceLevel()*(2-shared.Resources(c.internalParam.riskTaking)*scale)
	return shared.ForageDecision{
		Type:         shared.ForageType(ft),
		Contribution: shared.Resources(resources),
	}, nil
}

// ForageUpdate is called by the server upon completion of a foraging session. This handler can be used by clients to
// analyse their returns - resources returned to them, as well as number of fish/deer caught.
func (c *client) ForageUpdate(initialDecision shared.ForageDecision, resourceReturn shared.Resources, numberCaught uint) {
	//	dump("./forage.txt", "Call: %v\nResources: %v\n, numberCaught%v\n\n", initialDecision, resourceReturn, numberCaught)

}

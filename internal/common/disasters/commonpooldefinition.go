package disasters

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// DisasterMitigate mitigates the disaster's damage using CP before it hits the island
// 	+ If cp have enough resource to fully mitigate the disaster, island's personal resources won't be affected
//	+ Else, island will be affected with leftover damaged. Each island receives proportional helps from cp with respect to the total damage on 6 islands
//	+ EX: if island 1 gets hit the most then cp will help island 1 more than the rest
//	+ There is a threshold T that incentivizes contributing to the cp from islands
//	+ If threshold T is met, which mean islands have prepared well before the disaster hit, the damage that disaster have on cp and islands will be half
//	+ If not, cp and islands will feel the full effect of disaster
//	+ For now, each 1 effect of disaster is equivalent to DisasterConfig.MagnitudeResourceMulitplier resources
func (e Environment) DisasterMitigate(cpResources shared.Resources, effects DisasterEffects, dConf config.DisasterConfig) map[shared.ClientID]shared.Magnitude {
	// compute total effect of disaster on 6 islands
	totalEffect := 0.0
	updatedIndividualEffects := map[shared.ClientID]float64{}

	for _, effect := range effects.Absolute {
		totalEffect = totalEffect + effect
	}
	adjMagnitude := 0.0

	if cpResources >= dConf.CommonpoolThreshold { //exceeds cp threshold
		adjMagnitude = (totalEffect / 2) * dConf.MagnitudeResourceMultiplier // better prep means less effect
	} else {
		adjMagnitude = totalEffect * dConf.MagnitudeResourceMultiplier // no prep so goodluck
	}

	//compute remaining effect for each island
	if adjMagnitude <= float64(cpResources) { // Case when cp can fully mitigate disaster
		for islandID := range effects.Absolute {
			updatedIndividualEffects[islandID] = 0 // 0 means fully mitigated
		}
	} else { // Case when damage is too high, cp cannot fully mitigate
		leftOverEffect := adjMagnitude - float64(cpResources) // damage that has to be mitigated by islands
		for islandID, effect := range effects.Proportional {
			updatedIndividualEffects[islandID] = leftOverEffect * effect //leftover damage for each island is computed proportionally with respect to the damage on 6 islands
		}
	}
	return updatedIndividualEffects
}

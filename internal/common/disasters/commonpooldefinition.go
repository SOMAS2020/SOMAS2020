package disasters

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type Commonpool struct {															
	Resource, Threshold  uint
}

// DisasterMitigate mitigates the disaster's damage using CP before it hits the island
// 	+ If cp have enough resource to fully mitigate the disaster, island's personal resources won't be affected 
//	+ Else, island will be affected with leftover damaged. Each island receives proportional helps from cp with respect to the total damage on 6 islands 
//	+ EX: if island 1 gets hit the most then cp will help island 1 more than the rest 

//	+ There is a threshold T that incentivizes contributing to the cp from islands
//	+ If threshold T is met, which mean islands have prepared well before the disaster hit, the damage that disaster have on cp and islands will be half 
//	+ If not, cp and islands will feel the full effect of disaster
//	+ For now, each 1 effect of disaster is equivalent to 1000 resource 
func (e Environment) DisasterMitigate(individualEffect map[shared.ClientID]float64, proportionalEffect map[shared.ClientID]float64) map[shared.ClientID]float64 {				//input is the porportional effect of each island
	// compute total effect of disaster on 6 islands 
	totalEffect := 0.0				
	updatedIndividualEffect := map[shared.ClientID]float64{}

	for _, effect := range individualEffect {
		totalEffect = totalEffect + effect
	}

	// compute cp power towards disaster
	var newMagnitude float64
	tempResource := float64(e.CommonPool.Resource)
	if e.CommonPool.Resource >= e.CommonPool.Threshold { 				//exceeds cp threshold
		newMagnitude = (totalEffect / 2) * 1000 						// better prep means less effect
	} else {
		newMagnitude = totalEffect * 1000 								// no prep so goodluck
	}

	//compute remaining effect for each island
	if newMagnitude <= tempResource { 									// Case when cp can fully mitigate disaster
		tempResource = tempResource - newMagnitude
		for islandID, _ := range individualEffect {
			updatedIndividualEffect[islandID] = 0 						// 0 means fully mitigated
		}
	} else {															// Case when damage is too high, cp cannot fully mitigate 
		leftover_effect := newMagnitude-tempResource  					// damage that has to be mitigated by islands
		tempResource = 0                              					// nothing remains in the cp
		for islandID, effect := range proportionalEffect {
			updatedIndividualEffect[islandID] = leftover_effect * effect //leftover damage for each island is computed proportionally with respect to the damage on 6 islands
		}
	}
	e.CommonPool.Resource = uint(tempResource)
	return updatedIndividualEffect				
}



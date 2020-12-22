package disasters

import (

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type Commonpool struct {
	Resource, Threshold  uint
}

//islandContribution takes the resource donation from island to common pool
func (e *Environment) islandContribution(resource uint) {
	e.CommonPool.Resource += resource
}

//disasterMitigate mitigates the disaster's damage using CP efore it hits the island
func (e *Environment) DisasterMitigate(individualEffect map[shared.ClientID]float64, proportionalEffect map[shared.ClientID]float64) (map[shared.ClientID]float64, uint) {				//input is the porportional effect of each island
	// compute total effect of disaster
	totalEffect := 0.0
	updatedIndividualEffect := map[shared.ClientID]float64{}

	for _, effect := range individualEffect {
		totalEffect = totalEffect + effect
	}

	// compute cp power towards disaster
	var newmag float64
	tempResource := float64(e.CommonPool.Resource)
	if e.CommonPool.Resource >= e.CommonPool.Threshold { //exceeds cp threshold
		newmag = (totalEffect / 2) * 1000 // better prep means less effect
	} else {
		newmag = totalEffect * 1000 // no prep so goodluck
	}

	//compufte reminant effect for each island
	if newmag <= tempResource { //magn*1000 => effect in resource unit
		tempResource = tempResource - newmag
		for islandID, _ := range individualEffect {
			updatedIndividualEffect[islandID] = 0 // 0 means fully mitigated
		}
	} else {

		leftover_effect := newmag-tempResource  //disaster that has to be mitigated by islands
		tempResource = 0                              // nothing remains in the cp
		for islandID, effect := range proportionalEffect {
			updatedIndividualEffect[islandID] = leftover_effect * effect // 0 means fully mitigated
		}
	}
	e.CommonPool.Resource = uint(tempResource)
	//log.Printf("\nCommon pool resource the end of mitigation: %d\n", e.CommonPool.Resource)
	return updatedIndividualEffect, uint(tempResource) 					//the leftover effect that islands need to mitigate by their own resources
}



package team2

import (
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type ForagingResults struct {
	// no. teams that we know (they have shared the info with us)
	Hunters int
	Fishers int
	// total from all the teams we know
	Result shared.Resources
}

func (c *client) DecideForage() (shared.ForageDecision, error) {
	// implement a normal distribution which shifts closer to hunt or fish
	var Threshold float64 = c.decideHuntingLikelihood() //number from 0 to 1
	var forageDecision shared.ForageType

	//we fish when above the threshold
	if rand.Float64() > Threshold {
		forageDecision = shared.FishForageType
	} else {
		// we hunt when below the threshold
		forageDecision = shared.DeerForageType
	}

	c.Logf("[Our beautiful Forage decision]:", forageDecision)
	c.Logf("[Our beautiful Forage contribution]:", shared.Resources(c.DecideForageAmount(Threshold)))

	return shared.ForageDecision{
		Type:         forageDecision,
		Contribution: shared.Resources(c.DecideForageAmount(Threshold)),
	}, nil
}

//Decide amount of resources to put into foraging
func (c *client) DecideForageAmount(foragingDecisionThreshold float64) shared.Resources {
	// we have given to the pool already by this point in the turn

	if c.criticalStatus() || c.getAgentExcessResources() == 0 {
		return shared.Resources(0)
	}
	// otherwise we have excess to allocate to foraging
	multiplier := shared.Resources(0)
	switch c.setAgentStrategy() {
	case Altruist:
		multiplier = 0.1 //when the common pool is struggling we forage with 10% of our resources
	case FairSharer:
		multiplier = 0.2 //by default give 20% of our excess
	case Selfish:
		multiplier = 0.3 //when common pool is doing well we take risk and forage with 30% of our excess
	}

	resourcesForForaging := c.getAgentExcessResources() * multiplier //how much of our excess we forage with depends on the method of play
	c.Logf("Resources allocated to foraging:", resourcesForForaging)

	return shared.Resources(resourcesForForaging)
}

//being the only agent to hunt is undesirable, having one hunting partner is the desirable, the more hunters after that the less we want to hunt
func (c *client) decideHuntingLikelihood() float64 { //will move the threshold, higher value means more likely to hunt
	// Todo: this could be null
	hunters := c.otherHunters()
	if hunters > 1.0 && hunters < 1.2 { //in the case when one other person only is hunting
		return 0.95
	} else if hunters > 1 { //if no one is likely to hunt then we do default probability
		return 0.95 - (hunters * 0.15) //default hunt probability is 10%, the less people hunting the more likely we do it

	} else {
		return 0.1 //when no one is likely to hunt we have a default 10% chance of hunting just in the off chance another person hunts
	}
}

//EXTRA FUNCTIONALITY: find the probability based off of how agents act in specific circumstances not just the agents themselves
func (c *client) otherHunters() float64 { //will return a value of how many agents will likely hunt
	aliveClients := c.getAliveClients()
	HuntNum := 0.00 //this the average number of likely hunters

	for _, id := range aliveClients { //loop through every agent
		if islandForageHist, ok := c.foragingReturnsHist[id]; ok {
			for _, forageInfo := range islandForageHist { //loop through the agents array and add their average to HuntNum
				if len(islandForageHist) != 0 {
					HuntNum += float64(forageInfo.DecisionMade.Type) / float64(len(islandForageHist)) //add the agents decision to HuntNum and then average
				}
			}
		}
	}
	return HuntNum
}

//TODO: This function needs to be changed according to Eirik, I have no idea how
func (c *client) ReceiveForageInfo(forageInfos []shared.ForageShareInfo) {
	for _, info := range forageInfos {
		forageInfo := ForageInfo{
			DecisionMade:      info.DecisionMade,
			ResourcesObtained: info.ResourceObtained,
		}
		c.foragingReturnsHist[info.SharedFrom] = append(c.foragingReturnsHist[info.SharedFrom], forageInfo)
	}
}

// MakeForageInfo allows clients to share their most recent foraging DecisionMade, ResourceObtained from it to
// other clients.
// OPTIONAL. If this is not implemented then all values are nil.
func (c *client) MakeForageInfo() shared.ForageShareInfo {
	contribution := shared.ForageDecision{Type: shared.DeerForageType, Contribution: 0}
	return shared.ForageShareInfo{DecisionMade: contribution, ResourceObtained: 0, ShareTo: []shared.ClientID{}}
}

package team2

import (
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type ForagingResults struct {
	Hunters int //no. teams that we know (they have shared the info with us)
	Fishers int
	Result  shared.Resources //total from all the teams we know
}

func (c *client) DecideForage() (shared.ForageDecision, error) {
	// implement a normal distribution which shifts closer to hunt or fish
	var Threshold float64 = c.decideHuntingLikelihood() //number from 0 to 1
	var forageDecision shared.ForageType
	if rand.Float64() > Threshold { //we fish when above the threshold
		forageDecision = 1
	} else { // we hunt when below the threshold
		forageDecision = 0
	}
	return shared.ForageDecision{
		Type:         shared.ForageType(forageDecision),
		Contribution: shared.Resources(c.DecideForageAmount(Threshold)), //contribute fixed amount for now
	}, nil
}

//Decide amount of resources to put into foraging
func (c *client) DecideForageAmount(foragingDecisionThreshold float64) shared.Resources {
	ourResources := c.gameState().ClientInfo.Resources // we have given to the pool already by this point in the turn
	if criticalStatus(c) {
		return 0
	}
	if ourResources < c.agentThreshold() && foragingDecisionThreshold < 0.6 { //tune threshold (lower threshold = more likely to have better reward from foraging)
		return (ourResources - c.agentThreshold()) / 2 //tune divisor
	}
	resourcesForForaging := (ourResources - c.agentThreshold())
	return resourcesForForaging
}

//being the only agent to hunt is undesirable, having one hunting partner is the desirable, the more hunters after that the less we want to hunt
func (c *client) decideHuntingLikelihood() float64 { //will move the threshold, higher value means more likely to hunt
	hunters := c.otherHunters()
	if hunters == 1.0 { //in the case when one other person only is hunting
		return 0.95
	} else if hunters > 1 { //if no one is likely to hunt then we do default probability
		return 0.95 - (hunters * 0.15) //default hunt probability is 10%, the less people hunting the more likely we do it

	} else {
		return 0.1 //when no one is likely to hunt we have a default 10% chance of hunting just in the off chance another person hunts
	}
}

//EXTRA FUNCTIONALITY: find the probability based off of how agents act in specific circumstances not just the agents themselves
func (c *client) otherHunters() float64 { //will return a value of how many agents will likely hunt

	HuntNum := 0.00                                                //this the average number of likely hunters
	for id, lifeStatus := range c.gameState().ClientLifeStatuses { //loop through every agent
		if lifeStatus != shared.Dead { //client is dead ignore their decisions
			for _, forageInfo := range c.foragingReturnsHist[id] { //loop through the agents array and add their average to HuntNum
				HuntNum += float64(forageInfo.DecisionMade.Type) / float64(len(c.foragingReturnsHist[id])) //add the agents decision to HuntNum and then average
			}
		}
	}
	return HuntNum
}

//TODO: This function needs to be changed according to Eirik, I have no idea how
func (c *client) ReceiveForageInfo(neighbourForaging []shared.ForageShareInfo) {
	// updates our foragingReturnsHist with the decisions everyone made
	roi := map[shared.ClientID]shared.Resources{}
	for _, val := range neighbourForaging {
		/*c.foragingReturnsHist[val.SharedFrom].append(
		shared.ForageInfo{ DecisionMade: val.DecisionMade, //this needs to change
			ResourcesObtained: val.ResourcesObtained})
		*/
		if val.DecisionMade.Type == shared.DeerForageType {
			roi[val.SharedFrom] = val.ResourceObtained / shared.Resources(val.DecisionMade.Contribution) * 100
		}
	}
}

// MakeForageInfo allows clients to share their most recent foraging DecisionMade, ResourceObtained from it to
// other clients.
// OPTIONAL. If this is not implemented then all values are nil.
func (c *client) MakeForageInfo() shared.ForageShareInfo {
	contribution := shared.ForageDecision{Type: shared.DeerForageType, Contribution: 0}
	return shared.ForageShareInfo{DecisionMade: contribution, ResourceObtained: 0, ShareTo: []shared.ClientID{}}
}

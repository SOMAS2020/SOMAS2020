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
	var Threshold float64 = decideThreshold(c) //number from 0 to 1

	if rand.Float64() > Threshold { //we fish when above the threshold
		ft := 1
	} else { // we hunt when below the threshold
		ft := 0
	}
	return shared.ForageDecision{
		Type:         shared.ForageType(ft),
		Contribution: shared.Resources(20), //contribute fixed amount for now
	}, nil
}

//Decide amount of resources to put into foraging
func (c *client) DecideForageAmount(foragingDecisionThreshold) shared.Resources {
	ourResources := c.gameState().ClientInfo.Resources // we have given to the pool already by this point in the turn
	if criticalStatus(c) {
		return 0
	}
	if ourResources < internalThreshold(c) && foragingDecisionThreshold < 0.6 { //tune threshold (lower threshold = more likely to have better reward from foraging)
		return (ourResources - determineThreshold(c)) / 2 //tune divisor
	}
	resourcesForForaging := (ourResources - internalThreshold)
	return resourcesForForaging
}

// ReceiveForageInfo lets clients know what other clients have obtained from their most recent foraging attempt.
// Most recent foraging attempt includes information about: foraging DecisionMade and ResourceObtained as well
// as where this information came from.
// OPTIONAL.
func (c *client) ReceiveForageInfo(neighbourForaging []shared.ForageShareInfo) {
	// Return on Investment

	for _, val := range neighbourForaging {
		decisionMade := val.DecisionMade
		resourcesObtained := val.ResourceObtained

		newResult := ForageInfo{
			DecisionMade:      decisionMade,
			ResourcesObtained: resourcesObtained,
		}
		hist := append(c.foragingReturnsHist[val.SharedFrom], newResult)
		c.foragingReturnsHist[val.SharedFrom] = hist
	}
}

//being the only agent to hunt is undesirable, having one hunting partner is the desirable, the more hunters after that the less we want to hunt
func decideThreshold(c *client) float64 { //will move the threshold, higher value means more likely to hunt
	if Otheragentinfo(c) == 1 { //in the case when one other person only is hunting
		return 0.95
	} else if Otheragentinfo(c) > 1 {
		return 0.95 - (Otheragentinfo(c) * 0.15) //default hunt probability is 10%, the less people hunting the more likely we do it
	} else {
		return 0.1 //when no one is likely to hunt we have a default 10% chance of hunting just in the off chance another person hunts
	}
}

//TODO: LINE 83 foragingReturnHist how to access a clients foraging decision?
//EXTRA FUNCTIONALITY: find the probability based off of how agents act in specific circumstances not just the agents themselves
func Otheragentinfo(c *client) float64 { //will return a value of how many agents will likely hunt
	HuntNum := 0.00
	for _, id := range clientInfo {​​​​​​​
		if   {//client is dead ignore their decision
			for index, forageInfo := range c.foragingReturnsHist[id] {
			forageInfo.DecisionMade.ForageType
			}
		}
	}​​​​​​​
	return HuntNum
}

//TODO: This function needs to be changed according to Eirik, I have no idea why
func (c *client) ReceiveForageInfo(neighbourForaging []shared.ForageShareInfo) {
	// updates our foragingReturnsHist with the decisions everyone made
	roi := map[shared.ClientID]shared.Resources{}
	for _, val := range neighbourForaging {
		c.foragingReturnsHist[val.SharedFrom].append(shared.ForageInfo{DecisionMade: val.DecisionMade,
			ResourcesObtained: val.ResourcesObtained})

		if val.DecisionMade.Type == shared.DeerForageType {
			roi[val.SharedFrom] = val.ResourceObtained / shared.Resources(val.DecisionMade.Contribution) * 100
		}
	}
}

// MakeForageInfo allows clients to share their most recent foraging DecisionMade, ResourceObtained from it to
// other clients.
// OPTIONAL. If this is not implemented then all values are nil.
func (c *BaseClient) MakeForageInfo() shared.ForageShareInfo {
	contribution := shared.ForageDecision{Type: shared.DeerForageType, Contribution: 0}
	return shared.ForageShareInfo{DecisionMade: contribution, ResourceObtained: 0, ShareTo: []shared.ClientID{}}
}

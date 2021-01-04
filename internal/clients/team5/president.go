package team5

import (
	"math/rand"
	"time"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type president struct {
	*baseclient.BasePresident

	c *client
}

//This I'm not 100% certain

func (c *client) GetClientPresidentPointer() roles.President {
	c.Logf("Team 5 is now the President, Shalom to all")

	return &c.team5president
}

// EvaluateAllocationRequests sets allowed resource allocation based on each islands requests
// COMPULSORY: decide how much each island can take from the common pool.
func (pres *president) EvaluateAllocationRequests(resourceRequest map[shared.ClientID]shared.Resources, availCommonPool shared.Resources) (map[shared.ClientID]shared.Resources, bool) {
	var requestSum shared.Resources
	resourceAllocation := make(map[shared.ClientID]shared.Resources)

	for _, request := range resourceRequest { //sum of all available requests
		requestSum += request
	}

	if requestSum < 0.8*availCommonPool || requestSum == 0 { //If it's smaller than 0.8 of the CP or 0 we then allocate resources
		resourceAllocation = resourceRequest
	} else if requestSum > 0.8*availCommonPool && requestSum < availCommonPool { // If it's between 0.8 or 1 times the size of the CP, we divide requests by 2 and then allocate them
		for id, request := range resourceRequest {
			resourceAllocation[id] = shared.Resources(request * availCommonPool * 1 / (2 * requestSum)) //returns a fraction of requests
		}
	} else if requestSum > availCommonPool && requestSum < 2*availCommonPool {
		for id, request := range resourceRequest {
			resourceAllocation[id] = shared.Resources(request * availCommonPool * 1 / (5 * requestSum)) // requests smaller fraction of requests
		}
	} else {
		for id := range resourceRequest { //for any request above two times the CP, we reject it
			resourceAllocation[id] = 0
		}
	}
	return resourceAllocation, true
}

//TODO: NEED TO LOOK AT THE RULES
// PickRuleToVote chooses a rule proposal from all the proposals
// COMPULSORY: decide a method for picking a rule to vote on
//NOT SURE ABOUT THISSSSSSS
func (pres *president) PickRuleToVote(rulesProposals []string) (string, bool) {
	if len(rulesProposals) == 0 {
		// No rules were proposed by the islands
		return "", false
	}
	return rulesProposals[rand.Intn(len(rulesProposals))], true
}

// SetTaxationAmount sets taxation amount for all of the living islands
// islandsResources: map of all the living islands and their remaining resources
// COMPULSORY: decide how much to tax islands, using rules as a guide if you wish.

//Essentially the more the resources the greater the applied tax, as the point of the game is to survive
func (pres *president) SetTaxationAmount(islandsResources map[shared.ClientID]shared.Resources) (map[shared.ClientID]shared.Resources, bool) {
	var totalrecleft shared.Resources
	taxAmountMap := make(map[shared.ClientID]shared.Resources)

	for _, t := range islandsResources { //sum of all available resources
		totalrecleft += t
	}

	for id, resourceLeft := range islandsResources {
		taxAmountMap[id] = (resourceLeft * resourceLeft / totalrecleft) //instead of applying random taxation, the taxation is now directly proportional to each islands remaining resources/total remaining resources

	}
	return taxAmountMap, true
}

// PaySpeaker pays the speaker a salary.
// OPTIONAL: override to pay the Speaker less than the full amount.

//For this, whenever we are presidents, the salary of the speaker comes in hand with our own wealth state
//this is for the sake of everyone paying less, thus having a higher chance of our island to recover
func (pres *president) PaySpeaker(salary shared.Resources) (shared.Resources, bool) {
	if pres.c.wealth() == jeffBezos {
		return salary, true
	} else if pres.c.wealth() == middleClass {
		salary = salary * 0.8
		return salary, true
	} else {
		salary = salary * 0.5
		return salary, true
	}
}

//TODO: NEED TO LOOK AT THE RULES
// CallSpeakerElection is called by the executive to decide on power-transfer
// COMPULSORY: decide when to call an election following relevant rulesInPlay if you wish
func (pres *president) CallSpeakerElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) shared.ElectionSettings {
	// example implementation calls an election if monitoring was performed and the result was negative
	// or if the number of turnsInPower exceeds 3
	var electionsettings = shared.ElectionSettings{
		VotingMethod:  shared.Plurality,
		IslandsToVote: allIslands,
		HoldElection:  false,
	}
	if monitoring.Performed && !monitoring.Result {
		electionsettings.HoldElection = true
	}
	if turnsInPower >= 2 {
		electionsettings.HoldElection = true
	}
	return electionsettings
}

// DecideNextSpeaker returns the ID of chosen next Speaker
// OPTIONAL: override to manipulate the result of the election

//Spicing things up and randomly choosing the speaker
func (pres *president) DecideNextSpeaker(winner shared.ClientID) shared.ClientID {
	rand.Seed(time.Now().UnixNano())
	return shared.ClientID(rand.Intn(5))
}

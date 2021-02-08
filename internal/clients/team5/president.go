package team5

import (
	"math/rand"
	"time"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type president struct {
	*baseclient.BasePresident

	c *client
}

func (c *client) GetClientPresidentPointer() roles.President {
	c.Logf("Team 5 is now the President, Shalom to all")
	return &president{c: c, BasePresident: &baseclient.BasePresident{GameState: c.ServerReadHandle.GetGameState()}}
}

func (p *president) EvaluateAllocationRequests(resourceRequest map[shared.ClientID]shared.Resources, availCommonPool shared.Resources) shared.PresidentReturnContent {
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

	return shared.PresidentReturnContent{
		ContentType: shared.PresidentAllocation,
		ResourceMap: resourceAllocation,
		ActionTaken: true,
	}
}

//TODO: NEED TO LOOK AT THE RULES
// PickRuleToVote chooses a rule proposal from all the proposals
func (p *president) PickRuleToVote(rulesProposals []rules.RuleMatrix) shared.PresidentReturnContent {
	// DefaulContentType: No rules were proposed by the islands

	// proposedRule := ""
	actionTaken := false

	// if some rules were proposed
	if len(rulesProposals) != 0 {
		// proposedRule = rulesProposals[rand.Intn(len(rulesProposals))]
		actionTaken = true
	}

	return shared.PresidentReturnContent{
		ContentType: shared.PresidentRuleProposal,
		// ProposedRuleMatrix: proposedRule,
		ActionTaken: actionTaken,
	}
}

// SetTaxationAmount sets taxation amount for all of the living islands
// islandsResources: map of all the living islands and their remaining resources
// COMPULSORY: decide how much to tax islands, using rules as a guide if you wish.

//Essentially the more the resources the greater the applied tax, as the point of the game is to survive
func (p *president) SetTaxationAmount(islandsResources map[shared.ClientID]shared.ResourcesReport) shared.PresidentReturnContent {
	var totalrecleft shared.Resources
	taxAmountMap := make(map[shared.ClientID]shared.Resources)

	for _, t := range islandsResources { //sum of all available resources
		totalrecleft += t.ReportedAmount
	}

	for id, resourceLeft := range islandsResources {
		if resourceLeft.Reported {
			if totalrecleft != 0 {
				taxAmountMap[id] = (resourceLeft.ReportedAmount * resourceLeft.ReportedAmount / totalrecleft) //instead of applying random taxation, the taxation is now directly proportional to each islands remaining resources/total remaining resources
			} else {
				taxAmountMap[id] = 0 //If reported is 0 we avoid dividing by 0
			}
		} else {
			taxAmountMap[id] = 20 // Flat tax rate which high enough, in order to encourage islands to share their resources
		}
	}
	return shared.PresidentReturnContent{
		ContentType: shared.PresidentTaxation,
		ResourceMap: taxAmountMap,
		ActionTaken: true,
	}
}

// PaySpeaker pays the speaker a salary.

//For this, whenever we are presidents, the salary of the speaker comes in hand with our own wealth state
//this is for the sake of everyone paying less, thus having a higher chance of our island to recover
func (p *president) PaySpeaker() shared.PresidentReturnContent {
	// return p.BasePresident.PaySpeaker()

	SpeakerSalaryRule, ok := p.GameState.RulesInfo.CurrentRulesInPlay["salary_cycle_speaker"]
	var salary shared.Resources = 0
	if ok {
		salary = shared.Resources(SpeakerSalaryRule.ApplicableMatrix.At(0, 1))
	}
	if p.c.wealth() == jeffBezos {
		return shared.PresidentReturnContent{
			ContentType:   shared.PresidentSpeakerSalary,
			SpeakerSalary: salary,
			ActionTaken:   true,
		}
	} else if p.c.wealth() == middleClass {
		salary = salary * 0.8
	} else {
		salary = salary * 0.5
	}

	return shared.PresidentReturnContent{
		ContentType:   shared.PresidentSpeakerSalary,
		SpeakerSalary: salary,
		ActionTaken:   true,
	}
}

//TODO: NEED TO LOOK AT THE RULES
// CallSpeakerElection is called by the executive to decide on power-transfer
// COMPULSORY: decide when to call an election following relevant rulesInPlay if you wish
func (p *president) CallSpeakerElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) shared.ElectionSettings {
	// example implementation calls an election if monitoring was performed and the result was negative
	// or if the number of turnsInPower exceeds 3
	var electionsettings = shared.ElectionSettings{
		VotingMethod:  shared.BordaCount,
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

//Deciding next speaker revokes the voting procedure and it is selected using our opinion,
//the island with the highest opinion is selected as the speaker of next round

func (p *president) DecideNextSpeaker(winner shared.ClientID) shared.ClientID {
	rand.Seed(time.Now().UnixNano())
	oparray := []opinionScore{}
	for id, opinion := range p.c.opinions { //stores scores except team 5's in an array
		if id != p.c.GetID() {
			oparray = append(oparray, opinion.getScore())
		}
	}
	_, max := p.c.minmaxOpinion(oparray) //calculates mas score from the array of opinion scores

	if max > 0.0 {
		for client, op := range p.c.opinions { //matches highest score to the id of the island and returns that island as the speaker
			if op.getScore() == max {
				p.c.Logf("Opinion on team %v is %v", client, op.getScore())
				return shared.ClientID(client)
			}

		}
		return shared.ClientID(rand.Intn(5)) //this should never be trigerred, just here for completeness
	}
	return shared.ClientID(rand.Intn(5)) //triggered if max<0 and thus it's better to randomize speaker selection
}

//This function takes in an array of opinions when called and outputs the minimum and maximum scores
//that can then be matched to the id of an island.
func (c client) minmaxOpinion(o []opinionScore) (min opinionScore, max opinionScore) {
	min = o[0]
	max = o[0]
	for _, value := range o {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}
	return min, max
}

// Package team5 contains code for team 5's client implementation
package team5

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const id = shared.Team5

// Wealth defines how much money we have
type Wealth int

// ForageOutcome records the ROI on a foraging session
type ForageOutcome struct {
	turn   uint
	input  shared.Resources
	output shared.Resources
}

// ForageHistory stores history of foraging outcomes
type ForageHistory map[shared.ForageType][]ForageOutcome

type clientConfig struct {
	// Initial non planned foraging
	InitialForageTurns uint

	// Skip forage for x amount of returns if theres no return > 1* multiplier
	SkipForage uint

	// If resources go above this limit we are balling with money
	JBThreshold shared.Resources

	// Middle class:  Middle < Jeff bezos
	MiddleThreshold shared.Resources

	// Poor: Imperial student < Middle
	ImperialThreshold shared.Resources
}

// Client is the island number
type client struct {
	*baseclient.BaseClient

	forageHistory ForageHistory // Stores our previous foraging data

	taxAmount shared.Resources

	// allocation is the president's response to your last common pool resource request
	allocation shared.Resources

	config clientConfig
}

// Classes we are in
const (
	Dying           Wealth = iota // Sets values = 0
	ImperialStudent               // iota sets the folloing values =1
	MiddleClass                   // = 2
	JeffBezos                     // = 3
)

// NewClient get the base client
func NewClient(clientID shared.ClientID) baseclient.Client {
	return &client{
		BaseClient:    baseclient.NewClient(clientID),
		forageHistory: ForageHistory{},
		taxAmount:     0,
		allocation:    0,
		config: clientConfig{
			InitialForageTurns: 10,
			// SkipForage:         5,

			JBThreshold:       1000.0,
			MiddleThreshold:   60.0,
			ImperialThreshold: 30.0,
		},
	}
}

//================================================================
/*  Wealth class */
//================================================================

func (c client) wealth() Wealth {
	cData := c.gameState().ClientInfo
	switch {
	case cData.LifeStatus == shared.Critical:
		c.Logf("Critical") // Debugging
		return Dying
	case cData.Resources > c.config.ImperialThreshold && cData.Resources < c.config.MiddleThreshold:
		c.Logf("Student") // Debugging
		return ImperialStudent
	case cData.Resources > c.config.JBThreshold:
		c.Logf("Threshold %f", c.config.JBThreshold) // Debugging
		c.Logf("Resources %f", cData.Resources)      // Debugging
		return JeffBezos
	default:
		c.Logf("Middle class") // Debugging
		return MiddleClass
	}

}

func (st Wealth) String() string {
	strings := [...]string{"Dying", "ImperialStudent", "MiddleClass", "JeffBezos"}
	if st >= 0 && int(st) < len(strings) {
		return strings[st]
	}
	return fmt.Sprintf("Unkown internal state '%v'", int(st))
}

//================================================================
/*  Init */
//================================================================
func init() {
	baseclient.RegisterClient(
		id,
		&client{
			BaseClient:    baseclient.NewClient(id),
			forageHistory: ForageHistory{},
		},
	)
}

func (c *client) StartOfTurn() {
	c.Logf("Wealth: %v", c.wealth())
	c.Logf("Resources: %v", c.gameState().ClientInfo.Resources)

	for clientID, status := range c.gameState().ClientLifeStatuses { //if not dead then can start the turn, else no return
		if status != shared.Dead && clientID != c.GetID() {
			return
		}
	}

}

/*
================================================================
	FORAGING
		DecideForage() (shared.ForageDecision, error)
		ForageUpdate(shared.ForageDecision, shared.Resources)
================================================================
IntialForage() Gamble the first attempt to try to get rich (born rich)
If we fail we are in middle class and have a 50% chance to go up fast(Deer) or slow (Fish)
If we lose again then we are in the Imperial class and we just fish to try to get back to middle class
*/
func (c *client) DecideForage() (shared.ForageDecision, error) {
	if c.forageHistorySize() < c.config.InitialForageTurns { // Start with initial foraging turns
		return c.InitialForage(), nil
	} else if c.wealth() == Dying { // If dying go to last hope
		return c.lasthopeF(), nil
	} else {
		return c.normalForage(), nil // Else normally forage
	}
}

func (c *client) InitialForage() shared.ForageDecision {
	// Default contribution amount is a random amount between 0 -> 20%
	forageContribution := shared.Resources(0.2*rand.Float64()) * c.gameState().ClientInfo.Resources

	var forageType shared.ForageType

	switch {
	case c.wealth() == JeffBezos: // JB then we have so much might as well gamble 50% of it
		forageContribution = shared.Resources(0.5*rand.Float64()) * c.gameState().ClientInfo.Resources
		forageType = shared.DeerForageType

	case c.wealth() == ImperialStudent: // Imperial student
		forageType = shared.FishForageType
	case c.wealth() == Dying: // Dying
		c.lasthopeF()
	default: // Midle class
		if rand.Float64() < 0.50 {
			forageType = shared.DeerForageType
		} else {
			forageType = shared.FishForageType
		}
	}

	c.Logf("[Forage][Decision]:Initial")

	return shared.ForageDecision{
		Type:         forageType,
		Contribution: forageContribution,
	}
}

// bestForagingType indicates the best foraging method
func bestHistoryForaging(forageHistory ForageHistory) shared.ForageType {
	bestForagingMethod := shared.ForageType(-1)
	bestReturn := 0.0

	for forageType, outcomes := range forageHistory { // For each foraging type
		totalReturn := 0.0
		for _, returns := range outcomes {
			totalReturn += float64(returns.output / returns.input) // cumlative sum of the outputs over the inputs (multiplier)
		}

		totalReturn = totalReturn / float64(len(outcomes))

		if bestReturn < totalReturn { // Finds the method with the highest multiplier
			bestReturn = totalReturn
			bestForagingMethod = forageType
		}
	}

	if bestReturn < 1 {
		bestForagingMethod = shared.ForageType(-1)
	}

	return bestForagingMethod
}

//normalForage is based on the method that previously gave the best multiplier return
func (c *client) normalForage() shared.ForageDecision {
	// Find the forageType with the best average multiplier
	bestForageType := bestHistoryForaging(c.forageHistory)

	// If theres no return with multiplier greater than 1 then dont go foraging
	if bestForageType == shared.ForageType(-1) {
		return shared.ForageDecision{
			Type:         shared.FishForageType,
			Contribution: 0,
		}
	}

	// Find the value of resources that gave us the best return and add some
	// noise to it. Cap to 20% of our stockpile
	pastOutcomes := c.forageHistory[bestForageType]
	bestInput := shared.Resources(0)
	bestRoI := shared.Resources(0)

	// For all returns find the best return on investment ((output/input) -1 )
	for _, returns := range pastOutcomes { // Look at the returns of the previous
		if returns.output-returns.input < 10 { // less than 10 profit then continue
			continue
		}
		if returns.input != 0 { // if returns are not 0
			ROI := (returns.output / returns.input) - 1 // Find the input that gave the best return on investment
			if ROI > bestRoI {                          // If the return on investment is better than previous
				bestInput = returns.input // the best value input would be the one that returned the best return on investment
				bestRoI = ROI             // assign best return on investment to be the best value ROI
			}
		}
	}

	// Pick a the minimum value between the best value and 20%
	bestInput = shared.Resources(math.Min(
		float64(bestInput),
		float64(0.2*c.gameState().ClientInfo.Resources)),
	)
	// Add a random amount to the bestInput (max 5%)
	bestInput += shared.Resources(math.Min(
		rand.Float64(),
		float64(0.05*c.gameState().ClientInfo.Resources)),
	)

	// Now return the foraging decision
	forageDecision := shared.ForageDecision{
		Type:         bestForageType,
		Contribution: bestInput,
	}

	// Log the results
	c.Logf(
		"[Forage][Decision]: Decision %v | Expected ROI %v",
		forageDecision, bestRoI)

	return forageDecision
}

// updateForageHistory : Update the foraging history
func (c *client) ForageUpdate(forageDecision shared.ForageDecision, output shared.Resources) {
	c.forageHistory[forageDecision.Type] = append(c.forageHistory[forageDecision.Type], ForageOutcome{
		input:  forageDecision.Contribution,
		output: output,
		turn:   c.gameState().Turn,
	})

	c.Logf(
		"[Forage][Update]: ForageType %v | Profit %v | Contribution %v | Actual ROI %v",
		forageDecision.Type,
		output-forageDecision.Contribution,
		forageDecision.Contribution,
		(output/forageDecision.Contribution)-1,
	)
}

/*  Dying MODE, RISK IT ALL and ask for gifts
lasthopeF put everything in foraging for Deer */
func (c *client) lasthopeF() shared.ForageDecision {
	forageDecision := shared.ForageDecision{
		Type:         shared.DeerForageType,
		Contribution: 0.95 * c.gameState().ClientInfo.Resources,
	}
	c.Logf("[Forage][Decision]: Desperate | Decision %v", forageDecision)
	return forageDecision
}

func (c *client) CommonPoolResourceRequest() shared.Resources {
	switch c.wealth() {
	case Dying:
		c.Logf("Common pool request: 20")
		return 20
	default:
		return 0
	}
}

// func (c *client) RequestAllocation() shared.Resources {
// 	var allocation shared.Resources

// 	if c.wealth() == Dying {
// 		allocation = c.config.desperateStealAmount
// 	} else if c.allocation != 0 {
// 		allocation = c.allocation
// 		c.allocation = 0
// 	}

// 	if allocation != 0 {
// 		c.Logf("Taking %v from common pool", allocation)
// 	}
// 	return allocation
// }

//=============================================================
//=======GIFTS=================================================
//=============================================================

/*
	GetGiftRequests() shared.GiftRequestDict
	GetGiftOffers(receivedRequests shared.GiftRequestDict) shared.GiftOfferDict
	GetGiftResponses(receivedOffers shared.GiftOfferDict) shared.GiftResponseDict
	UpdateGiftInfo(receivedResponses shared.GiftResponseDict)
*/
// func (c* client) GetGiftRequests

// // GiftAccceptance
// func (c *client) GiftAcceptance() shared.GiftResponse {
// 	type AmountofResources int
// 	type AccReason int
// 	if c.wealth() == JeffBezos {
// 		c.Logf("I don't need any gifts you peasants")
// 	} else if c.wealth() == MiddleClass {
// 		c.Logf("We only accept Gifts from Team 1, 2, 4, 6")
// 	} else if c.wealth() == ImperialStudent {
// 		c.Logf("We accept Gifts from everyone")
// 	} else {
// 		c.Logf("We accept Gifts from everyone")
// 	}

// 	GiftResponse := shared.GiftResponse{
// 		AcceptedAmount: AmountofResources,
// 		Reason:         AccReason,
// 	}
// 	return GiftResponse
// }

// /*Gift Requests*/
// func (c *client) GiftRequests() shared.GiftRequest {
// 	var giftrequest shared.GiftRequest

// 	if c.wealth() == JeffBezos {
// 		c.Logf("I don't need any gifts you peasants")
// 	} else if c.wealth() == MiddleClass {
// 		c.Logf("We only accept Gifts from Team 1, 2, 4, 6")
// 	} else if c.wealth() == ImperialStudent {
// 		c.Logf("We accept Gifts from everyone")
// 	} else {
// 		c.Logf("We accept Gifts from everyone")
// 	}

// 	return giftrequest
// }

/********************/
/***    IIFO        */
/********************/

func (c *client) MakeForageInfo() shared.ForageShareInfo {
	var shareTo []shared.ClientID

	for id, status := range c.gameState().ClientLifeStatuses {
		if status != shared.Dead {
			shareTo = append(shareTo, id)
		}
	}

	lastDecisionTurn := -1
	var lastDecision shared.ForageDecision
	var lastRevenue shared.Resources

	for forageType, outcomes := range c.forageHistory {
		for _, outcome := range outcomes {
			if int(outcome.turn) > lastDecisionTurn {
				lastDecisionTurn = int(outcome.turn)
				lastDecision = shared.ForageDecision{
					Type:         forageType,
					Contribution: outcome.input,
				}
				lastRevenue = outcome.output
			}
		}
	}

	if lastDecisionTurn < 0 {
		shareTo = []shared.ClientID{}
	}

	forageInfo := shared.ForageShareInfo{
		ShareTo:          shareTo,
		ResourceObtained: lastRevenue,
		DecisionMade:     lastDecision,
	}

	c.Logf("Sharing forage info: %v", forageInfo)
	return forageInfo
}

func (c *client) ReceiveForageInfo(forageInfos []shared.ForageShareInfo) {
	for _, forageInfo := range forageInfos {
		c.forageHistory[forageInfo.DecisionMade.Type] =
			append(
				c.forageHistory[forageInfo.DecisionMade.Type],
				ForageOutcome{
					input:  forageInfo.DecisionMade.Contribution,
					output: forageInfo.ResourceObtained,
				},
			)
	}
}

/*Foraging History*/
func (c *client) forageHistorySize() uint {
	length := uint(0)
	for _, lst := range c.forageHistory {
		length += uint(len(lst))
	}
	return length // Return how many turns of foraging we have been on depending on the History
}

func (c *client) gameState() gamestate.ClientGameState {
	return c.BaseClient.ServerReadHandle.GetGameState()
}

package team5

import (
	"math"
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

/*
================================================================
	FORAGING
================================================================
	Foraging Functions:
		DecideForage() (shared.ForageDecision, error)
		ForageUpdate(shared.ForageDecision, shared.Resources)

	IIFO Foraging Functions:
		MakeForageInfo() shared.ForageShareInfo
		ReceiveForageInfo([]shared.ForageShareInfo)
================================================================
	TBA
	- Share foraging information with others
	- Look at others foraging information to determine out persuit
	
================================================================
*/

// DecideForage helps us pick the forage (First X turns are the initial setup / data gathering phases)
func (c *client) DecideForage() (shared.ForageDecision, error) {
	if c.forageHistorySize() < c.config.InitialForageTurns { // Start with initial foraging turns
		return c.InitialForage(), nil
	} else if c.wealth() == Dying { // If dying go to last hope
		return c.lastHopeForage(), nil
	} else {
		return c.normalForage(), nil // Else normally forage
	}
}

/*
IntialForage() Gamble the first attempt to try to get rich (born rich)
If we fail we are in middle class and have a 50% chance to go up fast(Deer) or slow (Fish)
If we lose again then we are in the Imperial class and we just fish to try to get back to middle class
*/
func (c *client) InitialForage() shared.ForageDecision {
	var forageType shared.ForageType

	// Default contribution amount is a random amount between 0 -> 10%
	forageContribution := shared.Resources(0.01 + rand.Float64()*(0.05 - 0.01))* c.gameState().ClientInfo.Resources

	switch {
	case c.wealth() == JeffBezos: // JB then we have so much might as well gamble 20% of it
		forageContribution = shared.Resources(0.025 + rand.Float64()*(0.10 - 0.025))* c.gameState().ClientInfo.Resources
		forageType = shared.DeerForageType
	case c.wealth() == ImperialStudent: // Imperial student (Need to save money so dont spent a lot)
		forageType = shared.FishForageType
	case c.wealth() == Dying: // Dying (Its all or nothing now )
		c.lastHopeForage()
	default: // Midle class (lets see where the coin takes us)
		if rand.Float64() < 0.50 {
			forageType = shared.DeerForageType
		} else {
			forageType = shared.FishForageType
		}
	}

	c.Logf("[Debug] - [Initial Forage]:[%v][%v]",forageType,forageContribution)

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

	c.Logf(
		"[Debug] - [Normal Forage]:[%v][%v][Expected Return %v]",
		bestForageType,bestInput,bestRoI)
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
lastHopeForage put everything in foraging for Deer */
func (c *client) lastHopeForage() shared.ForageDecision {
	forageDecision := shared.ForageDecision{
		Type:         shared.DeerForageType,
		Contribution: 0.95 * c.gameState().ClientInfo.Resources,
	}
	c.Logf("[Forage][Decision]: Desperate | Decision %v", forageDecision)
	return forageDecision
}

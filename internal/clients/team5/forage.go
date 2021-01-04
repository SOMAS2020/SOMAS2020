package team5

import (
	"math"
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"gonum.org/v1/gonum/stat/distuv"
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
	- Look at others foraging information to determine our persuit

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
	//=============================================================================
	// find out how much we foraged in the first
	if c.gameState().Turn == 2 {
		var turncost shared.Resources
		for _, outcomes := range c.forageHistory { // For each foraging type find the outcome
			for _, returns := range outcomes { // For each outcome find the returns
				turncost = returns.input // cumlative sum of the return on investment
				c.Logf("DEBUG %v - %v + %v", c.resourceHistory[1], c.resourceHistory[2], returns.input)
			}
		}
		costOfTurn1 := c.resourceHistory[1] - c.resourceHistory[2] + turncost
		c.Logf("[Debug] Cost per turn to live %v = %v - %v + %v", costOfTurn1, c.resourceHistory[1], c.resourceHistory[2], turncost)
	}

	// =============================================================================

	var forageType shared.ForageType

	// Default contribution amount is a random amount between 1% -> 5%
	forageContribution := shared.Resources(0.01+rand.Float64()*(0.05-0.01)) * c.gameState().ClientInfo.Resources

	switch {
	case c.wealth() == JeffBezos: // JB then we have so much might as well gamble 5%->10% of it
		forageContribution = shared.Resources(0.05+rand.Float64()*(0.10-0.05)) * c.gameState().ClientInfo.Resources
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

	c.Logf("[Debug] - [Initial Forage]:[%v][%v]", forageType, forageContribution)

	return shared.ForageDecision{
		Type:         forageType,
		Contribution: forageContribution,
	}
}

// bestForagingType indicates the best average foraging method
func (c *client) bestHistoryForaging(forageHistory ForageHistory) shared.ForageType {
	bestForagingMethod := shared.ForageType(-1)
	bestReturn := 0.0

	for forageType, outcomes := range forageHistory { // For each foraging type
		returnOI := 0.0
		for _, returns := range outcomes {
			returnOI += float64((returns.output / returns.input) - 1) // cumlative sum of the return on investment
		}

		returnOI = returnOI / float64(len(outcomes))

		if returnOI > bestReturn && returnOI > 0 { // Finds the method with the return on investment
			bestReturn = returnOI
			bestForagingMethod = forageType
		}
	}

	// ======================================= Testing seeing how many people went deer hunting
	deerHunters := int(0)
	probDeerHunting := float64(0.5)
	fishHunters := int(0)
	probFishHunting := float64(0)

	for forageType, FOutcome := range forageHistory { // For the whole foraging history
		for _, returns := range FOutcome {
			if forageType == shared.DeerForageType && returns.turn == c.gameState().Turn-1 {
				deerHunters++
				probDeerHunting += 0.1
			} else if forageType == shared.FishForageType && returns.turn == c.gameState().Turn-1 {
				fishHunters++
				probFishHunting += 0.1
			}
		}
	}

	c.Logf("Number of Deer Hunters from pervious turn %v", deerHunters)
	c.Logf("Number of Fish Hunters %v", fishHunters)

	if bestForagingMethod == shared.FishForageType { // Fishing but a lot of people went deer hunting last turn
		bDeer := distuv.Bernoulli{P: 1 - probDeerHunting}     // P = P fishing (1)
		bestForagingMethod *= shared.ForageType(bDeer.Rand()) // Multipy the 0 in or 1
	} else if bestForagingMethod == shared.DeerForageType { // Deer hunting is the best choice
		bFish := distuv.Bernoulli{P: probFishHunting}         // Add the randomness of how many people went fihsin last turn
		bestForagingMethod += shared.ForageType(bFish.Rand()) // + 1 if we did pick fishing above
	}

	// ================================================================

	return bestForagingMethod
}

//normalForage is based on the method that previously gave the best multiplier return
func (c *client) normalForage() shared.ForageDecision {
	// Find the forageType with the best average multiplier
	bestForagingMethod := c.bestHistoryForaging(c.forageHistory)

	// No good returns
	if bestForagingMethod == shared.ForageType(-1) && c.config.SkipForage > 0 {
		c.Logf("[Debug] - Skipping Foraging %v", c.config.SkipForage)
		c.config.SkipForage = c.config.SkipForage - 1
		return shared.ForageDecision{
			Type:         shared.FishForageType,
			Contribution: 0,
		}
	} else if bestForagingMethod == shared.ForageType(-1) && c.config.SkipForage == 0 {
		c.Logf("[Debug] - Force Foraging %v", c.config.SkipForage)
		c.config.SkipForage = 1

		// For now randomly choose either deer or fishing after skipping 3 foraging turns
		var forageMethod shared.ForageType
		forageContribution := shared.Resources(0.05+rand.Float64()*(0.10-0.05)) * c.gameState().ClientInfo.Resources // Invest between 5 to 10%
		if rand.Float64() < 0.50 {
			forageMethod = shared.DeerForageType
		} else {
			forageMethod = shared.FishForageType
		}

		return shared.ForageDecision{
			Type:         forageMethod,
			Contribution: forageContribution,
		}
	}

	// Find the value of resources that gave us the best return and add some
	// noise to it. Cap to 20% of our stockpile
	pastOutcomes := c.forageHistory[bestForagingMethod]
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
		Type:         bestForagingMethod,
		Contribution: bestInput,
	}

	c.Logf(
		"[Debug] - [Forage][Normal]:[%v][%v][Expected Return %v]",
		bestForagingMethod, bestInput, bestRoI)
	return forageDecision
}

/*  Dying MODE, RISK IT ALL and ask for gifts
lastHopeForage put everything in foraging for Deer */
func (c *client) lastHopeForage() shared.ForageDecision {
	forageDecision := shared.ForageDecision{
		Type:         shared.DeerForageType,
		Contribution: 0.95 * c.gameState().ClientInfo.Resources,
	}
	c.Logf("[Debug] - [Forage][LastHopeForage]: Decision %v | Amount %v",
		forageDecision, forageDecision.Contribution)
	return forageDecision
}

// updateForageHistory : Update the foraging history
func (c *client) ForageUpdate(forageDecision shared.ForageDecision, output shared.Resources) {
	c.forageHistory[forageDecision.Type] = append(c.forageHistory[forageDecision.Type], ForageOutcome{
		turn:   c.gameState().Turn,
		input:  forageDecision.Contribution,
		output: output,
	})

	c.Logf(
		"[Debug] - [Update Forage History]: Type %v | Profit %v | Contribution %v | Real RoI %v",
		forageDecision.Type,
		output-forageDecision.Contribution,
		forageDecision.Contribution,
		(output/forageDecision.Contribution)-1,
	)
}

// forageHistorySize gets the size of our history to tell us how many rounds we have foraged
func (c *client) forageHistorySize() uint {
	length := uint(0)
	for _, lst := range c.forageHistory {
		length += uint(len(lst))
	}
	return length // Return how many turns of foraging we have been on depending on the History
}

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
	Work in progress
	- Find out cost of living
================================================================
*/

// DecideForage helps us pick the foraging method
func (c *client) DecideForage() (shared.ForageDecision, error) {
	if c.forageHistorySize() < c.config.InitialForageTurns { // Start with initial foraging turns (semi - randomized)
		return c.InitialForage(), nil
	} else if c.wealth() == dying { // If dying go to last hope
		return c.lastHopeForage(), nil
	}
	return c.normalForage(), nil // else forage normally
}

//================================================================
/*	Foraging types
	Each of the types of foraging are below */
//=================================================================
/* InitialForage() (Risk for wealth if we have wealth or conserve if we dont)
Born in the middle class
Roll randomly
If we succeed we are JB and will risk more to gain more
If we lose again then we are in the Imperial class and we fish to try to get back to middle class
*/
func (c *client) InitialForage() shared.ForageDecision {
	// Figure out how much the cost of livin per turn is
	//=============================================================================
	// if c.getTurn() == 2 { // On turn 2
	// 	var turncost shared.Resources
	// 	for _, outcomes := range c.forageHistory { // For each foraging type find the outcome
	// 		for _, returns := range outcomes { // For each outcome find the returns
	// 			turncost += returns.input // cumlative sum of the return on investment (should only be 1 return from foraging)
	// 			c.Logf("[DEBUG] %v - %v + %v", c.resourceHistory[1], c.resourceHistory[2], returns.input)
	// 		}
	// 	}
	// 	costOfTurn1 := c.resourceHistory[1] - c.resourceHistory[2] + turncost
	// 	c.Logf("[Debug] Cost per turn to live %v = %v - %v + %v", costOfTurn1, c.resourceHistory[1], c.resourceHistory[2], turncost)
	// }
	// =============================================================================
	var forageType shared.ForageType
	forageContribution := shared.Resources(c.config.MinimumForagePercentage+
		rand.Float64()*
			(c.config.NormalForagePercentage-c.config.MinimumForagePercentage)) *
		c.gameState().ClientInfo.Resources // Random amount between Min and Normal contribution
	switch c.wealth() {
	case jeffBezos: // Rich
		forageContribution = shared.Resources(c.config.NormalForagePercentage+
			rand.Float64()*
				(c.config.JBForagePercentage-c.config.NormalForagePercentage)) *
			c.gameState().ClientInfo.Resources // JB then we have so much might as well gamble Normal % -> JB% of it
		forageType = shared.DeerForageType
	case imperialStudent:
		forageType = shared.FishForageType // Save money by fishing and hope that we can get some returns
	case dying:
		c.lastHopeForage() // Invest all our money into fishing to hope we can get some return
	case middleClass: // Midle class (lets see where the coin takes us)
		if rand.Float64() < 0.50 { // Coin
			forageType = shared.DeerForageType
		} else {
			forageType = shared.FishForageType
		}
	}

	c.Logf("[Initial Forage]:[%v][%v]", forageType, forageContribution)

	return shared.ForageDecision{
		Type:         forageType,
		Contribution: forageContribution,
	}
}

// bestForagingType indicates the best foraging method best of RoI (Return on Investment Output/Input - 1)
func (c *client) bestHistoryForaging(forageHistory forageHistory) shared.ForageType {
	bestForagingMethod := shared.ForageType(-1) // Default is that there is no good method
	bestReturn := 0.0

	for forageType, outcomes := range forageHistory { // For each foraging type
		returnOI := 0.0
		for _, returns := range outcomes {
			returnOI += float64((returns.output / returns.input) - 1) // Cumlative sum of the return on investment
		}
		returnOI = returnOI / float64(len(outcomes)) // Average RoI for the type

		if returnOI > bestReturn && returnOI > 0 { // Compares the type to the previous type and 0
			bestReturn = returnOI // If its greater than 0 then it has some return
			bestForagingMethod = forageType
		}
	}
	// We have the best foragine method according to the RoI History

	if bestForagingMethod != shared.ForageType(-1) { // If RoI < 0 then dont bother look at previous hunters
		//=============================================================================
		// Looking at other islands amount of hunters last turn

		probDeerHunting := c.config.RandomChanceToHunt // Base probaility to add some randomness
		probFishing := c.config.RandomChanceToFish     // So we dont always go for the best type

		for forageType, FOutcome := range forageHistory { // For the whole foraging history
			for _, returns := range FOutcome {
				if forageType == shared.DeerForageType && //Deer Hunters
					returns.turn == c.getTurn()-1 && // Last turn
					returns.team != shared.Team5 &&
					returns.input > 0 { // Not including us
					probDeerHunting += c.config.IncreasePerHunterLastTurn // Incremenet the probability we hunt
				} else if forageType == shared.FishForageType &&
					returns.turn == c.getTurn()-1 &&
					returns.team != shared.Team5 &&
					returns.input > 0 {
					probFishing += c.config.IncreasePerFisherMenLastTurn
				}
			}
		}
		// c.Logf("Chance to Hunt %v | Chance to Fish %v", probDeerHunting, probFishing)
		//=============================================================================
		// Check the previous X turns to see how many hunted deer
		// Start only when we have enough turns to look back at
		if c.getTurn() > c.config.DeerTurnsToLookBack {
			prevTurnsHunters := make(map[uint]uint)
			for _, returns := range forageHistory[shared.DeerForageType] { // finds Number of hunters  for each turn
				for i := c.getTurn() - c.config.DeerTurnsToLookBack; i < c.getTurn() && i >= c.getTurn()-c.config.DeerTurnsToLookBack; i++ {
					if returns.output > returns.input*1.2 && returns.turn == i { // Only count the ones that had positive returns (as they must have went hunting)
						prevTurnsHunters[i] = prevTurnsHunters[i] + 1
					}
				}
			}
			for i := c.getTurn() - c.config.DeerTurnsToLookBack; i < c.getTurn(); i++ {
				probDeerHunting -= (float64(prevTurnsHunters[i]) / float64(c.getTurn()-i)) * c.config.DecreasePerHunterInLookBack
			}

			// Logger
			c.Logf("[Best Foraging Method] PrevTurnsHunters %v | Prob Hunting %v | Prob Fishing %v", prevTurnsHunters, probDeerHunting, probFishing)
		}
		// ================================================================
		// If best foraging was none of the 2 above then return shared.ForageType(-1)
		if bestForagingMethod == shared.FishForageType { // Fishing is best but 3 Deer hunters last turn
			probFishing = math.Min(1, 1+probFishing-probDeerHunting)
			bFish := distuv.Bernoulli{P: probFishing}             // P(1)[Fishing]
			bestForagingMethod *= shared.ForageType(bFish.Rand()) // Multiply the 0 in if Deer Hunting was picked in randomness

		} else if bestForagingMethod == shared.DeerForageType { // Deer hunting is the best choice but 3 Fishers
			probDeerHunting = math.Min(1, 1+probDeerHunting-probFishing)
			bDeer := distuv.Bernoulli{P: 1 - probDeerHunting}     // P(1)[Fishing]= 0.1 + 0.1*3 = 0.4
			bestForagingMethod += shared.ForageType(bDeer.Rand()) // +1 [makes it fishing] if Fishing was picked in randomness
		}
		// c.Logf("Chance to Hunt %v | Chance to Fish %v", probDeerHunting, probFishing)
	} // If both methods are less than 0 RoI then return SharedType(-1)
	return bestForagingMethod
}

/* normalForage() (Past the initial, based on the history of our foraging and some randomness) */
func (c *client) normalForage() shared.ForageDecision {
	bestForagingMethod := c.bestHistoryForaging(c.forageHistory) // Find the best foragine type in based on history

	// No good returns all our history had RoI < 0 , Skip one turn
	//=============================================================================
	if bestForagingMethod == shared.ForageType(-1) && c.config.SkipForage > 0 {
		c.Logf("Skipping Foraging: %v", c.config.SkipForage)
		c.config.SkipForage--         // Count down the number of turns to skip
		return shared.ForageDecision{ // Dont go foraging
			Type:         shared.DeerForageType,
			Contribution: 0.01,
		}
	} else if bestForagingMethod == shared.ForageType(-1) && c.config.SkipForage == 0 { // Force Foraging
		c.Logf("Force Foraging: %v", c.config.SkipForage)
		c.config.SkipForage = 1 // Reassign the number of skips for next time we have RoI < 0

		// Randomly pick type and invest 1->3%
		var forageMethod shared.ForageType
		if rand.Float64() < 0.50 {
			forageMethod = shared.DeerForageType
		} else {
			forageMethod = shared.FishForageType // Maybe pick only fishing?
		}
		// Between 1->5%
		forageContribution := shared.Resources(c.config.MinimumForagePercentage+
			rand.Float64()*
				(c.config.NormalForagePercentage-c.config.MinimumForagePercentage)) *
			c.gameState().ClientInfo.Resources
		forageContribution = forageContribution * 2 // Double the amount we invested (possibly investing too little and we need returns)
		return shared.ForageDecision{
			Type:         forageMethod,
			Contribution: forageContribution,
		}
	}
	//=============================================================================

	// Foraging with previous history thats not -1
	pastOutcomes := c.forageHistory[bestForagingMethod]
	bestInput := shared.Resources(0)
	bestRoI := shared.Resources(0)

	// For all returns find the best return on investment ((output/input) -1 )
	for _, returns := range pastOutcomes { // Look at the returns of the previous
		if returns.input > 0 { // If returns are not 0
			RoI := (returns.output / returns.input) - 1 // Find the input that gave the best RoI
			if RoI > bestRoI {                          // RoI better than previous
				bestInput = returns.input // best amount to invest
				bestRoI = RoI             // best RoI so far
			}
		}
	}

	// Add a random amount -5% -> 5% to the bestInput (max +-X%)
	if rand.Float64() < 0.5 { // Increase or Decrease
		bestInput -= bestInput * shared.Resources(rand.Float64()*c.config.NormalRandomChange)
	} else {
		bestInput += bestInput * shared.Resources(rand.Float64()*c.config.NormalRandomChange)
	}
	bestInput = bestInput * shared.Resources(c.mapToRange(float64(len(c.getAliveTeams(true))), 6, 1, 1, 2))

	// Pick the minimum value between the best value and x% of our resources
	bestInput = shared.Resources(math.Min(
		float64(bestInput),
		float64(shared.Resources(c.config.MaxForagePercentage)*c.gameState().ClientInfo.Resources)),
	)

	// Now return the foraging decision
	forageDecision := shared.ForageDecision{
		Type:         bestForagingMethod,
		Contribution: bestInput,
	}

	c.Logf(
		"[Normal Forage]:Method: %v | Input: %v | Expected RoI: %v",
		bestForagingMethod, bestInput, bestRoI)

	return forageDecision
}

/*  dying MODE, RISK IT ALL, put everything in foraging for Fishing */
func (c *client) lastHopeForage() shared.ForageDecision {
	forageDecision := shared.ForageDecision{
		Type:         shared.FishForageType,
		Contribution: 0.95 * c.gameState().ClientInfo.Resources, // Almost everything we still want to be > 0 in case 0 means insta death
	}
	c.Logf("[Forage][Just let me die]: Decision %v | Amount %v (please)?",
		forageDecision, forageDecision.Contribution)
	return forageDecision
}

//================================================================
/*	Foraging History Functions */
//=================================================================

//ForageUpdate Updates the foraging history
func (c *client) ForageUpdate(forageDecision shared.ForageDecision, output shared.Resources, numberCaught uint) {
	c.forageHistory[forageDecision.Type] = append(c.forageHistory[forageDecision.Type], forageOutcome{ // Append new data
		team:   shared.Team5,
		turn:   c.getTurn(),
		input:  forageDecision.Contribution,
		output: output,
		caught: numberCaught,
	})
	c.Logf(
		"[Update Forage History]:[%v] Type %v | Input %v | Output %v | No.Caught %v | Real RoI %v",
		c.getTurn(),
		forageDecision.Type,
		forageDecision.Contribution,
		output-forageDecision.Contribution,
		numberCaught,
		(output/forageDecision.Contribution)-1,
	)
}

// forageHistorySize gets the size of our history to tell us how many rounds we have foraged for
func (c *client) forageHistorySize() uint {
	length := uint(0)
	for _, lst := range c.forageHistory {
		length += uint(len(lst))
	}
	return length // Return how many turns of foraging we have been on depending on the History
}

//======================= Part of IIFO ====================================

//ReceiveForageInfo get info from other teams
func (c *client) ReceiveForageInfo(forageInfos []shared.ForageShareInfo) {
	for _, forageInfo := range forageInfos { // for all foraging information from all islands (ignore the islands)
		c.forageHistory[forageInfo.DecisionMade.Type] = // all their information (based on method of foraging)
			append( // add to our history
				c.forageHistory[forageInfo.DecisionMade.Type], // Type of foraging
				forageOutcome{ // Outcome of their foraging
					team:   forageInfo.SharedFrom,
					turn:   c.getTurn(),                          // The current turn
					input:  forageInfo.DecisionMade.Contribution, // Contribution
					output: forageInfo.ResourceObtained,          // Resource obtained
				},
			)
	}

	c.Logf("[Forage Info From] %v", forageInfos)

	for _, forageInfo := range forageInfos {
		if forageInfo.DecisionMade.Contribution >= 1 { // has to be meaningful forage
			c.opinions[forageInfo.SharedFrom].updateOpinion(generalBasis, +0.2*c.getMood()) // Thanks for the information dude very usefull
		}
	}
}

//MakeForageInfo
func (c *client) MakeForageInfo() shared.ForageShareInfo {
	var shareTo []shared.ClientID
	if c.getTurn() > c.config.InitialForageTurns { // for the turns we are NOT doing initial forage
		for _, FOutcome := range c.forageHistory {
			for _, returns := range FOutcome {
				for _, team := range c.getAliveTeams(false) { // For all alive teams
					if c.getTurn() > c.config.DeerTurnsToLookBack { // Turns greater than look back
						if returns.turn > c.getTurn()-c.config.DeerTurnsToLookBack &&
							returns.turn <= c.getTurn() &&
							returns.team == team { // If a certain team within a certain range
							shareTo = append(shareTo, team) // add to shrae to list if they shared to us
						} else {
							c.opinions[team].updateOpinion(generalBasis, -0.05*c.getMood()) // booo give me your foraging data
						}
					} else {
						startOfGame := uint(0) // case that the deerTurnsTo look back is greater than turns
						if returns.turn > startOfGame &&
							returns.turn <= c.getTurn() &&
							returns.team == team {
							shareTo = append(shareTo, team)
						} else {
							c.opinions[team].updateOpinion(generalBasis, -0.05*c.getMood()) // booo give me your foraging data
						}
					}
				}
			}
		}
		// Delete the repeated teams in the shareTO list
		keys := make(map[shared.ClientID]bool)
		shareToPrevShare := []shared.ClientID{}
		for _, entry := range shareTo {
			if _, value := keys[entry]; !value {
				keys[entry] = true
				shareToPrevShare = append(shareToPrevShare, entry)
			}
		}
		shareTo = shareToPrevShare

	} else if c.getTurn() > 1 { // share info for all turns in initial forage
		for team, status := range c.gameState().ClientLifeStatuses { // Check the clients that are alive
			if status != shared.Dead { // if they are not dead then append the shareTo,id
				shareTo = append(shareTo, team)
			}
		}
	}

	lastTurn := c.getTurn()
	if c.getTurn() > 1 {
		lastTurn--
	}
	var contribution shared.ForageDecision
	var output shared.Resources
	for forageType, outcomes := range c.forageHistory { //For each type look at the outcome
		for _, outcome := range outcomes {
			if outcome.turn == lastTurn { // If the turn is the same as the last turn then return the result
				output = outcome.output               // output of the outcome
				contribution = shared.ForageDecision{ // Foraging Decision
					Type:         forageType,
					Contribution: outcome.input,
				}
			}
		}
	}

	forageInfo := shared.ForageShareInfo{
		DecisionMade:     contribution, // contribution and Resources obtained
		ResourceObtained: output,       // How much we got back
		ShareTo:          shareTo,      // []shared.ClientIDs
		SharedFrom:       shared.Team5,
	}

	c.Logf("[Sharing Foraging Info] %v", forageInfo)
	return forageInfo
}

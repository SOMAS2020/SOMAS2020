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

// ForageOutcome records the ROI on a foraging session
type ForageOutcome struct {
	turn    uint
	input   shared.Resources
	utility shared.Resources
}

// ForageHistory stores history of foraging outcomes
type ForageHistory map[shared.ForageType][]ForageOutcome

type InternalState int

const (
	Normal InternalState = iota
	Anxious
	Desperate
)

func (st InternalState) String() string {
	strings := [...]string{"Normal", "Anxious", "Desperate"}
	if st >= 0 && int(st) < len(strings) {
		return strings[st]
	}
	return fmt.Sprintf("Unkown internal state '%v'", int(st))
}

func init() {
	baseclient.RegisterClient(
		id,
		&client{
			BaseClient:    baseclient.NewClient(id),
			forageHistory: ForageHistory{},
		},
	)
}

type clientConfig struct {
	// At the start of the game forage randomly for this many turns. If true,
	// pay some initial tax to help get the first IIGO running
	randomForageTurns uint

	// If resources go below this limit, go into "desperation" mode
	anxietyThreshold shared.Resources

	// If true, ignore requests for taxes
	evadeTaxes bool

	// If true, pay some initial tax to help get the first IIGO running
	kickstartTaxPercent shared.Resources

	// desperateStealAmount is the amount the agent will steal from the commonPool
	desperateStealAmount shared.Resources
}

// client is Pittunia
type client struct {
	*baseclient.BaseClient

	forageHistory ForageHistory
	taxAmount     shared.Resources

	// allocation is the president's response to your last common pool resource request
	allocation shared.Resources

	config clientConfig
}

func NewClient(clientID shared.ClientID) baseclient.Client {
	return &client{
		BaseClient:    baseclient.NewClient(clientID),
		forageHistory: ForageHistory{},
		taxAmount:     0,
		allocation:    0,
		config: clientConfig{
			randomForageTurns:    5,
			anxietyThreshold:     20,
			desperateStealAmount: 30,
			evadeTaxes:           false,
			kickstartTaxPercent:  0.25,
		},
	}
}

func (c client) internalState() InternalState {
	ci := c.gameState().ClientInfo
	switch {
	case ci.LifeStatus == shared.Critical:
		return Desperate
	case ci.Resources < c.config.anxietyThreshold:
		return Anxious
	default:
		return Normal
	}
}

func (c *client) StartOfTurn() {
	c.Logf("Emotional state: %v", c.internalState())
	c.Logf("Resources: %v", c.gameState().ClientInfo.Resources)

	for clientID, status := range c.gameState().ClientLifeStatuses {
		if status != shared.Dead && clientID != c.GetID() {
			return
		}
	}
	c.Logf("I'm all alone :c")

}

/********************/
/***  Desperation   */
/********************/

func (c *client) desperateForage() shared.ForageDecision {
	forageDecision := shared.ForageDecision{
		Type:         shared.FishForageType,
		Contribution: c.gameState().ClientInfo.Resources,
	}
	c.Logf("[Forage][Decision]: Desperate | Decision %v", forageDecision)
	return forageDecision
}

func (c *client) CommonPoolResourceRequest() shared.Resources {
	switch c.internalState() {
	case Anxious:
		c.Logf("Common pool request: 20")
		return 20
	default:
		return 0
	}
}

func (c *client) RequestAllocation() shared.Resources {
	var allocation shared.Resources

	if c.internalState() == Desperate {
		allocation = c.config.desperateStealAmount
	} else if c.allocation != 0 {
		allocation = c.allocation
		c.allocation = 0
	}

	if allocation != 0 {
		c.Logf("Taking %v from common pool", allocation)
	}
	return allocation
}

/********************/
/***  Messages      */
/********************/

func (c *client) ReceiveCommunication(
	sender shared.ClientID,
	data map[shared.CommunicationFieldName]shared.CommunicationContent,
) {
	for field, content := range data {
		switch field {
		case shared.TaxAmount:
			// if !__isPresident__(sender) {
			// }
			c.taxAmount = shared.Resources(content.IntegerData)
		case shared.AllocationAmount:
			// if !__isPresident__(sender) {
			// }
			c.allocation = shared.Resources(content.IntegerData)
		}
	}
}

func (c *client) GetTaxContribution() shared.Resources {
	if c.config.evadeTaxes {
		return 0
	}
	if c.gameState().Turn == 1 {
		// Put in some initial resources to get IIGO to run
		return c.config.kickstartTaxPercent * c.gameState().ClientInfo.Resources
	}
	return c.taxAmount
}

/********************/
/***    Foraging    */
/********************/

func (c *client) randomForage() shared.ForageDecision {
	// Up to 10% of our current resources
	forageContribution := shared.Resources(0.1*rand.Float64()) * c.gameState().ClientInfo.Resources

	var forageType shared.ForageType
	if rand.Float64() < 0.5 {
		forageType = shared.DeerForageType
	} else {
		forageType = shared.FishForageType
	}

	c.Logf("[Forage][Decision]:Random")
	// TODO Add fish foraging when it's done
	return shared.ForageDecision{
		Type:         forageType,
		Contribution: forageContribution,
	}
}

func (c *client) normalForage() shared.ForageDecision {
	// Find the forageType with the best average returns
	bestForageType := shared.ForageType(-1)
	bestROI := 0.0

	// TODO: We could also do this in one go. So just repeat the exact decision
	// that gave the best ROI, instead of picking forageType based on its
	// average ROI

	for forageType, outcomes := range c.forageHistory {
		ROIsum := 0.0
		for _, outcome := range outcomes {
			if outcome.input != 0 {
				ROIsum += float64(outcome.utility/outcome.input) - 1
			}
		}

		averageROI := ROIsum / float64(len(outcomes))

		if averageROI > bestROI {
			bestROI = averageROI
			bestForageType = forageType
		}
	}

	// Not foraging is best
	if bestForageType == shared.ForageType(-1) {
		return shared.ForageDecision{
			Type:         shared.FishForageType,
			Contribution: 0,
		}
	}

	// Find the value of resources that gave us the best return and add some
	// noise to it. Cap to 20% of our stockpile
	pastOutcomes := c.forageHistory[bestForageType]
	bestValue := shared.Resources(0)
	bestValueROI := shared.Resources(0)
	for _, outcome := range pastOutcomes {
		if outcome.utility-outcome.input < 10 {
			continue
		}
		if outcome.input != 0 {
			ROI := (outcome.utility / outcome.input) - 1
			if ROI > bestValueROI {
				bestValue = outcome.input
				bestValueROI = ROI
			}
		}
	}
	bestValue = shared.Resources(math.Min(
		float64(bestValue),
		float64(0.2*c.gameState().ClientInfo.Resources)),
	)
	bestValue += shared.Resources(math.Min(
		rand.Float64(),
		float64(0.01*c.gameState().ClientInfo.Resources),
	))

	forageDecision := shared.ForageDecision{
		Type:         bestForageType,
		Contribution: bestValue,
	}
	c.Logf(
		"[Forage][Decision]: Decision %v | Expected ROI %v",
		forageDecision, bestROI)

	return forageDecision
}

func (c *client) DecideForage() (shared.ForageDecision, error) {
	if c.forageHistorySize() < c.config.randomForageTurns {
		return c.randomForage(), nil
	} else if c.internalState() == Desperate {
		return c.desperateForage(), nil
	} else {
		return c.normalForage(), nil
	}
}

func (c *client) ForageUpdate(forageDecision shared.ForageDecision, utility shared.Resources) {
	c.forageHistory[forageDecision.Type] = append(c.forageHistory[forageDecision.Type], ForageOutcome{
		input:   forageDecision.Contribution,
		utility: utility,
		turn:    c.gameState().Turn,
	})

	c.Logf(
		"[Forage][Update]: ForageType %v | Profit %v | Contribution %v | Actual ROI %v",
		forageDecision.Type,
		utility-forageDecision.Contribution,
		forageDecision.Contribution,
		(utility/forageDecision.Contribution)-1,
	)
}

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
				lastRevenue = outcome.utility
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
					input:   forageInfo.DecisionMade.Contribution,
					utility: forageInfo.ResourceObtained,
				},
			)
	}
}

/********************/
/***    Helpers     */
/********************/

func (c *client) forageHistorySize() uint {
	length := uint(0)
	for _, lst := range c.forageHistory {
		length += uint(len(lst))
	}
	return length
}

func (c *client) gameState() gamestate.ClientGameState {
	return c.BaseClient.ServerReadHandle.GetGameState()
}

package team2

import (
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// no. teams that we know (they have shared the info with us)
// total from all the teams we know
type ForagingResults struct {
	Hunters int
	Fishers int
	Result  shared.Resources
}

// Returns a forage decision and an error that includes our intended foraging type and foraging resource allocation
func (c *client) DecideForage() (shared.ForageDecision, error) {
	// implement a normal distribution which shifts closer to hunt or fish
	// number from 0 to 1
	var threshold float64 = c.decideHuntingLikelihood()
	forageDecision := shared.FishForageType

	if rand.Float64() > threshold {
		// we fish when above the threshold
		forageDecision = shared.FishForageType
	} else {
		// we hunt when below the threshold
		forageDecision = shared.DeerForageType
	}

	// Update last forage type and amount
	c.lastForageType = forageDecision
	c.lastForageAmount = shared.Resources(c.DecideForageAmount(threshold))

	return shared.ForageDecision{
		Type:         forageDecision,
		Contribution: shared.Resources(c.DecideForageAmount(threshold)),
	}, nil
}

// Returns decision on amount of resources to put into foraging
func (c *client) DecideForageAmount(foragingDecisionThreshold float64) shared.Resources {
	// Note: we have given to the pool already by this point in the turn so act cautiously if we're in danger
	if c.criticalStatus() || c.getAgentExcessResources() == 0 {
		return shared.Resources(0)
	}

	// Otherwise we have excess to allocate to foraging

	resourcesForForaging := c.getAgentExcessResources()

	return shared.Resources(resourcesForForaging)
}

// being the only agent to hunt is undesirable, having one hunting partner is the desirable, the more hunters after that the less we want to hunt
// will move the threshold, higher value means more likely to hunt
func (c *client) decideHuntingLikelihood() float64 {
	// Todo: this could be null
	hunters := c.otherHunters()
	multiplier := 1.0

	// We want to shift our likelihood of hunting up if it likely hunting will result in higher utility
	switch c.getAgentStrategy() {
	case Selfish:
		multiplier = 1.2
	case Altruist:
		multiplier = 1.1
	default:
		multiplier = 1.15
	}

	// ideal 2 hunters, 4 fishers

	// in the case when one other person only is hunting
	if hunters > 1.0 && hunters < 1.2 {
		return 0.8 * multiplier
	} else if hunters >= 1.2 {
		// If no one is likely to hunt then we do default probability
		// default hunt probability is 10%, the less people hunting the more likely we do it
		return (0.95 - (hunters * 0.15)) * multiplier
	}

	//when no one is likely to hunt we have a default 10% chance of hunting just in the off chance another person hunts
	return 0.1 * multiplier
}

// Returns the expected number of other hunters based on previous foraging behaviour
func (c *client) otherHunters() float64 {
	aliveClients := c.getAliveClients()
	HuntNum := 0.00

	for _, id := range aliveClients {
		if islandForageHist, ok := c.foragingReturnsHist[id]; ok {
			for _, forageInfo := range islandForageHist {
				if len(islandForageHist) != 0 {
					HuntNum += float64(forageInfo.DecisionMade.Type) / float64(len(islandForageHist))
				}
			}
		}
	}

	return HuntNum
}

// ForageUpdate is called by the server upon completion of a foraging session. This handler can be used by clients to
// analyse their returns - resources returned to them, as well as number of fish/deer caught.
func (c *client) ForageUpdate(initialDecision shared.ForageDecision, resourceReturn shared.Resources, numberCaught uint) {
	ourForageInfo := ForageInfo{
		DecisionMade:      initialDecision,
		ResourcesObtained: resourceReturn,
	}
	c.foragingReturnsHist[c.GetID()] = append(c.foragingReturnsHist[c.GetID()], ourForageInfo)
}

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

	if val, ok := c.foragingReturnsHist[c.GetID()]; ok {
		if len(val) != 0 {
			contribution := shared.ForageDecision{Type: c.lastForageType, Contribution: c.lastForageAmount}
			return shared.ForageShareInfo{
				DecisionMade:     contribution,
				ResourceObtained: c.foragingReturnsHist[c.GetID()][len(c.foragingReturnsHist[c.GetID()])-1].ResourcesObtained,
				ShareTo:          c.getIslandsToShareWith(),
			}
		}
	}

	noInfo := shared.ForageShareInfo{
		DecisionMade:     shared.ForageDecision{},
		ResourceObtained: shared.Resources(0),
		ShareTo:          []shared.ClientID{},
	}

	return noInfo
}

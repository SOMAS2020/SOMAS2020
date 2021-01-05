package team6

import (
	"math/rand"
	"sort"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func chooseForageType() int {
	return int(shared.DeerForageType)
}

func (c client) changeForageType() shared.ForageType { //TODO
	return shared.DeerForageType
}

func (c client) decideContribution() shared.Resources {
	var multiplier shared.Resources = 0.8
	var safetyBuffer shared.Resources = 10
	ourResources := c.ServerReadHandle.GetGameState().ClientInfo.Resources
	return multiplier * (ourResources - safetyBuffer)
}

func (c *client) randomForage() shared.ForageDecision {
	var resources shared.Resources
	var forageType shared.ForageType

	if c.ServerReadHandle.GetGameState().Turn == 2 {
		forageType = shared.FishForageType
	} else {
		forageType = shared.DeerForageType
	}
	tmp := rand.Float64()
	if tmp > 0.2 { //up to 20% resources
		resources = 0.2 * c.ServerReadHandle.GetGameState().ClientInfo.Resources
	} else {
		resources = shared.Resources(tmp) * c.ServerReadHandle.GetGameState().ClientInfo.Resources
	}

	return shared.ForageDecision{
		Type:         shared.ForageType(forageType),
		Contribution: shared.Resources(resources),
	}
}

func (c *client) noramlForage() shared.ForageDecision {
	ft := c.changeForageType()
	amt := c.decideContribution()
	return shared.ForageDecision{
		Type:         shared.ForageType(ft),
		Contribution: shared.Resources(amt),
	}
}

func (c *client) DecideForage() (shared.ForageDecision, error) {
	if c.ServerReadHandle.GetGameState().Turn < 3 { //the agent will randomly forage at the start
		return c.randomForage(), nil
	} else {
		return c.noramlForage(), nil
	}

}

func (c *client) ForageUpdate(forageDecision shared.ForageDecision, outcome shared.Resources) {
	c.forageHistory[forageDecision.Type] =
		append(
			c.forageHistory[forageDecision.Type],
			ForageResults{
				forageIn:     forageDecision.Contribution,
				forageReturn: outcome,
				forageROI:    float64(forageDecision.Contribution/outcome - 1),
			},
		)
}

// Pair is used for sorting the friendship
type Pair struct {
	Key   shared.ClientID
	Value uint
}

// PairList is a slice of Pairs that implements sort.Interface to sort by Value.
type PairList []Pair

func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

// A function to turn a map into a PairList, then sort and return it.
func sortMapByValue(m map[shared.ClientID]uint) PairList {
	p := make(PairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = Pair{k, v}
	}
	sort.Sort(p)
	return p
}

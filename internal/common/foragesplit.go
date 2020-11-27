package common

type forageSplitRuleFuncType = func(teamForageInvestments map[ClientID]uint, totalPayoff uint) (payoff map[ClientID]uint)

// ForageSplitRuleKey is a key that indexes forageSplitRuleMap for the correct split rule func to use
type ForageSplitRuleKey int

const (
	// ForageSplitEven : splits evenly
	ForageSplitEven ForageSplitRuleKey = iota
	// ForageSplitProportionate : splits proportionately
	ForageSplitProportionate
)

// DefaultForageSplitRuleKey denotes the default split rule used.
const DefaultForageSplitRuleKey = ForageSplitEven

// ComputeForageSplit computes the required split for the set rule and forage situation
func ComputeForageSplit(rule ForageSplitRuleKey, teamForageInvestments map[ClientID]uint, totalPayoff uint) (payoff map[ClientID]uint) {
	f := forageSplitRuleMap[rule]
	payoff = f(teamForageInvestments, totalPayoff)
	return
}

// forageSplitRuleMap maps from the forage split enum to the forage split function implementing
// the split.
var forageSplitRuleMap = map[ForageSplitRuleKey]forageSplitRuleFuncType{
	ForageSplitEven:          forageSplitEven,
	ForageSplitProportionate: forageSplitProprotionate,
}

// forageSplitEven evenly splits the total payoff between all teams.
// forageInfo maps from id of foraging team to resources invested to forage.
// If totalPayoff is not completely divisible by the number of teams, then the remained resources are discarded.
func forageSplitEven(teamForageInvestments map[ClientID]uint, totalPayoff uint) (payoff map[ClientID]uint) {
	payoff = map[ClientID]uint{}
	payoffPerTeam := uint(totalPayoff / uint(len(teamForageInvestments)))
	for c := range teamForageInvestments {
		payoff[c] = payoffPerTeam
	}
	return
}

// forageSplitProportionate splits the total payoff based the proportionate amount of resources
// invested into foraging by each team.
// Values are rounded down if they don't divide completely.
func forageSplitProprotionate(teamForageInvestments map[ClientID]uint, totalPayoff uint) (payoff map[ClientID]uint) {
	payoff = map[ClientID]uint{}
	totalInvestments := uint(0)
	for _, inv := range teamForageInvestments {
		totalInvestments += inv
	}

	for c, inv := range teamForageInvestments {
		payoffC := uint(0)
		if totalInvestments != 0 {
			frac := float64(inv) / float64(totalInvestments)
			if frac != 0 {
				payoffC = uint(float64(totalPayoff) * frac)
			}
		}
		payoff[c] = payoffC
	}
	return
}

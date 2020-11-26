package forage

type splitRuleFuncType = func(teamForageInvestments map[int]int, totalPayoff int) (payoff map[int]int)

// SplitRuleKey is a key that indexes SplitRuleMap for the correct split rule func to use
type SplitRuleKey int

const (
	// SplitEven : splits evenly
	SplitEven SplitRuleKey = iota
	// SplitProportionate : splits proportionately
	SplitProportionate
)

// DefaultSplitRuleKey denotes the default split rule used.
const DefaultSplitRuleKey = SplitEven

// ComputeSplit computes the required split for the set rule and forage situation
func ComputeSplit(rule SplitRuleKey, teamForageInvestments map[int]int, totalPayoff int) (payoff map[int]int) {
	f := splitRuleMap[rule]
	payoff = f(teamForageInvestments, totalPayoff)
	return
}

// splitRuleMap maps from the forage split enum to the forage split function implementing
// the split.
var splitRuleMap = map[SplitRuleKey]splitRuleFuncType{
	SplitEven:          forageSplitEven,
	SplitProportionate: forageSplitProprotionate,
}

// forageSplitEven evenly splits the total payoff between all teams.
// forageInfo maps from id of foraging team to resources invested to forage.
// If totalPayoff is not completely divisible by the number of teams, then the remained resources are discarded.
func forageSplitEven(teamForageInvestments map[int]int, totalPayoff int) (payoff map[int]int) {
	payoff = map[int]int{}
	payoffPerTeam := int(totalPayoff / len(teamForageInvestments))
	for c := range teamForageInvestments {
		payoff[c] = payoffPerTeam
	}
	return
}

// forageSplitProportionate splits the total payoff based the proportionate amount of resources
// invested into foraging by each team.
// Values are rounded down if they don't divide completely.
func forageSplitProprotionate(teamForageInvestments map[int]int, totalPayoff int) (payoff map[int]int) {
	payoff = map[int]int{}
	totalInvestments := 0
	for _, inv := range teamForageInvestments {
		totalInvestments += inv
	}

	for c, inv := range teamForageInvestments {
		payoffC := 0
		if totalInvestments != 0 {
			frac := float64(inv) / float64(totalInvestments)
			if frac != 0 {
				payoffC = int(float64(totalPayoff) * frac)
			}
		}
		payoff[c] = payoffC
	}
	return
}

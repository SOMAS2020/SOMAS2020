package common

type forageSplitRuleFuncType = func(teamForageInvestments map[int]int, totalPayoff int) (payoff map[int]int)

// ForageSplitRuleMap maps from the forage split name to the forage split function implementing
// the spli .
var ForageSplitRuleMap = map[string]forageSplitRuleFuncType{
	"even":          forageSplitEven,
	"proportionate": forageSplitProprotionate,
}

// forageSplitEven evenly splits the total payoff between all teams.
// forageInfo maps from id of foraging team to resources invested to forage.
// If totalPayoff is not completely divisible by the number of teams, then the remained resources are discarded.
func forageSplitEven(teamForageInvestments map[int]int, totalPayoff int) (payoff map[int]int) {
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
	totalInvestments := 0
	for _, inv := range teamForageInvestments {
		totalInvestments += inv
	}

	for c, inv := range teamForageInvestments {
		frac := int(inv / totalInvestments)
		payoff[c] = totalPayoff / frac
	}
	return
}

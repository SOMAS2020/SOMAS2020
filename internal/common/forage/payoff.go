package forage

type payoffRuleFuncType = func(input int) (output int)

// DefaultPayoffRuleKey denotes the default payoff for an instance of foraging
const DefaultPayoffRuleKey = PayoffSquare

// PayoffRuleKey is a key that indexes PayoffRuleMap for the correct payoff rule func to use
type PayoffRuleKey int

const (
	// PayoffSquare : pays off square of input
	PayoffSquare PayoffRuleKey = iota
	// PayoffCube : pays off cube of input
	PayoffCube
)

// ComputePayoff computes the required payoff for the given rule and forage situation
func ComputePayoff(rule PayoffRuleKey, input int) (output int) {
	f := payoffRuleMap[rule]
	output = f(input)
	return
}

// payoffRuleMap maps from the forage payoff rule enum to the forage payoff
// function implementing the payoff function.
var payoffRuleMap = map[PayoffRuleKey]payoffRuleFuncType{
	PayoffSquare: foragePayoffSquare,
	PayoffCube:   foragePayoffCube,
}

func foragePayoffSquare(input int) (output int) {
	output = input * input
	return
}

func foragePayoffCube(input int) (output int) {
	output = input * input * input
	return
}

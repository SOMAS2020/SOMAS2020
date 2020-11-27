package common

type foragePayoffRuleFuncType = func(input uint) (output uint)

// DefaultForagePayoffRuleKey denotes the default payoff for an instance of foraging
const DefaultForagePayoffRuleKey = ForagePayoffSquare

// ForagePayoffRuleKey is a key that indexes foragePayoffRuleMap for the correct payoff rule func to use
type ForagePayoffRuleKey int

const (
	// ForagePayoffSquare : pays off square of input
	ForagePayoffSquare ForagePayoffRuleKey = iota
	// ForagePayoffCube : pays off cube of input
	ForagePayoffCube
)

// ComputeForagePayoff computes the required payoff for the given rule and forage situation
func ComputeForagePayoff(rule ForagePayoffRuleKey, input uint) (output uint) {
	f := foragePayoffRuleMap[rule]
	output = f(input)
	return
}

// foragePayoffRuleMap maps from the forage payoff rule enum to the forage payoff
// function implementing the payoff function.
var foragePayoffRuleMap = map[ForagePayoffRuleKey]foragePayoffRuleFuncType{
	ForagePayoffSquare: foragePayoffSquare,
	ForagePayoffCube:   foragePayoffCube,
}

func foragePayoffSquare(input uint) (output uint) {
	output = input * input
	return
}

func foragePayoffCube(input uint) (output uint) {
	output = input * input * input
	return
}

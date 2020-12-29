package shared

// ResourceDistributionStrategy characterises the different ways in which resources can be distributed
// amongst concerned parties. This is required in splitting foraging returns, for example.
type ResourceDistributionStrategy int

const (
	// Equal split. Ignores input contributions
	Equal ResourceDistributionStrategy = iota
	// InputProportional splits proportionally to input contributions
	InputProportional
	// RankProportional splits according to some other external rank specification
	RankProportional
)

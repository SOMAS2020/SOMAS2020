package shared

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/pkg/miscutils"
	"github.com/pkg/errors"
)

// ResourceDistributionStrategy characterises the different ways in which resources can be distributed
// amongst concerned parties. This is required in splitting foraging returns, for example.
type ResourceDistributionStrategy int

const (
	// EqualSplit ignores input contributions
	EqualSplit ResourceDistributionStrategy = iota
	// InputProportionalSplit splits proportional to input contributions
	InputProportionalSplit
	// RankProportionalSplit splits according to some other external rank specification
	RankProportionalSplit

	// don't modify this
	_resourceDistrEnd
)

func (rd ResourceDistributionStrategy) String() string {
	strings := [...]string{"EqualSplit", "InputProportionalSplit", "RankProportionalSplit"}
	if rd >= 0 && int(rd) < len(strings) {
		return strings[rd]
	}
	return fmt.Sprintf("UNKNOWN ResourceDistributionStrategy '%v'", int(rd))
}

// GoString implements GoStringer
func (rd ResourceDistributionStrategy) GoString() string {
	return rd.String()
}

// MarshalText implements TextMarshaler
func (rd ResourceDistributionStrategy) MarshalText() ([]byte, error) {
	return miscutils.MarshalTextForString(rd.String())
}

// MarshalJSON implements RawMessage
func (rd ResourceDistributionStrategy) MarshalJSON() ([]byte, error) {
	return miscutils.MarshalJSONForString(rd.String())
}

// ParseResourceDistributionStrategy gets the ResourceDistributionStrategy based on an iota index
func ParseResourceDistributionStrategy(x int) (ResourceDistributionStrategy, error) {
	if x >= 0 && ResourceDistributionStrategy(x) < _resourceDistrEnd {
		return ResourceDistributionStrategy(x), nil
	}
	return EqualSplit, errors.Errorf("Unknown ResourceDistribution Strategy specified: '%v'.", x)
}

// HelpResourceDistributionStrategy returns a help string for ResourceDistributionStrategy
func HelpResourceDistributionStrategy() string {
	help := "Determine basis on which group resources are split amongst concerned parties\n"

	for i := 0; i < int(_resourceDistrEnd); i++ {
		help += fmt.Sprintf("%v: %v\n", i, ResourceDistributionStrategy(i))
	}

	return help
}

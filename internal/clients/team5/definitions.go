package team5

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/SOMAS2020/SOMAS2020/pkg/miscutils"
)

const id = shared.Team5

// WealthTier defines how much money we have
type WealthTier int

// ForageOutcome records the ROI on a foraging session
type ForageOutcome struct {
	turn   uint
	input  shared.Resources
	output shared.Resources
}

// ForageHistory stores history of foraging outcomes
type ForageHistory map[shared.ForageType][]ForageOutcome

type clientConfig struct {
	// Initial non planned foraging
	InitialForageTurns uint

	// Skip forage for x amount of returns if theres no return > 1* multiplier
	SkipForage uint

	// If resources go above this limit we are balling with money
	JBThreshold shared.Resources

	// Middle class:  Middle < Jeff bezos
	MiddleThreshold shared.Resources

	// Poor: Imperial student < Middle
	ImperialThreshold shared.Resources
}

// Client is the island number
type client struct {
	*baseclient.BaseClient

	forageHistory ForageHistory // Stores our previous foraging data

	taxAmount shared.Resources

	// allocation is the president's response to your last common pool resource request
	allocation shared.Resources

	config clientConfig
}

// Possible wealth classes
const (
	Dying           WealthTier = iota // Sets values = 0
	ImperialStudent                   // iota sets the folloing values =1
	MiddleClass                       // = 2
	JeffBezos                         // = 3
)

func (wt WealthTier) String() string {
	strings := [...]string{"Dying", "ImperialStudent", "MiddleClass", "JeffBezos"}
	if wt >= 0 && int(wt) < len(strings) {
		return strings[wt]
	}
	return fmt.Sprintf("Unkown wealth state '%v'", int(wt))
}

// GoString implements GoStringer
func (wt WealthTier) GoString() string {
	return wt.String()
}

// MarshalText implements TextMarshaler
func (wt WealthTier) MarshalText() ([]byte, error) {
	return miscutils.MarshalTextForString(wt.String())
}

// MarshalJSON implements RawMessage
func (wt WealthTier) MarshalJSON() ([]byte, error) {
	return miscutils.MarshalJSONForString(wt.String())
}

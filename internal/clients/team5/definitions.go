package team5

import (
	"fmt"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/SOMAS2020/SOMAS2020/pkg/miscutils"
)

type WealthTier int

// Array to keep tracks of CP requests and allocations history
type CPRequestHistory []shared.Resources
type CPAllocationHistory []shared.Resources 

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

const (  
	Dying            WealthTier = iota // Sets values = 0  
	Imperial_Student               // iota sets the folloing values =1  
	Middle_Class                   // = 2  
	Jeff_Bezos                     // = 3
)

func (st WealthTier) String() string {  
	strings := [...]string{"Dying", "Imperial_Student", "Middle_Class", "Jeff_Bezos"}  
	if st >= 0 && int(st) < len(strings) {    
		return strings[st]  
	}  
	return fmt.Sprintf("Unkown internal state '%v'", int(st))
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
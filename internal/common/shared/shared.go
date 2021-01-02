// Package shared is used to encapsulate items used by all other
// packages to prevent import cycles.
package shared

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/pkg/miscutils"
)

// ClientID is an enum for client IDs
type ClientID int

const (
	// Team1 ID
	Team1 ClientID = iota
	// Team2 ID
	Team2
	// Team3 ID
	Team3
	// Team4 ID
	Team4
	// Team5 ID
	Team5
	// Team6 ID
	Team6
)

const (
	// President Role
	President int = iota
	// Speaker Role
	Speaker
	// Judge Role
	Judge
)

// SortClientByID implements sort.Interface for []ClientID
type SortClientByID []ClientID

func (a SortClientByID) Len() int           { return len(a) }
func (a SortClientByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortClientByID) Less(i, j int) bool { return a[i] < a[j] }

// TeamIDs contain sequential IDs of all teams
var TeamIDs = [...]ClientID{Team1, Team2, Team3, Team4, Team5, Team6}

func (c ClientID) String() string {
	clientIDStrings := [...]string{"Team1", "Team2", "Team3", "Team4", "Team5", "Team6"}
	if c >= 0 && int(c) < len(clientIDStrings) {
		return clientIDStrings[c]
	}
	return fmt.Sprintf("UNKNOWN ClientID '%v'", int(c))
}

// GoString implements GoStringer
func (c ClientID) GoString() string {
	return c.String()
}

// MarshalText implements TextMarshaler
func (c ClientID) MarshalText() ([]byte, error) {
	return miscutils.MarshalTextForString(c.String())
}

// MarshalJSON implements RawMessage
func (c ClientID) MarshalJSON() ([]byte, error) {
	return miscutils.MarshalJSONForString(c.String())
}

// Resources represents amounts of resources.
// Used for foraging inputs and utility outputs (returns)
type Resources float64

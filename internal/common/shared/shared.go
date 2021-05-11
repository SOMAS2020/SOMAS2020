// Package shared is used to encapsulate items used by all other
// packages to prevent import cycles.
package shared

import (
	"fmt"
	// "os"
	"strconv"

	"github.com/SOMAS2020/SOMAS2020/pkg/miscutils"
)

// ClientID is an enum for client IDs
type ClientID int

// TeamIDs
// const (
// 	Team1 ClientID = iota
// 	Team2
// 	Team3
// 	Team4
// 	Team5
// 	Team6
// 	Team7
// 	Team8
// 	Team9
// 	Team10
// 	Team11
// 	Team12
// )

var TotalTeams int = 6

// Map of all IDs
var Teams map[string]ClientID = GenTeams()

// SortClientByID implements sort.Interface for []ClientID
type SortClientByID []ClientID

func (a SortClientByID) Len() int           { return len(a) }
func (a SortClientByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortClientByID) Less(i, j int) bool { return a[i] < a[j] }

// TeamIDs contain sequential IDs of all teams
// var TeamIDs = [...]ClientID{Team1, Team2, Team3, Team4, Team5, Team6, Team7, Team8, Team9, Team10, Team11, Team12}

// var TeamIDs []ClientID

var TeamIDs = make([]ClientID, TotalTeams)

func init() {

	// var args = os.Args

	// if len(args) > 1 {
	// 	TotalTeams, _ = strconv.Atoi(args[1])
	// } else {
	// 	TotalTeams = 6
	// }

	// TeamIDs = make([]ClientID, TotalTeams)

	for i := 0; i < TotalTeams; i++ {
		TeamIDs[i] = ClientID(i)
	}

	fmt.Println(TotalTeams)

}

func (c ClientID) String() string {
	// clientIDStrings := [...]string{"Team1", "Team2", "Team3", "Team4", "Team5", "Team6", "Team7", "Team8", "Team9", "Team10", "Team11", "Team12"}
	var clientIDStrings = make([]string, TotalTeams)

	for i := 0; i < TotalTeams; i++ {
		clientIDStrings[i] = "Team" + strconv.Itoa(i+1)
	}

	if c >= 0 && int(c) < len(clientIDStrings) {
		return clientIDStrings[c]
	}
	return fmt.Sprintf("UNKNOWN ClientID '%v'", int(c))
}

func GenTeams() map[string]ClientID {
	allTeams := make(map[string]ClientID)
	for i := 1; i < TotalTeams+1; i++ {
		allTeams["Team"+strconv.Itoa(i)] = ClientID(i - 1)
	}
	return allTeams
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

// Logger type for convenience in other definitions
type Logger func(format string, a ...interface{})

// Resources represents amounts of resources.
// Used for foraging inputs and utility outputs (returns)
type Resources float64

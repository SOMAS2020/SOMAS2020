package shared

import "fmt"

// ClientLifeStatus represents the three states a client's life can be in.
type ClientLifeStatus int

const (
	// Alive : able to perform all actions
	Alive ClientLifeStatus = iota
	// Critical : still alive (and hence can perform all actions, but will die when reach MaxCriticalConsecutiveTurns)
	Critical
	// Dead : no actions possible
	Dead
)

func (c ClientLifeStatus) String() string {
	clientLifeStatusStrings := [...]string{"Alive", "Critical", "Dead"}
	if c >= 0 && int(c) < len(clientLifeStatusStrings) {
		return clientLifeStatusStrings[c]
	}
	return fmt.Sprintf("UNKNOWN ClientLifeStatus '%v'", int(c))
}

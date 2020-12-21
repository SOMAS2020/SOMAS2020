package shared

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/pkg/miscutils"
)

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

// GoString implements GoStringer (for %#v printing)
func (c ClientLifeStatus) GoString() string {
	return c.String()
}

// MarshalText implements TextMarshaler
func (c ClientLifeStatus) MarshalText() ([]byte, error) {
	return miscutils.MarshalTextForString(c.String())
}

// MarshalJSON implements RawMessage
func (c ClientLifeStatus) MarshalJSON() ([]byte, error) {
	return miscutils.MarshalJSONForString(c.String())
}

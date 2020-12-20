package gamestate

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/pkg/miscutils"
)

func (a ActionType) String() string {
	if v, ok := actionTypeStringMap[a]; ok {
		return v
	}
	return fmt.Sprintf("UNKNOWN ActionType '%v'", int(a))
}

// GoString implements GoStringer (for %#v printing)
func (a ActionType) GoString() string {
	return a.String()
}

// MarshalText implements TextMarshaler
func (a ActionType) MarshalText() ([]byte, error) {
	return miscutils.MarshalTextForString(a.String())
}

// MarshalJSON implements RawMessage
func (a ActionType) MarshalJSON() ([]byte, error) {
	return miscutils.MarshalJSONForString(a.String())
}

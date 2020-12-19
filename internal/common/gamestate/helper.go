package gamestate

import "fmt"

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

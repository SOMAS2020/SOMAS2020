// Package miscutils provides miscellatenous utilities.
package miscutils

import (
	"fmt"
)

// MarshalTextForString returns the MarshalText function output required for a string.
func MarshalTextForString(s string) ([]byte, error) {
	return []byte(s), nil
}

// MarshalJSONForString returns the MarshalJSON function output required for a string.
func MarshalJSONForString(s string) ([]byte, error) {
	return []byte(fmt.Sprintf("\"%v\"", s)), nil
}

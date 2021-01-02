// Package testutils contains useful test utilities
// Based on https://github.com/facebook/openbmc/tools/flashy/tests/testutils.go
// Original author: lhl2617, Facebook, (GPLv2)
package testutils

import (
	"testing"
)

// CompareTestErrors is used to test and compare errors in testing.
func CompareTestErrors(want error, got error, t *testing.T) {
	if got == nil {
		if want != nil {
			t.Errorf("want '%v' got '%v'", want, got)
		}
	} else {
		if want == nil {
			t.Errorf("want '%v' got '%v'", want, got)
		} else if got.Error() != want.Error() {
			t.Errorf("want '%v' got '%v'", want.Error(), got.Error())
		}
	}
}

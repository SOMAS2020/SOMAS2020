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

// CompareTestIntegers is used to check wheter an expected int value is equal to what was got
func CompareTestIntegers(want int, got int, t *testing.T) {
	if want != got {
		t.Errorf("Wanted integer '%v' got '%v'", want, got)
	}
}

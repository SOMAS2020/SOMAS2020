package team5

import (
	"fmt"
	"testing"
)

func TestAbsoluteCap(t *testing.T) {
	var tests = []struct {
		input, thresh, want float64
	}{
		{0.5, 1.0, 0.5},
		{1.5, 1.0, 1.0},
		{-0.5, 1.0, -0.5},
		{-1.5, 1.0, -1.0},
		{0.5, 10.5, 0.5},
		{11, 10.5, 10.5},
		{-0.5, 12, -0.5},
		{-14, 12, -12},
	}
	for _, tc := range tests {
		testname := fmt.Sprintf("test input %.3f", tc.input)
		t.Run(testname, func(t *testing.T) {
			ans := absoluteCap(tc.input, tc.thresh)
			if ans != tc.want {
				t.Errorf("got %.1f, want %.1f", ans, tc.want)
			}
		})
	}
}

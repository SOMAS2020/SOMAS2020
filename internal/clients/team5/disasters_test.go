package team5

import (
	"testing"
)

func TestSortKeys(t *testing.T) {
	di := disasterInfo{}
	dh := disasterHistory{8: di, 3: di, 5: di, 1: di}

	sortedTurns := dh.sortKeys()
	for i, turn := range sortedTurns[1:] {
		if turn < sortedTurns[i] {
			t.Errorf("Expected sorted elements in ascending order, got %v", sortedTurns)
			return
		}
	}
}

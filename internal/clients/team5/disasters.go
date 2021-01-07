package team5

import (
	"sort"

	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
)

// stores all information pertaining to a disaster
type disasterInfo struct {
	report  disasters.DisasterReport
	effects disasters.DisasterEffects
	season  uint
}

type disasterHistory map[uint]disasterInfo

// effects contain abs magnitude, prop. mag relative to other islands and CP mitigated mag.
func (c *client) DisasterNotification(dR disasters.DisasterReport, effects disasters.DisasterEffects) {
	c.Logf("CRITICAL: Received notification of disaster: %v", dR.Display())
	c.disasterHistory[c.getTurn()] = disasterInfo{
		report:  dR,
		effects: effects,
		season:  c.getSeason(),
	}
}

func (d disasterHistory) sortKeys() []uint {
	keys := make([]int, 0)
	for k := range d {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	finalKeys := make([]uint, 0)
	for _, k := range keys {
		finalKeys = append(finalKeys, uint(k))
	}
	return finalKeys
}

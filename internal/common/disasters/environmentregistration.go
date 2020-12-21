package disasters

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

// InitEnvironment initialises environment according to definitions
func InitEnvironment(islandIDs []shared.ClientID, xBounds [2]float64, yBounds [2]float64, disasterParams DisasterParameters, cp Commonpool) (*Environment, error) {

	ag := ArchipelagoGeography{map[shared.ClientID]Island{}, xBounds, yBounds}

	for i, id := range islandIDs {
		island := Island{id, float64(i), float64(0)} // begin with points on x axis
		ag.islands[id] = island
	}
	return &Environment{ag, disasterParams, DisasterReport{}, cp}, nil // may want ability to return error in future
}

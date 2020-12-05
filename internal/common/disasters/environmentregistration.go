package disasters

// InitEnvironment initialises environment according to definitions
func InitEnvironment(islandNames []string, xBounds [2]float64, yBounds [2]float64, disasterParams DisasterParameters) (*Environment, error) {

	ag := ArchipelagoGeography{[]Island{}, xBounds, yBounds}

	for i, name := range islandNames {
		island := Island{name, float64(i), float64(0)} // begin with points on x axis
		ag.islands = append(ag.islands, island)
	}
	return &Environment{ag, disasterParams, DisasterReport{}}, nil // may want ability to return error in future
}

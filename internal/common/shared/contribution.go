package shared

// IntendedContribution is a struct containing the intended contribution to the common pool 
// of an island
type IntendedContribution struct {
	Contribution float64
	TeamsOfferedTo []ClientID
}

// ReceivedIntendedContribution is a struct containing the information describing 
// an intended contribution made by a given island
type ReceivedIntendedContribution struct {
	Contribution float64
	SharedFrom     ClientID
}

// IntendedContributionDict is a dictionary of IntendedContribution
type IntendedContributionDict = map[ClientID]IntendedContribution

// ReceivedIntendedContributionDict is a dictionary of ReceivedIntendedContribution
type ReceivedIntendedContributionDict = map[ClientID]ReceivedIntendedContribution

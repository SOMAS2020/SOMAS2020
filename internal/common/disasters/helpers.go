package disasters

import (
	"fmt"
	"strings"

	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// IslandLocation is a convenience method to extract an island's location given its index
func (a ArchipelagoGeography) IslandLocation(id shared.ClientID) (shared.Coordinate, shared.Coordinate) {
	island := a.Islands[id]
	return island.X, island.Y
}

// GetIslandIDs is a helper function to return the IDs of islands currently in env
func (env Environment) GetIslandIDs() []shared.ClientID {
	IDs := make([]shared.ClientID, 0, len(env.Geography.Islands))
	for k := range env.Geography.Islands {
		IDs = append(IDs, k)
	}
	return IDs
}

// Display is simply a string format method to viz a `DisasterReport` in console
func (report DisasterReport) Display() string {
	if report.Magnitude == 0 {
		return "No disaster reported."
	}
	return fmt.Sprintf("ALERT: Disaster of magnitude %.3f recorded at co-ordinates (%.2f, %.2f)\n", report.Magnitude, report.X, report.Y)
}

// DisplayReport is a string format method to viz a disaster report and its effect,
// as well as how the disaster is been mitigated by the common pool
func (env Environment) DisplayReport(cpResources shared.Resources, dConf config.DisasterConfig) string {
	disasterReport := env.LastDisasterReport
	if env.LastDisasterReport.Magnitude == 0 {
		return "No disaster reported. No disaster effects."
	}
	var sb strings.Builder
	sb.WriteString(disasterReport.Display())

	effects := env.ComputeDisasterEffects(cpResources, dConf)

	// display absolute effects for each island
	sb.WriteString("\n------------------------ Disaster Effects ------------------------\n")

	for islandID, absEffect := range effects.Absolute {
		island := env.Geography.Islands[islandID]
		sb.WriteString(fmt.Sprintf(
			"island ID: %d, \txy co-ords: (%.2f, %.2f), \tabsolute damage: %.2f \n",
			islandID, island.X, island.Y, absEffect*dConf.MagnitudeResourceMultiplier))
	}

	// display propotional effects relative to other islands and effects after CP mitigation
	sb.WriteString("\n------------------------ Disaster Effects after CP Mitigation ------------------------\n")
	for islandID, propEffect := range effects.Proportional {
		sb.WriteString(fmt.Sprintf("%v: proportional damage: %.2f, \t damage after common pool mitigation: %.2f \n", islandID, propEffect, effects.CommonPoolMitigated[islandID]))
	}
	return sb.String()
}

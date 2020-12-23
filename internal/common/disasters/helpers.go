package disasters

import (
	"fmt"
	"strings"

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

// DisplayReport is a string format method to viz a disaster report and its effect
func (env Environment) DisplayReport() string {
	disasterReport := env.LastDisasterReport.Display()
	if env.LastDisasterReport.Magnitude == 0 {
		return disasterReport // just return default no disaster message. Not necessary to report affected islands.
	}
	var sb strings.Builder
	sb.WriteString(disasterReport + "\n")
	sb.WriteString("------------------------ Disaster Effects ------------------------\n")
	for islandID, effect := range env.DisasterEffects() {
		island := env.Geography.Islands[islandID]
		sb.WriteString(fmt.Sprintf("island ID: %d, \txy co-ords: (%.2f, %.2f), \tdisaster effect: %.2f \n", islandID, island.X, island.Y, effect))
	}
	return sb.String()
}

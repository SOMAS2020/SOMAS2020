package disasters

import (
	"fmt"
	"strings"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// IslandLocation is a convenience method to extract an island's location given its index
func (a ArchipelagoGeography) IslandLocation(id shared.ClientID) []float64 {
	island := a.islands[id]
	return []float64{island.x, island.y}
}

// Display is simply a string format method to viz a `DisasterReport` in console
func (report DisasterReport) Display() string {
	if report.magnitude == 0 {
		return "No disaster reported."
	}
	return fmt.Sprintf("\nALERT: Disaster of magnitude %.3f recorded at co-ordinates (%.2f, %.2f)\n", report.magnitude, report.x, report.y)
}

// DisplayReport is a string format method to viz a disaster report and its effect
func (env Environment) DisplayReport() string {
	disasterReport := env.lastDisasterReport.Display()
	if env.lastDisasterReport.magnitude == 0 {
		return disasterReport // just return default no disaster message. Not necessary to report affected islands.
	}
	var sb strings.Builder
	sb.WriteString(disasterReport)
	sb.WriteString("\n------------------------ Disaster Effects ------------------------\n")
	for islandID, effect := range env.DisasterEffects() {
		island := env.geography.islands[islandID]
		sb.WriteString(fmt.Sprintf("island ID: %d, \txy co-ords: (%.2f, %.2f), \tdisaster effect: %.2f \n", islandID, island.x, island.y, effect))
	}
	return sb.String()
}

package disasters

import (
	"fmt"
	"strings"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// SpatialPDFType is an enum for xy prob. density function types for disaster simulation
type SpatialPDFType int

const (
	// Uniform xy distribution: disaster peak occurs uniformally over xy bounds of env
	Uniform SpatialPDFType = iota

	// add other PDFs here post-MVP
)

// IslandLocation is a convenience method to extract an island's location given its index
func (a ArchipelagoGeography) IslandLocation(id shared.ClientID) []float64 {
	island := a.islands[id]
	return []float64{island.x, island.y}
}

// GetIslandIDs is a helper function to return the IDs of islands currently in env
func (env Environment) GetIslandIDs() []shared.ClientID {
	IDs := make([]shared.ClientID, 0, len(env.Geography.islands))
	for k := range env.Geography.islands {
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
		island := env.Geography.islands[islandID]
		sb.WriteString(fmt.Sprintf("island ID: %d, \txy co-ords: (%.2f, %.2f), \tdisaster effect: %.2f \n", islandID, island.x, island.y, effect))
	}
	return sb.String()
}

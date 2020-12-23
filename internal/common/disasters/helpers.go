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
// DisplayReport also visualize how the disaster has been mitigated by the common pool
func (env *Environment) DisplayReport() (string, map[shared.ClientID]float64) {
	disasterReport := env.LastDisasterReport.Display()
	if env.LastDisasterReport.Magnitude == 0 {
		return disasterReport, nil // just return default no disaster message. Not necessary to report affected islands.
	}
	var sb strings.Builder
	sb.WriteString(disasterReport)
	sb.WriteString("\n------------------------ Disaster Effects ------------------------\n")

	proportionalEffect := map[shared.ClientID]float64{}
	individualEffect := map[shared.ClientID]float64{}
	individualEffect, proportionalEffect = env.DisasterEffects()

	for islandID, effect := range individualEffect {
	
		island := env.Geography.Islands[islandID]
		sb.WriteString(fmt.Sprintf("island ID: %d, \txy co-ords: (%.2f, %.2f), \tdisaster effect: %.2f \tActual damage: %.2f \n", islandID, island.X, island.Y, effect, individualEffect[islandID]*1000))
	}

	updatedProportionalEffect := map[shared.ClientID]float64{}
	updatedProportionalEffect = env.DisasterMitigate(individualEffect, proportionalEffect)
	damageDifference := map[shared.ClientID]float64{}
	for islandID, _ := range updatedProportionalEffect {
		damageDifference[islandID] = individualEffect[islandID] * 1000 - updatedProportionalEffect[islandID]
	}

	sb.WriteString("\n------------------------ Updated Disaster Effects ------------------------\n")
	for islandID, effect := range individualEffect {
		island := env.Geography.Islands[islandID]
		sb.WriteString(fmt.Sprintf("island ID: %d, \txy co-ords: (%.2f, %.2f), \tdisaster effect: %.2f, \tUpdated damage: %.2f, \t Common pool mitigated: %.2f \n", islandID, island.X, island.Y, effect, updatedProportionalEffect[islandID], damageDifference[islandID] ))
	}
	return sb.String(), updatedProportionalEffect
}

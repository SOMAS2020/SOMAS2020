package server

import (
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common"
)

// TestClientReg checks that all clients are registered
func TestClientReg(t *testing.T) {
	const numTeams = 6 // we have 6 teams
	numRegClients := len(common.RegisteredClients)
	if numRegClients != numTeams {
		t.Errorf("Are all teams registered? want '%v' got '%v'", numTeams, numRegClients)
	}
}

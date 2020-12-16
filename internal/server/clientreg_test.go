package server

import (
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// TestClientReg checks that all clients are registered
func TestNumClientReg(t *testing.T) {
	const numTeams = 6 // we have 6 teams
	numRegClients := len(common.RegisteredClients)
	if numRegClients != numTeams {
		t.Errorf("Are all teams registered? want '%v' got '%v'", numTeams, numRegClients)
	}
}

func TestClientReg(t *testing.T) {
	checkClientReg := func(id shared.ClientID, c common.Client) {
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("Client %v was registered with a nil common.Client!", id)
			}
		}()
		c.Echo("checking!")
	}

	for id, client := range common.RegisteredClients {
		checkClientReg(id, client)
	}
}

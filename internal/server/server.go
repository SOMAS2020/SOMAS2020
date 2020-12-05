// Package server contains server-side code
package server

import (
	"fmt"
	"log"

	"github.com/SOMAS2020/SOMAS2020/internal/common"
	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// Server represents the primary server interface exposed to the simulation.
type Server interface {
	GetEcho(s string) error
	Logf(format string, a ...interface{})
}

// SOMASServer implements Server.
type SOMASServer struct {
	gameState common.GameState
}

// SOMASServerFactory returns an instance of the main server we use.
func SOMASServerFactory() Server {
	return &SOMASServer{
		gameState: common.GameState{
			Day:         1,
			ClientInfos: getClientInfoFromRegisteredClients(common.RegisteredClients),
			Environment: initEnvironment(),
		},
	}
}

func getClientInfoFromRegisteredClients(registeredClients map[shared.ClientID]common.Client) map[shared.ClientID]common.ClientInfo {
	clientInfos := map[shared.ClientID]common.ClientInfo{}

	for id, c := range registeredClients {
		clientInfos[id] = common.ClientInfo{
			Client:    c,
			Resources: common.DefaultResources,
			Alive:     true,
		}
	}

	return clientInfos
}

// GetEcho retrieves an echo from all the clients and make sure they are the same.
func (s *SOMASServer) GetEcho(str string) error {
	cis := s.gameState.ClientInfos
	for _, id := range shared.TeamIDs {
		ci := cis[id]
		c := ci.Client
		got := c.Echo(str)
		if str != got {
			return fmt.Errorf("Echo error: want '%v' got '%v' from client %v",
				str, got, c.GetID())
		}
		s.Logf("Received echo `%v` from client %v", str, c.GetID())
	}
	return nil
}

// Logf is the server's default logger.
func (s *SOMASServer) Logf(format string, a ...interface{}) {
	log.Printf("[SERVER]: %v", fmt.Sprintf(format, a...))
}

func initEnvironment() *disasters.Environment {
	islandNames := []string{}
	for id := range common.RegisteredClients {
		islandNames = append(islandNames, fmt.Sprintf("Island %v", id))
	}
	xBounds := [2]float64{0, 10}
	yBounds := [2]float64{0, 10}
	dp := disasters.DisasterParameters{GlobalProb: 0.1, SpatialPDF: "uniform", MagnitudeLambda: 1.0}
	env, _ := disasters.InitEnvironment(islandNames, xBounds, yBounds, dp)
	return env
}

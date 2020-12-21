package server

import (
	"fmt"
	"log"
	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/pkg/errors"
)

// Server represents the primary server interface exposed to the simulation.
type Server interface {
	// EntryPoint function that returns a list of historic gamestate.ClientInfos until the
	// game ends.
	EntryPoint() ([]gamestate.GameState, error)
}

// SOMASServer implements Server.
type SOMASServer struct {
	gameState gamestate.GameState

	// ClientMap maps from the ClientID to the Client object.
	// We don't store this in gameState--gameState is shared to clients and should
	// not contain pointers to other clients!
	clientMap map[shared.ClientID]baseclient.Client
}

// SOMASServerFactory returns an instance of the main server we use.
func SOMASServerFactory() Server {
	clientInfos, clientMap := getClientInfosAndMapFromRegisteredClients(baseclient.RegisteredClients)

	return &SOMASServer{
		clientMap: clientMap,
		gameState: gamestate.GameState{
			Season:      1,
			Turn:        1,
			ClientInfos: clientInfos,
			Environment: initEnvironment(),
		},
	}
}

// EntryPoint function that returns a list of historic gamestate.GameState until the
// game ends.
func (s *SOMASServer) EntryPoint() ([]gamestate.GameState, error) {
	states := []gamestate.GameState{s.gameState.Copy()}

	for !s.gameOver(config.GameConfig().MaxTurns, config.GameConfig().MaxSeasons) {
		if err := s.runTurn(); err != nil {
			return states, err
		}
		states = append(states, s.gameState.Copy())
	}
	return states, nil
}

// getEcho retrieves an echo from all the clients and make sure they are the same.
func (s *SOMASServer) getEcho(str string) error {
	for _, c := range s.clientMap {
		got := c.Echo(str)
		if str != got {
			return errors.Errorf("Echo error: want '%v' got '%v' from %v",
				str, got, c.GetID())
		}
		s.logf("Received echo `%v` from %v", str, c.GetID())
	}
	return nil
}

// logf is the server's default logger.
func (s *SOMASServer) logf(format string, a ...interface{}) {
	log.Printf("[SERVER]: %v", fmt.Sprintf(format, a...))
}

func initEnvironment() *disasters.Environment {
	var cp = disasters.Commonpool{900, 1000}
	_, clientMap := getClientInfosAndMapFromRegisteredClients(baseclient.RegisteredClients)
	islandIDs := make([]shared.ClientID, 0, len(clientMap))
	for id := range clientMap {
		islandIDs = append(islandIDs, id)
	}
	xBounds := [2]float64{0, 10}
	yBounds := [2]float64{0, 10}
	dp := disasters.DisasterParameters{GlobalProb: 1, SpatialPDF: "uniform", MagnitudeLambda: 1.0}
	env, _ := disasters.InitEnvironment(islandIDs, xBounds, yBounds, dp, cp)
	return env
}

//islandDistribute distributes leftover resource after disaster's mitigation
func (s *SOMASServer) islandDistribute(leftover uint, disasterReport string) { //receives 
	if disasterReport != "No disaster reported." {
		var clientPortion int
		clientPortion = int(leftover / 6) // equally divide leftover
		for _, id := range shared.TeamIDs {
		ci := s.gameState.ClientInfos[id]
		ci.Resources = ci.Resources + clientPortion //adds to each client its portion
		s.gameState.ClientInfos[id] = ci
		s.gameState.Environment.CommonPool.Resource = 0;
	}
		
	}
	

}

//islandDeplete depletes island's resource based on the severity of the storm
func (s *SOMASServer) islandDeplete(porpotionEffect map[shared.ClientID]float64) {
	temp := 0.0
	for id, ci := range s.gameState.ClientInfos {
		temp = math.Max(float64(ci.Resources)-porpotionEffect[id], 0.0) //adds to each client its portion	
		ci.Resources = int(temp)
		s.gameState.ClientInfos[id] = ci
	}
}
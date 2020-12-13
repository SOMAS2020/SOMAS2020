// Package server contains server-side code
package server

import (
	"fmt"
	"log"

	"github.com/SOMAS2020/SOMAS2020/internal/common"
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
	"github.com/SOMAS2020/SOMAS2020/internal/common/foraging"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/pkg/errors"
)

// Server represents the primary server interface exposed to the simulation.
type Server interface {
	// EntryPoint function that returns a list of historic common.GameStates until the
	// game ends.
	EntryPoint() ([]common.GameState, error)
}

// SOMASServer implements Server.
type SOMASServer struct {
	gameState common.GameState

	// ClientMap maps from the ClientID to the Client object.
	// We don't store this in gameState--gameState is shared to clients and should
	// not contain pointers to other clients!
	clientMap map[shared.ClientID]common.Client
}

// SOMASServerFactory returns an instance of the main server we use.
func SOMASServerFactory() Server {
	clientInfos, clientMap := getClientInfosAndMapFromRegisteredClients(common.RegisteredClients)

	return &SOMASServer{
		clientMap: clientMap,
		gameState: common.GameState{
			Environment: initEnvironment(),
			Season:      1,
			Turn:        1,
			ClientInfos: clientInfos,
		},
	}
}

// EntryPoint function that returns a list of historic common.GameStates until the
// game ends.
func (s *SOMASServer) EntryPoint() ([]common.GameState, error) {
	states := []common.GameState{s.gameState.Copy()}

	for !s.gameOver(config.MaxTurns, config.MaxSeasons) {
		if err := s.runTurn(); err != nil {
			return states, err
		}
		states = append(states, s.gameState.Copy())
	}
	return states, nil
}

// runRound runs a round (day) of the game.
func (s *SOMASServer) runRound() error {
	if err := s.getEcho("HELLO WORLD!"); err != nil {
		return fmt.Errorf("getEcho failed with: %v", err)
	}
	s.gameState.Environment.SampleForDisaster()
	fmt.Println(s.gameState.Environment.DisplayReport())

	huntParticipants := map[shared.ClientID]float64{shared.Team1: 1.0, shared.Team2: 0.9} // just to test for now
	deerHunt := createDeerHunt(huntParticipants)
	fmt.Printf("\nResults of deer hunt: return of %.3f at cost of %.3f\n", deerHunt.Hunt(), deerHunt.TotalInput())

	dp := foraging.CreateBasicDeerPopulationModel()
	consumption := []int{0, 0, 2, 0, 0, 1, 0, 3, 0, 0} // simulate deer consumption (no. deer hunted each day) over 10 days
	dp.Simulate(consumption)
	s.killAllClients()
	return nil
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
	islandIDs := make([]shared.ClientID, 0, len(common.RegisteredClients))
	for id := range common.RegisteredClients {
		islandIDs = append(islandIDs, id)
	}
	xBounds := [2]float64{0, 10}
	yBounds := [2]float64{0, 10}
	dp := disasters.DisasterParameters{GlobalProb: 0.1, SpatialPDF: "uniform", MagnitudeLambda: 1.0}
	env, _ := disasters.InitEnvironment(islandIDs, xBounds, yBounds, dp)
	return env
}

func createDeerHunt(teamResourceInputs map[shared.ClientID]float64) foraging.DeerHunt {
	params := foraging.DeerHuntParams{P: 0.95, Lam: 1.0} // TODO: move to central config store
	return foraging.DeerHunt{Participants: teamResourceInputs, Params: params}
}

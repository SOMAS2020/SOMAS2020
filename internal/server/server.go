// Package server contains server-side code
package server

import (
	"fmt"
	"log"

	"github.com/SOMAS2020/SOMAS2020/internal/common"
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
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

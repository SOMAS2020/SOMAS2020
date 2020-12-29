package server

import (
	"fmt"
	"log"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
	"github.com/SOMAS2020/SOMAS2020/internal/common/foraging"
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

	gameConfig config.Config

	// ClientMap maps from the ClientID to the Client object.
	// We don't store this in gameState--gameState is shared to clients and should
	// not contain pointers to other clients!
	clientMap map[shared.ClientID]baseclient.Client
}

// NewSOMASServer returns an instance of the main server we use.
func NewSOMASServer(gameConfig config.Config) Server {
	clientInfos, clientMap := getClientInfosAndMapFromRegisteredClients(
		baseclient.RegisteredClients,
		gameConfig.InitialResources,
	)
	return createSOMASServer(clientInfos, clientMap, gameConfig)
}

// createSOMASServer creates the main server given initial data about the
// clients. Extracted from SOMASServerFactory for testing purposes.
func createSOMASServer(
	clientInfos map[shared.ClientID]gamestate.ClientInfo,
	clientMap map[shared.ClientID]baseclient.Client,
	gameConfig config.Config,
) Server {
	clientIDs := make([]shared.ClientID, 0, len(clientMap))
	for k := range clientMap {
		clientIDs = append(clientIDs, k)
	}

	server := &SOMASServer{
		clientMap:  clientMap,
		gameConfig: gameConfig,
		gameState: gamestate.GameState{
			Season:         1,
			Turn:           1,
			ClientInfos:    clientInfos,
			Environment:    disasters.InitEnvironment(clientIDs, gameConfig.DisasterConfig),
			DeerPopulation: foraging.CreateDeerPopulationModel(gameConfig.ForagingConfig),
			IIGOHistory:    []shared.Accountability{},
			SpeakerID:      shared.Team1,
			JudgeID:        shared.Team2,
			PresidentID:    shared.Team3,
		},
	}

	for _, client := range clientMap {
		client.Initialise(ServerForClient{
			clientID: client.GetID(),
			server:   server,
		})
	}

	return server
}

// EntryPoint function that returns a list of historic gamestate.GameState until the
// game ends.
func (s *SOMASServer) EntryPoint() ([]gamestate.GameState, error) {
	states := []gamestate.GameState{s.gameState.Copy()}

	for !s.gameOver(s.gameConfig.MaxTurns, s.gameConfig.MaxSeasons) {
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

// ServerForClient is a reference to the server for particular client. It is
// meant as an instance of baseclient.ServerReadHandle
type ServerForClient struct {
	clientID shared.ClientID
	server   *SOMASServer
}

// GetGameState gets the ClientGameState for the client matching s.clientID in
// s.server
func (s ServerForClient) GetGameState() gamestate.ClientGameState {
	return s.server.gameState.GetClientGameStateCopy(s.clientID)
}

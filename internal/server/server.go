// Package server contains server-side code
package server

import (
	"fmt"
	"log"

	"github.com/SOMAS2020/SOMAS2020/internal/common"
)

// Server represents the primary server interface exposed to the simulation.
type Server interface {
	GetEcho(s string) error
	Logf(format string, a ...interface{})
}

// SOMASServer implements Server.
type SOMASServer struct {
	clients   []common.Client
	gameState common.GameState
}

// SOMASServerFactory returns an instance of the main server we use.
func SOMASServerFactory() Server {
	return &SOMASServer{clients: common.RegisteredClients}
}

// GetEcho retrieves an echo from all the clients and make sure they are the same.
func (s *SOMASServer) GetEcho(str string) error {
	for _, c := range s.clients {
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

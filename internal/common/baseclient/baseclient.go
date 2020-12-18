// Package baseclient contains the Client interface as well as a base client implementation.
package baseclient

import (
	"fmt"
	"log"

	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// Client is a base interface to be implemented by each client struct.
type Client interface {
	Echo(s string) string
	GetID() shared.ClientID

	// StartOfTurnUpdate is where SOMASServer.updateIsland sends the game state over
	// at start of turn. Do whatever you like here :).
	StartOfTurnUpdate(gameState gamestate.GameState)

	// GameStateUpdate updates game state mid-turn.
	GameStateUpdate(gameState gamestate.GameState)

	Logf(format string, a ...interface{})
	// EndOfTurnActions should return all end of turn actions.
	EndOfTurnActions() []gamestate.Action
}

// NewClient produces a new client with the BaseClient already implemented.
// BASE: Do not overwrite in team client.
func NewClient(id shared.ClientID) Client {
	return &BaseClient{id: id}
}

// BaseClient provides a basic implementation for all functions of the client interface and should always the interface fully.
// All clients should be based off of this BaseClient to ensure that all clients implement the interface,
// even when new features are added.
type BaseClient struct {
	id shared.ClientID
}

// Echo prints a message to show that the client exists
// BASE: Do not overwrite in team client.
func (c *BaseClient) Echo(s string) string {
	c.Logf("Echo: '%v'", s)
	return s
}

// GetID gets returns the id of the client.
// BASE: Do not overwrite in team client.
func (c *BaseClient) GetID() shared.ClientID {
	return c.id
}

// Logf is the client's logger that prepends logs with your ID. This makes
// it easier to read logs. DO NOT use other loggers that will mess logs up!
// BASE: Do not overwrite in team client.
func (c *BaseClient) Logf(format string, a ...interface{}) {
	log.Printf("[%v]: %v", c.id, fmt.Sprintf(format, a...))
}

// StartOfTurnUpdate is updates the gamestate of the client at the start of each turn.
// The gameState is served by the server.
// OPTIONAL. Base should be able to handle it but feel free to implement your own.
func (c *BaseClient) StartOfTurnUpdate(gameState gamestate.GameState) {
	c.Logf("Received start of turn game state update: %v", gameState)
	// TODO
}

// GameStateUpdate updates game state mid-turn.
// The gameState is served by the server.
// OPTIONAL. Base should be able to handle it but feel free to implement your own.
func (c *BaseClient) GameStateUpdate(gameState gamestate.GameState) {
	c.Logf("Received game state update: %v", gameState)
}

// EndOfTurnActions executes and returns the actions done by the client that turn.
// OPTIONAL. Base should be able to handle it but feel free to implement your own.
func (c *BaseClient) EndOfTurnActions() []gamestate.Action {
	c.Logf("EndOfTurnActions")
	return nil
}

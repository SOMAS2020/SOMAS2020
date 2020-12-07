// Package team2 contains code for team 2's client implementation
package team2

import (
	"fmt"
	"log"

	"github.com/SOMAS2020/SOMAS2020/internal/common"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const id = shared.Team2

func init() {
	common.RegisterClient(id, &client{id: id})
}

type client struct {
	id shared.ClientID
}

func (c *client) Echo(s string) string {
	c.logf("Echo: '%v'", s)
	return s
}

func (c *client) GetID() shared.ClientID {
	return c.id
}

// logf is the client's logger that prepends logs with your ID. This makes
// it easier to read logs. DO NOT use other loggers that will mess logs up!
func (c *client) logf(format string, a ...interface{}) {
	log.Printf("[%v]: %v", c.id, fmt.Sprintf(format, a...))
}

func (c *client) ReceiveGameStateUpdate(gameState common.GameState) {
	c.logf("Received game state update: %v", gameState)
	// TODO
}

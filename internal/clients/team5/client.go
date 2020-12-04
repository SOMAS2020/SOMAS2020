// Package team5 contains code for team 5's client implementation
package team5

import (
	"fmt"
	"log"

	"github.com/SOMAS2020/SOMAS2020/internal/common"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const id = shared.Team5

func init() {
	common.RegisterClient(id, &client{id: id})
}

type client struct {
	id shared.ClientID
}

func (c *client) Echo(s string) string {
	c.Logf("Echo: '%v'", s)
	return s
}

func (c *client) GetID() shared.ClientID {
	return c.id
}

// Logf is the client's logger that prepends logs with your ID. This makes
// it easier to read logs. DO NOT use other loggers that will mess logs up!
func (c *client) Logf(format string, a ...interface{}) {
	log.Printf("[%v]: %v", c.id, fmt.Sprintf(format, a...))
}

// Package team3 contains code for team 3's client implementation
package team3

import (
	"fmt"
	"github.com/SOMAS2020/SOMAS2020/internal/common"
	"log"
)

const id = common.Team3

func init() {
	common.RegisterClient(id, &client{id: id})
}

type client struct {
	id common.ClientID
}

func (c *client) Echo(s string) string {

	//c.DemoEvaluation()
	c.Logf("Echo: '%v'", s)

	return s
}

func (c *client) DemoEvaluation() {
	evalResult, err := common.BasicRuleEvaluator("Kinda Complicated Rule")
	if err != nil {
		panic(err.Error())
	}
	c.Logf("Rule Eval: %t", evalResult)
}

func (c *client) GetID() common.ClientID {
	return c.id
}

// Logf is the client's logger that prepends logs with your ID. This makes
// it easier to read logs. DO NOT use other loggers that will mess logs up!
func (c *client) Logf(format string, a ...interface{}) {
	log.Printf("[%v]: %v", c.id, fmt.Sprintf(format, a...))
}

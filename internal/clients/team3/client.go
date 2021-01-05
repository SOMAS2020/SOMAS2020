// Package team3 contains code for team 3's client implementation
package team3

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const id = shared.Team3

func init() {
	baseclient.RegisterClient(id, &client{BaseClient: baseclient.NewClient(id)})
}

type client struct {
	*baseclient.BaseClient
}

func (c *client) DemoEvaluation() {
	ret := rules.EvaluateRule("Kinda Complicated Rule")
	if ret.EvalError != nil {
		panic(ret.EvalError.Error())
	}
	c.Logf("Rule Eval: %t", ret.RulePasses)
}

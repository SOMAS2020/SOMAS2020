// Package team3 contains code for team 3's client implementation
package team3

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const id = shared.Team3

func init() {
	common.RegisterClient(id, &client{Client: common.NewClient(id)})
}

type client struct {
	common.Client
}

func (c *client) DemoEvaluation() {
	evalResult, err := rules.BasicBooleanRuleEvaluator("Kinda Complicated Rule")
	if err != nil {
		panic(err.Error())
	}
	c.Logf("Rule Eval: %t", evalResult)
}

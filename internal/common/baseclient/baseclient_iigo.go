package baseclient

import (
	"math/rand"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
)

// CommonPoolResourceRequest is called by the President in IIGO to
// request an allocation of resources from the common pool.
func(c *BaseClient) CommonPoolResourceRequest() int{
	//TODO : Team 5 need access to the common pool gamestate...
	return 0
}

// ResourceReport is an island's self-report of its own resources.
func(c *BaseClient) ResourceReport() int{
	return c.clientGameState.ClientInfo.Resources
}

// RuleProposal is called by the President in IIGO to propose a
// rule to be voted on.
func(c *BaseClient) RuleProposal() string{
	allrules := make([]string, 0, len(rules.RulesInPlay))
    for r := range rules.RulesInPlay {
        allrules = append(allrules, r)
	}
	return allrules[rand.Intn(len(allrules))]
}
package baseclient

import "github.com/SOMAS2020/SOMAS2020/internal/common/rules"

// CommonPoolResourceRequest is called by the President in IIGO to
// request an allocation of resources from the common pool.
func (c *BaseClient) CommonPoolResourceRequest() int {
	return 20
}

// ResourceReport is an island's self-report of its own resources.
func (c *BaseClient) ResourceReport() int {
	return c.clientGameState.ClientInfo.Resources
}

// RuleProposal is called by the President in IIGO to propose a
// rule to be voted on.
func (c *BaseClient) RuleProposal() string {
	allRules := rules.AvailableRules
	for k, _ := range allRules {
		return k
	}
	return ""
}

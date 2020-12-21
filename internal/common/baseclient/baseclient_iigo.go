package baseclient

// CommonPoolResourceRequest is called by the President in IIGO to
// request an allocation of resources from the common pool.
func(c *BaseClient) CommonPoolResourceRequest(){

}

// ResourceReport is an islands self-report of its own resources.
func(c *BaseClient) ResourceReport()int{
	return c.clientGameState.ClientInfo.Resources
}

// RuleProposal is called by the President in IIGO to propose a
// rule to be voted on.
func(c *BaseClient) RuleProposal(){

}
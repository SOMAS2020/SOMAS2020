package team2

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
)

type judge struct {
	// Base implementation
	*baseclient.BaseJudge
	// Our client
	c *client
}

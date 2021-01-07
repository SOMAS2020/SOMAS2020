package team2

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
)

type president struct {
	// Base implementation
	*baseclient.BasePresident
	// Our client
	c *client
}

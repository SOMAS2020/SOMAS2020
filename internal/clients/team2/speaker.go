package team2

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
)

type speaker struct {
	// Base implementation
	*baseclient.BaseSpeaker
	// Our client
	c *client
}

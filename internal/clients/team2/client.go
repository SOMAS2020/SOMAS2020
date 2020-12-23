// Package team2 contains code for team 2's client implementation
package team2

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const id = shared.Team2

func init() {
	baseclient.RegisterClient(id, &client{Client: baseclient.NewClient(id)})
}

type client struct {
	baseclient.Client
}

func test() {
	c.Logf("I am here")
}

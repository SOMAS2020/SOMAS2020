// Package team1 contains code for team 1's client implementation
package team1

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const id = shared.Team1

func init() {
	baseclient.RegisterClient(id, &client{Client: baseclient.NewClient(id)})
}

type client struct {
	baseclient.Client
}

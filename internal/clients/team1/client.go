// Package team1 contains code for team 1's client implementation
package team1

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const id = shared.Team1

func init() {
	common.RegisterClient(id, &client{Client: common.NewClient(id)})
}

type client struct {
	common.Client
}

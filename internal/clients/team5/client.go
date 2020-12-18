// Package team5 contains code for team 5's client implementation
package team5

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const id = shared.Team5

func init() {
	common.RegisterClient(id, &client{Client: common.NewClient(id)})
}

type client struct {
	common.Client
}

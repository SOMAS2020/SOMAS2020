// Package team6 contains code for team 6's client implementation
package team6

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const id = shared.Team6

func init() {
	common.RegisterClient(id, &client{Client: common.NewClient(id)})
}

type client struct {
	common.Client
}

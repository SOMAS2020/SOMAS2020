// Package team4 contains code for team 4's client implementation
package team4

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const id = shared.Team4

func init() {
	baseclient.RegisterClient(id, &client{BaseClient: baseclient.NewClient(id)})
}

type client struct {
	*baseclient.BaseClient
}

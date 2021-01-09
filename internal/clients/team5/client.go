// Package team5 contains code for team 5's client implementation
package team5

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const id = shared.Team5

func init() {
	baseclient.RegisterClientFactory(id, func() baseclient.Client { return baseclient.NewClient(id) })
}

type client struct {
	*baseclient.BaseClient
}

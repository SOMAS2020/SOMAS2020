// Package team6 contains code for team 6's client implementation
package team6

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

const id = shared.Team6

func init() {
	baseclient.RegisterClient(id, &client{Client: baseclient.NewClient(id)})
}

type client struct {
	baseclient.Client
}

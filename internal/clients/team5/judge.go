package team5

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
)

type judge struct {
	*baseclient.BaseJudge
}

func (c *client) GetClientJudgePointer() roles.Judge {
	return &c.team5Judge
}

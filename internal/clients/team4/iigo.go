package team4

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
)

func (c *client) GetClientJudgePointer() roles.Judge {
	return &c.clientJudge
}

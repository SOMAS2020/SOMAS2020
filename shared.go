package main

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/pkg/gitinfo"
)

// output represents what is output into the output.json file
type output struct {
	GameStates []gamestate.GameState
	Config     config.Config
	GitInfo    gitinfo.GitInfo
}

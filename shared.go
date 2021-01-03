package main

import (
	"time"

	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/SOMAS2020/SOMAS2020/pkg/gitinfo"
)

type runInfo struct {
	TimeStart       time.Time
	TimeEnd         time.Time
	DurationSeconds float64
	Version         string
	GOOS            string
	GOARCH          string
}

// auxInfo contains other useful auxiliary information, mainly
// used to help visualisation
type auxInfo struct {
	TeamIDs []string
}

func getAuxInfo() auxInfo {
	teams := make([]string, len(shared.TeamIDs))
	for idx, teamID := range shared.TeamIDs {
		teams[idx] = teamID.String()
	}
	return auxInfo{
		TeamIDs: teams,
	}
}

// output represents what is output into the output.json file
type output struct {
	Config     config.Config
	GitInfo    gitinfo.GitInfo
	RunInfo    runInfo
	AuxInfo    auxInfo
	GameStates []gamestate.GameState
}

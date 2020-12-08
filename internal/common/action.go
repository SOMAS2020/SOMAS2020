package common

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/pkg/errors"
)

// ActionType is an enum that classifies actions by their type
type ActionType int

const (
	// GiveClientResource : give resource from a client to the other
	GiveClientResource ActionType = iota
)

// Action represents an action that modifies game state
type Action struct {
	ActionType               ActionType
	GiveClientResourceAction *GiveClientResourceAction
}

// GiveClientResourceAction describes an action of giving resources from a client to the other
type GiveClientResourceAction struct {
	SourceClientID shared.ClientID
	TargetClientID shared.ClientID
	Resources      int
}

func (g *GameState) DispatchActions(actions []Action) error {
	g.logf("start DispatchActions")
	defer g.logf("finish DispatchActions")
	for _, action := range actions {
		err := g.dispatchAction(action)
		if err != nil {
			g.logf("ERROR IGNORED: Cannot dispatch action '%v': %v",
				action, err)
		}
	}
	return nil
}

func (g *GameState) dispatchAction(action Action) error {
	switch action.ActionType {
	case GiveClientResource:
		return g.dispatchGiveClientResourceAction(action.GiveClientResourceAction)
	default:
		return errors.Errorf("Unknown action type: %v", action.ActionType)
	}
}

func (g *GameState) dispatchGiveClientResourceAction(giveClientResourceAction *GiveClientResourceAction) error {
	if giveClientResourceAction == nil {
		return errors.Errorf("giveClientResourceAction is nil!")
	}
	srcClientInfo := g.ClientInfos[giveClientResourceAction.SourceClientID]
	trgClientInfo := g.ClientInfos[giveClientResourceAction.TargetClientID]

	if srcClientInfo.Resources < giveClientResourceAction.Resources {
		return errors.Errorf("%v does not have %v resources to give to %v",
			giveClientResourceAction.SourceClientID, giveClientResourceAction.Resources,
			giveClientResourceAction.TargetClientID)
	}

	srcClientInfo.Resources -= giveClientResourceAction.Resources
	trgClientInfo.Resources += giveClientResourceAction.Resources

	g.ClientInfos[giveClientResourceAction.SourceClientID] = srcClientInfo
	g.ClientInfos[giveClientResourceAction.TargetClientID] = trgClientInfo

	g.logf("%v gave %v resources to %v", giveClientResourceAction.SourceClientID,
		giveClientResourceAction.Resources,
		giveClientResourceAction.TargetClientID)

	return nil
}

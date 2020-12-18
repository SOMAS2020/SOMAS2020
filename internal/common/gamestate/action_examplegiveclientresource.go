package gamestate

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/pkg/errors"
)

// THIS IS JUST AN EXAMPLE - lhl2617

func init() {
	registerAction(ExampleGiveClientResource, "ExampleGiveClientResource", dispatchExampleGiveClientResource)
}

// ExampleGiveClientResourceAction describes an action of giving resources from a client to the other
type ExampleGiveClientResourceAction struct {
	SourceClientID shared.ClientID
	TargetClientID shared.ClientID
	Resources      int
}

func dispatchExampleGiveClientResource(g GameState, action Action) (GameState, error) {
	giveClientResourceAction := action.ExampleGiveClientResourceAction
	if giveClientResourceAction == nil {
		return g, errors.Errorf("giveClientResourceAction is nil!")
	}
	srcClientInfo := g.ClientInfos[giveClientResourceAction.SourceClientID]
	trgClientInfo := g.ClientInfos[giveClientResourceAction.TargetClientID]

	if srcClientInfo.Resources < giveClientResourceAction.Resources {
		return g, errors.Errorf("%v does not have %v resources to give to %v",
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

	return g, nil
}

package gamestate

import (
	"fmt"

	"github.com/pkg/errors"
)

// Action represents an action that modifies game state
type Action struct {
	ActionType                      ActionType
	ExampleGiveClientResourceAction *ExampleGiveClientResourceAction
}

// ActionType is an enum that classifies actions by their type
type ActionType int

// actionDispatcher is a type for action dispatcher functions that take in the current GameState
// and action then attempt to apply the action and return the resulting GameState.
type actionDispatcher = func(g GameState, a Action) (GameState, error)

const (
	// ExampleGiveClientResource : give resource from a client to the other
	ExampleGiveClientResource ActionType = iota

	// ADD NEW ActionTypes HERE. REMEMBER TO CALL RegisterAction for your actions in your
	// action file's init function!
	// Each action MUST be in a file with naming convention action_<action-name>.go.
	// A TEST WILL FAIL IF YOU DO NOT FOLLOW THIS!

	// DO NOT TOUCH THIS! USED FOR TESTING!
	actionTypeEnd
)

var actionTypeStringMap = map[ActionType]string{}
var actionDispatcherMap = map[ActionType]actionDispatcher{}

// registerAction registers an action.
func registerAction(actionType ActionType, actionStringName string, dispatcher actionDispatcher) {
	if _, ok := actionTypeStringMap[actionType]; ok {
		panic(fmt.Sprintf("Action '%v' already registered", actionType))
	}
	if _, ok := actionDispatcherMap[actionType]; ok {
		panic(fmt.Sprintf("Action '%v' already registered", actionType))
	}
	actionTypeStringMap[actionType] = actionStringName
	actionDispatcherMap[actionType] = dispatcher
}

// DispatchActions runs all actions and alters GameState.
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
	if dispatcher, ok := actionDispatcherMap[action.ActionType]; ok {
		gotG, err := dispatcher(g.Copy(), action)
		if err != nil {
			return err
		}
		*g = gotG
		return nil
	}
	return errors.Errorf("Unknown action type--was this ActionType registered in GameState.dispatchAction? : %v", action.ActionType)
}

// Package team4 contains code for team 4's client implementation
package team4

import (
	"fmt"
	"reflect"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"gonum.org/v1/gonum/mat"
)

const id = shared.Team4

func init() {
	baseclient.RegisterClient(id, &client{BaseClient: baseclient.NewClient(id)})
}

type client struct {
	*baseclient.BaseClient //this line gives us access to the methods of the Client interface(not the same as the BaseClient struct) in baseclient.go
	// id                shared.ClientID
	// clientGameState   gamestate.ClientGameState
	// communications    map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent
	// serverReadHandle  baseclient.ServerReadHandle

	//custom fields
	yes                string //this field is just for testing
	internalStateField *internalState
}

type internalState struct {
	trustMatrix *mat.Dense
}

func (c *client) Initialise(serverReadHandle baseclient.ServerReadHandle) {
	c.ServerReadHandle = serverReadHandle

	//custom things below, just a random matrix for now
	v := []float64{1, 0, 0, 0, -4, 0, -1, -1, 0, 2, 0, 0, 0, 0, 2, 0, 0, 1, 0, -1}
	c.internalStateField = &internalState{
		trustMatrix: mat.NewDense(4, 5, v),
	}
}

func (c *client) StartOfTurn() {
	c.yes = "yes"
	c.Logf(`what are you doing?
	=========================================
	==============================================
	================================================`)
	c.Logf("this is a %v for you ", c.yes)
	fmt.Println(reflect.TypeOf(c))
}

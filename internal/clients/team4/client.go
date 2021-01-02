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
	*baseclient.BaseClient //client struct has access to methods and fields of the BaseClient struct which implements implicitly the Client interface.

	//custom fields
	yes         string                  //this field is just for testing
	obs         *observation            //observation is the raw input into our client
	internalRep *internalRepresentation //internal state is condensed from observation
}

type observation struct {
}

type internalRepresentation struct {
	trustMatrix *mat.Dense
}

//Overriding the Initialise method of the BaseClient to initilise the trust matrix too
func (c *client) Initialise(serverReadHandle baseclient.ServerReadHandle) {
	c.ServerReadHandle = serverReadHandle

	//custom things below, trust matrix initilised to values of 1
	numClient := len(shared.TeamIDs)
	v := make([]float64, numClient*numClient)
	for i := range v {
		v[i] = 1
	}
	c.internalRep = &internalRepresentation{
		trustMatrix: mat.NewDense(numClient, numClient, v),
	}
}

//Overriding the StartOfTurn method of the BaseClient
func (c *client) StartOfTurn() {
	c.yes = "yes"
	c.Logf(`what are you doing?
	=========================================
	==============================================
	================================================`)
	c.Logf("this is a %v for you ", c.yes)
	fmt.Println(reflect.TypeOf(c))
}

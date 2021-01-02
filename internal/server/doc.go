/*
Package server contains server-side code.

It runs sequential turns until the end of the game.
Note that turn means a "day", and a Season ends with a disaster.

The server's EntryPoint function returns a slice of historic GameStates of the game
until the end of the game.

The current structure of the turn is as follows:

	runTurn
		startOfTurn
		runOrgs
			runIIGO
			runIIFO
			runIITO
		endOfTurn
			runOrgsEndOfTurn
				runIIGOEndOfTurn
				runIIFOEndOfTurn
				runIITOEndOfTurn
			runForage
			probeDisaster
			incrementTurnAndSeason
			deductCostOfLiving
			updateIslandLivingStatus
*/
package server

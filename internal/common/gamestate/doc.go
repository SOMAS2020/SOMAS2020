/*
Package gamestate contains information about the current game state.
It also contains action code that mutates gameState.

How to: Add an action
- See [`gamestate/action_examplegiveclientresource.go`](gamestate/action_examplegiveclientresource.go) for an example.
1. Define a new `ActionType` for your action in [`gamestate/action.go:ActionType`](gamestate/action.go).
2. Make a new file for your action following the naming convention `gamestate/action_<action-name>.go`.
3. Make a new struct in your new file with the payload required for that action.
4. Extend the `Action` struct in [`gamestate/action.go:Action`](gamestate/action.go) with your payload created in 3. Make sure it is a pointer so that it is `nil`-able.
4. Make a dispatcher function in your file.
5. Register your action in the `init` function.
6. If you did all the above correctly, running `go test ./...` should give no errors.
*/
package gamestate

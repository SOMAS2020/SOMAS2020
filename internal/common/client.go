package common

// Client is a base interface to be implemented by each client struct.
type Client interface {
	Echo(s string) string
	GetID() int
	Logf(format string, a ...interface{})
}

// RegisteredClients contain all registered clients, exposed for the server.
var RegisteredClients = []Client{}

// RegisterClient registers clients into RegisteredClients
func RegisterClient(c Client) {
	RegisteredClients = append(RegisteredClients, c)
}

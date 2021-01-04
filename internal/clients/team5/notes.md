## Team 5 Client Notes

### Foraging

Currently, the server returns very limited information - only hunt type, client's contribution and client's return. Doesn't currently seem like there's any way to determine who the other participants were?

```go
s.giveResources(participantID, participantReturn, retReason)
    s.clientMap[participantID].ForageUpdate(shared.ForageDecision{
        Type:         huntReport.ForageType,
        Contribution: contribution,
    }, participantReturn)
```

Similarly, a foraging decision does not offer any insight into other participants. The server merely requests your desired forage type and the contribution:

```go
// DecideForage makes a foraging decision
// the forageContribution can not be larger than the total resources available
func (c *BaseClient) DecideForage() (shared.ForageDecision, error) {
	ft := int(math.Round(rand.Float64())) // 0 or 1 with equal prob.
	return shared.ForageDecision{
		Type:         shared.ForageType(ft),
		Contribution: shared.Resources(rand.Float64() * 5),
	}, nil
}
```

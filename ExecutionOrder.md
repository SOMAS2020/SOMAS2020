
# Execution Order

Sequence of function calls from start to end of game.

## Initialisations

| Filename | Function | Description |
| ---- | ---- | ---- |
| main.go |  main | Begins game |
| internal/server/server.go | EntryPoint| Creates deep copy of list of gamestates. Then, while game is not over, calls runTurn and keeps track of the states during the game. |
|internal/server/turn.go|gameOver| Checks at least one client is alive and we haven't reached maximum number of turns or seasons. |
|internal/server/turn.go |runTurn| Gets a start of turn update, runs the organisations and runs the end of turn procedures. |
|internal/server/turn.go|startOfTurnUpdate| Sends update of gameState to alive clients. |
|internal/common/baseclient/baseclient.go|Client.StartOfTurnUpdate| Where the client receives the updated gameState from the server. |
|internal/server/turn.go|runOrgs| Runs IIGO, IIFO, IITO. |
|internal/server/iigo.go|runIIGO| Runs IIGO with alive clients, then updates the alive islands variables. |
|internal/server/iigointernal/orchestration.go| RunIIGO| Gets available clients and current roles. If there are sufficient funds, pays the president, then judge, then speaker. If there are insufficient funds, IIGO cannot be run. The judge performs its duties (inspections), followed by the president (resource reports, taxation, allocations) and then the speaker (rules). Then, the judge declares performance. New judge, speaker, president appointed.  |
|internal/server/iigointernal/base_judge.go | withdrawPresidentSalary | Gets the presidentSalary amount from the rules and withdraws the amount from the common pool and stores the amount if there is enough, else, returns an error. *|
|internal/server.iigointernal/utilities.go| WithdrawFromCommonPool| If there are enough resources iin the common pool, subtract the resources required. Else, return error. *|
|internal/server/iigointernal/base_judge.go|sendPresidentSalary| If there is a judge, the judge pays the president, send the salary to president.  *|
|internal/server/iigointernal/base_judge.go|InspectHistory| Budget decreased by 10 (MVP amount). If there is a judge, perform historical evaluation of clients behaviour (e.g. checking contributions to the common pool). |
| internal/server/iigointernal/basejudge.go | inspectHistoryInternal | MVP implementation of inspecting clients' behaviour. Picks up all the rules currently in play and iterates through and evaluates clients' behaviour. For non-MVP, this could be a starting point for each groups' implementation. |
|internal/server/iigointernal/base_president.go|broadcastTaxation|Subtract 10 as cost of action. Calculates tax to pay with getTaxMap (this calls setTaxation amount which is to be implemented by agents) and communicates this to islands. |
|internal/server/iigointernal/base_president.go|requestAllocationRequests|Creates a map of island : common pool resource request|
|internal/server/iigointernal/base_president.go|replyAllocationRequest|Subtracts 10 as cost of action. For each alive island, allocate an amount of resouces (MVP: this is the requested amount) and communicate this to the islands. |
|internal/server/iigointernal/base_president.go|requestRuleProposal|Go through and collect all agents' rule proposals. setRuleProposals stores the rules. |
|internal/server/iigointernal/base_president.go|getRuleForSpeaker|MVP: All rules proposed passed to speaker. Non-MVP: each agent's PickRuleToVote selects rule to vote.|
|internal/server/iigointernal/base_speaker.go|setVotingResult|Ballot is opened and iigo clients vote. Then the ballot is closed. |
|internal/server/iigointernal/base_speaker.go|announceVotingResult|Can change (non-MVP) and announce result by broadcasting to islands. Then, update the rules. |
|internal/server/iigointernal/base_judge.go|declarePresidentPerformanceWrapped| Uses DelcarePresidentPerformance (overwrite for non-MVP for own implementation) to evaluate president performance and broadcast result to all islands. |
|internal/server/iigointernal/base_judge.go| declareSpeakerPerformanceWrapped| Uses DeclareSpeakerPerformance (overwrite for non-MVP for own implementation) to evaluate speaker performance and broadcast result to all islands. |
|internal/server/iigointernal/base_speaker.go|appointNextJudge|Subtract 10 as cost of action. Propose and carry out election of next judge and returns results. * |
|internal/server/iifo.go|runIIFO|Runs prediction session. |
|internal/server/iifo.go|runPredictionSession| Get islands' predictions and distribute predictions to all the islands. |
|internal/server/iifo.go|getPredictions|For alive clients, get their prediction for disaster (and foraging - non-MVP) using MakePrediction which is implemented by each client.|
|internal/server/iifo.go|distributePreictions|Share predictions with islands of your choice. Islands who receive a prediction(s) do not know who else receieved the same prediciton(s).|
|internal/server/iito.go|runIITO|MVP: gifts session only.|
|internal/server/iito.go|runGiftSession|Gets all fiftr requests and offers. Gets the result (accept/reject gift(s)) and distributes gifts accordingly. |
|internal/server/iito.go|getGiftAcceptance|Matches gift requests and offers per client, then allows clients to accept/reject offers. Then, returns result (accept/reject) of offer/request pairs. |
|internal/server/turn.go|endOfTurn|Runs the organisations end of turn, forage, rpobe for disasters, increase turn and season count, deduct cost of living from clients and check clients' alive/dead status. |
|internal/server/turn.go|runOrgsEndOfTurn|Runs end of turn for IIGO then IIFO then IITO.|
|internal/server/iigo.go|runIIGOEndOfTurn|Get tax from each alive client and add it to the common pool. Update the game state. Logs how much tax each client contributed. |
|internal/server/iifo.go|runIIFOEndOfTurn| To be implemented. |
|internal/server/iito.go|runIITOEndOfTurn| To be implemented. |
|internal/server/forage.go|runForage|Calls getForagingDecisions (this calls DecideForage which is to be implemented by the client) to get foraging participants. MVP: deer hunting only via runDeerHunt.|
|internal/server/forage.go|runDeerHunt|Creates a deer hunt for the participants and returns amount generated by foraging.|
|internal/server/disaster.go|probeDisaster|Checks if a disaster occurs this turn. |
|internal/server/turn.go|incrementTurnAndSeason|Increments turn. Increments season if a disaster occurs.|
|internal/server/turn.go|deductCostOfLiving|Deduct cost of living from each alive client then update the gamestate. |
|internal/server/turn.go|updateIslandLivingStatus|Update island status to dead is needed. |

\* The same process is repeated for the other roles .

## Questions

- internal/server/iigointernal/orchestration.go - what happens if president and judge can be paid but speaker can't be? How are the salaries that are taken from the pool and stored returned back to the common pool? IIGO is exited before they can do their job.
- internal/server/iigointernal/base_judge.go  - error(?) in withdrawPresidentSalary - if withdrawError == nil then pay (currently !=nil then pay) line 36
- internal/server/iigo_internal/base_president replyAllocationRequest - where do we put in our agent's implementation for allocating resources?
- Why do you need to overwrite UpdateGiftInfo in baseclient.go ?
- from deductCostOfLiving - how are critical islands accounted for?

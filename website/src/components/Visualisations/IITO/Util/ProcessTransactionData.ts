import {
    Link,
    Node,
    Transaction,
    OutputJSONType,
    Team,
} from '../../../../consts/types'

export const getIIGOTransactions = (data: OutputJSONType) => {
    const acc: Transaction[] = []

    // Since IIGOHistories is repeated, take the one from the LAST GameState and
    // do Object.entries to make it iterable. List of array'ed tuples.
    const IIGOHistories = Object.entries(
        data.GameStates[data.GameStates.length - 1].IIGOHistory
    )
    // For each of these arrayed tuples, we have [turnNumber: <"pair events">[]]
    IIGOHistories.forEach(([_, exchanges]) => {
        if (exchanges) {
            exchanges.forEach((teamAction) => {
                const type = teamAction.Pairs[0].VariableName
                let transaction: Transaction
                // There are three types of transactions
                // the target could be the client id depending on the type of team action
                // else accounts for SanctionPaid and IslandTaxContribution
                if (type === 'IslandAllocation') {
                    transaction = {
                        from: Team.CommonPool,
                        to: Team[teamAction.ClientID as keyof typeof Team],
                        amount: teamAction.Pairs[0].Values[0],
                    }
                } else {
                    transaction = {
                        from: Team[teamAction.ClientID as keyof typeof Team],
                        to: Team.CommonPool,
                        amount: teamAction.Pairs[0].Values[0],
                    }
                }
                if (transaction.amount) acc.push(transaction)
            })
        }
    })
    return acc
}

// islandGifts should get a list of IITO Transactions that happened in that turn.
// TODO: Try getting the newly written types to fit these functions
export const getIITOTransactions = (data: OutputJSONType) => {
    return data.GameStates.map((turnState) => {
        // Guard to prevent crashing on the first turn where it's undefined
        if (turnState.IITOTransactions) {
            // map over the IITO transactions in this turn
            return (
                Object.entries(turnState.IITOTransactions)
                    .map(
                        // map over each giftResponseDict for this team's offers
                        ([toTeam, giftResponseDict]) => {
                            const transactionsForThisIsland: Transaction[] = []
                            // iterate over the giftResponseDict and push Transactions to an accumulator
                            Object.entries(giftResponseDict).forEach(
                                ([fromTeam, response]) => {
                                    if (response) {
                                        transactionsForThisIsland.push({
                                            from:
                                                Team[
                                                    fromTeam as keyof typeof Team
                                                ],
                                            to:
                                                Team[
                                                    toTeam as keyof typeof Team
                                                ],
                                            amount: response.AcceptedAmount,
                                        })
                                    }
                                }
                            )

                            return transactionsForThisIsland
                        }
                    )
                    // fold the island transaction lists together for this turn
                    .reduce((acc, nextLst) => acc.concat(nextLst))
            )
        }

        return []
        // fold all turns together once more to get the whole game
    }).reduce((acc, nextLst) => acc.concat(nextLst))
}

function processTransactionData(data: OutputJSONType) {
    let nodes: Node[] = []
    let links: Link[] = []

    // let allTransactions = getIIGOTransactions(data).concat(getIITOTransactions(data));
    const allTransactions = getIIGOTransactions(data)

    if (allTransactions) {
        links = allTransactions.map((item) => {
            return {
                source: item.from,
                target: item.to,
            }
        })
    }

    // TODO: fix logic error here
    // TODO: absolute it and make them green or red
    // input range: 0 - 1000
    // output range: 10 - 100
    function normaliseMag(val, xMax, xMin, yMax, yMin) {
        return ((val - yMin) / (yMax - yMin)) * (xMax - xMin) + xMin
    }

    const bubbleIds = [0, 1, 2, 3, 4, 5, 6]

    nodes = bubbleIds
        .map((team) =>
            allTransactions
                .filter(
                    (transaction) =>
                        transaction.from === team || transaction.to === team
                )
                .reduce<number>(
                    (acc, curr) =>
                        curr.to === team
                            ? acc + curr.amount
                            : acc - curr.amount,
                    0
                )
        )
        .map((mag, teamNo) => {
            return {
                id: teamNo,
                color: mag < 0 ? 'red' : 'green',
                magnitude: normaliseMag(Math.abs(mag), 150, 5, 1000, 0),
            }
        })

    // remove the duplicate elements of zero magnitude
    return { nodes, links }
}

export default processTransactionData

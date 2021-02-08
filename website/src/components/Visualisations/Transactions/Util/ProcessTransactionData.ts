import { Transaction, OutputJSONType, TeamName } from '../../../../consts/types'
import { Link, Node } from './ForceGraph'
import { teamColors } from '../../utils'

export const getIIGOTransactions = (data: OutputJSONType) => {
    const acc: Transaction[] = []

    // Since IIGOHistories is repeated, take the one from the LAST GameState and
    // do Object.entries to make it iterable. List of array'ed tuples.
    const IIGOHistory = Object.entries(
        data.GameStates[data.GameStates.length - 1].IIGOHistory
    )
    // For each of these arrayed tuples, we have [turnNumber: <"pair events">[]]
    IIGOHistory.forEach(([turnNumber, exchanges]) => {
        if (exchanges) {
            exchanges.forEach((teamAction) => {
                const type = teamAction.Pairs[0].VariableName
                let transaction: Transaction | undefined
                // There are three types of transactions
                // the target could be the client id depending on the type of team action
                // else accounts for SanctionPaid and IslandTaxContribution
                switch (type) {
                    case 'IslandAllocation':
                    case 'AllocationMade':
                        transaction = {
                            from: TeamName.CommonPool,
                            to:
                                TeamName[
                                    teamAction.ClientID as keyof typeof TeamName
                                ],
                            amount: teamAction.Pairs[0].Values[0],
                        }
                        break
                    case 'SpeakerPayment':
                        transaction = {
                            from:
                                TeamName[
                                    teamAction.ClientID as keyof typeof TeamName
                                ],
                            to:
                                TeamName[
                                    data.GameStates[Number(turnNumber)]
                                        .SpeakerID as keyof typeof TeamName
                                ],
                            amount: teamAction.Pairs[0].Values[0],
                        }
                        break
                    case 'JudgePayment':
                        transaction = {
                            from:
                                TeamName[
                                    teamAction.ClientID as keyof typeof TeamName
                                ],
                            to:
                                TeamName[
                                    data.GameStates[Number(turnNumber)]
                                        .JudgeID as keyof typeof TeamName
                                ],
                            amount: teamAction.Pairs[0].Values[0],
                        }
                        break
                    case 'PresidentPayment':
                        transaction = {
                            from:
                                TeamName[
                                    teamAction.ClientID as keyof typeof TeamName
                                ],
                            to:
                                TeamName[
                                    data.GameStates[Number(turnNumber)]
                                        .PresidentID as keyof typeof TeamName
                                ],
                            amount: teamAction.Pairs[0].Values[0],
                        }
                        break
                    case 'IslandTaxContribution':
                    case 'SanctionPaid':
                        transaction = {
                            from:
                                TeamName[
                                    teamAction.ClientID as keyof typeof TeamName
                                ],
                            to: TeamName.CommonPool,
                            amount: teamAction.Pairs[0].Values[0],
                        }
                        break
                    default:
                        transaction = undefined
                        break
                }
                if (transaction?.amount) acc.push(transaction)
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
                                                TeamName[
                                                    fromTeam as keyof typeof TeamName
                                                ],
                                            to:
                                                TeamName[
                                                    toTeam as keyof typeof TeamName
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
                    .reduce((acc, nextLst) => acc.concat(nextLst), [])
            )
        }

        return []
        // fold all turns together once more to get the whole game
    }).reduce((acc, nextLst) => acc.concat(nextLst), [])
}

function processTransactionData(data: OutputJSONType) {
    let nodes: Node[] = []
    let links: Link[] = []

    const allTransactions = getIIGOTransactions(data).concat(
        getIITOTransactions(data)
    )
    // const allTransactions = getIIGOTransactions(data)

    if (allTransactions) {
        links = allTransactions.map((item) => {
            return {
                source: item.from,
                target: item.to,
                amount: item.amount,
            }
        })
    }

    // input range: 0 - 1000
    // output range: 10 - 100
    function normaliseMag(val, xMax, xMin, yMax, yMin) {
        return ((val - yMin) / (yMax - yMin)) * (xMax - xMin) + xMin
    }

    function getRandomColor() {
        const letters = '0123456789ABCDEF'
        let color = '#'
        for (let i = 0; i < 6; i++) {
            color += letters[Math.floor(Math.random() * 16)]
        }
        return color
    }

    const bubbleIds = [0, 1, 2, 3, 4, 5, 6]
    let maxMag = 1000
    const magnitudes = bubbleIds.map((team) =>
        allTransactions
            .filter(
                (transaction) =>
                    transaction.from === team || transaction.to === team
            )
            .reduce<number>(
                (acc, curr) =>
                    curr.to === team ? acc + curr.amount : acc - curr.amount,
                0
            )
    )

    magnitudes.forEach((mag) => {
        maxMag = mag > maxMag ? mag : maxMag
    })

    nodes = magnitudes.map((mag, teamNo) => {
        const thisTeamColor = teamColors.get(`Team${teamNo}`)
        return {
            id: teamNo,
            colorStatus: mag < 0 ? 'red' : 'green',
            islandColor: thisTeamColor ?? '#031927',
            magnitude: normaliseMag(Math.abs(mag), 125, 15, maxMag, 0),
        }
    })

    // remove the duplicate elements of zero magnitude
    return { nodes, links }
}

export default processTransactionData

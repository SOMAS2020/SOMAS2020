import { Transaction, OutputJSONType, TeamName } from '../../../../consts/types'
import { Link, Node } from './ForceGraph'
import { teamColors } from '../../utils'

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

    const allTransactions = getIITOTransactions(data)
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

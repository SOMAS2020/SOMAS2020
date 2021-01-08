import { OutputJSONType, Team, GiftResponse } from '../../../../consts/types'

// TODO: transactions can also be with the common pool so they should not be teams - maybe more appropriate to rename entities
type Transaction = {
    from: Team
    to: Team
    amount: number
}

// TODO: Decide on node structure (i.e. what determines bubble size)
// TODO: Extract summary metric for bubble size from transactions[] and islandGifts[]
// TODO: might be cool to have max and min resources of each entity as a summary metric in the tooltip
type Node = {
    id: number
    magnitude: number
    color: string
}

type Link = {
    source: number
    target: number
}

export const getIIGOTransactions = (data: OutputJSONType) => {
    let acc: Transaction[] = []

    // Since IIGOHistories is repeated, take the one from the LAST GameState and
    // do Object.entries to make it iterable. List of array'ed tuples.
    let IIGOHistories = Object.entries(
        data.GameStates[data.GameStates.length - 1].IIGOHistory
    )
    // For each of these arrayed tuples, we have [turnNumber: <"pair events">[]]
    IIGOHistories.forEach(([_, exchanges]) => {
        if (exchanges) {
            exchanges.forEach((teamAction) => {
                let type = teamAction.Pairs[0].VariableName
                let transaction: Transaction
                // There are three types of transactions
                // the target could be the client id depending on the type of team action
                // else accounts for SanctionPaid and IslandTaxContribution
                if (type === "IslandAllocation") {
                    transaction = {
                        from: Team.CommonPool,
                        to: Team[teamAction.ClientID as keyof typeof Team],
                        amount: teamAction.Pairs[0].Values[0],
                    };
                } else {
                    transaction = {
                        from: Team[teamAction.ClientID as keyof typeof Team],
                        to: Team.CommonPool,
                        amount: teamAction.Pairs[0].Values[0],
                    };
                }
                transaction.amount ? acc.push(transaction) : console.log("no bueno");
            });
        }
    });
    return acc;
};

// islandGifts should get a list of IITO Transactions that happened in that turn.
// TODO: Try getting the newly written types to fit these functions
export const getIITOTransactions = (data: OutputJSONType) => {
    return data.GameStates.map((turnState) => {
        // Guard to prevent crashing on the first turn where it's undefined
        if (turnState.IITOTransactions) {
            // map over the IITO transactions in this turn
            return Object.entries(turnState.IITOTransactions)
                .map(
                    // map over each giftResponseDict for this team's offers
                    ([toTeam, giftResponseDict]) => {
                        let transactionsForThisIsland: Transaction[] = [];
                        // iterate over the giftResponseDict and push Transactions to an accumulator
                        Object.entries(giftResponseDict).forEach(([fromTeam, response]) => {
                            if (response) {
                                transactionsForThisIsland.push({
                                    from: Team[fromTeam as keyof typeof Team],
                                    to: Team[toTeam as keyof typeof Team],
                                    amount: response.AcceptedAmount,
                                });
                            }
                        });
                        return transactionsForThisIsland
                    }
                )
                // fold the island transaction lists together for this turn
                .reduce((acc, nextLst) => acc.concat(nextLst));
        } else return [];
        // fold all turns together once more to get the whole game
    }).reduce((acc, nextLst) => acc.concat(nextLst));
};

// TODO: Decide on link representation (do we show all links, do we collate them and use colours or thickness)
function processTransactionData(data: OutputJSONType) {
    // map over turns, map through IIGOHistories, map through sub turns, extract allocations and TaxContributions

    // We want to construct the node array of Teams and their total resources traded (in/out)
    // For now this is being used to determine bubble size  
    var nodes: Node[] = [];
    var nodesNew: Node[] = [];
    var links: Link[] = [];

    // sum of all transactions for each team
    // construct the nodes array and links array (source and target)

    // let allTransactions = getIIGOTransactions(data).concat(getIITOTransactions(data));
    let allTransactions = getIIGOTransactions(data);
    console.log(allTransactions);
    if (allTransactions) {
        links = allTransactions.map(item => {
            return {
                source: item.from,
                target: item.to
            }
        })
    }


    // TODO: fix logic error here
    // TODO: absolute it and make them green or red
    // input range: 0 - 1000
    // output range: 10 - 100

    function normaliseMag(val, xMax, xMin, yMax, yMin) {
        return ((val - yMin) / (yMax - yMin)) * (xMax - xMin) + xMin;
    }

    let bubbleIds = [0, 1, 2, 3, 4, 5, 6]

    // Object.values(Team);
    console.log({ bubbleIds });

    nodes = bubbleIds.map(team =>
        allTransactions.filter(
            transaction => (transaction.from === team) || (transaction.to === team)
        ).reduce<number>((acc, curr) =>
            curr.to === team ? (acc + curr.amount) : (acc - curr.amount), 0)
    ).map((mag, teamNo) => {
        return {
            id: teamNo,
            color: mag < 0 ? "red" : "green",
            magnitude: normaliseMag(Math.abs(mag), 150, 5, 1000, 0)
        }
    })

    // remove the duplicate elements of zero magnitude
    return { nodes, links }
}

export default processTransactionData;

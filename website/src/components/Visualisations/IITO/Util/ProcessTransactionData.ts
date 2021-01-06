// Process Output JSON to extract Transaction data (inc. Gifts and common pool interactions)
// JSON parse the output here
import { OutputJSONType, Team, GiftResponse } from "../../../../consts/types";


// TODO: transactions can also be with the common pool so they should not be teams
type Transaction = {
    from: Team,
    to: Team,
    amount: number,
}

// TODO: Decide on node structure (i.e. what determines bubble size)
// TODO: Extract summary metric for bubble size from transactions[] and islandGifts[]
type Node = {
    entity: Team,
    amount: number
}

// TODO: Decide on link representation (do we show all links, do we collate them and use colours or thickness)
function processTransactionData(data: OutputJSONType) {
    // make one for the common pool and then have the N islands
    let bubbleIds = new Set(data.AuxInfo.TeamIDs.map(teamName => teamName));
    // map over turns, map through IIGOHistories, map through sub turns, extract allocations and TaxContributions

    let transactions: Transaction[] = Object.entries(data.GameStates[data.GameStates.length - 1].IIGOHistory).map(([turnNum, exchanges]: [string, Array<any>]) => {
        exchanges.map(teamAction => {
            let type = teamAction.Pairs[0].VariableName
            let transaction: Transaction;
            // There are three types of transactions
            // the target could be the client id depending on the type of team action
            // else accounts for SanctionPaid and IslandTaxContribution
            if (type === "IslandAllocation") {
                transaction = {
                    from: "CommonPool",
                    to: teamAction.ClientID,
                    amount: teamAction.Pairs[0].Values[0]
                };
            } else {
                transaction = {
                    from: teamAction.ClientID,
                    to: "CommonPool",
                    amount: teamAction.Pairs[0].Values[0]
                };
            }

            transactions.push(transaction);
        })
    })

    // islandGifts should get a list of IITO Transactions that happened in that turn.
    // TODO: Try getting the newly written types to fit these functions
    let islandGifts = data.GameStates.map(turnState => {
        if (turnState.IITOTransactions) {
            let thisIslandTransactions =
                Object.entries(turnState.IITOTransactions).map(([fromTeam, giftResponse]: [Team, any], toTeam: Team) => {
                    return {
                        from: fromTeam, to: toTeam, amount: giftResponse.AcceptedAmount,
                    }
                })

        } else {
            let arr: Array<Transaction> = []
            return arr;
        }

        Object.entries(turnState.IITOTransactions).map(team => {
            // each team will have a series of transactions

        })
    })

    // We want to construct the node array of Teams and their total resources traded
    let nodes = {};

    islandTaxes.map(([fromTeam, amount]) => {
        nodes.push(
            id:
        )
    })

}

module.exports = {
    constructNetwork,
};

// let islandTaxes = data.GameStates.map(turnState => {
    //     Object.entries(turnState.IIGOHistory).map(([turnNum, exchanges]: [string, Array<any>]) => {
    //         exchanges.map(teamAction => {
    //             let team = teamAction.ClientID
    //             let transactionPair = teamAction.Pairs[0].Values[0]

    //         })
    //     })
    // });
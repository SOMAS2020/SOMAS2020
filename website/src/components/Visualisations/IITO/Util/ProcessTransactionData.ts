// Process Output JSON to extract Transaction data (inc. Gifts and common pool interactions)
// JSON parse the output here
import { OutputJSONType, Team, GiftResponse } from "../../../../consts/types";

// TODO: transactions can also be with the common pool so they should not be teams - maybe more appropriate to rename entities
type Transaction = {
  from: Team;
  to: Team;
  amount: number;
};

// TODO: Decide on node structure (i.e. what determines bubble size)
// TODO: Extract summary metric for bubble size from transactions[] and islandGifts[]
type Node = {
  id: number;
  magnitude: number;
  color: string;
};

type Link = {
  source: number;
  target: number;
};

export const getIIGOTransactions = (data: OutputJSONType) => {
  let acc: Transaction[] = [];

  // Since IIGOHistories is repeated, take the one from the LAST GameState and
  // do Object.entries to make it iterable. List of array'ed tuples.
  let IIGOHistories = Object.entries(
    data.GameStates[data.GameStates.length - 1].IIGOHistory
  );
  // For each of these arrayed tuples, we have [turnNumber: <"pair events">[]]
  IIGOHistories.forEach(([_, exchanges]) => {
    if (exchanges) {
      exchanges.forEach((teamAction) => {
        let type = teamAction.Pairs[0].VariableName;
        let transaction: Transaction;
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
      return Object.entries(turnState.IITOTransactions).map(
        ([fromTeam, giftResponse], toTeam: Team) => {
          return {
            from: fromTeam,
            to: toTeam,
            amount: giftResponse.AcceptedAmount,
          };
        }
      );
    } else return [];
  }).reduce((acc, nextLst) => acc.concat(nextLst));
};

// TODO: Decide on link representation (do we show all links, do we collate them and use colours or thickness)
function processTransactionData(data: OutputJSONType) {
  // map over turns, map through IIGOHistories, map through sub turns, extract allocations and TaxContributions

  // We want to construct the node array of Teams and their total resources traded (in/out)
  // For now this is being used to determine bubble size
  var nodes: Node[] = [];
  var links: Link[] = [];

  // sum of all transactions for each team
  // construct the nodes array and links array (source and target)

  // let allTransactions = getIIGOTransactions(data).concat(getIITOTransactions(data));
  let allTransactions = getIIGOTransactions(data);
  console.log(allTransactions);
  if (allTransactions) {
    links = allTransactions.map((item) => {
      return {
        source: item.from,
        target: item.to,
      };
    });
  }

  let bubbleIds = Object.values(Team);

  // First we add each of the islands to the list of nodes
  // make one for the common pool and then have the N islands

  // let allMagsPerTeam = bubbleIds.map(team =>
  //     allTransactions.filter(
  //         transaction => (transaction.from === team) || (transaction.to === team)
  //     ).reduce<number>((acc, curr) =>
  //         curr.to === team ? (acc + curr.amount) : (acc - curr.amount), 0)
  // )

  // TODO: fix logic error here
  // TODO: absolute it and make them green or red
  nodes = bubbleIds
    .map((team) =>
      allTransactions
        .filter(
          (transaction) => transaction.from === team || transaction.to === team
        )
        .reduce<number>(
          (acc, curr) =>
            curr.to === team ? acc + curr.amount : acc - curr.amount,
          0
        )
    )
    .map((mag, teamNo) => {
      return {
        id: teamNo,
        color: mag < 0 ? "red" : "green",
        magnitude: Math.max(0, mag),
      };
    });

  console.log({ nodes });
  return { nodes, links };
}

export default processTransactionData;

import outputJSONData from '../output/output.json'

// TODO: what if there are more islands - dynamic typing
export enum TeamName {
    'CommonPool',
    'Team1',
    'Team2',
    'Team3',
    'Team4',
    'Team5',
    'Team6',
}

// TODO: transactions can also be with the common pool so they should not be teams - maybe more appropriate to rename entities
export type Transaction = {
    from: TeamName
    to: TeamName
    amount: number
}

// TODO: Decide on link representation (do we show all links, do we collate them and use colours or thickness)
// TODO: Decide on node structure (i.e. what determines bubble size)
// TODO: Extract summary metric for bubble size from transactions[] and islandGifts[]
// TODO: might be cool to have max and min resources of each entity as a summary metric in the tooltip
export type Node = {
    id: number
    magnitude: number
    color: string
}

export type Link = {
    source: number
    target: number
}

export type OutputJSONType = GameStatesType & typeof outputJSONData

type GameStatesType = {
    GameStates: {
        IIGOHistory: {
            [turn: string]:
                | {
                      ClientID: string
                      Pairs: {
                          VariableName: string
                          Values: number[]
                      }[]
                  }[]
                | undefined
        }
        IITOTransactions:
            | {
                  [offerTeam in TeamName]: {
                      [receiveTeam in TeamName]?: {
                          AcceptedAmount: number
                          Reason: number
                      }
                  }
              }
            | undefined
    }[]
}

type Overwrite<T, U> = Pick<T, Exclude<keyof T, keyof U>> & U

type StrongTypedGameState = Overwrite<
    typeof outputJSONData.GameStates,
    {
        IIGOHistory: IIGOHistory
        IITOTransactions: IITOTransactions
    }
>

type StrongTypedOutputJSONData = Overwrite<OutputJSONType, StrongTypedGameState>

// IIGO TYPES

// IIGOHistory will be at most data.Config.Maxturns long, containing an "Accountability" occurrence for a given client.
// Returns undefined if the accessing an unavailable key
export type IIGOHistory = {
    [key: number]: Accountability
}

export type Accountability = {
    ClientID: TeamName
    Pairs: VariableValuePair[]
}

export type VariableValuePair = {
    VariableName: string
    Values: number[]
}

// IITO Types
export type IITOTransactions = {
    [key: number]: GiftResponseDict
}

export type GiftResponseDict = {
    [team: number]: GiftResponse
}

export type GiftResponse = {
    AcceptedAmount: number
    Reason: number
}

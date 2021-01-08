import outputJSONData from '../output/output.json'

// TODO: what if there are more islands - dynamic typing
// eslint-disable-next-line no-shadow
export enum Team {
    'CommonPool',
    'Team1',
    'Team2',
    'Team3',
    'Team4',
    'Team5',
    'Team6',
}

export type Transaction = {
    from: Team
    to: Team
    amount: number
}

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

// Custom utility type
type Overwrite<T, U> = Pick<T, Exclude<keyof T, keyof U>> & U
// export type OutputJSONType = GameStatesType & typeof outputJSONData

export type OverwriteGameState = Overwrite<
    typeof outputJSONData.GameStates,
    GameStatesType
>

export type OutputJSONType = Overwrite<
    typeof outputJSONData,
    OverwriteGameState
>

type GameStatesType = {
    GameStates: {
        IIGOHistory: IIGOHistory | undefined
        IITOTransactions:
            | {
                  [offerTeam in Team]: {
                      [receiveTeam in Team]?: {
                          AcceptedAmount: number
                          Reason: number
                      }
                  }
              }
            | undefined
    }[]
}

// IIGOHistory will be at most data.Config.Maxturns long, containing an "Accountability" occurrence for a given client.
// Returns undefined if the accessing an unavailable key
export type IIGOHistory = {
    [key: number]: Accountability | undefined
}

export type Accountability = {
    ClientID: Team
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

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

export type Transaction = {
    from: TeamName
    to: TeamName
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
// type Overwrite<T, U> = Pick<T, Exclude<keyof T, keyof U>> & U
// export type OutputJSONType = GameStatesType & typeof outputJSONData
type ChangeKeyType<T, K extends keyof T, U> = Omit<T, K> & { [k in K]?: U }

export type OutputJSONType = ChangeKeyType<
    typeof outputJSONData,
    'GameStates',
    GameStatesType
>

// export type OverwriteGameState = Overwrite<
//     typeof outputJSONData.GameStates,
//     GameStatesType
// >

// export type OutputJSONType = Overwrite<
//     typeof outputJSONData,
//     OverwriteGameState
// >

type GameStatesType = {
    GameStates: {
        Season: number
        Turn: number
        ClientInfos: any
        Environment: any
        DeerPopulation: any
        ForagingHistory: any
        CurrentRulesInPlay: any
        IIGORolesBudget: any
        IIGOTurnsInPower: any
        IIGOCache: any
        SpeakerID: string
        JudgeID: string
        PresidentID: string
        IIGOHistory: IIGOHistory | undefined
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

// IIGOHistory will be at most data.Config.Maxturns long, containing an "Accountability" occurrence for a given client.
// Returns undefined if the accessing an unavailable key
export type IIGOHistory = {
    [key: number]: Accountability | undefined
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

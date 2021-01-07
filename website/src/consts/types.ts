import outputJSONData from '../output/output.json'

export type OutputJSONType = typeof outputJSONData

// TODO: what if there are more islands - dynamic typing
export enum Team {
    "CommonPool",
    "Team1",
    "Team2",
    "Team3",
    "Team4",
    "Team5",
    "Team6",
}

type Overwrite<T, U> = Pick<T, Exclude<keyof T, keyof U>> & U;

type StrongTypedGameState = Overwrite<typeof outputJSONData.GameStates, { IIGOHistory: IIGOHistory, IITOTransactions: IITOTransactions }>

type StrongTypedOutputJSONData = Overwrite<OutputJSONType, StrongTypedGameState>

// IIGO TYPES

// IIGOHistory will be at most data.Config.Maxturns long, containing an "Accountability" occurrence for a given client.
// Returns undefined if the accessing an unavailable key
export type IIGOHistory = {
    [key: number]: Accountability,
}

export type Accountability = {
    ClientID: Team,
    Pairs: VariableValuePair[]
}

export type VariableValuePair = {
    VariableName: string,
    Values: number[],
}

// IITO Types
export type IITOTransactions = {
    [key: number]: GiftResponseDict,
}

export type GiftResponseDict = {
    [team: number]: GiftResponse
}

export type GiftResponse = {
    AcceptedAmount: number,
    Reason: number,
}

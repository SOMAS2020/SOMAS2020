import outputJSONData from '../output/output.json'

// TODO: what if there are more islands - dynamic typing
// eslint-disable-next-line no-shadow
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
// Custom utility type
// type Overwrite<T, U> = Pick<T, Exclude<keyof T, keyof U>> & U
// export type OutputJSONType = GameStatesType & typeof outputJSONData
type ChangeKeyType<T, K extends keyof T, U> = Omit<T, K> & { [k in K]: U }

export type OutputJSONType = ChangeKeyType<
    typeof outputJSONData,
    'GameStates',
    GameStates
>
type Config = {
    MaxSeasons: number
    MaxTurns: number
    InitialResources: number
    InitialCommonPool: number
    CostOfLiving: number
    MinimumResourceThreshold: number
    MaxCriticalConsecutiveTurns: number
    ForagingConfig: ForagingConfig
    DisasterConfig: DisasterConfig
    IIGOConfig: any
}

type ForagingConfig = {
    DeerHuntConfig: DeerHuntConfig
    FishingConfig: FishingConfig
}

type DeerHuntConfig = {
    MaxDeerPerHunt: number
    IncrementalInputDecay: number
    BernoulliProb: number
    ExponentialRate: number
    InputScaler: number
    OutputScaler: number
    DistributionStrategy: string
    ThetaCritical: number
    ThetaMax: number
    MaxDeerPopulation: number
    DeerGrowthCoefficient: number
}
type FishingConfig = {
    MaxFishPerHunt: number
    IncrementalInputDecay: number
    Mean: number
    Variance: number
    InputScaler: number
    OutputScaler: number
    DistributionStrategy: string
}
type DisasterConfig = {
    XMin: number
    XMax: number
    YMin: number
    YMax: number
    Period: number
    SpatialPDFType: string
    MagnitudeLambda: number
    MagnitudeResourceMultiplier: number
    CommonpoolThreshold: number
    StochasticPeriod: boolean
    CommonpoolThresholdVisible: boolean
    PeriodVisible: boolean
    StochasticPeriodVisible: boolean
}
type GameStates = GameState[]
type GameState = {
    Season: number
    Turn: number
    CommonPool: number
    ClientInfos: any
    Environment: Environment
    DeerPopulation: any
    ForagingHistory: ForagingHistory
    RulesInfo: any
    IIGOHistory: IIGOHistory
    IIGORolesBudget: any
    IIGOTurnsInPower: any
    IIGOTaxAmount: any
    IIGOAllocationMap: any
    IIGOSanctionMap: any
    IIGOSanctionCache: any
    IIGORunStatus: any
    IIGOHistoryCache: any
    IIGORoleMonitoringCache: any
    IITOTransactions: IITOTransactions
    SpeakerID: string
    JudgeID: string
    PresidentID: string
    RulesBrokenByIslands: RulesBrokenByIslands | null
}

// IIGOHistory will be at most data.Config.Maxturns long, containing an "Accountability" occurrence for a given client.
// Returns undefined if the accessing an unavailable key
export type IIGOHistory = {
    [key: string]: Accountability[] | undefined
}

export type Accountability = {
    ClientID: string
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

export type ForagingHistory = {
    DeerForageType: ForagingInfo[]
    FishForageType: ForagingInfo[]
}

export type ForagingInfo = {
    ForageType: string
    InputResources: number
    ParticipantContributions: ParticipantContributions
    NumberCaught: number
    TotalUtility: number
    CatchSizes: number[]
    Turn: number
}
export type ParticipantContributions = {
    [team: number]: number
}

export type RulesBrokenByIslands = {
    [team: number]: string[]
}

export type Environment = {
    Geography: any
    LastDisasterReport: LastDisasterReport
}

export type LastDisasterReport = {
    Magnitude: number
    X: number
    Y: number
    Effects: DisasterEffects
}
export type DisasterEffects = {
    Absolute: IslandMap | null
    Proportional: IslandMap | null
    CommonPoolMitigated: IslandMap | null
}

export type IslandMap = {
    [team: number]: number
}

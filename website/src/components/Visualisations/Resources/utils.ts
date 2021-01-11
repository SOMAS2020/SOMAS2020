import { OutputJSONType } from '../../../consts/types'
import { notUndefined } from '../../../utils/filter'

export type ResourceLevel = {
    team1: number
    team2: number
    team3: number
    team4: number
    team5: number
    team6: number
    CommonPool: number
    TotalResources: number
    Turn: number
}

export const outputToResourceLevels = (
    output: OutputJSONType
): ResourceLevel[] => {
    return output.GameStates.map((g) => {
        const cis = g.ClientInfos
        const resLevel: ResourceLevel = {
            team1: cis.Team1.Resources,
            team2: cis.Team2.Resources,
            team3: cis.Team3.Resources,
            team4: cis.Team4.Resources,
            team5: cis.Team5.Resources,
            team6: cis.Team6.Resources,
            CommonPool: g.CommonPool,
            TotalResources: 0, // set 0 for now, calculate below
            Turn: 0, // set 0 for now, set below
        }
        resLevel.TotalResources = Object.values(resLevel).reduce(
            (a, b) => a + b,
            0
        )
        resLevel.Turn = g.Turn
        return resLevel
    })
}

export const getSeasonEnds = (output: OutputJSONType): number[] => {
    return output.GameStates.map((g) => {
        if (g.Environment.LastDisasterReport.Magnitude !== 0) return g.Turn - 1
        return undefined
    }).filter(notUndefined)
}

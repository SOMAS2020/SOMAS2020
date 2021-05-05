import { OutputJSONType } from '../../../consts/types'
import { notUndefined } from '../../../utils/filter'

export const outputToResourceLevels = (
    output: OutputJSONType
): Record<string, number>[] => {
    return output.GameStates.map((g) => {
        const cis = g.ClientInfos
        const numTeams = Object.keys(cis).length

        const resLevel: Record<string, number> = {
            CommonPool: g.CommonPool,
            TotalResources: 0, // set 0 for now, calculate below
            Turn: 0, // set 0 for now, set below
        }

        for (let i = 0; i < numTeams; i++) {
            resLevel[`team${i + 1}`] = cis[`Team${i + 1}`].Resources
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

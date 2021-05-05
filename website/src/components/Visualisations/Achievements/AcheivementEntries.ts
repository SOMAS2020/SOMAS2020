import { OutputJSONType } from '../../../consts/types'
import { numAgents } from '../utils'

const emptyMetrics = (totalAgents: number): Record<string, number> => {
    const output: Record<string, number> = {}
    for (let i = 0; i < totalAgents; i++) {
        output[`Team${i + 1}`] = 0
    }
    return output
}

const teamNames = (totalAgents: number): string[] => {
    const output: string[] = []
    for (let i = 0; i < totalAgents; i++) {
        output.push(`Team${i + 1}`)
    }
    return output
}

export type AcheivementEntry = {
    title: string
    description: string
    collectMetrics: (data: OutputJSONType) => Record<string, number>
    evalLargest: boolean
}

const addMetrics = (totalAgents: number, metArr: Record<string, number>[]) =>
    metArr.reduce((metrics, currMet) => {
        teamNames(totalAgents).forEach((team) => {
            metrics[team] += currMet[team]
        })
        return metrics
    }, emptyMetrics(totalAgents))

const turnsAlive = (data: OutputJSONType, team: string): number =>
    data.GameStates.filter(
        (gameState) => gameState.ClientInfos[team].LifeStatus !== 'Dead'
    ).length

export const evaluateMetrics = (
    data: OutputJSONType,
    aEntry: AcheivementEntry
): string[] => {
    const totalAgents = numAgents(data)
    const metrics = aEntry.collectMetrics(data)
    const ret = teamNames(totalAgents).reduce(
        (maxTeams: { val: number; teams: string[] }, team: string) => {
            if (aEntry.evalLargest && metrics[team] > maxTeams.val)
                return { val: metrics[team], teams: [team] }

            if (!aEntry.evalLargest && metrics[team] < maxTeams.val)
                return { val: metrics[team], teams: [team] }

            if (metrics[team] === maxTeams.val) {
                maxTeams.val = metrics[team]
                maxTeams.teams.push(team)
            }
            return maxTeams
        },
        { val: aEntry.evalLargest ? 0 : Number.MAX_VALUE, teams: [] }
    ).teams
    return ret
}

const peakResourcesMetricCollection = (
    data: OutputJSONType
): Record<string, number> => {
    const totalAgents = numAgents(data)
    return data.GameStates.reduce(
        (metrics: Record<string, number>, gameState) =>
            teamNames(totalAgents).reduce((metAcc, teamName) => {
                const teamResources: number =
                    gameState.ClientInfos[teamName].Resources
                metAcc[teamName] =
                    teamResources > metAcc[teamName]
                        ? teamResources
                        : metAcc[teamName]
                return metAcc
            }, metrics),
        emptyMetrics(totalAgents)
    )
}

const averageResourcesMetricCollection = (
    data: OutputJSONType
): Record<string, number> => {
    const totalAgents = numAgents(data)
    const retMetrics = data.GameStates.reduce(
        (metrics: Record<string, number>, gameState) =>
            teamNames(totalAgents).reduce((metAcc, teamName) => {
                metAcc[teamName] += gameState.ClientInfos[teamName].Resources
                return metAcc
            }, metrics),
        emptyMetrics(totalAgents)
    )
    teamNames(totalAgents).forEach((team) => {
        retMetrics[team] /= turnsAlive(data, team)
    })
    return retMetrics
}

const returnsFromCriticalMetricCollection = (data: OutputJSONType) => {
    const totalAgents = numAgents(data)
    const metrics = emptyMetrics(totalAgents)
    teamNames(totalAgents).forEach((team) => {
        metrics[team] = data.GameStates.reduce(
            (status, gameState) => {
                const lifeStatus = gameState.ClientInfos[team].LifeStatus
                if (status.prev === 'Critical' && lifeStatus === 'Alive')
                    return {
                        prev: lifeStatus,
                        occurred: status.occurred + 1,
                    }
                return {
                    ...status,
                    prev: lifeStatus,
                }
            },
            {
                prev: 'Alive',
                occurred: 0,
            }
        ).occurred
    })
    return metrics
}

const turnsAliveMetricCollection = (data: OutputJSONType) => {
    const totalAgents = numAgents(data)
    const metrics = emptyMetrics(totalAgents)
    teamNames(totalAgents).forEach((team) => {
        metrics[team] = turnsAlive(data, team)
    })
    return metrics
}

const turnsAsRoleMetricCollection = (
    data: OutputJSONType,
    role: 'PresidentID' | 'SpeakerID' | 'JudgeID'
) => {
    const totalAgents = numAgents(data)
    const metrics = emptyMetrics(totalAgents)
    data.GameStates.forEach((gameState) => {
        metrics[gameState[role]]++
    })
    return metrics
}

const turnsInPowerMetricCollection = (data: OutputJSONType) => {
    const totalAgents = numAgents(data)
    return addMetrics(totalAgents, [
        turnsAsRoleMetricCollection(data, 'PresidentID'),
        turnsAsRoleMetricCollection(data, 'JudgeID'),
        turnsAsRoleMetricCollection(data, 'SpeakerID'),
    ])
}

const acheivementList: AcheivementEntry[] = [
    {
        title: 'I Will Survive',
        description: 'Island alive the longest',
        collectMetrics: turnsAliveMetricCollection,
        evalLargest: true,
    },
    {
        title: 'F',
        description: 'First island to die',
        collectMetrics: turnsAliveMetricCollection,
        evalLargest: false,
    },
    {
        title: 'Baller',
        description: 'Island with the highest average resources',
        collectMetrics: averageResourcesMetricCollection,
        evalLargest: true,
    },
    {
        title: 'Broke',
        description: 'Island with the lowest average resources',
        collectMetrics: averageResourcesMetricCollection,
        evalLargest: false,
    },
    {
        title: 'Jackpot!',
        description: 'Island with the highest peak resources',
        collectMetrics: peakResourcesMetricCollection,
        evalLargest: true,
    },
    {
        title: 'Back to Life',
        description: 'Island who returned from critical the most',
        collectMetrics: returnsFromCriticalMetricCollection,
        evalLargest: true,
    },
    {
        title: 'The Donald',
        description: 'Island who spent the most time as President',
        collectMetrics: (data) =>
            turnsAsRoleMetricCollection(data, 'PresidentID'),
        evalLargest: true,
    },
    {
        title: 'Judge Judy',
        description: 'Island who spent the most time as Judge',
        collectMetrics: (data) => turnsAsRoleMetricCollection(data, 'JudgeID'),
        evalLargest: true,
    },
    {
        title: 'Speak Now or Forever Hold Your Peace',
        description: 'Island who spent the most time as Speaker',
        collectMetrics: (data) =>
            turnsAsRoleMetricCollection(data, 'SpeakerID'),
        evalLargest: true,
    },
    {
        title: 'Power Hungry',
        description: 'Island who spent the most time in power',
        collectMetrics: turnsInPowerMetricCollection,
        evalLargest: true,
    },
]

export default acheivementList

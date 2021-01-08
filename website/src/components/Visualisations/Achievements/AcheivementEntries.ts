import { OutputJSONType } from '../../../consts/types'

export type TeamName = 'Team1' | 'Team2' | 'Team3' | 'Team4' | 'Team5' | 'Team6'

type MetricsType = {
    [key in TeamName]: number
}

const emptyMetrics = (): MetricsType => ({
    Team1: 0,
    Team2: 0,
    Team3: 0,
    Team4: 0,
    Team5: 0,
    Team6: 0,
})

const teamNames = (): TeamName[] => [
    'Team1',
    'Team2',
    'Team3',
    'Team4',
    'Team5',
    'Team6',
]

export type AcheivementEntry = {
    title: string
    description: string
    collectMetrics: (data: OutputJSONType) => MetricsType
    evalLargest: boolean
}

const addMetrics = (metArr: MetricsType[]) =>
    metArr.reduce((metrics, currMet) => {
        teamNames().forEach((team) => {
            metrics[team] += currMet[team]
        })
        return metrics
    }, emptyMetrics())

const turnsAlive = (data: OutputJSONType, team: TeamName): number =>
    data.GameStates.filter(
        (gameState) => gameState.ClientInfos[team].LifeStatus !== 'Dead'
    ).length

export const evaluateMetrics = (
    data: OutputJSONType,
    aEntry: AcheivementEntry
): TeamName[] => {
    const metrics = aEntry.collectMetrics(data)
    const ret = teamNames().reduce(
        (maxTeams: { val: number; teams: TeamName[] }, team: TeamName) => {
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

const peakResourcesMetricCollection = (data: OutputJSONType): MetricsType =>
    data.GameStates.reduce(
        (metrics: MetricsType, gameState) =>
            teamNames().reduce((metAcc, teamName) => {
                const teamResources: number =
                    gameState.ClientInfos[teamName].Resources
                metAcc[teamName] =
                    teamResources > metAcc[teamName]
                        ? teamResources
                        : metAcc[teamName]
                return metAcc
            }, metrics),
        emptyMetrics()
    )

const averageResourcesMetricCollection = (
    data: OutputJSONType
): MetricsType => {
    const retMetrics = data.GameStates.reduce(
        (metrics: MetricsType, gameState) =>
            teamNames().reduce((metAcc, teamName) => {
                metAcc[teamName] += gameState.ClientInfos[teamName].Resources
                return metAcc
            }, metrics),
        emptyMetrics()
    )
    teamNames().forEach((team) => {
        retMetrics[team] /= turnsAlive(data, team)
    })
    return retMetrics
}

const returnsFromCriticalMetricCollection = (data: OutputJSONType) => {
    const metrics = emptyMetrics()
    teamNames().forEach((team) => {
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
    const metrics = emptyMetrics()
    teamNames().forEach((team) => {
        metrics[team] = turnsAlive(data, team)
    })
    return metrics
}

const turnsAsRoleMetricCollection = (
    data: OutputJSONType,
    role: 'PresidentID' | 'SpeakerID' | 'JudgeID'
) => {
    const metrics = emptyMetrics()
    data.GameStates.forEach((gameState) => {
        metrics[gameState[role]]++
    })
    return metrics
}

const turnsInPowerMetricCollection = (data: OutputJSONType) =>
    addMetrics([
        turnsAsRoleMetricCollection(data, 'PresidentID'),
        turnsAsRoleMetricCollection(data, 'JudgeID'),
        turnsAsRoleMetricCollection(data, 'SpeakerID'),
    ])

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

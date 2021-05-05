import { OutputJSONType } from '../../../consts/types'
import { numAgents } from '../utils'

export type Metric = {
    teamName: string
    value: number
}

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

export type MetricEntry = {
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

const seasonsAlive = (data: OutputJSONType, team: string): number =>
    data.GameStates.filter(
        (gameState) =>
            gameState.ClientInfos[team].LifeStatus !== 'Dead' &&
            gameState.Environment.LastDisasterReport.Magnitude !== 0
    ).length

export const evaluateMetrics = (
    data: OutputJSONType,
    aEntry: MetricEntry
): Metric[] => {
    const metrics = Object.entries(aEntry.collectMetrics(data))
    let metric: Metric | undefined
    const acc: Metric[] = []
    metrics.forEach((team) => {
        metric = {
            teamName: team[0],
            value: team[1],
        }
        acc.push(metric)
    })
    return acc
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

const rulesBrokenPerTeamCollection = (data: OutputJSONType) => {
    const totalAgents = numAgents(data)
    const metrics = emptyMetrics(totalAgents)
    // Since IIGOHistories is repeated, take the one from the LAST GameState and
    // do Object.entries to make it iterable. List of array'ed tuples.
    data.GameStates.forEach((gameState) => {
        if (gameState.RulesBrokenByIslands !== null) {
            Object.entries(gameState.RulesBrokenByIslands).forEach((island) => {
                metrics[island[0]] += island[1].length
            })
        }
    })
    return metrics
}
const totalDisasterImpactCollection = (data: OutputJSONType) => {
    const totalAgents = numAgents(data)
    const metrics = emptyMetrics(totalAgents)
    // Since IIGOHistories is repeated, take the one from the LAST GameState and
    // do Object.entries to make it iterable. List of array'ed tuples.
    data.GameStates.forEach((gameState) => {
        if (
            gameState.Environment.LastDisasterReport.Effects
                .CommonPoolMitigated !== null
        ) {
            Object.entries(
                gameState.Environment.LastDisasterReport.Effects
                    .CommonPoolMitigated
            ).forEach((island) => {
                metrics[island[0]] += island[1]
            })
        }
    })
    return metrics
}

const totalDisasterImpactMitigatedCollection = (data: OutputJSONType) => {
    const totalAgents = numAgents(data)
    const metrics = emptyMetrics(totalAgents)
    // Since IIGOHistories is repeated, take the one from the LAST GameState and
    // do Object.entries to make it iterable. List of array'ed tuples.
    data.GameStates.forEach((gameState) => {
        if (
            gameState.Environment.LastDisasterReport.Effects.Absolute !== null
        ) {
            Object.entries(
                gameState.Environment.LastDisasterReport.Effects.Absolute
            ).forEach((island) => {
                metrics[island[0]] +=
                    island[1] *
                    data.Config.DisasterConfig.MagnitudeResourceMultiplier
            })
        }
        if (
            gameState.Environment.LastDisasterReport.Effects
                .CommonPoolMitigated !== null
        ) {
            Object.entries(
                gameState.Environment.LastDisasterReport.Effects
                    .CommonPoolMitigated
            ).forEach((island) => {
                metrics[island[0]] -= island[1]
            })
        }
    })
    return metrics
}

const deerForagingResultCollection = (data: OutputJSONType) => {
    const totalAgents = numAgents(data)
    const metrics = emptyMetrics(totalAgents)
    // Since IIGOHistories is repeated, take the one from the LAST GameState and
    // do Object.entries to make it iterable. List of array'ed tuples.
    const DeerForageHistory = Object.entries(
        data.GameStates[data.GameStates.length - 1].ForagingHistory
            .DeerForageType
    )
    // For each of these arrayed tuples, we have [turnNumber: <"pair events">[]]
    DeerForageHistory.forEach(([turn, forageRound]) => {
        Object.entries(forageRound.ParticipantContributions).forEach(
            ([team, value]) => {
                metrics[team] += value
            }
        )
    })
    return metrics
}
const fishForagingResultCollection = (data: OutputJSONType) => {
    const totalAgents = numAgents(data)
    const metrics = emptyMetrics(totalAgents)
    // Since IIGOHistories is repeated, take the one from the LAST GameState and
    // do Object.entries to make it iterable. List of array'ed tuples.
    const FishForageHistory = Object.entries(
        data.GameStates[data.GameStates.length - 1].ForagingHistory
            .FishForageType
    )
    // For each of these arrayed tuples, we have [turnNumber: <"pair events">[]]
    FishForageHistory.forEach(([turn, forageRound]) => {
        Object.entries(forageRound.ParticipantContributions).forEach(
            ([team, value]) => {
                metrics[team] += value
            }
        )
    })
    return metrics
}

const fishForagingEfficiencyCollection = (data: OutputJSONType) => {
    const totalAgents = numAgents(data)
    const metrics = emptyMetrics(totalAgents)
    const totalSpent = emptyMetrics(totalAgents)

    // Since IIGOHistories is repeated, take the one from the LAST GameState and
    // do Object.entries to make it iterable. List of array'ed tuples.
    const FishForageHistory = Object.entries(
        data.GameStates[data.GameStates.length - 1].ForagingHistory
            .FishForageType
    )
    const FishingDistributionStrategy =
        data.Config.ForagingConfig.FishingConfig.DistributionStrategy
    // For each of these arrayed tuples, we have [turnNumber: <"pair events">[]]
    FishForageHistory.forEach(([turn, forageRound]) => {
        const Participants = Object.entries(
            forageRound.ParticipantContributions
        )
        Participants.forEach(([team, value]) => {
            if (FishingDistributionStrategy === 'EqualSplit') {
                metrics[team] += forageRound.TotalUtility / Participants.length
            } else {
                metrics[team] +=
                    forageRound.TotalUtility /
                    (forageRound.InputResources / value)
            }
            totalSpent[team] += value
        })
    })
    Object.entries(totalSpent).forEach(([team, spent]) => {
        if (spent > 0) {
            metrics[team] /= spent
        }
    })
    return metrics
}

const deerForagingEfficiencyCollection = (data: OutputJSONType) => {
    const totalAgents = numAgents(data)
    const metrics = emptyMetrics(totalAgents)
    const totalSpent = emptyMetrics(totalAgents)

    // Since IIGOHistories is repeated, take the one from the LAST GameState and
    // do Object.entries to make it iterable. List of array'ed tuples.
    const DeerForageHistory = Object.entries(
        data.GameStates[data.GameStates.length - 1].ForagingHistory
            .DeerForageType
    )
    const HuntingDistributionStrategy =
        data.Config.ForagingConfig.DeerHuntConfig.DistributionStrategy
    // For each of these arrayed tuples, we have [turnNumber: <"pair events">[]]
    DeerForageHistory.forEach(([turn, forageRound]) => {
        const Participants = Object.entries(
            forageRound.ParticipantContributions
        )
        Participants.forEach(([team, value]) => {
            if (HuntingDistributionStrategy === 'EqualSplit') {
                metrics[team] += forageRound.TotalUtility / Participants.length
            } else {
                metrics[team] +=
                    forageRound.TotalUtility /
                    (forageRound.InputResources / value)
            }
            totalSpent[team] += value
        })
    })
    Object.entries(totalSpent).forEach(([team, spent]) => {
        if (spent > 0) {
            metrics[team] /= spent
        }
    })
    return metrics
}

const totalForagingEfficiencyCollection = (data: OutputJSONType) => {
    const totalAgents = numAgents(data)
    const metrics = emptyMetrics(totalAgents)
    const totalSpent = emptyMetrics(totalAgents)

    // Since IIGOHistories is repeated, take the one from the LAST GameState and
    // do Object.entries to make it iterable. List of array'ed tuples.
    const DeerForageHistory = Object.entries(
        data.GameStates[data.GameStates.length - 1].ForagingHistory
            .DeerForageType
    )
    const HuntingDistributionStrategy =
        data.Config.ForagingConfig.DeerHuntConfig.DistributionStrategy
    // For each of these arrayed tuples, we have [turnNumber: <"pair events">[]]
    DeerForageHistory.forEach(([turn, forageRound]) => {
        const Participants = Object.entries(
            forageRound.ParticipantContributions
        )
        Participants.forEach(([team, value]) => {
            if (HuntingDistributionStrategy === 'EqualSplit') {
                metrics[team] += forageRound.TotalUtility / Participants.length
            } else {
                metrics[team] +=
                    forageRound.TotalUtility /
                    (forageRound.InputResources / value)
            }
            totalSpent[team] += value
        })
    })
    // Since IIGOHistories is repeated, take the one from the LAST GameState and
    // do Object.entries to make it iterable. List of array'ed tuples.
    const FishForageHistory = Object.entries(
        data.GameStates[data.GameStates.length - 1].ForagingHistory
            .FishForageType
    )
    const FishingDistributionStrategy =
        data.Config.ForagingConfig.FishingConfig.DistributionStrategy
    // For each of these arrayed tuples, we have [turnNumber: <"pair events">[]]
    FishForageHistory.forEach(([turn, forageRound]) => {
        const Participants = Object.entries(
            forageRound.ParticipantContributions
        )
        Participants.forEach(([team, value]) => {
            if (FishingDistributionStrategy === 'EqualSplit') {
                metrics[team] += forageRound.TotalUtility / Participants.length
            } else {
                metrics[team] +=
                    forageRound.TotalUtility /
                    (forageRound.InputResources / value)
            }
            totalSpent[team] += value
        })
    })
    Object.entries(totalSpent).forEach(([team, spent]) => {
        if (spent > 0) {
            metrics[team] /= spent
        }
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

const seasonsAliveMetricCollection = (data: OutputJSONType) => {
    const totalAgents = numAgents(data)
    const metrics = emptyMetrics(totalAgents)
    teamNames(totalAgents).forEach((team) => {
        metrics[team] = seasonsAlive(data, team)
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

const metricList: MetricEntry[] = [
    {
        title: 'Island longevity',
        description: 'Turns island is alive',
        collectMetrics: turnsAliveMetricCollection,
        evalLargest: true,
    },
    {
        title: 'Island disaster longevity',
        description: 'Island survived n disasters',
        collectMetrics: seasonsAliveMetricCollection,
        evalLargest: true,
    },
    {
        title: 'Island disaster impact',
        description: 'Island suffered this much from disasters',
        collectMetrics: totalDisasterImpactCollection,
        evalLargest: true,
    },
    {
        title: 'Island disaster mitigated',
        description: 'Common pool disaster mitigation for each team',
        collectMetrics: totalDisasterImpactMitigatedCollection,
        evalLargest: true,
    },
    {
        title: 'Resources Collected',
        description: 'Island with the highest average resources',
        collectMetrics: averageResourcesMetricCollection,
        evalLargest: true,
    },
    {
        title: 'Max resources',
        description: 'Island with the highest peak resources',
        collectMetrics: peakResourcesMetricCollection,
        evalLargest: true,
    },
    {
        title: 'Fish Foraging Spent',
        description: 'Total fish foraging resources spent',
        collectMetrics: fishForagingResultCollection,
        evalLargest: true,
    },
    {
        title: 'Fish Foraging Efficiency',
        description: 'Total fish foraging resources efficency',
        collectMetrics: fishForagingEfficiencyCollection,
        evalLargest: true,
    },
    {
        title: 'Deer Foraging Spent',
        description: 'Total deer foraging resources spent',
        collectMetrics: deerForagingResultCollection,
        evalLargest: true,
    },
    {
        title: 'Deer Foraging Efficiency',
        description: 'Total deer foraging resources efficency',
        collectMetrics: deerForagingEfficiencyCollection,
        evalLargest: true,
    },
    {
        title: 'Total Foraging Efficieny',
        description: 'Total deer foraging resources spent',
        collectMetrics: totalForagingEfficiencyCollection,
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
    {
        title: 'Criminal',
        description: 'Number of rules broken',
        collectMetrics: rulesBrokenPerTeamCollection,
        evalLargest: true,
    },
]

export default metricList

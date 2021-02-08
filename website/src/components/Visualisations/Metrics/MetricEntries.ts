import { OutputJSONType } from '../../../consts/types'

export type TeamName = 'Team1' | 'Team2' | 'Team3' | 'Team4' | 'Team5' | 'Team6'

export type MetricsType = {
    [key in TeamName]: number
}
export type Metric = {
    teamName: string
    value: number
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

export type MetricEntry = {
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

const seasonsAlive = (data: OutputJSONType, team: TeamName): number =>
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

const rulesBrokenPerTeamCollection = (data: OutputJSONType) => {
    const metrics = emptyMetrics()
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
    const metrics = emptyMetrics()
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
    const metrics = emptyMetrics()
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
    const metrics = emptyMetrics()
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
    const metrics = emptyMetrics()
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
    const metrics = emptyMetrics()
    const totalSpent = emptyMetrics()

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
    const metrics = emptyMetrics()
    const totalSpent = emptyMetrics()

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
    const metrics = emptyMetrics()
    const totalSpent = emptyMetrics()

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
    const metrics = emptyMetrics()
    teamNames().forEach((team) => {
        metrics[team] = turnsAlive(data, team)
    })
    return metrics
}

const seasonsAliveMetricCollection = (data: OutputJSONType) => {
    const metrics = emptyMetrics()
    teamNames().forEach((team) => {
        metrics[team] = seasonsAlive(data, team)
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

import { getIITOTransactions, getIIGOTransactions } from "./ProcessTransactionData"

const testInput = {
    "Config": {
        "MaxSeasons": 10,
        "MaxTurns": 10,
        "InitialResources": 100,
        "InitialCommonPool": 100,
        "CostOfLiving": 10,
        "MinimumResourceThreshold": 5,
        "MaxCriticalConsecutiveTurns": 3,
        "ForagingConfig": {
            "DeerHuntConfig": {
                "MaxDeerPerHunt": 4,
                "IncrementalInputDecay": 0.8,
                "BernoulliProb": 0.95,
                "ExponentialRate": 1,
                "InputScaler": 12,
                "OutputScaler": 18,
                "DistributionStrategy": "InputProportionalSplit",
                "ThetaCritical": 0.8,
                "ThetaMax": 0.95,
                "MaxDeerPopulation": 12,
                "DeerGrowthCoefficient": 0.2
            },
            "FishingConfig": {
                "MaxFishPerHunt": 6,
                "IncrementalInputDecay": 0.8,
                "Mean": 0.9,
                "Variance": 0.2,
                "InputScaler": 10,
                "OutputScaler": 12,
                "DistributionStrategy": "EqualSplit"
            }
        },
        "DisasterConfig": {
            "XMin": 0,
            "XMax": 10,
            "YMin": 0,
            "YMax": 10,
            "Period": 15,
            "SpatialPDFType": "Uniform",
            "MagnitudeLambda": 1,
            "MagnitudeResourceMultiplier": 500,
            "CommonpoolThreshold": 50,
            "StochasticPeriod": false,
            "CommonpoolThresholdVisible": false
        },
        "IIGOConfig": {
            "GetRuleForSpeakerActionCost": 10,
            "BroadcastTaxationActionCost": 10,
            "ReplyAllocationRequestsActionCost": 10,
            "RequestAllocationRequestActionCost": 10,
            "RequestRuleProposalActionCost": 10,
            "AppointNextSpeakerActionCost": 10,
            "InspectHistoryActionCost": 10,
            "InspectBallotActionCost": 10,
            "InspectAllocationActionCost": 10,
            "AppointNextPresidentActionCost": 10,
            "SanctionCacheDepth": 3,
            "HistoryCacheDepth": 3,
            "AssumedResourcesNoReport": 500,
            "SanctionLength": 2,
            "SetVotingResultActionCost": 10,
            "SetRuleToVoteActionCost": 10,
            "AnnounceVotingResultActionCost": 10,
            "UpdateRulesActionCost": 10,
            "AppointNextJudgeActionCost": 10,
            "StartWithRulesInPlay": true
        }
    },
    "GitInfo": {
        "Hash": "f7f1382ae3dbdd9c5b8526ad2259d7691c86767d",
        "GithubURL": "https://github.com/SOMAS2020/SOMAS2020.git/tree/f7f1382ae3dbdd9c5b8526ad2259d7691c86767d"
    },
    "RunInfo": {
        "TimeStart": "2021-01-06T09:08:43.831209Z",
        "TimeEnd": "2021-01-06T09:08:43.845011Z",
        "DurationSeconds": 0.01380157,
        "Version": "go1.15.6",
        "GOOS": "darwin",
        "GOARCH": "amd64"
    },
    "AuxInfo": {
        "TeamIDs": [
            "Team1",
            "Team2",
            "Team3",
            "Team4",
            "Team5",
            "Team6"
        ]
    },
    "GameStates": [
        {
            "Season": 1,
            "Turn": 2,
            "CommonPool": 100,
            "ClientInfos": {
                "Team1": {
                    "Resources": 100,
                    "LifeStatus": "Alive",
                    "CriticalConsecutiveTurnsCounter": 0
                },
                "Team2": {
                    "Resources": 100,
                    "LifeStatus": "Alive",
                    "CriticalConsecutiveTurnsCounter": 0
                },
                "Team3": {
                    "Resources": 100,
                    "LifeStatus": "Alive",
                    "CriticalConsecutiveTurnsCounter": 0
                },
                "Team4": {
                    "Resources": 100,
                    "LifeStatus": "Alive",
                    "CriticalConsecutiveTurnsCounter": 0
                },
                "Team5": {
                    "Resources": 100,
                    "LifeStatus": "Alive",
                    "CriticalConsecutiveTurnsCounter": 0
                },
                "Team6": {
                    "Resources": 100,
                    "LifeStatus": "Alive",
                    "CriticalConsecutiveTurnsCounter": 0
                }
            },
            "Environment": {
                "Geography": {
                    "Islands": {
                        "Team1": {
                            "ID": "Team1",
                            "X": 8,
                            "Y": 0
                        },
                        "Team2": {
                            "ID": "Team2",
                            "X": 10,
                            "Y": 0
                        },
                        "Team3": {
                            "ID": "Team3",
                            "X": 0,
                            "Y": 0
                        },
                        "Team4": {
                            "ID": "Team4",
                            "X": 2,
                            "Y": 0
                        },
                        "Team5": {
                            "ID": "Team5",
                            "X": 4,
                            "Y": 0
                        },
                        "Team6": {
                            "ID": "Team6",
                            "X": 6,
                            "Y": 0
                        }
                    },
                    "XMin": 0,
                    "XMax": 10,
                    "YMin": 0,
                    "YMax": 10
                },
                "LastDisasterReport": {
                    "Magnitude": 0,
                    "X": 0,
                    "Y": 0
                }
            },
            "DeerPopulation": {
                "Population": 12,
                "T": 0
            },
            "ForagingHistory": {
                "DeerForageType": [],
                "FishForageType": []
            },
            "IIGOHistory": {
                "1": [
                    {
                        "ClientID": "Team5",
                        "Pairs": [
                            {
                                "VariableName": "IslandAllocation",
                                "Values": [
                                    10
                                ]
                            },
                            {
                                "VariableName": "ExpectedAllocation",
                                "Values": [
                                    20
                                ]
                            }
                        ]
                    },
                    {
                        "ClientID": "Team6",
                        "Pairs": [
                            {
                                "VariableName": "IslandAllocation",
                                "Values": [
                                    10
                                ]
                            },
                            {
                                "VariableName": "ExpectedAllocation",
                                "Values": [
                                    20
                                ]
                            }
                        ]
                    },
                    {
                        "ClientID": "Team6",
                        "Pairs": [
                            {
                                "VariableName": "IslandTaxContribution",
                                "Values": [
                                    42
                                ]
                            },
                            {
                                "VariableName": "ExpectedTaxContribution",
                                "Values": [
                                    0
                                ]
                            }
                        ]
                    },
                    {
                        "ClientID": "Team6",
                        "Pairs": [
                            {
                                "VariableName": "SanctionPaid",
                                "Values": [
                                    0
                                ]
                            },
                            {
                                "VariableName": "SanctionExpected",
                                "Values": [
                                    0
                                ]
                            }
                        ]
                    },
                    {
                        "ClientID": "Team1",
                        "Pairs": [
                            {
                                "VariableName": "IslandTaxContribution",
                                "Values": [
                                    68
                                ]
                            },
                            {
                                "VariableName": "ExpectedTaxContribution",
                                "Values": [
                                    0
                                ]
                            }
                        ]
                    },
                    {
                        "ClientID": "Team1",
                        "Pairs": [
                            {
                                "VariableName": "SanctionPaid",
                                "Values": [
                                    0
                                ]
                            },
                            {
                                "VariableName": "SanctionExpected",
                                "Values": [
                                    0
                                ]
                            }
                        ]
                    },
                    {
                        "ClientID": "Team2",
                        "Pairs": [
                            {
                                "VariableName": "IslandTaxContribution",
                                "Values": [
                                    60
                                ]
                            },
                            {
                                "VariableName": "ExpectedTaxContribution",
                                "Values": [
                                    0
                                ]
                            }
                        ]
                    },
                    {
                        "ClientID": "Team2",
                        "Pairs": [
                            {
                                "VariableName": "SanctionPaid",
                                "Values": [
                                    0
                                ]
                            },
                            {
                                "VariableName": "SanctionExpected",
                                "Values": [
                                    0
                                ]
                            }
                        ]
                    },
                    {
                        "ClientID": "Team3",
                        "Pairs": [
                            {
                                "VariableName": "IslandTaxContribution",
                                "Values": [
                                    94
                                ]
                            },
                            {
                                "VariableName": "ExpectedTaxContribution",
                                "Values": [
                                    0
                                ]
                            }
                        ]
                    },
                    {
                        "ClientID": "Team3",
                        "Pairs": [
                            {
                                "VariableName": "SanctionPaid",
                                "Values": [
                                    0
                                ]
                            },
                            {
                                "VariableName": "SanctionExpected",
                                "Values": [
                                    0
                                ]
                            }
                        ]
                    },
                    {
                        "ClientID": "Team4",
                        "Pairs": [
                            {
                                "VariableName": "IslandTaxContribution",
                                "Values": [
                                    66
                                ]
                            },
                            {
                                "VariableName": "ExpectedTaxContribution",
                                "Values": [
                                    0
                                ]
                            }
                        ]
                    },
                    {
                        "ClientID": "Team4",
                        "Pairs": [
                            {
                                "VariableName": "SanctionPaid",
                                "Values": [
                                    0
                                ]
                            },
                            {
                                "VariableName": "SanctionExpected",
                                "Values": [
                                    0
                                ]
                            }
                        ]
                    },
                    {
                        "ClientID": "Team5",
                        "Pairs": [
                            {
                                "VariableName": "IslandTaxContribution",
                                "Values": [
                                    43
                                ]
                            },
                            {
                                "VariableName": "ExpectedTaxContribution",
                                "Values": [
                                    0
                                ]
                            }
                        ]
                    },
                    {
                        "ClientID": "Team5",
                        "Pairs": [
                            {
                                "VariableName": "SanctionPaid",
                                "Values": [
                                    0
                                ]
                            },
                            {
                                "VariableName": "SanctionExpected",
                                "Values": [
                                    0
                                ]
                            }
                        ]
                    }
                ],
            },
            "IIGORolesBudget": {
                "Judge": 0,
                "President": 0,
                "Speaker": 0
            },
            "IITOTransactions": null,
            "SpeakerID": "Team1",
            "JudgeID": "Team2",
            "PresidentID": "Team3"
        },
    ]
}

test('test getIIGOTransactions', () => {
    const want = [
        { from: 0, to: 5, amount: 10 },
        { from: 0, to: 6, amount: 10 },
        { from: 6, to: 0, amount: 42 },
        { from: 1, to: 0, amount: 68 },
        { from: 2, to: 0, amount: 60 },
        { from: 3, to: 0, amount: 94 },
        { from: 4, to: 0, amount: 66 },
        { from: 5, to: 0, amount: 43 },
    ];
    expect(getIIGOTransactions(testInput)).toEqual(want);
});

test('test getIITOTransactions', () => {
    // expect(sum(1, 2)).toBe(3);
});
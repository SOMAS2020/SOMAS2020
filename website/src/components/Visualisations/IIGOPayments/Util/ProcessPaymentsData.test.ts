import { processPaymentsData } from './ProcessedPaymentsData'

const testInput = {
    Config: {
        MaxSeasons: 100,
        MaxTurns: 100,
        InitialResources: 100,
        InitialCommonPool: 100,
        CostOfLiving: 10,
        MinimumResourceThreshold: 5,
        MaxCriticalConsecutiveTurns: 3,
        ForagingConfig: {
            DeerHuntConfig: {
                MaxDeerPerHunt: 4,
                IncrementalInputDecay: 0.8,
                BernoulliProb: 0.95,
                ExponentialRate: 1,
                InputScaler: 12,
                OutputScaler: 18,
                DistributionStrategy: 'InputProportionalSplit',
                ThetaCritical: 0.8,
                ThetaMax: 0.95,
                MaxDeerPopulation: 12,
                DeerGrowthCoefficient: 0.2,
            },
            FishingConfig: {
                MaxFishPerHunt: 6,
                IncrementalInputDecay: 0.8,
                Mean: 0.9,
                Variance: 0.2,
                InputScaler: 10,
                OutputScaler: 12,
                DistributionStrategy: 'EqualSplit',
            },
        },
        DisasterConfig: {
            XMin: 0,
            XMax: 10,
            YMin: 0,
            YMax: 10,
            Period: 15,
            SpatialPDFType: 'Uniform',
            MagnitudeLambda: 1,
            MagnitudeResourceMultiplier: 500,
            CommonpoolThreshold: 50,
            StochasticPeriod: false,
            CommonpoolThresholdVisible: false,
            PeriodVisible: true,
            StochasticPeriodVisible: true,
        },
        IIGOConfig: {
            IIGOTermLengths: {
                Judge: 4,
                President: 4,
                Speaker: 4,
            },
            GetRuleForSpeakerActionCost: 10,
            BroadcastTaxationActionCost: 10,
            ReplyAllocationRequestsActionCost: 10,
            RequestAllocationRequestActionCost: 10,
            RequestRuleProposalActionCost: 10,
            AppointNextSpeakerActionCost: 10,
            InspectHistoryActionCost: 10,
            HistoricalRetributionActionCost: 10,
            InspectBallotActionCost: 10,
            InspectAllocationActionCost: 10,
            AppointNextPresidentActionCost: 10,
            SanctionCacheDepth: 3,
            HistoryCacheDepth: 3,
            AssumedResourcesNoReport: 500,
            SanctionLength: 2,
            SetVotingResultActionCost: 10,
            SetRuleToVoteActionCost: 10,
            AnnounceVotingResultActionCost: 10,
            UpdateRulesActionCost: 10,
            AppointNextJudgeActionCost: 10,
            StartWithRulesInPlay: true,
        },
    },
    GitInfo: {
        Hash: 'f7f1382ae3dbdd9c5b8526ad2259d7691c86767d',
        GithubURL:
            'https://github.com/SOMAS2020/SOMAS2020.git/tree/f7f1382ae3dbdd9c5b8526ad2259d7691c86767d',
    },
    RunInfo: {
        TimeStart: '2021-01-06T09:08:43.831209Z',
        TimeEnd: '2021-01-06T09:08:43.845011Z',
        DurationSeconds: 0.01380157,
        Version: 'go1.15.6',
        GOOS: 'darwin',
        GOARCH: 'amd64',
    },
    AuxInfo: {
        TeamIDs: ['Team1', 'Team2', 'Team3', 'Team4', 'Team5', 'Team6'],
    },
    GameStates: [
        {
            Season: 1,
            Turn: 2,
            CommonPool: 55.60116950066092,
            ClientInfos: {
                Team1: {
                    Resources: 142.25088034365214,
                    LifeStatus: 'Alive',
                    CriticalConsecutiveTurnsCounter: 0,
                },
                Team2: {
                    Resources: 60.52656738758088,
                    LifeStatus: 'Alive',
                    CriticalConsecutiveTurnsCounter: 0,
                },
                Team3: {
                    Resources: 78.50412501698553,
                    LifeStatus: 'Alive',
                    CriticalConsecutiveTurnsCounter: 0,
                },
                Team4: {
                    Resources: 114.25376816777894,
                    LifeStatus: 'Alive',
                    CriticalConsecutiveTurnsCounter: 0,
                },
                Team5: {
                    Resources: 113.95670179599756,
                    LifeStatus: 'Alive',
                    CriticalConsecutiveTurnsCounter: 0,
                },
                Team6: {
                    Resources: 121.50316879276909,
                    LifeStatus: 'Alive',
                    CriticalConsecutiveTurnsCounter: 0,
                },
            },
            Environment: {
                Geography: {
                    Islands: {
                        Team1: {
                            ID: 'Team1',
                            X: 4,
                            Y: 0,
                        },
                        Team2: {
                            ID: 'Team2',
                            X: 6,
                            Y: 0,
                        },
                        Team3: {
                            ID: 'Team3',
                            X: 8,
                            Y: 0,
                        },
                        Team4: {
                            ID: 'Team4',
                            X: 10,
                            Y: 0,
                        },
                        Team5: {
                            ID: 'Team5',
                            X: 0,
                            Y: 0,
                        },
                        Team6: {
                            ID: 'Team6',
                            X: 2,
                            Y: 0,
                        },
                    },
                    XMin: 0,
                    XMax: 10,
                    YMin: 0,
                    YMax: 10,
                },
                LastDisasterReport: {
                    Magnitude: 0,
                    X: -1,
                    Y: -1,
                },
            },
            DeerPopulation: {
                Population: 8.72507698680008,
                T: 1,
            },
            ForagingHistory: {
                DeerForageType: [
                    {
                        ForageType: 'DeerForageType',
                        InputResources: 67.51314378599372,
                        NumberParticipants: 4,
                        NumberCaught: 4,
                        TotalUtility: 181.52860448572852,
                        CatchSizes: [
                            64.01619328509351,
                            18.19793571378622,
                            80.239234314334,
                            19.075241172514765,
                        ],
                        Turn: 1,
                    },
                ],
                FishForageType: [
                    {
                        ForageType: 'FishForageType',
                        InputResources: 36.27279921980194,
                        NumberParticipants: 2,
                        NumberCaught: 5,
                        TotalUtility: 58.85371952549224,
                        CatchSizes: [
                            9.219552904973892,
                            10.491508043284016,
                            12.250361935103452,
                            11.716990010074612,
                            15.175306632056264,
                        ],
                        Turn: 1,
                    },
                ],
            },
            RulesInfo: {
                VariableMap: {
                    AllocationMade: {
                        VariableName: 'AllocationMade',
                        Values: [0],
                    },
                    AllocationRequestsMade: {
                        VariableName: 'AllocationRequestsMade',
                        Values: [1],
                    },
                    AnnouncementResultMatchesVote: {
                        VariableName: 'AnnouncementResultMatchesVote',
                        Values: [0],
                    },
                    AnnouncementRuleMatchesVote: {
                        VariableName: 'AnnouncementRuleMatchesVote',
                        Values: [0],
                    },
                    AppointmentMatchesVote: {
                        VariableName: 'AppointmentMatchesVote',
                        Values: [0],
                    },
                    ConstSanctionAmount: {
                        VariableName: 'ConstSanctionAmount',
                        Values: [0],
                    },
                    ElectionHeld: {
                        VariableName: 'ElectionHeld',
                        Values: [0],
                    },
                    ExpectedAllocation: {
                        VariableName: 'ExpectedAllocation',
                        Values: [0],
                    },
                    ExpectedTaxContribution: {
                        VariableName: 'ExpectedTaxContribution',
                        Values: [0],
                    },
                    HasIslandReportPrivateResources: {
                        VariableName: 'HasIslandReportPrivateResources',
                        Values: [0],
                    },
                    IslandActualPrivateResources: {
                        VariableName: 'IslandActualPrivateResources',
                        Values: [0],
                    },
                    IslandAllocation: {
                        VariableName: 'IslandAllocation',
                        Values: [0],
                    },
                    IslandReportedPrivateResources: {
                        VariableName: 'IslandReportedPrivateResources',
                        Values: [0],
                    },
                    IslandReportedResources: {
                        VariableName: 'IslandReportedResources',
                        Values: [0],
                    },
                    IslandTaxContribution: {
                        VariableName: 'IslandTaxContribution',
                        Values: [0],
                    },
                    IslandsAlive: {
                        VariableName: 'IslandsAlive',
                        Values: [1, 2, 3, 4, 5, 0],
                    },
                    IslandsAllowedToVote: {
                        VariableName: 'IslandsAllowedToVote',
                        Values: [0],
                    },
                    IslandsProposedRules: {
                        VariableName: 'IslandsProposedRules',
                        Values: [0],
                    },
                    JudgeBudgetIncrement: {
                        VariableName: 'JudgeBudgetIncrement',
                        Values: [100],
                    },
                    JudgeHistoricalRetributionPerformed: {
                        VariableName: 'JudgeHistoricalRetributionPerformed',
                        Values: [0],
                    },
                    JudgeInspectionPerformed: {
                        VariableName: 'JudgeInspectionPerformed',
                        Values: [0],
                    },
                    JudgeLeftoverBudget: {
                        VariableName: 'JudgeLeftoverBudget',
                        Values: [0],
                    },
                    JudgePaid: {
                        VariableName: 'JudgePaid',
                        Values: [0],
                    },
                    JudgePayment: {
                        VariableName: 'JudgePayment',
                        Values: [50],
                    },
                    JudgeSalary: {
                        VariableName: 'JudgeSalary',
                        Values: [50],
                    },
                    MaxSeverityOfSanctions: {
                        VariableName: 'MaxSeverityOfSanctions',
                        Values: [2],
                    },
                    MonitorRoleAnnounce: {
                        VariableName: 'MonitorRoleAnnounce',
                        Values: [0],
                    },
                    MonitorRoleDecideToMonitor: {
                        VariableName: 'MonitorRoleDecideToMonitor',
                        Values: [0],
                    },
                    MonitorRoleEvalResult: {
                        VariableName: 'MonitorRoleEvalResult',
                        Values: [0],
                    },
                    MonitorRoleEvalResultDecide: {
                        VariableName: 'MonitorRoleEvalResultDecide',
                        Values: [0],
                    },
                    NumberOfAllocationsSent: {
                        VariableName: 'NumberOfAllocationsSent',
                        Values: [6],
                    },
                    NumberOfBallotsCast: {
                        VariableName: 'NumberOfBallotsCast',
                        Values: [6],
                    },
                    NumberOfBrokenAgreements: {
                        VariableName: 'NumberOfBrokenAgreements',
                        Values: [1],
                    },
                    NumberOfFailedForages: {
                        VariableName: 'NumberOfFailedForages',
                        Values: [0.5],
                    },
                    NumberOfIslandsAlive: {
                        VariableName: 'NumberOfIslandsAlive',
                        Values: [6],
                    },
                    NumberOfIslandsContributingToCommonPool: {
                        VariableName: 'NumberOfIslandsContributingToCommonPool',
                        Values: [5],
                    },
                    PresidentBudgetIncrement: {
                        VariableName: 'PresidentBudgetIncrement',
                        Values: [100],
                    },
                    PresidentLeftoverBudget: {
                        VariableName: 'PresidentLeftoverBudget',
                        Values: [0],
                    },
                    PresidentPaid: {
                        VariableName: 'PresidentPaid',
                        Values: [0],
                    },
                    PresidentPayment: {
                        VariableName: 'PresidentPayment',
                        Values: [50],
                    },
                    PresidentRuleProposal: {
                        VariableName: 'PresidentRuleProposal',
                        Values: [0],
                    },
                    PresidentSalary: {
                        VariableName: 'PresidentSalary',
                        Values: [50],
                    },
                    RuleChosenFromProposalList: {
                        VariableName: 'RuleChosenFromProposalList',
                        Values: [0],
                    },
                    RuleSelected: {
                        VariableName: 'RuleSelected',
                        Values: [0],
                    },
                    SanctionExpected: {
                        VariableName: 'SanctionExpected',
                        Values: [0],
                    },
                    SanctionPaid: {
                        VariableName: 'SanctionPaid',
                        Values: [0],
                    },
                    SpeakerBudgetIncrement: {
                        VariableName: 'SpeakerBudgetIncrement',
                        Values: [100],
                    },
                    SpeakerLeftoverBudget: {
                        VariableName: 'SpeakerLeftoverBudget',
                        Values: [0],
                    },
                    SpeakerPaid: {
                        VariableName: 'SpeakerPaid',
                        Values: [0],
                    },
                    SpeakerPayment: {
                        VariableName: 'SpeakerPayment',
                        Values: [50],
                    },
                    SpeakerProposedPresidentRule: {
                        VariableName: 'SpeakerProposedPresidentRule',
                        Values: [0],
                    },
                    SpeakerSalary: {
                        VariableName: 'SpeakerSalary',
                        Values: [50],
                    },
                    TaxDecisionMade: {
                        VariableName: 'TaxDecisionMade',
                        Values: [1],
                    },
                    TermEnded: {
                        VariableName: 'TermEnded',
                        Values: [0],
                    },
                    TurnsLeftOnSanction: {
                        VariableName: 'TurnsLeftOnSanction',
                        Values: [0],
                    },
                    VoteCalled: {
                        VariableName: 'VoteCalled',
                        Values: [0],
                    },
                    VoteResultAnnounced: {
                        VariableName: 'VoteResultAnnounced',
                        Values: [0],
                    },
                },
                AvailableRules: {
                    'Kinda Complicated Rule': {
                        RuleName: 'Kinda Complicated Rule',
                        RequiredVariables: [
                            'NumberOfIslandsContributingToCommonPool',
                            'NumberOfFailedForages',
                            'NumberOfBrokenAgreements',
                            'MaxSeverityOfSanctions',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    allocation_decision: {
                        RuleName: 'allocation_decision',
                        RequiredVariables: ['AllocationMade'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: true,
                            LinkType: 0,
                            LinkedRule: 'check_allocation_rule',
                        },
                    },
                    allocations_made_rule: {
                        RuleName: 'allocations_made_rule',
                        RequiredVariables: [
                            'AllocationRequestsMade',
                            'AllocationMade',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    announcement_matches_vote: {
                        RuleName: 'announcement_matches_vote',
                        RequiredVariables: [
                            'AnnouncementRuleMatchesVote',
                            'AnnouncementResultMatchesVote',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    check_allocation_rule: {
                        RuleName: 'check_allocation_rule',
                        RequiredVariables: [
                            'IslandAllocation',
                            'ExpectedAllocation',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    check_sanction_rule: {
                        RuleName: 'check_sanction_rule',
                        RequiredVariables: ['SanctionPaid', 'SanctionExpected'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    check_taxation_rule: {
                        RuleName: 'check_taxation_rule',
                        RequiredVariables: [
                            'IslandTaxContribution',
                            'ExpectedTaxContribution',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    iigo_economic_sanction_1: {
                        RuleName: 'iigo_economic_sanction_1',
                        RequiredVariables: [
                            'IslandReportedResources',
                            'ConstSanctionAmount',
                            'TurnsLeftOnSanction',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: true,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    iigo_economic_sanction_2: {
                        RuleName: 'iigo_economic_sanction_2',
                        RequiredVariables: [
                            'IslandReportedResources',
                            'ConstSanctionAmount',
                            'TurnsLeftOnSanction',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: true,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    iigo_economic_sanction_3: {
                        RuleName: 'iigo_economic_sanction_3',
                        RequiredVariables: [
                            'IslandReportedResources',
                            'ConstSanctionAmount',
                            'TurnsLeftOnSanction',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: true,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    iigo_economic_sanction_4: {
                        RuleName: 'iigo_economic_sanction_4',
                        RequiredVariables: [
                            'IslandReportedResources',
                            'ConstSanctionAmount',
                            'TurnsLeftOnSanction',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: true,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    iigo_economic_sanction_5: {
                        RuleName: 'iigo_economic_sanction_5',
                        RequiredVariables: [
                            'IslandReportedResources',
                            'ConstSanctionAmount',
                            'TurnsLeftOnSanction',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: true,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    iigo_monitor_rule_permission_1: {
                        RuleName: 'iigo_monitor_rule_permission_1',
                        RequiredVariables: [
                            'MonitorRoleDecideToMonitor',
                            'MonitorRoleAnnounce',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: true,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    iigo_monitor_rule_permission_2: {
                        RuleName: 'iigo_monitor_rule_permission_2',
                        RequiredVariables: [
                            'MonitorRoleEvalResult',
                            'MonitorRoleEvalResultDecide',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: true,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    increment_budget_judge: {
                        RuleName: 'increment_budget_judge',
                        RequiredVariables: ['JudgeBudgetIncrement'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: true,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    increment_budget_president: {
                        RuleName: 'increment_budget_president',
                        RequiredVariables: ['PresidentBudgetIncrement'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: true,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    increment_budget_speaker: {
                        RuleName: 'increment_budget_speaker',
                        RequiredVariables: ['SpeakerBudgetIncrement'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: true,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    inspect_ballot_rule: {
                        RuleName: 'inspect_ballot_rule',
                        RequiredVariables: [
                            'NumberOfIslandsAlive',
                            'NumberOfBallotsCast',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    island_must_report_actual_private_resource: {
                        RuleName: 'island_must_report_actual_private_resource',
                        RequiredVariables: [
                            'IslandActualPrivateResources',
                            'IslandReportedPrivateResources',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    island_must_report_private_resource: {
                        RuleName: 'island_must_report_private_resource',
                        RequiredVariables: ['HasIslandReportPrivateResources'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    islands_allowed_to_vote_rule: {
                        RuleName: 'islands_allowed_to_vote_rule',
                        RequiredVariables: [
                            'NumberOfIslandsAlive',
                            'IslandsAllowedToVote',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    judge_historical_retribution_permission: {
                        RuleName: 'judge_historical_retribution_permission',
                        RequiredVariables: [
                            'JudgeHistoricalRetributionPerformed',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    judge_inspection_rule: {
                        RuleName: 'judge_inspection_rule',
                        RequiredVariables: ['JudgeInspectionPerformed'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    judge_over_budget: {
                        RuleName: 'judge_over_budget',
                        RequiredVariables: ['JudgeLeftoverBudget'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    must_appoint_elected_island: {
                        RuleName: 'must_appoint_elected_island',
                        RequiredVariables: ['AppointmentMatchesVote'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    obl_to_propose_rule_if_some_are_given: {
                        RuleName: 'obl_to_propose_rule_if_some_are_given',
                        RequiredVariables: [
                            'IslandsProposedRules',
                            'RuleSelected',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    president_over_budget: {
                        RuleName: 'president_over_budget',
                        RequiredVariables: ['PresidentLeftoverBudget'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    roles_must_hold_election: {
                        RuleName: 'roles_must_hold_election',
                        RequiredVariables: ['TermEnded', 'ElectionHeld'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    rule_chosen_from_proposal_list: {
                        RuleName: 'rule_chosen_from_proposal_list',
                        RequiredVariables: ['RuleChosenFromProposalList'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    rule_to_vote_on_rule: {
                        RuleName: 'rule_to_vote_on_rule',
                        RequiredVariables: ['SpeakerProposedPresidentRule'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    salary_cycle_judge: {
                        RuleName: 'salary_cycle_judge',
                        RequiredVariables: ['JudgePayment'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: true,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    salary_cycle_president: {
                        RuleName: 'salary_cycle_president',
                        RequiredVariables: ['PresidentPayment'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: true,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    salary_cycle_speaker: {
                        RuleName: 'salary_cycle_speaker',
                        RequiredVariables: ['SpeakerPayment'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: true,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    salary_paid_judge: {
                        RuleName: 'salary_paid_judge',
                        RequiredVariables: ['JudgePaid'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    salary_paid_president: {
                        RuleName: 'salary_paid_president',
                        RequiredVariables: ['PresidentPaid'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    salary_paid_speaker: {
                        RuleName: 'salary_paid_speaker',
                        RequiredVariables: ['SpeakerPaid'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    speaker_over_budget: {
                        RuleName: 'speaker_over_budget',
                        RequiredVariables: ['SpeakerLeftoverBudget'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    tax_decision: {
                        RuleName: 'tax_decision',
                        RequiredVariables: ['TaxDecisionMade'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: true,
                            LinkType: 0,
                            LinkedRule: 'check_taxation_rule',
                        },
                    },
                    vote_called_rule: {
                        RuleName: 'vote_called_rule',
                        RequiredVariables: ['RuleSelected', 'VoteCalled'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    vote_result_rule: {
                        RuleName: 'vote_result_rule',
                        RequiredVariables: [
                            'VoteResultAnnounced',
                            'VoteCalled',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                },
                CurrentRulesInPlay: {
                    'Kinda Complicated Rule': {
                        RuleName: 'Kinda Complicated Rule',
                        RequiredVariables: [
                            'NumberOfIslandsContributingToCommonPool',
                            'NumberOfFailedForages',
                            'NumberOfBrokenAgreements',
                            'MaxSeverityOfSanctions',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    allocation_decision: {
                        RuleName: 'allocation_decision',
                        RequiredVariables: ['AllocationMade'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: true,
                            LinkType: 0,
                            LinkedRule: 'check_allocation_rule',
                        },
                    },
                    allocations_made_rule: {
                        RuleName: 'allocations_made_rule',
                        RequiredVariables: [
                            'AllocationRequestsMade',
                            'AllocationMade',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    announcement_matches_vote: {
                        RuleName: 'announcement_matches_vote',
                        RequiredVariables: [
                            'AnnouncementRuleMatchesVote',
                            'AnnouncementResultMatchesVote',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    check_allocation_rule: {
                        RuleName: 'check_allocation_rule',
                        RequiredVariables: [
                            'IslandAllocation',
                            'ExpectedAllocation',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    check_sanction_rule: {
                        RuleName: 'check_sanction_rule',
                        RequiredVariables: ['SanctionPaid', 'SanctionExpected'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    check_taxation_rule: {
                        RuleName: 'check_taxation_rule',
                        RequiredVariables: [
                            'IslandTaxContribution',
                            'ExpectedTaxContribution',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    iigo_economic_sanction_1: {
                        RuleName: 'iigo_economic_sanction_1',
                        RequiredVariables: [
                            'IslandReportedResources',
                            'ConstSanctionAmount',
                            'TurnsLeftOnSanction',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: true,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    iigo_economic_sanction_2: {
                        RuleName: 'iigo_economic_sanction_2',
                        RequiredVariables: [
                            'IslandReportedResources',
                            'ConstSanctionAmount',
                            'TurnsLeftOnSanction',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: true,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    iigo_economic_sanction_3: {
                        RuleName: 'iigo_economic_sanction_3',
                        RequiredVariables: [
                            'IslandReportedResources',
                            'ConstSanctionAmount',
                            'TurnsLeftOnSanction',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: true,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    iigo_economic_sanction_4: {
                        RuleName: 'iigo_economic_sanction_4',
                        RequiredVariables: [
                            'IslandReportedResources',
                            'ConstSanctionAmount',
                            'TurnsLeftOnSanction',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: true,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    iigo_economic_sanction_5: {
                        RuleName: 'iigo_economic_sanction_5',
                        RequiredVariables: [
                            'IslandReportedResources',
                            'ConstSanctionAmount',
                            'TurnsLeftOnSanction',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: true,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    iigo_monitor_rule_permission_1: {
                        RuleName: 'iigo_monitor_rule_permission_1',
                        RequiredVariables: [
                            'MonitorRoleDecideToMonitor',
                            'MonitorRoleAnnounce',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: true,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    iigo_monitor_rule_permission_2: {
                        RuleName: 'iigo_monitor_rule_permission_2',
                        RequiredVariables: [
                            'MonitorRoleEvalResult',
                            'MonitorRoleEvalResultDecide',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: true,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    increment_budget_judge: {
                        RuleName: 'increment_budget_judge',
                        RequiredVariables: ['JudgeBudgetIncrement'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: true,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    increment_budget_president: {
                        RuleName: 'increment_budget_president',
                        RequiredVariables: ['PresidentBudgetIncrement'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: true,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    increment_budget_speaker: {
                        RuleName: 'increment_budget_speaker',
                        RequiredVariables: ['SpeakerBudgetIncrement'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: true,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    inspect_ballot_rule: {
                        RuleName: 'inspect_ballot_rule',
                        RequiredVariables: [
                            'NumberOfIslandsAlive',
                            'NumberOfBallotsCast',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    island_must_report_actual_private_resource: {
                        RuleName: 'island_must_report_actual_private_resource',
                        RequiredVariables: [
                            'IslandActualPrivateResources',
                            'IslandReportedPrivateResources',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    island_must_report_private_resource: {
                        RuleName: 'island_must_report_private_resource',
                        RequiredVariables: ['HasIslandReportPrivateResources'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    islands_allowed_to_vote_rule: {
                        RuleName: 'islands_allowed_to_vote_rule',
                        RequiredVariables: [
                            'NumberOfIslandsAlive',
                            'IslandsAllowedToVote',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    judge_historical_retribution_permission: {
                        RuleName: 'judge_historical_retribution_permission',
                        RequiredVariables: [
                            'JudgeHistoricalRetributionPerformed',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    judge_inspection_rule: {
                        RuleName: 'judge_inspection_rule',
                        RequiredVariables: ['JudgeInspectionPerformed'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    judge_over_budget: {
                        RuleName: 'judge_over_budget',
                        RequiredVariables: ['JudgeLeftoverBudget'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    must_appoint_elected_island: {
                        RuleName: 'must_appoint_elected_island',
                        RequiredVariables: ['AppointmentMatchesVote'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    obl_to_propose_rule_if_some_are_given: {
                        RuleName: 'obl_to_propose_rule_if_some_are_given',
                        RequiredVariables: [
                            'IslandsProposedRules',
                            'RuleSelected',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    president_over_budget: {
                        RuleName: 'president_over_budget',
                        RequiredVariables: ['PresidentLeftoverBudget'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    roles_must_hold_election: {
                        RuleName: 'roles_must_hold_election',
                        RequiredVariables: ['TermEnded', 'ElectionHeld'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    rule_chosen_from_proposal_list: {
                        RuleName: 'rule_chosen_from_proposal_list',
                        RequiredVariables: ['RuleChosenFromProposalList'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    rule_to_vote_on_rule: {
                        RuleName: 'rule_to_vote_on_rule',
                        RequiredVariables: ['SpeakerProposedPresidentRule'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    salary_cycle_judge: {
                        RuleName: 'salary_cycle_judge',
                        RequiredVariables: ['JudgePayment'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: true,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    salary_cycle_president: {
                        RuleName: 'salary_cycle_president',
                        RequiredVariables: ['PresidentPayment'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: true,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    salary_cycle_speaker: {
                        RuleName: 'salary_cycle_speaker',
                        RequiredVariables: ['SpeakerPayment'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: true,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    salary_paid_judge: {
                        RuleName: 'salary_paid_judge',
                        RequiredVariables: ['JudgePaid'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    salary_paid_president: {
                        RuleName: 'salary_paid_president',
                        RequiredVariables: ['PresidentPaid'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    salary_paid_speaker: {
                        RuleName: 'salary_paid_speaker',
                        RequiredVariables: ['SpeakerPaid'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    speaker_over_budget: {
                        RuleName: 'speaker_over_budget',
                        RequiredVariables: ['SpeakerLeftoverBudget'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    tax_decision: {
                        RuleName: 'tax_decision',
                        RequiredVariables: ['TaxDecisionMade'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: true,
                            LinkType: 0,
                            LinkedRule: 'check_taxation_rule',
                        },
                    },
                    vote_called_rule: {
                        RuleName: 'vote_called_rule',
                        RequiredVariables: ['RuleSelected', 'VoteCalled'],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                    vote_result_rule: {
                        RuleName: 'vote_result_rule',
                        RequiredVariables: [
                            'VoteResultAnnounced',
                            'VoteCalled',
                        ],
                        ApplicableMatrix: {},
                        AuxiliaryVector: {},
                        Mutable: false,
                        Link: {
                            Linked: false,
                            LinkType: 0,
                            LinkedRule: '',
                        },
                    },
                },
            },
            IIGOHistory: {
                '1': [
                    {
                        ClientID: 'Team1',
                        Pairs: [
                            {
                                VariableName: 'IslandAllocation',
                                Values: [1],
                            },
                            {
                                VariableName: 'ExpectedAllocation',
                                Values: [1],
                            },
                        ],
                    },
                    {
                        ClientID: 'Team6',
                        Pairs: [
                            {
                                VariableName: 'IslandAllocation',
                                Values: [4],
                            },
                            {
                                VariableName: 'ExpectedAllocation',
                                Values: [2],
                            },
                        ],
                    },
                    {
                        ClientID: 'Team1',
                        Pairs: [
                            {
                                VariableName: 'IslandAllocation',
                                Values: [4],
                            },
                            {
                                VariableName: 'ExpectedAllocation',
                                Values: [5],
                            },
                        ],
                    },
                    {
                        ClientID: 'Team2',
                        Pairs: [
                            {
                                VariableName: 'IslandAllocation',
                                Values: [3],
                            },
                            {
                                VariableName: 'ExpectedAllocation',
                                Values: [1],
                            },
                        ],
                    },
                    {
                        ClientID: 'Team3',
                        Pairs: [
                            {
                                VariableName: 'IslandAllocation',
                                Values: [43],
                            },
                            {
                                VariableName: 'ExpectedAllocation',
                                Values: [56],
                            },
                        ],
                    },
                    {
                        ClientID: 'Team4',
                        Pairs: [
                            {
                                VariableName: 'IslandAllocation',
                                Values: [0],
                            },
                            {
                                VariableName: 'ExpectedAllocation',
                                Values: [6],
                            },
                        ],
                    },
                    {
                        ClientID: 'Team5',
                        Pairs: [
                            {
                                VariableName: 'IslandTaxContribution',
                                Values: [14],
                            },
                            {
                                VariableName: 'ExpectedTaxContribution',
                                Values: [13],
                            },
                        ],
                    },
                    {
                        ClientID: 'Team5',
                        Pairs: [
                            {
                                VariableName: 'SanctionPaid',
                                Values: [0],
                            },
                            {
                                VariableName: 'SanctionExpected',
                                Values: [0],
                            },
                        ],
                    },
                    {
                        ClientID: 'Team6',
                        Pairs: [
                            {
                                VariableName: 'IslandTaxContribution',
                                Values: [0],
                            },
                            {
                                VariableName: 'ExpectedTaxContribution',
                                Values: [13],
                            },
                        ],
                    },
                    {
                        ClientID: 'Team6',
                        Pairs: [
                            {
                                VariableName: 'SanctionPaid',
                                Values: [0],
                            },
                            {
                                VariableName: 'SanctionExpected',
                                Values: [0],
                            },
                        ],
                    },
                    {
                        ClientID: 'Team1',
                        Pairs: [
                            {
                                VariableName: 'IslandTaxContribution',
                                Values: [14],
                            },
                            {
                                VariableName: 'ExpectedTaxContribution',
                                Values: [14],
                            },
                        ],
                    },
                    {
                        ClientID: 'Team1',
                        Pairs: [
                            {
                                VariableName: 'SanctionPaid',
                                Values: [0],
                            },
                            {
                                VariableName: 'SanctionExpected',
                                Values: [0],
                            },
                        ],
                    },
                    {
                        ClientID: 'Team2',
                        Pairs: [
                            {
                                VariableName: 'IslandTaxContribution',
                                Values: [51],
                            },
                            {
                                VariableName: 'ExpectedTaxContribution',
                                Values: [1],
                            },
                        ],
                    },
                    {
                        ClientID: 'Team2',
                        Pairs: [
                            {
                                VariableName: 'SanctionPaid',
                                Values: [0],
                            },
                            {
                                VariableName: 'SanctionExpected',
                                Values: [0],
                            },
                        ],
                    },
                    {
                        ClientID: 'Team3',
                        Pairs: [
                            {
                                VariableName: 'IslandTaxContribution',
                                Values: [0],
                            },
                            {
                                VariableName: 'ExpectedTaxContribution',
                                Values: [2],
                            },
                        ],
                    },
                    {
                        ClientID: 'Team3',
                        Pairs: [
                            {
                                VariableName: 'SanctionPaid',
                                Values: [0],
                            },
                            {
                                VariableName: 'SanctionExpected',
                                Values: [0],
                            },
                        ],
                    },
                    {
                        ClientID: 'Team4',
                        Pairs: [
                            {
                                VariableName: 'IslandTaxContribution',
                                Values: [14],
                            },
                            {
                                VariableName: 'ExpectedTaxContribution',
                                Values: [14],
                            },
                        ],
                    },
                    {
                        ClientID: 'Team4',
                        Pairs: [
                            {
                                VariableName: 'SanctionPaid',
                                Values: [0],
                            },
                            {
                                VariableName: 'SanctionExpected',
                                Values: [0],
                            },
                        ],
                    },
                ],
            },
            IIGORolesBudget: {
                Judge: 90,
                President: 50,
                Speaker: 70,
            },
            IIGOTurnsInPower: {
                Judge: 1,
                President: 1,
                Speaker: 1,
            },
            IIGOTaxAmount: {
                Team1: 1.4002923751652314,
                Team2: 1.4002923751652314,
                Team3: 2.0989766869216955,
                Team4: 1.4002923751652314,
                Team5: 1.4002923751652314,
                Team6: 1.4002923751652314,
            },
            IIGOAllocationMap: {
                Team1: -0,
                Team2: -0,
                Team3: 56.00000001259784,
                Team4: -6.2989156481913814e-9,
                Team5: -6.2989156481913814e-9,
                Team6: -0,
            },
            IIGOSanctionMap: {},
            IIGOSanctionCache: {
                '0': [
                    {
                        ClientID: 'Team5',
                        SanctionTier: 'NoSanction',
                        TurnsLeft: 1,
                    },
                    {
                        ClientID: 'Team6',
                        SanctionTier: 'NoSanction',
                        TurnsLeft: 1,
                    },
                    {
                        ClientID: 'Team1',
                        SanctionTier: 'NoSanction',
                        TurnsLeft: 1,
                    },
                    {
                        ClientID: 'Team2',
                        SanctionTier: 'NoSanction',
                        TurnsLeft: 1,
                    },
                    {
                        ClientID: 'Team3',
                        SanctionTier: 'NoSanction',
                        TurnsLeft: 1,
                    },
                    {
                        ClientID: 'Team4',
                        SanctionTier: 'NoSanction',
                        TurnsLeft: 1,
                    },
                ],
                '1': [],
                '2': [],
            },
            IIGOHistoryCache: {
                '0': [],
                '1': [],
                '2': [],
            },
            IIGORoleMonitoringCache: [
                {
                    ClientID: 'Team2',
                    Pairs: [
                        {
                            VariableName: 'JudgeInspectionPerformed',
                            Values: [1],
                        },
                    ],
                },
                {
                    ClientID: 'Team2',
                    Pairs: [
                        {
                            VariableName: 'JudgeLeftoverBudget',
                            Values: [90],
                        },
                    ],
                },
                {
                    ClientID: 'Team2',
                    Pairs: [
                        {
                            VariableName: 'JudgeHistoricalRetributionPerformed',
                            Values: [0],
                        },
                    ],
                },
                {
                    ClientID: 'Team3',
                    Pairs: [
                        {
                            VariableName: 'PresidentLeftoverBudget',
                            Values: [90],
                        },
                    ],
                },
                {
                    ClientID: 'Team3',
                    Pairs: [
                        {
                            VariableName: 'PresidentLeftoverBudget',
                            Values: [80],
                        },
                    ],
                },
                {
                    ClientID: 'Team3',
                    Pairs: [
                        {
                            VariableName: 'PresidentLeftoverBudget',
                            Values: [70],
                        },
                    ],
                },
                {
                    ClientID: 'Team3',
                    Pairs: [
                        {
                            VariableName: 'PresidentLeftoverBudget',
                            Values: [60],
                        },
                    ],
                },
                {
                    ClientID: 'Team3',
                    Pairs: [
                        {
                            VariableName: 'IslandsProposedRules',
                            Values: [1],
                        },
                        {
                            VariableName: 'PresidentRuleProposal',
                            Values: [1],
                        },
                    ],
                },
                {
                    ClientID: 'Team3',
                    Pairs: [
                        {
                            VariableName: 'PresidentLeftoverBudget',
                            Values: [50],
                        },
                    ],
                },
                {
                    ClientID: 'Team3',
                    Pairs: [
                        {
                            VariableName: 'RuleChosenFromProposalList',
                            Values: [1],
                        },
                    ],
                },
                {
                    ClientID: 'Team3',
                    Pairs: [
                        {
                            VariableName: 'AllocationMade',
                            Values: [1],
                        },
                    ],
                },
                {
                    ClientID: 'Team1',
                    Pairs: [
                        {
                            VariableName: 'SpeakerLeftoverBudget',
                            Values: [90],
                        },
                    ],
                },
                {
                    ClientID: 'Team1',
                    Pairs: [
                        {
                            VariableName: 'SpeakerLeftoverBudget',
                            Values: [80],
                        },
                    ],
                },
                {
                    ClientID: 'Team1',
                    Pairs: [
                        {
                            VariableName: 'IslandsAllowedToVote',
                            Values: [6],
                        },
                    ],
                },
                {
                    ClientID: 'Team1',
                    Pairs: [
                        {
                            VariableName: 'SpeakerProposedPresidentRule',
                            Values: [1],
                        },
                    ],
                },
                {
                    ClientID: 'Team1',
                    Pairs: [
                        {
                            VariableName: 'SpeakerLeftoverBudget',
                            Values: [70],
                        },
                    ],
                },
                {
                    ClientID: 'Team1',
                    Pairs: [
                        {
                            VariableName: 'AnnouncementRuleMatchesVote',
                            Values: [1],
                        },
                        {
                            VariableName: 'AnnouncementResultMatchesVote',
                            Values: [1],
                        },
                    ],
                },
                {
                    ClientID: 'Team1',
                    Pairs: [
                        {
                            VariableName: 'RuleSelected',
                            Values: [1],
                        },
                        {
                            VariableName: 'VoteCalled',
                            Values: [1],
                        },
                    ],
                },
                {
                    ClientID: 'Team1',
                    Pairs: [
                        {
                            VariableName: 'VoteCalled',
                            Values: [1],
                        },
                        {
                            VariableName: 'VoteResultAnnounced',
                            Values: [1],
                        },
                    ],
                },
                {
                    ClientID: 'Team2',
                    Pairs: [
                        {
                            VariableName: 'PresidentPayment',
                            Values: [10],
                        },
                    ],
                },
                {
                    ClientID: 'Team1',
                    Pairs: [
                        {
                            VariableName: 'JudgePaid',
                            Values: [1],
                        },
                    ],
                },
                {
                    ClientID: 'Team3',
                    Pairs: [
                        {
                            VariableName: 'SpeakerPaid',
                            Values: [1],
                        },
                    ],
                },
            ],
            IITOTransactions: {
                Team1: {
                    Team3: {
                        AcceptedAmount: 5,
                        Reason: 0,
                    },
                    Team6: {
                        AcceptedAmount: 4.5125,
                        Reason: 0,
                    },
                },
                Team2: {},
                Team3: {
                    Team1: {
                        AcceptedAmount: 12.5,
                        Reason: 0,
                    },
                    Team2: {
                        AcceptedAmount: 12.5,
                        Reason: 0,
                    },
                    Team4: {
                        AcceptedAmount: 12.5,
                        Reason: 0,
                    },
                    Team5: {
                        AcceptedAmount: 12.5,
                        Reason: 0,
                    },
                    Team6: {
                        AcceptedAmount: 4.999999999999999,
                        Reason: 0,
                    },
                },
                Team4: {},
                Team5: {},
                Team6: {
                    Team3: {
                        AcceptedAmount: 11.785113019775793,
                        Reason: 0,
                    },
                },
            },
            SpeakerID: 'Team1',
            JudgeID: 'Team2',
            PresidentID: 'Team3',
        },
    ],
}

test('test getIIGOTransactions', () => {
    const want = [
        {
            name: 1,
            expectedAlloc: 6,
            actualAlloc: 5,
            expectedTax: 14,
            actualTax: 14,
            expectedSanction: 0,
            actualSanction: 0,
        },
        {
            name: 2,
            expectedAlloc: 1,
            actualAlloc: 3,
            expectedTax: 1,
            actualTax: 51,
            expectedSanction: 0,
            actualSanction: 0,
        },
        {
            name: 3,
            expectedAlloc: 56,
            actualAlloc: 43,
            expectedTax: 2,
            actualTax: 0,
            expectedSanction: 0,
            actualSanction: 0,
        },
        {
            name: 4,
            expectedAlloc: 6,
            actualAlloc: 0,
            expectedTax: 14,
            actualTax: 14,
            expectedSanction: 0,
            actualSanction: 0,
        },
        {
            name: 5,
            expectedAlloc: 0,
            actualAlloc: 0,
            expectedTax: 13,
            actualTax: 14,
            expectedSanction: 0,
            actualSanction: 0,
        },
        {
            name: 6,
            expectedAlloc: 2,
            actualAlloc: 4,
            expectedTax: 13,
            actualTax: 0,
            expectedSanction: 0,
            actualSanction: 0,
        },
    ]
    expect(processPaymentsData(testInput)).toEqual(want)
})

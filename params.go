package main

import (
	"flag"

	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/pkg/errors"
)

var (
	// config.Config
	maxSeasons = flag.Uint(
		"maxSeasons",
		100,
		"The maximum number of 1-indexed seasons to run the game.",
	)
	maxTurns = flag.Uint(
		"maxTurns",
		50,
		"The maximum numbers of 1-indexed turns to run the game.",
	)
	initialResources = flag.Float64(
		"initialResources",
		50,
		"The default number of resources at the start of the game.",
	)
	initialCommonPool = flag.Float64(
		"initialCommonPool",
		0,
		"The default number of resources in the common pool at the start of the game.",
	)
	costOfLiving = flag.Float64(
		"costOfLiving",
		10,
		"Subtracted from an islands pool before the next turn.\n"+
			"This is the simulation-level equivalent to using resources to stay \n"+
			"alive (e.g. food consumed). These resources are permanently consumed and do \n"+
			" NOT go into the common pool. Note: this is NOT the same as the tax",
	)
	minimumResourceThreshold = flag.Float64(
		"minimumResourceThreshold",
		5,
		"The minimum resources required for an island to not be in Critical state.",
	)
	maxCriticalConsecutiveTurns = flag.Uint(
		"maxCriticalConsecutiveTurns",
		3,
		"The maximum consecutive turns an island can be in the critical state.",
	)

	// config.ForagingConfig.DeerHuntConfig
	foragingDeerMaxPerHunt = flag.Uint(
		"foragingMaxDeerPerHunt",
		5,
		"Max possible number of deer on a single hunt (regardless of number of participants). ** should be strictly less than max deer population.",
	)
	foragingDeerIncrementalInputDecay = flag.Float64(
		"foragingDeerIncrementalInputDecay",
		0.9,
		"Determines decay of incremental input cost of hunting more deer.",
	)
	foragingDeerBernoulliProb = flag.Float64(
		"foragingDeerBernoulliProb",
		0.95,
		"`p` param in D variable (see foraging README). Controls prob of catching a deer or not.",
	)
	foragingDeerExponentialRate = flag.Float64(
		"foragingDeerExponentialRate",
		0.3,
		"`lambda` param in W variable (see foraging README). Controls distribution of deer sizes.",
	)
	foragingDeerInputScaler = flag.Float64(
		"foragingDeerInputScaler",
		18,
		"scalar value that adjusts deer input resources to be in a range that is commensurate with cost of living, salaries etc.",
	)
	foragingDeerOutputScaler = flag.Float64(
		"foragingDeerOutputScaler",
		18,
		"scalar value that adjusts deer returns to be in a range that is commensurate with cost of living, salaries etc.",
	)
	foragingDeerDistributionStrategy = flag.Int(
		"foragingDeerDistributionStrategy",
		int(shared.InputProportionalSplit),
		shared.HelpResourceDistributionStrategy(),
	)
	foragingDeerThetaCritical = flag.Float64(
		"foragingDeerThetaCritical",
		0.97,
		"Bernoulli prob of catching deer when population ratio = running population/max deer per hunt = 1",
	)
	foragingDeerThetaMax = flag.Float64(
		"foragingDeerThetaMax",
		0.99,
		"Bernoulli prob of catching deer when population is at carrying capacity (max population)",
	)
	foragingDeerMaxPopulation = flag.Uint(
		"foragingDeerMaxPopulation",
		20,
		"Max possible deer population. ** Should be strictly greater than max deer per hunt.",
	)
	foragingDeerGrowthCoefficient = flag.Float64(
		"foragingDeerGrowthCoefficient",
		0.4,
		"Scaling parameter used in the population model. Larger coeff => deer pop. regenerates faster.",
	)

	// config.ForagingConfig.FishingConfig
	foragingFishMaxPerHunt = flag.Uint(
		"foragingMaxFishPerHunt",
		12,
		"Max possible catch (num. fish) on a single expedition (regardless of number of participants).",
	)
	foragingFishingIncrementalInputDecay = flag.Float64(
		"foragingFishingIncrementalInputDecay",
		0.95,
		"Determines decay of incremental input cost of catching more fish.",
	)
	foragingFishingMean = flag.Float64(
		"foragingFishingMean",
		1.45,
		"Determines mean of normal distribution of fishing return (see foraging README)",
	)
	foragingFishingVariance = flag.Float64(
		"foragingFishingVariance",
		0.1,
		"Determines variance of normal distribution of fishing return (see foraging README)",
	)
	foragingFishingInputScaler = flag.Float64(
		"foragingFishingInputScaler",
		18,
		"scalar value that adjusts input resources to be in a range that is commensurate with cost of living, salaries etc.",
	)
	foragingFishingOutputScaler = flag.Float64(
		"foragingFishingOutputScaler",
		18,
		"scalar value that adjusts returns to be in a range that is commensurate with cost of living, salaries etc.",
	)
	foragingFishingDistributionStrategy = flag.Int(
		"foragingFishingDistributionStrategy",
		int(shared.EqualSplit),
		shared.HelpResourceDistributionStrategy(),
	)

	// config.DisasterConfig
	disasterXMin = flag.Float64(
		"disasterXMin",
		0,
		"Min x bound of archipelago (bounds for possible disaster).",
	)
	disasterXMax = flag.Float64(
		"disasterXMax",
		10,
		"Max x bound of archipelago (bounds for possible disaster).",
	)
	disasterYMin = flag.Float64(
		"disasterYMin",
		0,
		"Min y bound of archipelago (bounds for possible disaster).",
	)
	disasterYMax = flag.Float64(
		"disasterYMax",
		10,
		"Max y bound of archipelago (bounds for possible disaster).",
	)
	disasterPeriod = flag.Uint(
		"disasterPeriod",
		5,
		"Period T between disasters in deterministic case and E[T] in stochastic case.",
	)
	disasterSpatialPDFType = flag.Int(
		"disasterSpatialPDFType",
		0,
		shared.HelpSpatialPDFType(),
	)
	disasterMagnitudeLambda = flag.Float64(
		"disasterMagnitudeLambda",
		1,
		"Exponential rate param for disaster magnitude",
	)
	disasterMagnitudeResourceMultiplier = flag.Float64(
		"disasterMagnitudeResourceMultiplier",
		85,
		"Multiplier to map disaster magnitude to CP resource deductions",
	)
	disasterCommonpoolThreshold = flag.Float64(
		"disasterCommonpoolThreshold",
		200,
		"Common pool threshold value for disaster to be mitigated",
	)
	disasterStochasticPeriod = flag.Bool(
		"disasterStochasticPeriod",
		false,
		"If true, period between disasters becomes random. If false, it will be consistent (deterministic)",
	)
	disasterCommonpoolThresholdVisible = flag.Bool(
		"disasterCommonpoolThresholdVisible",
		false,
		"Whether disasterCommonpoolThreshold is visible to agents",
	)
	disasterPeriodVisible = flag.Bool(
		"disasterPeriodVisible",
		true,
		"Whether disasterPeriod is visible to agents",
	)
	disasterStochasticPeriodVisible = flag.Bool(
		"disasterStochasticPeriodVisible",
		true,
		"Whether stochasticPeriod is visible to agents",
	)

	// config.IIGOConfig - Executive branch
	iigoGetRuleForSpeakerActionCost = flag.Float64(
		"iigoGetRuleForSpeakerActionCost",
		2,
		"IIGO action cost for getRuleForSpeaker action",
	)
	iigoBroadcastTaxationActionCost = flag.Float64(
		"iigoBroadcastTaxationActionCost",
		2,
		"IIGO action cost for broadcastTaxation action",
	)
	iigoReplyAllocationRequestsActionCost = flag.Float64(
		"iigoReplyAllocationRequestsActionCost",
		2,
		"IIGO action cost for replyAllocationRequests action",
	)
	iigoRequestAllocationRequestActionCost = flag.Float64(
		"iigoRequestAllocationRequestActionCost",
		2,
		"IIGO action cost for requestAllocationRequest action",
	)
	iigoRequestRuleProposalActionCost = flag.Float64(
		"iigoRequestRuleProposalActionCost",
		2,
		"IIGO action cost for requestRuleProposal action",
	)
	iigoAppointNextSpeakerActionCost = flag.Float64(
		"iigoAppointNextSpeakerActionCost",
		2,
		"IIGO action cost for appointNextSpeaker action",
	)

	// config.IIGOConfig - Judiciary branch
	iigoInspectHistoryActionCost = flag.Float64(
		"iigoInspectHistoryActionCost",
		2,
		"IIGO action cost for inspectHistory",
	)

	historicalRetributionActionCost = flag.Float64(
		"historicalRetributionActionCost",
		2,
		"IIGO action cost for inspectHistory retroactively (in turns before the last one)",
	)

	iigoInspectBallotActionCost = flag.Float64(
		"iigoInspectBallotActionCost",
		2,
		"IIGO action cost for inspectBallot",
	)

	iigoInspectAllocationActionCost = flag.Float64(
		"iigoInspectAllocationActionCost",
		2,
		"IIGO action cost for inspectAllocation",
	)

	iigoAppointNextPresidentActionCost = flag.Float64(
		"iigoAppointNextPresidentActionCost",
		2,
		"IIGO action cost for appointNextPresident",
	)

	iigoDefaultSanctionScore = flag.Uint(
		"iigoDefaultSanctionScore",
		2,
		"Default penalty score for breaking a rule",
	)

	iigoSanctionCacheDepth = flag.Uint(
		"iigoSanctionCacheDepth",
		3,
		"Turn depth of sanctions to be applied or pardoned",
	)

	iigoHistoryCacheDepth = flag.Uint(
		"iigoHistoryCacheDepth",
		3,
		"Turn depth of history cache for events to be evaluated",
	)

	iigoAssumedResourcesNoReport = flag.Uint(
		"iigoAssumedResourcesNoReport",
		100,
		"If an island doesn't report usaged this value is assumed for sanction calculations",
	)

	iigoSanctionLength = flag.Uint(
		"iigoSanctionLength",
		5,
		"Sanction length for all sanctions",
	)

	// config.IIGOConfig - Legislative branch
	iigoSetVotingResultActionCost = flag.Float64(
		"iigoSetVotingResultActionCost",
		2,
		"IIGO action cost for setVotingResult",
	)

	iigoSetRuleToVoteActionCost = flag.Float64(
		"iigoSetRuleToVoteActionCost",
		2,
		"IIGO action cost for setRuleToVote action",
	)

	iigoAnnounceVotingResultActionCost = flag.Float64(
		"iigoAnnounceVotingResultActionCost",
		2,
		"IIGO action cost for announceVotingResult action",
	)

	iigoUpdateRulesActionCost = flag.Float64(
		"iigoUpdateRulesActionCost",
		2,
		"IIGO action cost for updateRules action",
	)

	iigoAppointNextJudgeActionCost = flag.Float64(
		"iigoAppointNextJudgeActionCost",
		2,
		"IIGO action cost for appointNextJudge action",
	)

	iigoTermLengthPresident = flag.Uint(
		"iigoTermLengthPresident",
		4,
		"Length of the term for the President",
	)

	iigoTermLengthSpeaker = flag.Uint(
		"iigoTermLengthSpeaker",
		4,
		"Length of the term for the Speaker",
	)

	iigoTermLengthJudge = flag.Uint(
		"iigoTermLengthJudge",
		4,
		"Length of the term for the Judge",
	)

	startWithRulesInPlay = flag.Bool(
		"startWithRulesInPlay",
		true,
		"Pull all available rules into play at start of run",
	)
)

func parseConfig() (config.Config, error) {
	flag.Parse()

	parsedForagingDeerDistributionStrategy, err := shared.ParseResourceDistributionStrategy(*foragingDeerDistributionStrategy)
	if err != nil {
		return config.Config{}, errors.Errorf("Error parsing foragingDeerDistributionStrategy: %v", err)
	}

	parsedDeerMaxPerHunt, parseDeerMaxPopulation, err := shared.ParseDeerPopulationParams(
		*foragingDeerMaxPerHunt,
		*foragingDeerMaxPopulation,
	)
	if err != nil {
		return config.Config{}, errors.Errorf("Error parsing foragingDeerMaxPerHunt and/or foragingDeerMaxPopulation: %v", err)
	}

	parsedForagingFishingDistributionStrategy, err := shared.ParseResourceDistributionStrategy(*foragingFishingDistributionStrategy)
	if err != nil {
		return config.Config{}, errors.Errorf("Error parsing foragingFishingDistributionStrategy: %v", err)
	}

	parsedDisasterSpatialPDFType, err := shared.ParseSpatialPDFType(*disasterSpatialPDFType)
	if err != nil {
		return config.Config{}, errors.Errorf("Error parsing disasterSpatialPDFType: %v", err)
	}

	deerConf := config.DeerHuntConfig{
		//Deer parameters
		MaxDeerPerHunt:        parsedDeerMaxPerHunt,
		IncrementalInputDecay: *foragingDeerIncrementalInputDecay,
		BernoulliProb:         *foragingDeerBernoulliProb,
		ExponentialRate:       *foragingDeerExponentialRate,
		InputScaler:           *foragingDeerInputScaler,
		OutputScaler:          *foragingDeerOutputScaler,
		DistributionStrategy:  parsedForagingDeerDistributionStrategy,
		ThetaCritical:         *foragingDeerThetaCritical,
		ThetaMax:              *foragingDeerThetaMax,
		MaxDeerPopulation:     parseDeerMaxPopulation,
		DeerGrowthCoefficient: *foragingDeerGrowthCoefficient,
	}
	fishingConf := config.FishingConfig{
		// Fish parameters
		MaxFishPerHunt:        *foragingFishMaxPerHunt,
		IncrementalInputDecay: *foragingFishingIncrementalInputDecay,
		Mean:                  *foragingFishingMean,
		Variance:              *foragingFishingVariance,
		InputScaler:           *foragingFishingInputScaler,
		OutputScaler:          *foragingFishingOutputScaler,
		DistributionStrategy:  parsedForagingFishingDistributionStrategy,
	}
	foragingConf := config.ForagingConfig{
		DeerHuntConfig: deerConf,
		FishingConfig:  fishingConf,
	}
	disasterConf := config.DisasterConfig{
		XMin:                        *disasterXMin,
		XMax:                        *disasterXMax,
		YMin:                        *disasterYMin,
		YMax:                        *disasterYMax,
		Period:                      *disasterPeriod,
		SpatialPDFType:              parsedDisasterSpatialPDFType,
		MagnitudeLambda:             *disasterMagnitudeLambda,
		StochasticPeriod:            *disasterStochasticPeriod,
		MagnitudeResourceMultiplier: *disasterMagnitudeResourceMultiplier,
		CommonpoolThreshold:         shared.Resources(*disasterCommonpoolThreshold),
		CommonpoolThresholdVisible:  *disasterCommonpoolThresholdVisible,
		PeriodVisible:               *disasterPeriodVisible,
		StochasticPeriodVisible:     *disasterStochasticPeriodVisible,
	}

	iigoConf := config.IIGOConfig{
		IIGOTermLengths: map[shared.Role]uint{shared.President: *iigoTermLengthPresident,
			shared.Speaker: *iigoTermLengthSpeaker,
			shared.Judge:   *iigoTermLengthJudge},
		// Executive branch
		GetRuleForSpeakerActionCost:        shared.Resources(*iigoGetRuleForSpeakerActionCost),
		BroadcastTaxationActionCost:        shared.Resources(*iigoBroadcastTaxationActionCost),
		ReplyAllocationRequestsActionCost:  shared.Resources(*iigoReplyAllocationRequestsActionCost),
		RequestAllocationRequestActionCost: shared.Resources(*iigoRequestAllocationRequestActionCost),
		RequestRuleProposalActionCost:      shared.Resources(*iigoRequestRuleProposalActionCost),
		AppointNextSpeakerActionCost:       shared.Resources(*iigoAppointNextSpeakerActionCost),
		// Judiciary branch
		InspectHistoryActionCost:        shared.Resources(*iigoInspectHistoryActionCost),
		HistoricalRetributionActionCost: shared.Resources(*historicalRetributionActionCost),
		InspectBallotActionCost:         shared.Resources(*iigoInspectBallotActionCost),
		InspectAllocationActionCost:     shared.Resources(*iigoInspectAllocationActionCost),
		AppointNextPresidentActionCost:  shared.Resources(*iigoAppointNextPresidentActionCost),
		DefaultSanctionScore:            shared.IIGOSanctionsScore(*iigoDefaultSanctionScore),
		SanctionCacheDepth:              *iigoSanctionCacheDepth,
		HistoryCacheDepth:               *iigoHistoryCacheDepth,
		AssumedResourcesNoReport:        shared.Resources(*iigoAssumedResourcesNoReport),
		SanctionLength:                  *iigoSanctionLength,
		// Legislative branch
		SetVotingResultActionCost:      shared.Resources(*iigoSetVotingResultActionCost),
		SetRuleToVoteActionCost:        shared.Resources(*iigoSetRuleToVoteActionCost),
		AnnounceVotingResultActionCost: shared.Resources(*iigoAnnounceVotingResultActionCost),
		UpdateRulesActionCost:          shared.Resources(*iigoUpdateRulesActionCost),
		AppointNextJudgeActionCost:     shared.Resources(*iigoAppointNextJudgeActionCost),
		StartWithRulesInPlay:           *startWithRulesInPlay,
	}

	return config.Config{
		MaxSeasons:                  *maxSeasons,
		MaxTurns:                    *maxTurns,
		InitialResources:            shared.Resources(*initialResources),
		InitialCommonPool:           shared.Resources(*initialCommonPool),
		CostOfLiving:                shared.Resources(*costOfLiving),
		MinimumResourceThreshold:    shared.Resources(*minimumResourceThreshold),
		MaxCriticalConsecutiveTurns: *maxCriticalConsecutiveTurns,
		ForagingConfig:              foragingConf,
		DisasterConfig:              disasterConf,
		IIGOConfig:                  iigoConf,
	}, nil
}

package main

import (
	"flag"

	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
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
		100,
		"The maximum numbers of 1-indexed turns to run the game.",
	)
	initialResources = flag.Float64(
		"initialResources",
		100,
		"The default number of resources at the start of the game.",
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
	// config.ForagingConfig
	foragingMaxDeerPerHunt = flag.Uint(
		"foragingMaxDeerPerHunt",
		4,
		"Max possible number of deer on a single hunt (regardless of number of participants).",
	)
	foragingIncrementalInputDecay = flag.Float64(
		"foragingIncrementalInputDecay",
		0.8,
		"Determines decay of incremental input cost of hunting more deer.",
	)
	foragingBernoulliProb = flag.Float64(
		"foragingBernoulliProb",
		0.95,
		"`p` param in D variable (see foraging README). Controls prob of catching a deer or not.",
	)
	foragingExponentialRate = flag.Float64(
		"foragingExponentialRate",
		1,
		"`lambda` param in W variable (see foraging README). Controls distribution of deer sizes.",
	)
	foragingMaxDeerPopulation = flag.Uint(
		"maxDeerPopulation",
		12,
		"Max possible deer population.",
	)
	foragingDeerGrowthCoefficient = flag.Float64(
		"foragingDeerGrowthCoefficient",
		0.4,
		"Scaling parameter used in the population model. Larger coeff => deer pop. regenerates faster.",
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
	disasterGlobalProb = flag.Float64(
		"disasterGlobalProb",
		0.1,
		"Bernoulli 'p' param. Chance of a disaster occurring.",
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

	// config.IIGOConfig - Executive branch
	getRuleForSpeakerActionCost = flag.Float64(
		"getRuleForSpeakerActionCost",
		10,
		"IIGO action cost for getRuleForSpeaker action",
	)
	broadcastTaxationActionCost = flag.Float64(
		"broadcastTaxationActionCost",
		10,
		"IIGO action cost for broadcastTaxation action",
	)
	replyAllocationRequestsActionCost = flag.Float64(
		"replyAllocationRequestsActionCost",
		10,
		"IIGO action cost for replyAllocationRequests action",
	)
	requestAllocationRequestActionCost = flag.Float64(
		"requestAllocationRequestActionCost",
		10,
		"IIGO action cost for requestAllocationRequest action",
	)
	requestRuleProposalActionCost = flag.Float64(
		"requestRuleProposalActionCost",
		10,
		"IIGO action cost for requestRuleProposal action",
	)
	appointNextSpeakerActionCost = flag.Float64(
		"appointNextSpeakerActionCost",
		10,
		"IIGO action cost for appointNextSpeaker action",
	)

	// config.IIGOConfig - Judiciary branch
	inspectHistoryActionCost = flag.Float64(
		"inspectHistoryActionCost",
		10,
		"IIGO action cost for inspectHistory",
	)

	inspectBallotActionCost = flag.Float64(
		"inspectBallotActionCost",
		10,
		"IIGO action cost for inspectBallot",
	)

	inspectAllocationActionCost = flag.Float64(
		"inspectAllocationActionCost",
		10,
		"IIGO action cost for inspectAllocation",
	)

	appointNextPresidentActionCost = flag.Float64(
		"appointNextPresidentActionCost",
		10,
		"IIGO action cost for appointNextPresident",
	)

	// config.IIGOConfig - Legislative branch
	setVotingResultActionCost = flag.Float64(
		"setVotingResultActionCost",
		10,
		"IIGO action cost for setVotingResult",
	)

	setRuleToVoteActionCost = flag.Float64(
		"setRuleToVoteActionCost",
		10,
		"IIGO action cost for setRuleToVote action",
	)

	announceVotingResultActionCost = flag.Float64(
		"announceVotingResultActionCost",
		10,
		"IIGO action cost for announceVotingResult action",
	)

	updateRulesActionCost = flag.Float64(
		"updateRulesActionCost",
		10,
		"IIGO action cost for updateRules action",
	)

	appointNextJudgeActionCost = flag.Float64(
		"appointNextJudgeActionCost",
		10,
		"IIGO action cost for appointNextJudge action",
	)
)

func parseConfig() config.Config {
	flag.Parse()
	foragingConf := config.ForagingConfig{
		MaxDeerPerHunt:        *foragingMaxDeerPerHunt,
		IncrementalInputDecay: *foragingIncrementalInputDecay,
		BernoulliProb:         *foragingBernoulliProb,
		ExponentialRate:       *foragingExponentialRate,

		MaxDeerPopulation:     *foragingMaxDeerPopulation,
		DeerGrowthCoefficient: *foragingDeerGrowthCoefficient,
	}
	disasterConf := config.DisasterConfig{
		XMin:            *disasterXMin,
		XMax:            *disasterXMax, // chosen quite arbitrarily for now
		YMin:            *disasterYMin,
		YMax:            *disasterYMax,
		GlobalProb:      *disasterGlobalProb,
		SpatialPDFType:  shared.ParseSpatialPDFType(*disasterSpatialPDFType),
		MagnitudeLambda: *disasterMagnitudeLambda,
	}

	iigoConf := config.IIGOConfig{

		// Executive branch
		GetRuleForSpeakerActionCost:        shared.Resources(*getRuleForSpeakerActionCost),
		BroadcastTaxationActionCost:        shared.Resources(*broadcastTaxationActionCost),
		ReplyAllocationRequestsActionCost:  shared.Resources(*replyAllocationRequestsActionCost),
		RequestAllocationRequestActionCost: shared.Resources(*requestAllocationRequestActionCost),
		RequestRuleProposalActionCost:      shared.Resources(*requestRuleProposalActionCost),
		AppointNextSpeakerActionCost:       shared.Resources(*appointNextSpeakerActionCost),
		// Judiciary branch
		InspectHistoryActionCost:       shared.Resources(*inspectHistoryActionCost),
		InspectBallotActionCost:        shared.Resources(*inspectBallotActionCost),
		InspectAllocationActionCost:    shared.Resources(*inspectAllocationActionCost),
		AppointNextPresidentActionCost: shared.Resources(*appointNextPresidentActionCost),
		// Legislative branch
		SetVotingResultActionCost:      shared.Resources(*setVotingResultActionCost),
		SetRuleToVoteActionCost:        shared.Resources(*setRuleToVoteActionCost),
		AnnounceVotingResultActionCost: shared.Resources(*announceVotingResultActionCost),
		UpdateRulesActionCost:          shared.Resources(*updateRulesActionCost),
		AppointNextJudgeActionCost:     shared.Resources(*appointNextJudgeActionCost),
	}

	return config.Config{
		MaxSeasons:                  *maxSeasons,
		MaxTurns:                    *maxTurns,
		InitialResources:            shared.Resources(*initialResources),
		CostOfLiving:                shared.Resources(*costOfLiving),
		MinimumResourceThreshold:    shared.Resources(*minimumResourceThreshold),
		MaxCriticalConsecutiveTurns: *maxCriticalConsecutiveTurns,
		ForagingConfig:              foragingConf,
		DisasterConfig:              disasterConf,
		IIGOConfig:                  iigoConf,
	}
}

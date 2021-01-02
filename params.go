package main

import (
	"flag"

	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/pkg/errors"
)

var (
	//config.Config
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
	initialCommonPool = flag.Float64(
		"initialCommonPool",
		100,
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
	foragingDeerMaxPerHunt = flag.Uint(
		"foragingMaxDeerPerHunt",
		4,
		"Max possible number of deer on a single hunt (regardless of number of participants).",
	)
	foragingDeerIncrementalInputDecay = flag.Float64(
		"foragingDeerIncrementalInputDecay",
		0.8,
		"Determines decay of incremental input cost of hunting more deer.",
	)
	foragingDeerBernoulliProb = flag.Float64(
		"foragingDeerBernoulliProb",
		0.95,
		"`p` param in D variable (see foraging README). Controls prob of catching a deer or not.",
	)
	foragingDeerExponentialRate = flag.Float64(
		"foragingDeerExponentialRate",
		1,
		"`lambda` param in W variable (see foraging README). Controls distribution of deer sizes.",
	)
	foragingDeerResourceMultiplier = flag.Float64(
		"foragingDeerResourceMultiplier",
		1,
		"scalar value that adjusts returns to be in a range that is commensurate with cost of living, salaries etc.",
	)
	foragingDeerDistributionStrategy = flag.Int(
		"foragingDeerDistributionStrategy",
		int(shared.InputProportionalSplit),
		shared.HelpResourceDistributionStrategy(),
	)
	foragingDeerMaxPopulation = flag.Uint(
		"foragingDeerMaxPopulation",
		12,
		"Max possible deer population.",
	)
	foragingDeerGrowthCoefficient = flag.Float64(
		"foragingDeerGrowthCoefficient",
		0.4,
		"Scaling parameter used in the population model. Larger coeff => deer pop. regenerates faster.",
	)
	foragingFishMaxPerHunt = flag.Uint(
		"foragingMaxFishPerHunt",
		6,
		"Max possible catch (num. fish) on a single expedition (regardless of number of participants).",
	)
	foragingFishingIncrementalInputDecay = flag.Float64(
		"foragingFishingIncrementalInputDecay",
		0.8,
		"Determines decay of incremental input cost of catching more fish.",
	)
	foragingFishingMean = flag.Float64(
		"foragingFishingMean",
		0.9,
		"Determines mean of normal distribution of fishing return (see foraging README)",
	)
	foragingFishingVariance = flag.Float64(
		"foragingFishingVariance",
		0.2,
		"Determines variance of normal distribution of fishing return (see foraging README)",
	)
	foragingFishingResourceMultiplier = flag.Float64(
		"foragingFishingResourceMultiplier",
		1,
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
)

func parseConfig() (config.Config, error) {
	flag.Parse()

	parsedForagingDeerDistributionStrategy, err := shared.ParseResourceDistributionStrategy(*foragingDeerDistributionStrategy)
	if err != nil {
		return config.Config{}, errors.Errorf("Error parsing foragingDeerDistributionStrategy: %v", err)
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
		MaxDeerPerHunt:        *foragingDeerMaxPerHunt,
		IncrementalInputDecay: *foragingDeerIncrementalInputDecay,
		BernoulliProb:         *foragingDeerBernoulliProb,
		ExponentialRate:       *foragingDeerExponentialRate,
		ResourceMultiplier:    *foragingDeerResourceMultiplier,
		DistributionStrategy:  parsedForagingDeerDistributionStrategy,

		MaxDeerPopulation:     *foragingDeerMaxPopulation,
		DeerGrowthCoefficient: *foragingDeerGrowthCoefficient,
	}
	fishingConf := config.FishingConfig{
		// Fish parameters
		MaxFishPerHunt:        *foragingFishMaxPerHunt,
		IncrementalInputDecay: *foragingFishingIncrementalInputDecay,
		Mean:                  *foragingFishingMean,
		Variance:              *foragingFishingVariance,
		ResourceMultiplier:    *foragingFishingResourceMultiplier,
		DistributionStrategy:  parsedForagingFishingDistributionStrategy,
	}
	foragingConf := config.ForagingConfig{
		DeerHuntConfig: deerConf,
		FishingConfig:  fishingConf,
	}
	disasterConf := config.DisasterConfig{
		XMin:            *disasterXMin,
		XMax:            *disasterXMax, // chosen quite arbitrarily for now
		YMin:            *disasterYMin,
		YMax:            *disasterYMax,
		GlobalProb:      *disasterGlobalProb,
		SpatialPDFType:  parsedDisasterSpatialPDFType,
		MagnitudeLambda: *disasterMagnitudeLambda,
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
	}, nil
}

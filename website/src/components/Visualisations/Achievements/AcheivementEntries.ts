import { OutputJSONType } from "../../../consts/types";

type metricsType = {
    Team1: number,
    Team2: number,
    Team3: number,
    Team4: number,
    Team5: number,
    Team6: number,
};

const emptyMetrics = (): metricsType => ({
    Team1: 0,
    Team2: 0,
    Team3: 0,
    Team4: 0,
    Team5: 0,
    Team6: 0,
});

const teamNames = [ "Team1", "Team2", "Team3", "Team4", "Team5", "Team6" ];

export type AcheivementEntry = {
    title: string,
    description: string,
    collectMetrics: (data: OutputJSONType) => metricsType,
    evalLargest: boolean,
};

const peakResourcesMetricCollection = (data: OutputJSONType): metricsType => (
    data.GameStates.reduce((metrics: metricsType, gameState) => (
        teamNames.reduce((metAcc, teamName) => {
            const teamResources: number = gameState.ClientInfos[teamName].Resources
            metAcc[teamName] = teamResources > metAcc[teamName] ? teamResources : metAcc[teamName]
            return metAcc
        }, metrics)
    ), emptyMetrics())
);

const acheivementList: AcheivementEntry[] = [
    {
        title: "Jackpot!",
        description: "Island with the highest peak resources",
        collectMetrics: peakResourcesMetricCollection,
        evalLargest: true,
    },
    {
        title: "Baller",
        description: "Island with the highest average resources",
        collectMetrics: (data) => emptyMetrics(), //TODO: implement
        evalLargest: true,
    },
    {
        title: "Broke",
        description: "Island with the lowest average resources",
        collectMetrics: (data) => emptyMetrics(), //TODO: implement
        evalLargest: false,
    },
    {
        title: "Biggest G",
        description: "Island who gifted the most in total",
        collectMetrics: (data) => emptyMetrics(), //TODO: implement
        evalLargest: true,
    },
    {
        title: "Back to Life",
        description: "Island who returned from critical the most",
        collectMetrics: (data) => emptyMetrics(), //TODO: implement
        evalLargest: true,
    },
    {
        title: "Saving Grace",
        description: "Island who gifted the most to critical islands",
        collectMetrics: (data) => emptyMetrics(), //TODO: implement
        evalLargest: true,
    },
    {
        title: "Hoarder",
        description: "Island with lowest total resource outflow",
        collectMetrics: (data) => emptyMetrics(), //TODO: implement
        evalLargest: false,
    },
    {
        title: "Spenny",
        description: "Island with largest total resource outflow",
        collectMetrics: (data) => emptyMetrics(), //TODO: implement
        evalLargest: true,
    },
    {
        title: "Greedy",
        description: "Island with largest total resource inflow",
        collectMetrics: (data) => emptyMetrics(), //TODO: implement
        evalLargest: true,
    },
    {
        title: "I Will Survive",
        description: "Island alive the longest",
        collectMetrics: (data) => emptyMetrics(), //TODO: implement
        evalLargest: true,
    },
    {
        title: "F",
        description: "First island to die",
        collectMetrics: (data) => emptyMetrics(), //TODO: implement
        evalLargest: false,
    },
    {
        title: "Mother Teresa",
        description: "Island who gives the most with low resources",
        collectMetrics: (data) => emptyMetrics(), //TODO: implement
        evalLargest: true,
    },
    {
        title: "Hunter Gatherer",
        description: "Island who put in the highest average resources for foraging",
        collectMetrics: (data) => emptyMetrics(), //TODO: implement
        evalLargest: true,
    },
    {
        title: "Lazy",
        description: "Island who put in the lowest average resources for foraging",
        collectMetrics: (data) => emptyMetrics(), //TODO: implement
        evalLargest: false,
    },
    {
        title: "Scammed",
        description: "Island who got the lowest average returns from foraging",
        collectMetrics: (data) => emptyMetrics(), //TODO: implement
        evalLargest: false,
    },
    {
        title: "Lucky",
        description: "Island with least average damage from disasters",
        collectMetrics: (data) => emptyMetrics(), //TODO: implement
        evalLargest: false,
    },
    {
        title: "Unlucky",
        description: "Island with most average damage from disasters",
        collectMetrics: (data) => emptyMetrics(), //TODO: implement
        evalLargest: true,
    },
    {
        title: "Power Hungry",
        description: "Island who spent the most time in power",
        collectMetrics: (data) => emptyMetrics(), //TODO: implement
        evalLargest: true,
    },
    {
        title: "The Donald",
        description: "Island who spent the most time as President",
        collectMetrics: (data) => emptyMetrics(), //TODO: implement
        evalLargest: true,
    },
    {
        title: "Judge Judy",
        description: "Island who spent the most time as Judge",
        collectMetrics: (data) => emptyMetrics(), //TODO: implement
        evalLargest: true,
    },
    {
        title: "Speak Now or Forever Hold Your Peace",
        description: "Island who spent the most time as Speaker",
        collectMetrics: (data) => emptyMetrics(), //TODO: implement
        evalLargest: true,
    },
];

export default acheivementList;
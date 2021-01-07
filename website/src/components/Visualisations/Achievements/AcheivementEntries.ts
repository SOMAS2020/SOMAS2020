import { OutputJSONType } from "../../../consts/types";

export type AcheivementEntry = {
    title: string,
    description: string,
    collectMetrics: (data: OutputJSONType) => { team: string, metric: number }[],
    evalLargest: boolean,
}

const acheivementList: AcheivementEntry[] = [
    {
        title: "Jackpot!",
        description: "Island with the highest peak resources",
        collectMetrics: (data) => [], //TODO: implement
        evalLargest: true,
    },
    {
        title: "Baller",
        description: "Island with the highest average resources",
        collectMetrics: (data) => [], //TODO: implement
        evalLargest: true,
    },
    {
        title: "Broke",
        description: "Island with the lowest average resources",
        collectMetrics: (data) => [], //TODO: implement
        evalLargest: false,
    },
    {
        title: "Biggest G",
        description: "Island who gifted the most in total",
        collectMetrics: (data) => [], //TODO: implement
        evalLargest: true,
    },
    {
        title: "Back to Life",
        description: "Island who returned from critical the most",
        collectMetrics: (data) => [], //TODO: implement
        evalLargest: true,
    },
    {
        title: "Saving Grace",
        description: "Island who gifted the most to critical islands",
        collectMetrics: (data) => [], //TODO: implement
        evalLargest: true,
    },
    {
        title: "Hoarder",
        description: "Island with lowest total resource outflow",
        collectMetrics: (data) => [], //TODO: implement
        evalLargest: false,
    },
    {
        title: "Spenny",
        description: "Island with largest total resource outflow",
        collectMetrics: (data) => [], //TODO: implement
        evalLargest: true,
    },
    {
        title: "Greedy",
        description: "Island with largest total resource inflow",
        collectMetrics: (data) => [], //TODO: implement
        evalLargest: true,
    },
    {
        title: "I Will Survive",
        description: "Island alive the longest",
        collectMetrics: (data) => [], //TODO: implement
        evalLargest: true,
    },
    {
        title: "F",
        description: "First island to die",
        collectMetrics: (data) => [], //TODO: implement
        evalLargest: false,
    },
    {
        title: "Mother Teresa",
        description: "Island who gives the most with low resources",
        collectMetrics: (data) => [], //TODO: implement
        evalLargest: true,
    },
    {
        title: "Hunter Gatherer",
        description: "Island who put in the highest average resources for foraging",
        collectMetrics: (data) => [], //TODO: implement
        evalLargest: true,
    },
    {
        title: "Lazy",
        description: "Island who put in the lowest average resources for foraging",
        collectMetrics: (data) => [], //TODO: implement
        evalLargest: false,
    },
    {
        title: "Scammed",
        description: "Island who got the lowest average returns from foraging",
        collectMetrics: (data) => [], //TODO: implement
        evalLargest: false,
    },
    {
        title: "Lucky",
        description: "Island with least average damage from disasters",
        collectMetrics: (data) => [], //TODO: implement
        evalLargest: false,
    },
    {
        title: "Unlucky",
        description: "Island with most average damage from disasters",
        collectMetrics: (data) => [], //TODO: implement
        evalLargest: true,
    },
    {
        title: "Power Hungry",
        description: "Island who spent the most time in power",
        collectMetrics: (data) => [], //TODO: implement
        evalLargest: true,
    },
    {
        title: "The Donald",
        description: "Island who spent the most time as President",
        collectMetrics: (data) => [], //TODO: implement
        evalLargest: true,
    },
    {
        title: "Judge Judy",
        description: "Island who spent the most time as Judge",
        collectMetrics: (data) => [], //TODO: implement
        evalLargest: true,
    },
    {
        title: "Speak Now or Forever Hold Your Peace",
        description: "Island who spent the most time as Speaker",
        collectMetrics: (data) => [], //TODO: implement
        evalLargest: true,
    },
];

export default acheivementList;
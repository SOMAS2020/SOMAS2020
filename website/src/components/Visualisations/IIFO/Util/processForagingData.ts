import { OutputJSONType } from '../../../../consts/types'

const processForagingData = (data: OutputJSONType) => {
    let ForagingTurn: {
        turn: number
        inputResources: number
        numParticipants: number
        numCaught: number
        totalUtility: number
    }

    const ForagingHistory: {
        FishForagingHistory: typeof ForagingTurn[]
        DeerForagingHistory: typeof ForagingTurn[]
    } = { FishForagingHistory: [], DeerForagingHistory: [] }

    // TODO: strongly type the output json
    // Have to check the type of Foraging History Matches the expected type
    // TODO: remove intermediate fish/deer histories and do it all in one step
    // Cannot show individual catch sizes for now
    Object.entries(
        data.GameStates[data.GameStates.length - 1].ForagingHistory
            .DeerForageType
    ).forEach((el) => {
        ForagingHistory.DeerForagingHistory.push({
            turn: el[1].Turn,
            inputResources: el[1].InputResources,
            numParticipants: el[1].NumberParticipants,
            numCaught: el[1].NumberCaught,
            totalUtility: el[1].TotalUtility,
        })
    })

    Object.entries(
        data.GameStates[data.GameStates.length - 1].ForagingHistory
            .FishForageType
    ).forEach((el) => {
        ForagingHistory.FishForagingHistory.push({
            turn: el[1].Turn,
            inputResources: el[1].InputResources,
            numParticipants: el[1].NumberParticipants,
            numCaught: el[1].NumberCaught,
            totalUtility: el[1].TotalUtility,
        })
    })

    return ForagingHistory
}

export default processForagingData

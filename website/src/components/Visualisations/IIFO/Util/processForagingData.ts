import { OutputJSONType } from '../../../../consts/types'

const processForagingData = (data: OutputJSONType) => {
    // TODO: strongly type the output json
    // Have to check the type of Foraging History Matches the expected type
    const deerForageHistory = Object.entries(
        data.GameStates[data.GameStates.length - 1].ForagingHistory
            .DeerForageType
    )

    const fishForageHistory = Object.entries(
        data.GameStates[data.GameStates.length - 1].ForagingHistory
            .DeerForageType
    )

    // TODO: turn numbers are likely to always be the same
    // TODO: remove intermediate fish/deer histories and do it all in one step
    const deerTurnNumbers: number[] = []
    const fishTurnNumbers: number[] = []
    const deerInputResources: number[] = []
    const fishInputResources: number[] = []
    const deerNumParticipants: number[] = []
    const fishNumParticipants: number[] = []
    const deerNumCaught: number[] = []
    const fishNumCaught: number[] = []
    const deerTotalUtility: number[] = []
    const fishTotalUtility: number[] = []
    // Cannot show individual catch sizes for now

    // Map through foraging history, extract all the data
    // Add entries in expected format for vis ForagingPlot
    // const deerForageHistory = foragingHistory.DeerForageType
    // const fishForageHistory = foragingHistory.FishForageType

    // console.log({ deerForageHistory })
    // console.log({ fishForageHistory })

    deerForageHistory.forEach((el) => {
        deerTurnNumbers.push(el[1].Turn)
        deerInputResources.push(el[1].InputResources)
        deerNumParticipants.push(el[1].NumberParticipants)
        deerNumCaught.push(el[1].NumberCaught)
        deerTotalUtility.push(el[1].TotalUtility)
    })

    fishForageHistory.forEach((el) => {
        fishTurnNumbers.push(el[1].Turn)
        fishInputResources.push(el[1].InputResources)
        fishNumParticipants.push(el[1].NumberParticipants)
        fishNumCaught.push(el[1].NumberCaught)
        fishTotalUtility.push(el[1].TotalUtility)
    })

    console.log({ deerTurnNumbers })
    console.log({ fishTurnNumbers })
    return 0
}

export default processForagingData

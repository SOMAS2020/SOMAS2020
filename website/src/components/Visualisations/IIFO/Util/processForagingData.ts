import { OutputJSONType } from '../../../../consts/types'
import { ForagingTurn, ForagingHistory } from './ForagingTypes'

const processForagingData: (data: OutputJSONType) => ForagingHistory = (
    data: OutputJSONType
) => {
    // TODO: poor form: probably don't need all these interemediate datastructures
    const fishTurns: number[] = []
    const deerTurns: number[] = []
    const deerInputResources: number[] = []
    const deerNumParticipants: number[] = []
    const deerNumCaught: number[] = []
    const deerTotalUtility: number[] = []
    const fishInputResources: number[] = []
    const fishNumParticipants: number[] = []
    const fishNumCaught: number[] = []
    const fishTotalUtility: number[] = []

    // reduce the full foraging hist into single foragingTurn objects
    // map through each foraging type and merge the result of that map all into one object

    Object.entries(
        data.GameStates[data.GameStates.length - 1].ForagingHistory
    ).forEach((foragingType) => {
        // this happens twice so we get twice the no of turns
        const foragingTypeData = foragingType[1]
        foragingTypeData.forEach((turn) => {
            if (turn.ForageType === 'DeerForageType') {
                deerInputResources.push(turn.InputResources)
                deerNumParticipants.push(turn.NumberParticipants)
                deerNumCaught.push(turn.NumberCaught)
                deerTotalUtility.push(turn.TotalUtility)
                deerTurns.push(turn.Turn)
            } else if (turn.ForageType === 'FishForageType') {
                fishInputResources.push(turn.InputResources)
                fishNumParticipants.push(turn.NumberParticipants)
                fishNumCaught.push(turn.NumberCaught)
                fishTotalUtility.push(turn.TotalUtility)
                fishTurns.push(turn.Turn)
            }
        })
    })

    // TODO: strongly type the output json
    // Have to check the type of Foraging History Matches the expected type
    // TODO: remove intermediate fish/deer histories and do it all in one step
    // Cannot show individual catch sizes for now
    const acc: ForagingTurn[] = []

    deerTurns.forEach((turn) => {
        acc.push({
            turn,
            deerInputResources: deerInputResources[turn - 1],
            deerNumParticipants: deerNumParticipants[turn - 1],
            deerNumCaught: deerNumCaught[turn - 1],
            deerTotalUtility: deerTotalUtility[turn - 1],
            fishNumParticipants: fishNumParticipants[turn - 1],
            fishInputResources: fishInputResources[turn - 1],
            fishNumCaught: fishNumCaught[turn - 1],
            fishTotalUtility: fishTotalUtility[turn - 1],
        })
    })
    return acc
}

export default processForagingData

import _ from 'lodash'
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
            switch (turn.ForageType) {
                case 'DeerForageType':
                    deerInputResources.push(turn.InputResources)
                    deerNumParticipants.push(turn.NumberParticipants)
                    deerNumCaught.push(turn.NumberCaught)
                    deerTotalUtility.push(turn.TotalUtility)
                    deerTurns.push(turn.Turn)
                    break
                case 'FishForageType':
                    fishInputResources.push(turn.InputResources)
                    fishNumParticipants.push(turn.NumberParticipants)
                    fishNumCaught.push(turn.NumberCaught)
                    fishTotalUtility.push(turn.TotalUtility)
                    fishTurns.push(turn.Turn)
                    break
                default:
                // shouldn't happen
            }
        })
    })

    const acc: ForagingTurn[] = []

    deerTurns.forEach((turn) => {
        acc.push({
            turn,
            deerInputResources: _.isUndefined(deerInputResources[turn - 1])
                ? 0
                : deerInputResources[turn - 1],
            deerNumParticipants: _.isUndefined(deerNumParticipants[turn - 1])
                ? 0
                : deerNumParticipants[turn - 1],
            deerNumCaught: _.isUndefined(deerNumCaught[turn - 1])
                ? 0
                : deerNumCaught[turn - 1],
            deerTotalUtility: _.isUndefined(deerTotalUtility[turn - 1])
                ? 0
                : deerTotalUtility[turn - 1],
            fishNumParticipants: _.isUndefined(fishNumParticipants[turn - 1])
                ? 0
                : fishNumParticipants[turn - 1],
            fishInputResources: _.isUndefined(fishInputResources[turn - 1])
                ? 0
                : fishInputResources[turn - 1],
            fishNumCaught: _.isUndefined(fishNumCaught[turn - 1])
                ? 0
                : fishNumCaught[turn - 1],
            fishTotalUtility: _.isUndefined(fishTotalUtility[turn - 1])
                ? 0
                : fishTotalUtility[turn - 1],
        })
    })
    return acc
}

export default processForagingData

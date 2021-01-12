import { OutputJSONType } from '../../../../consts/types'
import { ForagingTurn } from './ForagingTypes'

const processForagingData: (data: OutputJSONType) => ForagingTurn[] = (
    data: OutputJSONType
) => {
    const turns = data.GameStates.length - 1
    let res: ForagingTurn[] = new Array(turns)

    const fillRange = (start, end) => {
        return [...Array(end - start + 1)].map((item, index) => ({
            turn: index,
            deerInputResources: 0,
            deerNumCaught: 0,
            deerTotalUtility: 0,
            fishInputResources: 0,
            fishNumCaught: 0,
            fishTotalUtility: 0,
        }))
    }

    // Fill in the blanks and turns
    res = fillRange(0, turns)

    Object.entries(
        data.GameStates[data.GameStates.length - 1].ForagingHistory
    ).forEach((foragingType) => {
        const foragingTypeData = foragingType[1]
        foragingTypeData.forEach((forageInfo) => {
            switch (forageInfo.ForageType) {
                case 'DeerForageType':
                    res[forageInfo.Turn - 1] = {
                        turn: res[forageInfo.Turn - 1].turn,
                        fishInputResources:
                            res[forageInfo.Turn - 1].fishInputResources,
                        fishNumCaught: res[forageInfo.Turn - 1].fishNumCaught,
                        fishTotalUtility:
                            res[forageInfo.Turn - 1].fishTotalUtility,
                        deerInputResources: forageInfo.InputResources,
                        deerNumCaught: forageInfo.NumberCaught,
                        deerTotalUtility: forageInfo.TotalUtility,
                    }
                    break
                case 'FishForageType':
                    res[forageInfo.Turn - 1] = {
                        turn: res[forageInfo.Turn - 1].turn,
                        fishInputResources: forageInfo.InputResources,
                        fishNumCaught: forageInfo.NumberCaught,
                        fishTotalUtility: forageInfo.TotalUtility,
                        deerInputResources:
                            res[forageInfo.Turn - 1].deerInputResources,
                        deerNumCaught: res[forageInfo.Turn - 1].deerNumCaught,
                        deerTotalUtility:
                            res[forageInfo.Turn - 1].deerTotalUtility,
                    }

                    break
                default:
                // this shouldn't happen
            }
        })
    })

    return res
}

export default processForagingData

import { OutputJSONType } from '../../../../consts/types'
import { IIGOInfos, IIGOInfo } from './RoleTypes'

export const getIIGOStatuses = (data: OutputJSONType): IIGOInfos => {
    if (data.GameStates.length <= 1) return []

    const retData: IIGOInfos = []

    return data.GameStates.slice(1).map<IIGOInfo>((gameState) => {
        const info: IIGOInfo = {
            turn: gameState.Turn - 1,
            status: gameState.IIGORunStatus,
        }

        return info
    }, retData)
}

export default getIIGOStatuses

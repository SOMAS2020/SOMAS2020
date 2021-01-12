import { OutputJSONType } from '../../../../consts/types'
import { IIGOStatuses, IIGOStatus } from './RoleTypes'

export const getIIGOStatuses = (data: OutputJSONType): IIGOStatuses => {
    if (data.GameStates.length <= 1) return []

    const retData: IIGOStatuses = []

    return data.GameStates.slice(1).map<IIGOStatus>((gameState) => {
        const info: IIGOStatus = {
            turn: gameState.Turn - 1,
            status: gameState.IIGORunStatus,
        }

        return info
    }, retData)
}

export default getIIGOStatuses

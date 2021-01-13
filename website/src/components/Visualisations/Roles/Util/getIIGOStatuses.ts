import { OutputJSONType } from '../../../../consts/types'
import { IIGOInfos } from './RoleTypes'

export const getIIGOStatuses = (data: OutputJSONType): IIGOInfos => {
    if (data.GameStates.length <= 1) return []

    return data.GameStates.slice(1).map((gameState) => {
        return {
            turn: gameState.Turn - 1,
            status: gameState.IIGORunStatus,
        }
    })
}

export default getIIGOStatuses

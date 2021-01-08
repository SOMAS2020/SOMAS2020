import * as localForage from 'localforage'
import { VIS_OUTPUT } from '../../consts/localForage'
import { OutputJSONType } from '../../consts/types'

export const loadLocalVisOutput = async () => {
    const got: OutputJSONType | null | undefined = await localForage.getItem(
        VIS_OUTPUT
    )
    if (got) return got
    return undefined
}

export const storeLocalVisOutput = async (o: OutputJSONType) => {
    await localForage.setItem(VIS_OUTPUT, o)
}

export const clearLocalVisOutput = async () => {
    await localForage.removeItem(VIS_OUTPUT)
}

export const teamColors = new Map([
    ['Team0', '#101D42'],
    ['Team1', '#0095FF'],
    ['Team2', '#FF0000'],
    ['Team3', '#802FF0'],
    ['Team4', '#00C49F'],
    ['Team5', '#FFBB28'],
    ['Team6', '#FF8042'],
])

import * as localForage from 'localforage'
import { VIS_OUTPUT } from '../../consts/localForage'
import { OutputJSONType } from '../../consts/types'

export const loadLocalVisOutput = async () => {
    const got: OutputJSONType | null | undefined = await localForage.getItem(VIS_OUTPUT)
    if (got) return got
    return undefined
}

export const storeLocalVisOutput = async (o: OutputJSONType) => {
    await localForage.setItem(VIS_OUTPUT, o)
}

export const clearLocalVisOutput = async () => {
    await localForage.removeItem(VIS_OUTPUT)
}
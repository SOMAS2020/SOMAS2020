import * as localForage from 'localforage'
import { NEW_RUN_FLAGS, NEW_RUN_OUTPUT } from '../../consts/localForage'
import {
    Flag,
    getFlagsFormats,
    RunGameReturnType,
    runGame,
} from '../../wasmAPI'

export const loadFlags = async (
    loadLocal = true
): Promise<Map<string, Flag>> => {
    const goFlags = await getFlagsFormats()
    const fs: Map<string, Flag> = new Map()
    goFlags.forEach((f) => {
        fs.set(f.Name, { ...f, Value: f.DefValue })
    })

    if (loadLocal) {
        // try to load local flags' values if present
        const localFs:
            | Map<string, Flag>
            | undefined
            | null = await localForage.getItem(NEW_RUN_FLAGS)
        if (localFs) {
            localFs.forEach((f) => {
                const origFlag = fs.get(f.Name)
                if (origFlag) {
                    fs.set(f.Name, { ...origFlag, Value: f.Value })
                }
            })
        }
    }

    return fs
}

export const loadLocalOutput = async () => {
    const localOutput:
        | RunGameReturnType
        | undefined
        | null = await localForage.getItem(NEW_RUN_OUTPUT)
    return localOutput
}

export const runGameHelper = async (
    flags: Map<string, Flag>
): Promise<RunGameReturnType> => {
    const flagArr = Array.from(flags, ([_, value]) => value)
    const output = await runGame(flagArr)

    // async-ally save the flags and output in localForage
    // best effort, so we don't really care if it's not oK
    localForage
        .setItem(NEW_RUN_FLAGS, flags)
        .then(() => console.debug(`Set local flags`))
        .catch((err: any) => console.error(err))
    localForage
        .setItem(NEW_RUN_OUTPUT, output)
        .then(() => console.debug(`Set local output`))
        .catch((err: any) => console.error(err))

    return output
}

export const setFlagHelper = async (
    flags: Map<string, Flag> | undefined,
    flagName: string,
    val: string
): Promise<Map<string, Flag>> => {
    if (!flags) {
        // should not happen
        throw new Error(`Flags not loaded`)
    }
    const currFlag = flags.get(flagName)

    if (!currFlag) {
        throw new Error(`Unknown flag name ${flagName}`)
    }
    const newCurrFlag = { ...currFlag, Value: val }
    return new Map(flags.set(flagName, newCurrFlag))
}

export const clearLocalOutput = async () => {
    localForage
        .removeItem(NEW_RUN_OUTPUT)
        .then(() => console.debug(`Clear local output`))
        .catch((err: any) => console.error(err))
}

export const clearLocalFlags = async () => {
    localForage
        .removeItem(NEW_RUN_FLAGS)
        .then(() => console.debug(`Clear local flags`))
        .catch((err: any) => console.error(err))
}

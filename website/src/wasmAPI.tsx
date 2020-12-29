import outputJSONData from './output/output.json'
import { notEmpty } from './utils'

// @ts-ignore
const go: any = window.go

export type GoFlag = {
    Name: string,
    Usage: string,
    DefValue: string,
    Type: string,
}

export type Flag = GoFlag & { Value: string }

type RunGameReturnTypeWASM = {
    output: string,
    logs: string,
    error: string,
}

export type RunGameReturnType = {
    output: typeof outputJSONData,
    logs: string,
}

type GetFlagsFormatsReturnTypeWASM = {
    output: string,
    error: string,
}

export type GetFlagsFormatsReturnType = GoFlag[]

let loaded = false

let runGameWASM: ((args: string[]) => RunGameReturnTypeWASM) | undefined;
let getFlagsFormatsWASM: (() => GetFlagsFormatsReturnTypeWASM) | undefined;

const load = async () => {
    const { instance } = await WebAssembly.instantiateStreaming(
        fetch(`${process.env.PUBLIC_URL}/SOMAS2020.wasm`),
        go.importObject
    )
    go.run(instance)

    // @ts-ignore
    runGameWASM = window.RunGame
    // @ts-ignore
    getFlagsFormatsWASM = window.GetFlagsFormats

    loaded = true
}

export const runGame = async (flags: Flag[]): Promise<RunGameReturnType> => {
    if (!loaded) {
        await load()
    }
    if (!runGameWASM) {
        throw new Error("Game not loaded properly")
    }

    const args = await prepareFlags(flags)

    const result = runGameWASM(args)
    if (result.error.length > 0) {
        throw new Error(result.error)
    }

    const processedOutput = JSON.parse(result.output) as typeof outputJSONData

    // we need to patch git info
    processedOutput.GitInfo = outputJSONData.GitInfo

    return {
        output: processedOutput,
        logs: result.logs,
    }
}

const prepareFlags = async (flags: Flag[]): Promise<string[]> => {
    return flags
        .map(f => {
            if (f.Value !== f.DefValue) {
                return [f.Name, f.Value]
            }
            return undefined
        })
        .filter(notEmpty)
        .flat()
}

export const getFlagsFormats = async (): Promise<GetFlagsFormatsReturnType> => {
    if (!loaded) {
        await load()
    }
    if (!getFlagsFormatsWASM) {
        throw new Error("Game not loaded properly")
    }

    const result = getFlagsFormatsWASM()
    if (result.error.length > 0) {
        throw new Error(result.error)
    }

    const processedOutput = JSON.parse(result.output) as GetFlagsFormatsReturnType

    return processedOutput
}
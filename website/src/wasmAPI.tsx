// Hacks to properly type global objects set by wasm 
import outputJSONData from './output/output.json'

// @ts-ignore
export const go: any = window.go

let loaded = false

type RunGameReturnTypeWASM = {
    output: string,
    logs: string,
    error: string,
}

type RunGameReturnType = {
    output: typeof outputJSONData | undefined,
    logs: string | undefined,
}

const load = async () => {
    const { instance } = await WebAssembly.instantiateStreaming(
        fetch(`${process.env.PUBLIC_URL}/SOMAS2020.wasm`),
        go.importObject
    )
    go.run(instance)
    loaded = true
}

export const runGame = async (): Promise<RunGameReturnType> => {
    if (!loaded) {
        await load()
    }

    // @ts-ignore
    const runGameWASM: () => RunGameReturnTypeWASM = window.RunGame

    const result = runGameWASM()
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

// Hacks to properly type global objects set by wasm 
import outputJSONData from './output/output.json'

// @ts-ignore
export const go: any = window.go

type RunGameReturnTypeWASM = {
    output: string,
    logs: string,
    error: string,
}

type RunGameReturnType = {
    output: typeof outputJSONData | undefined,
    logs: string | undefined,
}

export const runGame = async (): Promise<RunGameReturnType> => {
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

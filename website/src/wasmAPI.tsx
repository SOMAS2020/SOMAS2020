import { OutputJSONType } from './consts/types'
import outputJSONData from './output/output.json'
import './wasmjs/wasm_exec'

// @ts-ignore
const go: any = new Go()

export type GoFlag = {
  Name: string
  Usage: string
  DefValue: string
  Type: string
}

export type Flag = GoFlag & { Value: string }

type RunGameReturnTypeWASM = {
  output?: string
  logs?: string
  error: string
}

export type RunGameReturnType = {
  output: OutputJSONType
  logs: string
}

type GetFlagsFormatsReturnTypeWASM = {
  output?: string
  error: string
}

export type GetFlagsFormatsReturnType = GoFlag[]

let loaded = false

let runGameWASM: ((args: string) => RunGameReturnTypeWASM) | undefined
let getFlagsFormatsWASM: (() => GetFlagsFormatsReturnTypeWASM) | undefined

// Safari polyfill
if (!WebAssembly.instantiateStreaming) {
  WebAssembly.instantiateStreaming = async (resp, importObject) => {
    const source = await (await resp).arrayBuffer()
    // eslint-disable-next-line no-return-await
    return await WebAssembly.instantiate(source, importObject)
  }
}

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

/**
 * Take all the flags and make them into the string argument required by runGame
 * (`arg1=value,arg2=value,...`)
 *
 * @param flags all input flags with information initially gotten from getFlagsFormats
 */
const prepareFlags = async (flags: Flag[]): Promise<string> => {
  return flags.map((f) => `${f.Name}=${f.Value}`).join(`,`)
}
export const runGame = async (flags: Flag[]): Promise<RunGameReturnType> => {
  if (!loaded) {
    await load()
  }
  if (!runGameWASM) {
    throw new Error('Game not loaded properly')
  }

  const args = await prepareFlags(flags)

  const result = runGameWASM(args)
  if (result.error.length > 0) {
    if (result.output) {
      console.debug(result.output)
    }
    if (result.logs) {
      console.debug(result.logs)
    }
    throw new Error(result.error)
  }
  if (result.output === undefined || result.logs === undefined) {
    throw new Error(`Can't get output or logs`)
  }

  const processedOutput = JSON.parse(result.output) as OutputJSONType

  // we need to patch git info
  processedOutput.GitInfo = outputJSONData.GitInfo

  return {
    output: processedOutput,
    logs: result.logs,
  }
}

export const getFlagsFormats = async (): Promise<GetFlagsFormatsReturnType> => {
  if (!loaded) {
    await load()
  }
  if (!getFlagsFormatsWASM) {
    throw new Error('Game not loaded properly')
  }

  const result = getFlagsFormatsWASM()
  if (result.error.length > 0) {
    throw new Error(result.error)
  }
  if (result.output === undefined) {
    throw new Error(`Can't get output`)
  }
  const processedOutput = JSON.parse(result.output) as GetFlagsFormatsReturnType

  return processedOutput
}

import React, { useState, useEffect } from "react"
import { runGame, getFlagsFormats, RunGameReturnType, Flag } from "../../wasmAPI"
import { Alert, Button, Row } from 'react-bootstrap'
import { useLoadingState, initialLoadingState } from "../../contexts/loadingState"
import Artifacts from '../Artifacts/Artifacts'
import FlagForm from "./FlagForm"
import { setFlagWithValidation } from "./utils"


const NewRun = () => {
  const [, setLoading] = useLoadingState()
  const [output, setOutput] = useState<RunGameReturnType | undefined>(undefined)
  const [runError, setRunError] = useState<string | undefined>(undefined)
  const [flags, setFlags] = useState<Map<string, Flag> | undefined>(undefined)

  useEffect(() => {
    loadFlags()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const loadFlags = async () => {
    setLoading({ loading: true, loadingText: `We're finding the archipelago...` })
    try {
      const goFlags = await getFlagsFormats()
      let fs: Map<string, Flag> = new Map()
      goFlags.forEach(f => { fs.set(f.Name, { ...f, Value: f.DefValue, InvalidReason: undefined }) })
      setFlags(fs)
    }
    catch (err) {
      setRunError(err.message)
    }
    setLoading(initialLoadingState)
  }

  const setFlag = async (flagName: string, val: string) => {
    if (!flags) {
      // should not happen
      throw new Error(`Flags not loaded`)
    }
    const currFlag = flags.get(flagName)

    if (!currFlag) {
      throw new Error(`Unknown flag name ${flagName}`)
    }
    const newCurrFlag = await setFlagWithValidation(currFlag, val)

    setFlags(new Map(flags.set(flagName, newCurrFlag)))
  }

  const reset = async () => {
    setOutput(undefined)
  }

  const run = async () => {
    setLoading({ loading: true, loadingText: `Minions in your computer are running the agents!` })

    // sleep for 100ms for stability
    await new Promise(r => setTimeout(r, 100));

    try {
      if (!flags) {
        throw new Error(`Flags unset!`)
      }
      const flagArr = Array.from(flags, ([, value]) => value)
      const res = await runGame(flagArr)
      setOutput(res)
    }
    catch (err) {
      setRunError(err.message)
    }
    setLoading(initialLoadingState)
  }

  const runDisabled = () => {
    if (!flags) return true

    const flagsHasInvalid = 
      Array.from(flags, ([, value]) => value.InvalidReason)
        .map(a => a !== undefined) // true if invalid
        .reduce((a, b) => (a || b), false)

    if (flagsHasInvalid) return false

    return false
  }

  const getFlagForms = (fs: Map<string, Flag>): JSX.Element[] => {
    let ret: JSX.Element[] = []
    fs.forEach((f, n) => {
      ret.push(
        <FlagForm setFlag={(val) => setFlag(n, val)} flag={f} disabled={output !== undefined} />
      )
    })
    return ret
  }

  return <div>
    <h1>New Run</h1>

    {
      !output &&
      <div>
        <Button variant="success" size="lg" onClick={run} disabled={runDisabled()}>Run</Button>
      </div>
    }
    {
      flags &&
      <div>
        <Button variant="danger" size="lg" onClick={loadFlags} disabled={output !== undefined} style={{ marginTop: 24 }}>Reset Flags</Button>
        <Row style={{ marginLeft: `5vw`, marginRight: `5vw`, marginTop: 48, marginBottom: 48 }}>
          {getFlagForms(flags)}
        </Row>
      </div>
    }
    {
      runError &&
      <Alert variant="danger" onClose={() => setRunError(undefined)} dismissible style={{ maxWidth: `90vw`, margin: `auto`, marginTop: 24 }}>
        <Alert.Heading>Oh reeeeeeeeee!</Alert.Heading>
        <p>{runError}</p>
      </Alert>
    }
    {
      output &&
      <div style={{ marginTop: 48 }}>
        <div>
          <Button variant="danger" size="lg" onClick={reset}>Reset Run</Button>
        </div>

        <h3 style={{ marginTop: 24 }}>Artifacts</h3>
        <Artifacts output={output.output} logs={output.logs} />
      </div>
    }
  </div>

}

export default NewRun
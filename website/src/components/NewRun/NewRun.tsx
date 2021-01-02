import React, { useState, useEffect } from "react"
import { RunGameReturnType, Flag } from "../../wasmAPI"
import { Alert, Button, Row, Col, OverlayTrigger, Tooltip, Form } from 'react-bootstrap'
import { useLoadingState, initialLoadingState } from "../../contexts/loadingState"
import Artifacts from '../Artifacts/Artifacts'
import { clearLocalOutput, loadFlags, loadLocalOutput, runGameHelper, setFlagHelper, clearLocalFlags } from './utils'


type FlagFormProps = {
  flag: Flag,
  setFlag: (val: string) => Promise<void>,
  disabled: boolean,
}

const FlagForm = (props: FlagFormProps) => {
  const { flag, setFlag, disabled } = props

  const handleChange = async (event: React.ChangeEvent<any>) => {
    await setFlag(event.target.value)
  }
  return <Col xs={4}>
    <Form>
      <Form.Group>
        <Form.Label>
          <OverlayTrigger
            placement="top"
            overlay={
              <Tooltip id={flag.Name}>
                {flag.Usage} (Type: {flag.Type}, Default: {flag.DefValue})
              </Tooltip>
            }
          >
            <span style={{ wordBreak: `break-all` }}>{flag.Name}</span>
          </OverlayTrigger >
        </Form.Label>
        <Form.Control value={flag.Value} onChange={handleChange} readOnly={disabled} />
      </Form.Group>
    </Form>
  </Col >
}

const NewRun = () => {
  const [, setLoading] = useLoadingState()
  const [output, setOutput] = useState<RunGameReturnType | undefined>(undefined)
  const [runError, setRunError] = useState<string | undefined>(undefined)
  const [flags, setFlags] = useState<Map<string, Flag> | undefined>(undefined)

  useEffect(() => {
    onDidMount()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const onDidMount = async () => {
    setLoading({ loading: true, loadingText: `We're finding the archipelago...` })
    try {
      const fs = await loadFlags()
      setFlags(fs)

      const localOutput = await loadLocalOutput()
      if (localOutput) {
        setOutput(localOutput)
      }

      setRunError(undefined)
    }
    catch (err) {
      setRunError(err.message)
    }
    setLoading(initialLoadingState)
  }

  const setFlag = async (flagName: string, val: string) => {
    try {
      const newFlags = await setFlagHelper(flags, flagName, val)
      setFlags(newFlags)
      setRunError(undefined)
    }
    catch (err) {
      setRunError(err.message)
    }
  }

  const reset = async () => {
    setOutput(undefined)
    clearLocalOutput()
  }

  const run = async () => {
    setLoading({ loading: true, loadingText: `Minions in your computer are running the agents!` })

    // sleep for 100ms for stability
    await new Promise(r => setTimeout(r, 100));

    try {
      if (!flags) {
        throw new Error(`Flags unset!`)
      }
      const res = await runGameHelper(flags)
      setOutput(res)
      setRunError(undefined)
    }
    catch (err) {
      setRunError(err.message)
    }
    setLoading(initialLoadingState)
  }

  const getFlagForms = (fs: Map<string, Flag>): JSX.Element[] => {
    let ret: JSX.Element[] = []
    fs.forEach((f, n) => {
      ret.push(
        <FlagForm key={n} setFlag={(val) => setFlag(n, val)} flag={f} disabled={output !== undefined} />
      )
    })
    return ret
  }

  const resetFlags = async () => {
    const fs = await loadFlags(false)
    clearLocalFlags()
    setFlags(fs)
  }

  return <div>
    <h1>New Run</h1>

    {
      !output &&
      <div>
        <Button variant="success" size="lg" onClick={run} disabled={flags === undefined}>Run</Button>
      </div>
    }
    {
      flags &&
      <div>
        <Button variant="danger" size="lg" onClick={resetFlags} disabled={output !== undefined} style={{ marginTop: 24 }}>Reset Flags</Button>
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
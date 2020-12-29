import React, { useState, useEffect } from "react"
import { CodeBlock } from 'react-code-blocks'
import { runGame, getFlagsFormats, RunGameReturnType, GoFlag } from "../../wasmAPI"
import { Alert, Button, Row, Col, OverlayTrigger, Tooltip, Form } from 'react-bootstrap'
import { useLoadingState, initialLoadingState } from "../../contexts/loadingState"

import styles from '../CIOutput/CIOutput.module.css'

type Flag = GoFlag & { Value: string }

type flagFormProps = {
  flag: Flag,
  setFlag: (val: string) => Promise<void>,
  disabled: boolean,
}

const FlagForm = (props: flagFormProps) => {
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
                {flag.Usage}
              </Tooltip>
            }
          >
            <span>{flag.Name}</span>
          </OverlayTrigger >
        </Form.Label>
        <Form.Control value={flag.Value} onChange={handleChange} readOnly={disabled} />
      </Form.Group>
    </Form>
  </Col >
}

const NewRun = () => {
  const [_, setLoading] = useLoadingState()
  const [output, setOutput] = useState<RunGameReturnType | undefined>(undefined)
  const [runError, setRunError] = useState<string | undefined>(undefined)
  const [flags, setFlags] = useState<Map<string, Flag> | undefined>(undefined)

  useEffect(() => {
    loadFlags()
  }, [])

  const loadFlags = async () => {
    setLoading({ loading: true, loadingText: `We're finding the archipelago...` })
    try {
      const goFlags = await getFlagsFormats()
      let fs: Map<string, Flag> = new Map()
      goFlags.forEach(f => { fs.set(f.Name, { ...f, Value: f.DefValue }) })
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
    const newCurrFlag = { ...currFlag, Value: val }

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
      const res = await runGame()
      setOutput(res)
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
        <FlagForm setFlag={(val) => setFlag(n, val)} flag={f} disabled={output !== undefined} />
      )
    })
    return ret
  }

  return <div className={styles.root}>
    <h1>{output ? `Run Output` : `New Run`}</h1>

    {
      !output &&
      <div>
        <Button variant="success" size="lg" onClick={run} disabled={flags === undefined}>Run</Button>
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

        <div style={{ textAlign: `left`, padding: `0 3vw` }}>
          <div style={{ marginBottom: 100 }}>
            <h2><code>output.json</code></h2>
            <CodeBlock text={JSON.stringify(output.output, null, "\t")} wrapLines showLineNumbers language="json" theme="dracula" />
          </div>
          <div style={{ marginBottom: 100 }}>
            <h2><code>log.txt</code></h2>
            <CodeBlock text={output.logs} wrapLines showLineNumbers language="text" theme="dracula" />
          </div>
        </div>
      </div>
    }
  </div>

}

export default NewRun
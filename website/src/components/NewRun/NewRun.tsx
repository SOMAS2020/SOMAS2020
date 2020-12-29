import React, { useState, useEffect } from "react"
import { CodeBlock } from 'react-code-blocks'
import { runGame, getFlagsFormats, RunGameReturnType, GoFlag } from "../../wasmAPI"
import { Alert, Button, Row, Col, OverlayTrigger, Tooltip } from 'react-bootstrap'
import { useLoadingState, initialLoadingState } from "../../contexts/loadingState"

import styles from '../CIOutput/CIOutput.module.css'

type Flag = GoFlag & { Value: string }

type flagFormProps = {
  flag: GoFlag,
  setFlag: (val: string) => Promise<void>,
}

const FlagForm = (props: flagFormProps) => {
  const { flag, setFlag } = props

  return <Col xs={4}>
    <OverlayTrigger
      placement="top"
      overlay={
        <Tooltip id={flag.Name}>
          {flag.Usage}
        </Tooltip>
      }
    >
      <p style={{ fontWeight: 600 }}>{flag.Name}</p>
    </OverlayTrigger >
  </Col >
}

const NewRun = () => {
  const [_, setLoading] = useLoadingState()
  const [output, setOutput] = useState<RunGameReturnType | undefined>(undefined)
  const [runError, setRunError] = useState<string | undefined>(undefined)
  const [flags, setFlags] = useState<Record<string, Flag> | undefined>(undefined)

  useEffect(() => {
    loadFlags()
  }, [])

  const loadFlags = async () => {
    setLoading({ loading: true, loadingText: `We're finding the archipelago...` })
    try {
      const goflags = await getFlagsFormats()
      let fs: Record<string, Flag> = {}
      goflags.forEach(f => { fs[f.Name] = { ...f, Value: f.DefValue } })
      setFlags(fs)
    }
    catch (err) {
      setRunError(err.message)
    }
    setLoading(initialLoadingState)
  }

  const setFlag = async (flagName: string, val: string) => {
    if (!flags) {
      throw new Error(`flags undefined!`)
    }
    if (!(flagName in flags)) {
      throw new Error(`Unknown flag name ${flagName}`)
    }
    const newFlags = flags
    newFlags[flagName].Value = val
    setFlags(newFlags)
  }

  const reset = async () => {
    setOutput(undefined)
  }

  const run = async () => {
    setLoading({ loading: true, loadingText: `Minions in your computer are running the agents!` })
    try {
      const res = await runGame()
      setOutput(res)
    }
    catch (err) {
      setRunError(err.message)
    }
    setLoading(initialLoadingState)
  }

  return <div className={styles.root}>
    <h1>{output ? `Run Output` : `New Run`}</h1>

    {
      output ?
        <div>
          <Button variant="danger" size="lg" onClick={reset}>Reset</Button>
        </div>
        :
        <div>
          <Button variant="success" size="lg" onClick={run} disabled={flags === undefined}>Run</Button>
        </div>
    }
    {
      flags &&
      <Row style={{ marginLeft: `5vw`, marginRight: `5vw`, marginTop: 48 }}>
        {Object.entries(flags).map(([n, f]) => <FlagForm setFlag={(val) => setFlag(n, val)} flag={f} />)}
      </Row>
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
    }
  </div>

}

export default NewRun
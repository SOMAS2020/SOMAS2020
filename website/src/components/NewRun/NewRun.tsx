import React, { useState } from "react"
import { CodeBlock } from 'react-code-blocks'
import { runGame, RunGameReturnType } from "../../wasmAPI"
import { Alert, Button } from 'react-bootstrap'
import { useLoadingState, initialLoadingState } from "../../contexts/loadingState"

import styles from '../CIOutput/CIOutput.module.css'

const NewRun = () => {
  const [_, setLoading] = useLoadingState()
  const [output, setOutput] = useState<RunGameReturnType | undefined>(undefined)
  const [runError, setRunError] = useState<string|undefined>(undefined)

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
          <Button variant="success" size="lg" onClick={run}>Run</Button>
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
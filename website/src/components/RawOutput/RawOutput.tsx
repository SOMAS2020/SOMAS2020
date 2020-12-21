import React from 'react'
import { CodeBlock } from 'react-code-blocks'

import styles from './RawOutput.module.css'

import outputJSON from '../../output/output.json'
import outputLog from '../../output/log.txt.json'

const processedOutputLog = outputLog.join(`\n`)

const RawOutput = () => {
    return <div className={styles.root}>
        <h1>Raw Output</h1>

        <div style={{ textAlign: `left`, padding: `0 3vw` }}>
            <div style={{ marginBottom: 100 }}>
                <h2><code>output.json</code></h2>
                <CodeBlock text={JSON.stringify(outputJSON, null, "\t")} wrapLines showLineNumbers language="json" theme="dracula"/>
            </div>
            <div style={{ marginBottom: 100 }}>
                <h2><code>log.txt</code></h2>
                <CodeBlock text={processedOutputLog} wrapLines showLineNumbers language="text" theme="dracula"/>
            </div>
        </div>
    </div>
}

export default RawOutput
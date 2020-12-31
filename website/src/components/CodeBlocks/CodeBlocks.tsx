import React from "react"
import { CodeBlock } from 'react-code-blocks'

import styles from './CodeBlocks.module.css'

type Props = {
    output: string,
    logs: string,
}

const CodeBlocks = (props: Props) => {
    const { output, logs } = props
    return <div style={{ textAlign: `left`, padding: `0 3vw` }} className={styles.root}>
        <div style={{ marginBottom: 100 }}>
            <h2><code>output.json</code></h2>
            <CodeBlock text={output} wrapLines showLineNumbers language="json" theme="dracula" />
        </div>
        <div style={{ marginBottom: 100 }}>
            <h2><code>log.txt</code></h2>
            <CodeBlock text={logs} wrapLines showLineNumbers language="text" theme="dracula" />
        </div>
    </div>
}

export default CodeBlocks 
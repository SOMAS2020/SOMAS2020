import React from 'react'
import outputJSON from '../../output/output.json'
import outputLog from '../../output/log.txt.json'
import CodeBlocks from '../CodeBlocks/CodeBlocks'

import styles from './CIOutput.module.css'

const processedOutputLog = outputLog.join(`\n`)

const CIOutput = () => {
    return <div className={styles.root}>
        <h1>CI Output</h1>

        <CodeBlocks output={JSON.stringify(outputJSON, null, "\t")} logs={processedOutputLog} />
    </div>
}

export default CIOutput
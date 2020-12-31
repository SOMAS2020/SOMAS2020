import React from 'react'

import outputJSON from '../../output/output.json'
import outputLog from '../../output/log.txt.json'
import Artifacts from '../Artifacts/Artifacts'

const processedOutputLog = outputLog.join(`\n`)

const CIOutput = () => {
  return <div>
    <h1>CI Output</h1>
    <h3 style={{ marginTop: 24 }}>Artifacts</h3>
    <Artifacts output={outputJSON} logs={processedOutputLog} />
  </div>
}

export default CIOutput
import React from 'react'

import outputJSON from '../../output/output.json'
import outputLog from '../../output/log.txt.json'
import Artifacts from '../Artifacts/Artifacts'
import VisualiseButton from '../VisualiseButton/VisualiseButton'

const processedOutputLog = outputLog.join(`\n`)

const CIOutput = () => {
  return (
    <div style={{ paddingTop: 24 }}>
      <h1>CI Output</h1>
      <p>Note that max turns is set to 20</p>
      <h3 style={{ marginTop: 24 }}>Artifacts</h3>
      <Artifacts output={outputJSON} logs={processedOutputLog} />

      <h3 style={{ marginTop: 24 }}>Visualise</h3>
      <VisualiseButton output={outputJSON} />
    </div>
  )
}

export default CIOutput

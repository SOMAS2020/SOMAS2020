import React from 'react'
import styles from './IIFO.module.css'

import { OutputJSONType } from '../../../consts/types'
import ForagingPlot from './Util/ForagingPlot'
import processForagingData from './Util/processForagingData'
// import { ForagingHistory } from './Util/ForagingTypes'

const IIFO = (props: { output: OutputJSONType }) => {
  const foragingHistory = processForagingData(props.output)

  console.log({ foragingHistory })
  return (
    <div className={styles.root}>
      <h2 className={styles.text}>Foraging Visualisation</h2>
      <div style={{ textAlign: 'center' }}>
        <ForagingPlot data={Object.values(foragingHistory)} />
      </div>
    </div>
  )
}

export default IIFO

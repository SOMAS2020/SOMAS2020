import React from 'react'
import logo from '../../../assets/logo/logo512.png'
import styles from './IIFO.module.css'

import { OutputJSONType } from '../../../consts/types'
import ForagingPlot from './Util/ForagingPlot'
import processForagingData from './Util/processForagingData'

// TODO: sketch out what the plot should look like and what you need

const IIFO = (props: { output: OutputJSONType }) => {
  const foragingHistory = processForagingData(props.output)

  return (
    <div className={styles.root}>
      <h2 className={styles.text}>Foraging Visualisation</h2>
      <div style={{ textAlign: 'center' }}>
        <ForagingPlot data={foragingHistory} />
      </div>
    </div>
  )
}

export default IIFO

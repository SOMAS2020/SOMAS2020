import React from 'react'
import logo from '../../../assets/logo/logo512.png'
import styles from './IIFO.module.css'

import { OutputJSONType } from '../../../consts/types'
import ForagingPlot from './Util/ForagingPlot'

const IIFO = (props: { output: OutputJSONType }) => {
  return (
    <div className={styles.root}>
      <h2 className={styles.text}>Foraging Visualisation</h2>
      <ForagingPlot />
    </div>
  )
}

export default IIFO

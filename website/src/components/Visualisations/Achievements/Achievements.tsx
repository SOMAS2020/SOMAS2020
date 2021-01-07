import React from 'react'
import logo from '../../../assets/logo/logo512.png'
import styles from './Achievements.module.css'
import { OutputJSONType } from '../../../consts/types'
import acheivementList, { evaluateMetrics } from './AcheivementEntries'

const AcheivementList = (props: { data: OutputJSONType }) => (
  <div>
    {acheivementList.map((achievement) => (
      <p key={achievement.title}>
        {achievement.title}:{' '}
        {evaluateMetrics(props.data, achievement).join(',')}
      </p>
    ))}
  </div>
)

const Achievements = (props: { output: OutputJSONType }) => {
  return (
    <div className={styles.root}>
      <p className={styles.text}>Achievements Visualisation</p>
      <AcheivementList data={props.output} />
    </div>
  )
}

export default Achievements

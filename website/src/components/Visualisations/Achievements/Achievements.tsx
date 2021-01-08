import React from 'react'
import logo from '../../../assets/logo/logo512.png'
import styles from './Achievements.module.css'
import { OutputJSONType } from '../../../consts/types'
import acheivementList, {
  evaluateMetrics,
  TeamName,
} from './AcheivementEntries'

type AchievementBarProps = {
  title: string
  desc: string | unknown
  winArr: TeamName[]
}

const AchievementBar = ({ title, desc, winArr }: AchievementBarProps) => {
  return (
    <p>
      {title}: {winArr.join(',')}
    </p>
  )
}

const AcheivementList = (props: { data: OutputJSONType }) => (
  <div>
    {acheivementList.map((achievement) => (
      <AchievementBar
        key={achievement.title}
        title={achievement.title}
        desc={achievement.description}
        winArr={evaluateMetrics(props.data, achievement)}
      />
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

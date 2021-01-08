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

const IndivAchievement = ({ title, desc, winArr }: AchievementBarProps) => {
  return (
    <p>
      {title}: {winArr.join(',')}
    </p>
  )
}

const Achievements = (props: { output: OutputJSONType }) => {
  return (
    <div className={styles.root}>
      <p className={styles.text}>Achievements Visualisation</p>
      {acheivementList.map((achievement) => (
        <IndivAchievement
          key={achievement.title}
          title={achievement.title}
          desc={achievement.description}
          winArr={evaluateMetrics(props.output, achievement)}
        />
      ))}
    </div>
  )
}

export default Achievements

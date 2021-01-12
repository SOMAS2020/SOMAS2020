import React, { useEffect, useState } from 'react'
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts'
import { teamColors } from '../utils'
import styles from './Roles.module.css'

import { ProcessedRoleData, TeamAndTurns, RoleName } from './Util/RoleTypes'
import { processRoleData } from './Util/ProcessedRoleData'
import { OutputJSONType } from '../../../consts/types'
import IIGOStatus from './IIGOStatus'

type CustomTooltipProps = {
  active: boolean
  payload: [{ name: string; value: number; unit: string }]
  label: string
  data: ProcessedRoleData
}

const CustomTooltip = ({ active, label, data }: CustomTooltipProps) => {
  const getTurnsAsTeams = (role: RoleName): TeamAndTurns =>
    data
      .find((elem) => elem.role === role)
      ?.occupied?.reduce((acc, tAndT) => acc.add(tAndT), new TeamAndTurns()) ??
    new TeamAndTurns()

  if (active && data.length > 0) {
    const turnsAsTeams = getTurnsAsTeams(label as RoleName)
    const totalTurns = data[0].occupied.reduce(
      (acc, elem) => acc + elem.turns(),
      0
    )

    const newLabel = label === 'Pres' ? 'President' : label

    return (
      <div className={styles.customTooltip}>
        <p className={styles.label}>{newLabel}</p>
        {turnsAsTeams.map((team, turns) => (
          <p
            className={styles.content}
            key={team}
            style={{ color: teamColors.get(team) }}
          >
            {`Turns as ${team}: ${turns} (${(
              (turns * 100) /
              totalTurns
            ).toFixed(1)} %)`}
          </p>
        ))}
      </div>
    )
  }

  return null
}

const Roles = (props: { output: OutputJSONType }) => {
  const [data, setData] = useState(processRoleData(props.output))

  useEffect(() => {
    setData(processRoleData(props.output))
  }, [props.output])

  const teams = ['Team1', 'Team2', 'Team3', 'Team4', 'Team5', 'Team6', 'NotRun']

  const localTeamColor: Map<string, string> = teamColors
  localTeamColor.set('NotRun', '#787878')

  return (
    <div className={styles.root}>
      <p className={styles.text}>Role Visualisation</p>
      <ResponsiveContainer height={460} width="100%">
        <BarChart data={data} layout="vertical">
          <YAxis type="category" dataKey="role" />
          <XAxis
            type="number"
            domain={[0, 'dataMax']}
            tickCount={20}
            allowDecimals={false}
          />
          <Tooltip
            content={(p: CustomTooltipProps) => CustomTooltip({ ...p, data })}
          />
          <Legend
            verticalAlign="top"
            payload={teams.map((team, i) => ({
              value: team,
              type: 'square',
              id: `${team}${i}`,
              color: localTeamColor.get(team),
            }))}
          />
          {data[0].occupied.map((a, i) => [
            teams.map((team) => (
              <Bar
                dataKey={`occupied[${i}].${team}`}
                stackId="a"
                fill={localTeamColor.get(team)}
                key={`${i.toString()}${team}`}
              />
            )),
          ])}
        </BarChart>
      </ResponsiveContainer>
      <p className={styles.graphLabel}>Turns</p>
      <IIGOStatus output={props.output} />
    </div>
  )
}

export default Roles

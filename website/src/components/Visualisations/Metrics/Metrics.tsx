import React from 'react'
import { IconButton, List, ListItemText, Collapse } from '@material-ui/core'
import { Table } from 'react-bootstrap'
import KeyboardArrowDownIcon from '@material-ui/icons/KeyboardArrowDown'
import KeyboardArrowUpIcon from '@material-ui/icons/KeyboardArrowUp'
import confetti from 'canvas-confetti'
import styles from './Metrics.module.css'
import { OutputJSONType } from '../../../consts/types'
import metricList, { evaluateMetrics, Metric } from './MetricEntries'

type MetricBarProps = {
  title: string
  desc: string
  metrics: Metric[]
}

const Metrics = (props: { output: OutputJSONType }) => {
  return (
    <>
      <h2 className={styles.text}>Metrics Summary</h2>
      <Table striped bordered hover size="sm">
        <thead>
          <th>Metric</th>
          <th>Description</th>
          <th>Team 1</th>
          <th>Team 2</th>
          <th>Team 3</th>
          <th>Team 4</th>
          <th>Team 5</th>
          <th>Team 6</th>
        </thead>
        <tbody>
          {metricList.map((metric) => (
            <tr>
              <td style={{ textAlign: 'left' }}>{metric.title}</td>
              <td style={{ textAlign: 'left' }}>{metric.description}</td>
              {evaluateMetrics(props.output, metric).map((team) => (
                <td style={{ textAlign: 'right' }} key={team.value.toString()}>
                  {team.value.toFixed(2).toString()}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </Table>
    </>
  )
}

export default Metrics

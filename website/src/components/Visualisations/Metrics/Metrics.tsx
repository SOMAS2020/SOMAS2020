import React from 'react'
import { Table } from 'react-bootstrap'
import styles from './Metrics.module.css'
import { OutputJSONType } from '../../../consts/types'
import metricList, { evaluateMetrics, Metric } from './MetricEntries'
import { numAgents } from '../utils'

const Metrics = (props: { output: OutputJSONType }) => {
  const totalAgents = numAgents(props.output)
  const teamHeadings: string[] = []

  for (let i = 0; i < totalAgents; i++) {
    teamHeadings.push(`Team${i + 1}`)
  }

  return (
    <>
      <h2 className={styles.text}>Metrics Summary</h2>
      <Table striped bordered hover size="sm">
        <thead>
          <tr>
            <th>Metric</th>
            <th>Description</th>
            {teamHeadings.map((val) => (
              <th>{val}</th>
            ))}
          </tr>
        </thead>
        <tbody>
          {metricList.map((metric) => (
            <tr key={metric.title}>
              <td style={{ textAlign: 'left' }}>
                <b>{metric.title}</b>
              </td>
              <td style={{ textAlign: 'left' }}>{metric.description}</td>
              {evaluateMetrics(props.output, metric).map((team) => (
                <td
                  style={{ textAlign: 'right' }}
                  // hacky but heyo if it works
                  key={`${team.teamName}+${metric.title}`}
                >
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

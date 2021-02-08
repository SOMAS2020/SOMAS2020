import React from 'react'
import { Table } from 'react-bootstrap'
import styles from './Metrics.module.css'
import { OutputJSONType } from '../../../consts/types'
import metricList, { evaluateMetrics, Metric } from './MetricEntries'

const Metrics = (props: { output: OutputJSONType }) => {
  return (
    <>
      <h2 className={styles.text}>Metrics Summary</h2>
      <Table striped bordered hover size="sm">
        <thead>
          <tr>
            <th>Metric</th>
            <th>Description</th>
            <th>Team 1</th>
            <th>Team 2</th>
            <th>Team 3</th>
            <th>Team 4</th>
            <th>Team 5</th>
            <th>Team 6</th>
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

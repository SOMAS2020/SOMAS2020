import React from 'react'
import { List, ListItemText } from '@material-ui/core'
import { Col, Container, Row } from 'react-bootstrap'
import confetti from 'canvas-confetti'
import styles from './Metrics.module.css'
import { OutputJSONType } from '../../../consts/types'
import metricList, { evaluateMetrics, TeamName, Metric } from './MetricEntries'

type MetricBarProps = {
  title: string
  desc: string
  metrics: Metric[]
}

const IndivMetric = ({ title, desc, metrics }: MetricBarProps) => {
  function handleAchievementClick() {
    confetti({
      particleCount: 300,
      spread: 100,
      origin: { y: 0.6 },
    })
  }

  return (
    <div className={styles.metricContainer}>
      <Container
        className={styles.innerContainer}
        onClick={handleAchievementClick}
      >
        <Row>
          <Col className={styles.leftColumn}>
            <h4 style={{ textAlign: 'left' }}>{title}</h4>
            <p style={{ textAlign: 'left' }}>{desc}</p>
          </Col>
        </Row>
        <Row>
          <Col className={styles.centerColumn}>
            <List component="nav">
              {metrics.map((metric) => (
                <ListItemText
                  primary={metric.teamName}
                  secondary={metric.value.toFixed(2)}
                />
              ))}
            </List>
          </Col>
        </Row>
      </Container>
    </div>
  )
}

const Metrics = (props: { output: OutputJSONType }) => {
  return (
    <div className={styles.root}>
      <p className={styles.text} style={{ marginBottom: 30 }}>
        Metrics
      </p>
      {metricList.map((metric) => (
        <IndivMetric
          key={metric.title}
          title={metric.title}
          desc={metric.description}
          metrics={evaluateMetrics(props.output, metric)}
        />
      ))}
    </div>
  )
}

export default Metrics

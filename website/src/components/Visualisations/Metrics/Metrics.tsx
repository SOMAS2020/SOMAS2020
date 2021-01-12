import React from 'react'
import { IconButton, List, ListItemText, Collapse } from '@material-ui/core'
import { Col, Container, Row } from 'react-bootstrap'
import KeyboardArrowDownIcon from '@material-ui/icons/KeyboardArrowDown'
import KeyboardArrowUpIcon from '@material-ui/icons/KeyboardArrowUp'
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
  const [open, setOpen] = React.useState(false)

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
        <IconButton
          aria-label="expand row"
          size="small"
          onClick={() => setOpen(!open)}
        >
          {open ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
        </IconButton>
        <Row>
          <Col className={styles.centerColumn}>
            <Collapse in={open} timeout="auto" unmountOnExit>
              <List component="nav">
                {metrics.map((metric) => (
                  <ListItemText
                    primary={[metric.teamName, ': ', metric.value.toFixed(2)]}
                    key={metric.teamName.toString()}
                  />
                ))}
              </List>
            </Collapse>
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

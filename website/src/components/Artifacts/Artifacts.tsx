import React from "react"
import { Table } from "react-bootstrap"

import styles from './Artifacts.module.css'
import outputJSON from '../../output/output.json'

type Props = {
  output: typeof outputJSON,
  logs: string,
}

type Item = {
  fileName: string
  description: string
  content: string
}

const DownloadLink = (props: { item: Item }) => {
  const { item } = props
  const blob = new Blob([item.content], { type: `text/plain` })
  return <a 
    download={item.fileName} href={URL.createObjectURL(blob)}
    >
      Download
  </a>
}

const Artifacts = (props: Props) => {
  const { output, logs } = props
  const outputTxt = JSON.stringify(output, null, `\t`)


  const items: Item[] = [
    {
      fileName: `output.json`,
      description: `Output JSON containing game states, configuration and other information.`,
      content: outputTxt,
    },
    {
      fileName: `log.txt`,
      description: `Logs`,
      content: logs,
    },
  ]

  return <div style={{ textAlign: `left`, padding: `0 3vw` }} className={styles.root}>
    <Table striped bordered hover>
      <thead>
        <th>Name</th>
        <th>Description</th>
        <th>Download</th>
      </thead>
      <tbody>
        {
          items.map(item =>
            <tr>
              <td>
                <span className={styles.code}>{item.fileName}</span>
              </td>
              <td>{item.description}</td>
              <td><DownloadLink item={item}/></td>
            </tr>
          )
        }
      </tbody>
    </Table>
  </div>
}

export default Artifacts 
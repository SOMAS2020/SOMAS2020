import React, { useEffect, useState } from 'react'
import { Table } from 'react-bootstrap'
import styles from './Roles.module.css'

import { getIIGOStatuses } from './Util/getIIGOStatuses'
import { OutputJSONType } from '../../../consts/types'

const IIGOStatus = (props: { output: OutputJSONType }) => {
  const data = getIIGOStatuses(props.output)
  return (
    <div>
      <p className={styles.text}>IIGO Run Status</p>
      <div style={{ overflow: 'scroll', height: '10em' }}>
        <Table striped bordered hover size="sm">
          <thead>
            <tr>
              <th>Turn</th>
              <th>Run Status</th>
            </tr>
          </thead>
          <tbody>
            {data.map((item) => {
              return (
                <tr>
                  <td>{item.turn}</td>
                  <td>{item.status}</td>
                </tr>
              )
            })}
          </tbody>
        </Table>
      </div>
    </div>
  )
}

export default IIGOStatus

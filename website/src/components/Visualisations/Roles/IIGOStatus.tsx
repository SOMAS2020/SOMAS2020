import React, { useEffect, useState } from 'react'
import { Table } from 'react-bootstrap'
import styles from './Roles.module.css'

import { getIIGOStatuses } from './Util/getIIGOStatuses'
import { OutputJSONType } from '../../../consts/types'

const IIGOStatus = (props: { output: OutputJSONType }) => {
  const data = getIIGOStatuses(props.output)
  return (
    <div>
      <h3> IIGO Status </h3>
      <div style={{ overflow: 'scroll', height: '10em' }}>
        <Table striped bordered hover size="sm" responsive>
          {data.map((item) => {
            return (
              <tr>
                <td>{item.turn}</td>
                <td className="text-left">{item.status}</td>
              </tr>
            )
          })}
        </Table>
      </div>
    </div>
  )
}

export default IIGOStatus

import React, { useCallback } from 'react'
import styles from './Transactions.module.css'
import { OutputJSONType } from '../../../consts/types'

import processTransactionData from './Util/ProcessTransactionData'
import ForceGraph from './Util/ForceGraph'

const Transactions = (props: { output: OutputJSONType }) => {
  const nodeHoverTooltip = useCallback((node) => {
    return `<div>${node.id === 0 ? 'Common Pool' : `Team ${node.id}`}</div>`
  }, [])

  const { links, nodes } = processTransactionData(props.output)

  return (
    <div
      className={styles.root}
      style={{
        border: 'black',
        borderWidth: '2px',
        textAlign: 'center',
      }}
    >
      <h2>Transactions Visualisation</h2>
      <p style={{ textAlign: 'left' }}>
        The following plot visualises the transactions between islands and with
        the common pool. The size of a bubble represents the total magnitude of
        resources traded by each island. The width of each connecting edge
        between two entities represents the magnitude of transactions between
        the two islands. Entities that gave more resources than they receive
        have a red border, while islands that received more than they gave have
        a green border. Unlike the IITO Visualisation, this plot includes
        transaction information from Sanctions, Role salary payments, gifting,
        and donations to/requests from the common pool.
      </p>
      <ForceGraph
        linksData={links}
        nodesData={nodes}
        nodeHoverTooltip={nodeHoverTooltip}
      />
    </div>
  )
}

export default Transactions

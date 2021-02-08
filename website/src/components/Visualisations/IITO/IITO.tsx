import React, { useCallback } from 'react'
import styles from './IITO.module.css'
import { OutputJSONType } from '../../../consts/types'

import processTransactionData from './Util/ProcessTradeData'
import ForceGraph from './Util/ForceGraph'

const IITO = (props: { output: OutputJSONType }) => {
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
      <h2>IITO Visualisation</h2>
      <p style={{ textAlign: 'left' }}>
        The following plot visualises the transactions between islands in the
        IITO. The size of a bubble represents the total magnitude of resources
        traded by each island. The width of each connecting edge between two
        islands represents the magnitude of transactions between the two
        islands. Islands that gave more resources than they receive have a red
        border, while islands that received more than they gave have a green
        border.
      </p>
      <ForceGraph
        linksData={links}
        nodesData={nodes}
        nodeHoverTooltip={nodeHoverTooltip}
      />
    </div>
  )
}

export default IITO

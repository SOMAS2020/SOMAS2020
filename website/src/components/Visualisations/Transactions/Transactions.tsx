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
      <ForceGraph
        linksData={links}
        nodesData={nodes}
        nodeHoverTooltip={nodeHoverTooltip}
      />
    </div>
  )
}

export default Transactions

import React from 'react'
import styles from './IITO.module.css'
import { OutputJSONType } from '../../../consts/types'

import processTransactionData from './Util/ProcessTransactionData'
import ForceGraph from './Util/ForceGraph'

const IITO = (props: { output: OutputJSONType }) => {
  const nodeHoverTooltip = React.useCallback((node) => {
    return `<div>${node.id}</div>`
  }, [])

  const { links, nodes } = processTransactionData(props.output)

  return (
    <div className={styles.root} style={{ width: '90%', textAlign: 'center' }}>
      <ForceGraph
        linksData={links}
        nodesData={nodes}
        nodeHoverTooltip={nodeHoverTooltip}
      />
    </div>
  )
}

export default IITO

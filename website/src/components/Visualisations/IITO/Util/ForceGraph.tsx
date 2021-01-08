import React, { useEffect, useRef } from 'react'
import runForceGraph from './ForceGraphGenerator'
import styles from '../IITO.module.css'
import { Link, Node } from '../../../../consts/types'

const ForceGraph = ({
  linksData,
  nodesData,
  nodeHoverTooltip,
}: {
  linksData: Link[]
  nodesData: Node[]
  nodeHoverTooltip: any
}) => {
  const containerRef = useRef(null)

  useEffect(() => {
    let destroyFn

    if (containerRef.current) {
      const { destroy } = runForceGraph(
        containerRef.current,
        linksData,
        nodesData,
        nodeHoverTooltip
      )
      destroyFn = destroy
    }

    return destroyFn
  }, [])

  return <div ref={containerRef} className={styles.container} />
}

export default ForceGraph

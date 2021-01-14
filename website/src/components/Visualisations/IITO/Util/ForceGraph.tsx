import React, { useEffect, useRef } from 'react'
import runForceGraph from './ForceGraphGenerator'
import styles from '../IITO.module.css'

export type Node = {
  id: number
  magnitude: number
  colorStatus: string
  islandColor: string
}

export type Link = {
  source: number
  target: number
  amount: number
}

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

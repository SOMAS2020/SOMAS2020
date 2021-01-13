// @ts-nocheck
// TypeScript does not work nicely at all with d3 so need to come back and fix these
import * as d3 from 'd3'
import styles from '../IITO.module.css'
import { Link, Node } from './ForceGraph'

export default function runForceGraph(
  container: any,
  linksData: Link[],
  nodesData: Node[],
  nodeHoverTooltip: any
) {
  // Assuming we get the links and nodes as expected
  const links = linksData.map((d) => {
    return { ...d }
  })
  const nodes = nodesData.map((d) => {
    return { ...d }
  })

  const containerRect = container.getBoundingClientRect()
  const { height, width } = containerRect

  const borderColor = (d: Node) => {
    return d.colorStatus
  }

  const fillColor = (d: Node) => {
    return d.islandColor
  }

  // size the bubbles by their magnitude
  // TODO: scale the bubble sizes for the visualisation here
  const bubbleSize = (d: Node) => {
    return d.magnitude
  }

  const getClass = (d: Node) => {
    return styles.bubble
  }

  const drag = (simulation) => {
    const dragstarted = (event: any, d: any) => {
      if (!event.active) simulation.alphaTarget(0.3).restart()
      d.fx = d.x
      d.fy = d.y
    }

    const dragged = (event: any, d: any) => {
      d.fx = event.x
      d.fy = event.y
    }

    const dragended = (event: any, d: any) => {
      if (!event.active) simulation.alphaTarget(0)
      d.fx = null
      d.fy = null
    }

    return d3
      .drag()
      .on('start', dragstarted)
      .on('drag', dragged)
      .on('end', dragended)
  }

  // Add the tooltip element to the graph
  const tooltip = document.querySelector('#graph-tooltip')

  if (!tooltip) {
    const tooltipDiv = document.createElement('div')
    tooltipDiv.classList.add(styles.tooltip)
    tooltipDiv.style.opacity = '0'
    tooltipDiv.id = 'graph-tooltip'
    document.body.appendChild(tooltipDiv)
  }

  const div = d3.select('#graph-tooltip')
  const addTooltip = (hoverTooltip, d, x, y) => {
    div.transition().duration(200).style('opacity', 0.9)
    div
      .html(hoverTooltip(d))
      .style('left', `${x}px`)
      .style('top', `${y - 28}px`)
  }

  const removeTooltip = () => {
    div.transition().duration(200).style('opacity', 0)
  }

  const simulation = d3
    .forceSimulation(nodes)
    .force(
      'link',
      d3.forceLink(links).id((d) => d.id)
    )
    .force('charge', d3.forceManyBody().strength(-15000)) // changes the central force
    .force('x', d3.forceX())
    .force('y', d3.forceY())

  const svg = d3
    .select(container)
    .append('svg')
    .attr('viewBox', [-width / 2, -height / 2, width, height])

  const borderPath = svg
    .append('rect')
    .attr('x', -width / 2)
    .attr('y', -height / 2)
    .attr('height', height)
    .attr('width', width)
    .style('stroke', '#B1B1B1')
    .style('fill', 'none')
    .style('stroke-width', 5)

  const link = svg
    .append('g')
    .attr('stroke', '#999')
    .attr('stroke-opacity', 0.6)
    .selectAll('line')
    .data(links)
    .join('line')
    .attr('stroke-width', (d: Link) => Math.sqrt(d.amount))

  const node = svg
    .append('g')
    .selectAll('circle')
    .data(nodes)
    .join('circle')
    .attr('r', bubbleSize)
    .attr('stroke', borderColor)
    .attr('stroke-width', 5)
    .attr('fill', fillColor)
    .call(drag(simulation))

  const label = svg
    .append('g')
    .attr('class', 'labels')
    .selectAll('text')
    .data(nodes)
    .enter()
    .append('text')
    .attr('text-anchor', 'middle')
    .attr('dominant-baseline', 'central')
    .attr('class', (d) => `fa ${getClass(d)}`)
    .style('fill', (d: Node) => (d.id === 0 ? 'white' : 'black'))
    .text((d: Node) => (d.id === 0 ? 'Common Pool' : `${d.id}`))
    .call(drag(simulation))

  label
    .on('mouseover', (event, d) => {
      addTooltip(nodeHoverTooltip, d, event.pageX, event.pageY)
    })

    .on('mouseout', () => {
      removeTooltip()
    })

  simulation.on('tick', () => {
    // update link positions
    link
      .attr('x1', (d: any) => d.source.x)
      .attr('y1', (d: any) => d.source.y)
      .attr('x2', (d: any) => d.target.x)
      .attr('y2', (d: any) => d.target.y)

    node.attr('cx', (d: any) => d.x).attr('cy', (d: any) => d.y)

    borderPath
      .attr('x', -width / 2)
      .attr('y', -height / 2)
      .attr('height', height)
      .attr('width', width)

    label.attr('x', (d: any) => d.x).attr('y', (d: any) => d.y)
  })

  return {
    destroy: () => {
      simulation.stop()
    },
    nodes: () => {
      return svg.node()
    },
  }
}

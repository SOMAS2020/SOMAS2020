import React, { useState, useEffect } from 'react'
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
} from 'recharts'
import { ForagingHistory } from './ForagingTypes'

// TODO: Add axis labels
// plot: deer population, fish population, foraging return, foraging input
const ForagingPlot = ({ data }: { data: ForagingHistory }) => {
  const [result, setResult] = useState<ForagingHistory>([])

  useEffect(() => {
    console.log('charts props:', data)
    setResult(data)
    console.log('results updated')
  })

  console.log({ data })
  return (
    <LineChart
      width={900}
      height={400}
      data={data}
      margin={{
        top: 5,
        right: 30,
        left: 20,
        bottom: 5,
      }}
    >
      <CartesianGrid strokeDasharray="3 3" />

      {/* Name of the x axis */}
      <XAxis dataKey="turn" />

      {/* name of the value that corresponds to left axis */}
      <YAxis yAxisId="left" />

      {/* name of the value that corresponds to right axis */}
      <YAxis yAxisId="right" orientation="right" />
      <Tooltip />
      <Legend />
      <Line
        yAxisId="right"
        type="monotone"
        dataKey="deerNumParticipants"
        stroke="#F1538A"
        activeDot={{ r: 8 }}
      />
      <Line
        yAxisId="right"
        type="monotone"
        dataKey="fishNumParticipants"
        stroke="#F1238F"
        activeDot={{ r: 8 }}
      />
      <Line
        yAxisId="left"
        type="monotone"
        dataKey="deerInputResources"
        stroke="#8884d8"
        activeDot={{ r: 8 }}
      />
      <Line
        yAxisId="left"
        type="monotone"
        dataKey="fishInputResources"
        stroke="#82ca9d"
        activeDot={{ r: 8 }}
      />
    </LineChart>
  )
}

export default ForagingPlot

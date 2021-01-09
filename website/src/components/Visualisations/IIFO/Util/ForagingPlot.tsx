import React, { PureComponent } from 'react'
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
} from 'recharts'

// TODO: X axis turns
// TODO: Y axis resources
// plot: deer population, fish population, foraging return, foraging input
const ForagingPlot = (data) => {
  const foragingHist = [
    {
      turn: 1,
      deerCount: 10,
      fishCount: 25,
      foragingReturn: 10,
      foragingInput: 4,
    },
    {
      turn: 2,
      deerCount: 20,
      fishCount: 50,
      foragingReturn: 25,
      foragingInput: 3,
    },
    {
      turn: 3,
      deerCount: 15,
      fishCount: 35,
      foragingReturn: 20,
      foragingInput: 5,
    },
  ]

  const testData = [
    {
      name: 'Page A',
      uv: 4000,
      pv: 2400,
      amt: 2400,
    },
    {
      name: 'Page B',
      uv: 3000,
      pv: 1398,
      amt: 2210,
    },
    {
      name: 'Page C',
      uv: 2000,
      pv: 9800,
      amt: 2290,
    },
    {
      name: 'Page D',
      uv: 2780,
      pv: 3908,
      amt: 2000,
    },
    {
      name: 'Page E',
      uv: 1890,
      pv: 4800,
      amt: 2181,
    },
    {
      name: 'Page F',
      uv: 2390,
      pv: 3800,
      amt: 2500,
    },
    {
      name: 'Page G',
      uv: 3490,
      pv: 4300,
      amt: 2100,
    },
  ]
  return (
    <LineChart
      width={900}
      height={400}
      data={testData}
      margin={{
        top: 5,
        right: 30,
        left: 20,
        bottom: 5,
      }}
    >
      <CartesianGrid strokeDasharray="3 3" />

      {/* Name of the x axis */}
      <XAxis dataKey="name" />

      {/* name of the value that corresponds to left axis */}
      <YAxis yAxisId="left" />

      {/* name of the value that corresponds to right axis */}
      <YAxis yAxisId="right" orientation="right" />
      <Tooltip />
      <Legend />
      <Line
        yAxisId="left"
        type="monotone"
        dataKey="pv"
        stroke="#8884d8"
        activeDot={{ r: 8 }}
      />
      <Line
        yAxisId="right"
        type="monotone"
        dataKey="uv"
        stroke="#82ca9d"
        activeDot={{ r: 8 }}
      />
    </LineChart>
  )
}

export default ForagingPlot

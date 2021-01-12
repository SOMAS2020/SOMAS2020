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
import { ForagingTurn } from './ForagingTypes'

const ForagingPlot = ({ data }: { data: ForagingTurn[] }) => {
  const [result, setResult] = useState<ForagingTurn[]>([])

  useEffect(() => {
    console.log('charts props:', data)
    setResult(data)
    console.log('results updated')
  })

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

      <XAxis
        dataKey="turn"
        label={{
          value: 'Turn',
          position: 'bottom',
          dy: 0,
        }}
      />

      <YAxis yAxisId="left" />

      <YAxis
        yAxisId="right"
        orientation="right"
        label={{ value: 'Number Participants', angle: 90 }}
      />
      <Tooltip />
      <Legend
        verticalAlign="bottom"
        align="center"
        height={20}
        wrapperStyle={{
          paddingTop: '15px',
        }}
      />
      <Line
        yAxisId="left"
        type="monotone"
        dataKey="deerInputResources"
        stroke="#0095FF"
        activeDot={{ r: 8 }}
      />
      <Line
        yAxisId="right"
        type="monotone"
        dataKey="deerNumCaught"
        stroke="#ff0000"
        activeDot={{ r: 8 }}
      />
      <Line
        yAxisId="left"
        type="monotone"
        dataKey="deerTotalUtility"
        stroke="#802FF0"
        activeDot={{ r: 8 }}
      />
      <Line
        yAxisId="right"
        type="monotone"
        dataKey="fishNumCaught"
        stroke="#00c49f"
        activeDot={{ r: 8 }}
      />
      <Line
        yAxisId="left"
        type="monotone"
        dataKey="fishInputResources"
        stroke="#ffbb29"
        activeDot={{ r: 8 }}
      />
      <Line
        yAxisId="left"
        type="monotone"
        dataKey="fishTotalUtility"
        stroke="#ff8042"
        activeDot={{ r: 8 }}
      />
    </LineChart>
  )
}

export default ForagingPlot

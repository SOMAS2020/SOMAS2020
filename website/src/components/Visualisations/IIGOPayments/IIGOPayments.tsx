import React, { useEffect, useState } from 'react'
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts'
import styles from './IIGOPayments.module.css'

import { ProcessedTaxData } from './Util/IIGOPaymentsTypes'
import { processPaymentsData } from './Util/ProcessedPaymentsData'
import { OutputJSONType, TeamName } from '../../../consts/types'

const Payments = (props: { output: OutputJSONType }) => {
  const data = processPaymentsData(props.output)

  const legend = new Map([
    ['actualTax', '#094fdb'],
    ['expectedTax', '#507cd4'],
    ['actualAlloc', '#cf1763'],
    ['expectedAlloc', '#c76f94'],
  ])

  return (
    <div className={styles.root}>
      <p className={styles.text}>IIGO Payments</p>
      <ResponsiveContainer height={460} width="100%">
        <BarChart data={data} layout="horizontal">
          <XAxis type="category" dataKey="name" />
          <YAxis
            type="number"
            domain={['dataMin', 'dataMax']}
            tickCount={20}
            allowDecimals={false}
          />
          <Tooltip />
          <Legend verticalAlign="top" />
          <Bar dataKey="expectedTax" fill="#094fbd" />
          <Bar dataKey="actualTax" fill="#507cd4" />
          <Bar dataKey="expectedAlloc" fill="#cf1763" />
          <Bar dataKey="actualAlloc" fill="#c76f94" />
          <Bar dataKey="expectedSanction" fill="#EBA421" />
          <Bar dataKey="actualSanction" fill="#e6c891" />
        </BarChart>
      </ResponsiveContainer>
      <p className={styles.graphLabel}>Team</p>
    </div>
  )
}

export default Payments

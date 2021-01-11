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
    ['Actual Tax', '#094fdb'],
    ['Expected Tax', '#507cd4'],
    ['Actual Allocation', '#cf1763'],
    ['expected Allocation', '#c76f94'],
  ])

  return (
    <div className={styles.root}>
      <p className={styles.text}>IIGO Payments</p>
      <ResponsiveContainer height={460} width="100%">
        <BarChart data={data} layout="horizontal">
          <XAxis type="category" dataKey="name" />
          <YAxis type="number" tickCount={20} allowDecimals={false} />
          <Tooltip />
          <Legend
            verticalAlign="top"
            wrapperStyle={{
              paddingLeft: '10px',
            }}
          />
          <Bar dataKey="expectedTax" fill="#094fbd" name="Expected Tax" />
          <Bar dataKey="actualTax" fill="#507cd4" name="Actual Tax Paid" />
          <Bar
            dataKey="expectedAlloc"
            fill="#cf1763"
            name="Allocation Given by President"
          />
          <Bar
            dataKey="actualAlloc"
            fill="#c76f94"
            name="Actual Allocation Taken"
          />
          <Bar
            dataKey="expectedSanction"
            fill="#EBA421"
            name="Expected Sanction Charged"
          />
          <Bar
            dataKey="actualSanction"
            fill="#e6c891"
            name="Actual Sanction Paid"
          />
        </BarChart>
      </ResponsiveContainer>
      <p className={styles.graphLabel}>Team</p>
    </div>
  )
}

export default Payments

/* eslint-disable react/no-access-state-in-setstate */
/* eslint-disable jsx-a11y/no-static-element-interactions */
/* eslint-disable jsx-a11y/click-events-have-key-events */
import React from 'react'
import {
  LegendProps,
  Brush,
  ResponsiveContainer,
  LineChart,
  Line,
  YAxis,
  XAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  Surface,
  Symbols,
  ReferenceLine,
  ReferenceArea,
  TooltipProps,
} from 'recharts'
import { OutputJSONType } from '../../../consts/types'
import styles from './Resources.module.css'
import { getSeasonEnds, outputToResourceLevels, ResourceLevel } from './utils'

const CustomTooltip = ({ active, label, payload }: TooltipProps) => {
  return (
    active &&
    label &&
    payload && (
      <div className={styles.customTooltip}>
        <p className={styles.label}>{label}</p>
        {payload.map((pl) => (
          <p
            style={{ color: pl.color }}
            className={styles.content}
            key={`${pl.name}${pl.value}`}
          >
            {pl.name}: {(pl.value as number).toFixed(1)}
          </p>
        ))}
      </div>
    )
  )
}

interface IProps {
  output: OutputJSONType
}

interface IState {
  disabled: string[]
}

const legendColours = {
  team1: '#0095FF',
  team2: '#FF0000',
  team3: '#802FF0',
  team4: '#00C49F',
  team5: '#FFBB28',
  team6: '#FF8042',
  CommonPool: '#ACE600',
  TotalResources: '	#FF69B4',
  CriticalThreshold: '#B7B4B0',
}
class ResourceGraph extends React.Component<IProps, IState> {
  private chartData: ResourceLevel[]

  constructor(props) {
    super(props)
    const { output } = this.props
    this.state = {
      disabled: [],
    }
    this.chartData = outputToResourceLevels(output)
  }

  handleClick = (dataKey: string) => {
    this.setState({
      disabled: this.state.disabled.includes(dataKey)
        ? this.state.disabled.filter((obj: string) => obj !== dataKey)
        : this.state.disabled.concat(dataKey),
    })
  }

  renderCustomizedLegend = ({ payload }: LegendProps) => {
    return (
      <div className="customized-legend">
        {payload?.map((entry) => {
          const { value, color } = entry
          const active = this.state.disabled.includes(value)
          const style = {
            marginRight: 10,
            colour: active ? '#AAA' : '#000',
          }

          return (
            <span
              className="legend-item"
              onClick={() => this.handleClick(value)}
              style={style}
            >
              <Surface
                width={10}
                height={10}
                viewBox={{ x: 0, y: 0, width: 10, height: 10 }}
              >
                <Symbols cx={5} cy={5} type="circle" size={50} fill={color} />
                {active && (
                  <Symbols cx={5} cy={5} type="circle" size={25} fill="#FFF" />
                )}
              </Surface>
              <span>{value}</span>
            </span>
          )
        })}
      </div>
    )
  }

  render() {
    const { disabled } = this.state
    const { output } = this.props
    return (
      <ResponsiveContainer height={330} width="100%">
        <LineChart
          data={this.chartData}
          margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
        >
          {Object.entries(legendColours)
            .filter(([id]) => !disabled.includes(id))
            .map(([id, colour]) => (
              <Line name={id} type="monotone" dataKey={id} stroke={colour} />
            ))}
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis
            dataKey="Turn"
            height={50}
            label={{
              value: 'Turn',
              angle: 0,
              position: 'insideBottom',
              dy: -8,
            }}
          />
          <YAxis
            label={{ value: 'Resources', angle: -90, position: 'insideLeft' }}
          />
          <Tooltip content={CustomTooltip} />
          {getSeasonEnds(output).map((seasonEnd) => (
            <ReferenceLine
              x={seasonEnd}
              label="Season End"
              stroke="black"
              strokeDasharray="3 3"
            />
          ))}
          {!disabled.includes('CriticalThreshold') && (
            <ReferenceArea
              y1={0}
              y2={output.Config.MinimumResourceThreshold}
              label="CriticalThreshold"
              stroke="#e6eeff"
              strokeOpacity={0.1}
            />
          )}
          <Legend
            verticalAlign="top"
            align="center"
            height={20}
            wrapperStyle={{ top: 0, left: 25, right: 0, width: 'auto' }}
            payload={Object.entries(legendColours).map(([id, colour]) => ({
              value: id,
              color: colour,
              type: 'circle',
              id: `${id}${colour}`,
            }))}
            content={this.renderCustomizedLegend}
          />
          <Brush dataKey="Turn" height={25} stroke="#2fa1c6" />
        </LineChart>
      </ResponsiveContainer>
    )
  }
}

export default ResourceGraph

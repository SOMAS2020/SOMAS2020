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
import _ from 'lodash'
import { OutputJSONType } from '../../../consts/types'
import styles from './Resources.module.css'

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

class LineRechartComponent extends React.Component<IProps, any> {
  constructor(props: IProps | Readonly<IProps>) {
    super(props)
    this.state = {
      disabled: [],
      legendColours: {
        team1: '#0095FF',
        team2: '#FF0000',
        team3: '#802FF0',
        team4: '#00C49F',
        team5: '#FFBB28',
        team6: '#FF8042',
        CommonPool: '#ACE600',
        CriticalThreshold: '#B7B4B0',
      },
      datapaths: {
        team1: 'ClientInfos.Team1.Resources',
        team2: 'ClientInfos.Team2.Resources',
        team3: 'ClientInfos.Team3.Resources',
        team4: 'ClientInfos.Team4.Resources',
        team5: 'ClientInfos.Team5.Resources',
        team6: 'ClientInfos.Team6.Resources',
        CommonPool: 'CommonPool',
      },
      chartData: this.props.output,
    }
  }

  handleClick = (dataKey: string) => {
    this.setState({
      disabled: this.state.disabled.includes(dataKey)
        ? this.state.disabled.filter((obj: string) => obj !== dataKey)
        : this.state.disabled.concat(dataKey),
    })
  }

  getSeasonEnds() {
    const seasonEnds: number[] = []
    let i
    // eslint-disable-next-line no-restricted-syntax
    for (i in this.state.chartData.GameStates) {
      if (
        this.state.chartData.GameStates[i].Environment.LastDisasterReport
          .Magnitude !== 0
      ) {
        seasonEnds.push(this.state.chartData.GameStates[i].Turn - 1)
      }
    }
    return seasonEnds
  }

  renderCustomizedLegend = ({ payload }: LegendProps) => {
    return (
      <div className="customized-legend">
        {payload?.map((entry) => {
          const { value, color } = entry
          const { disabled } = this.state
          const active = disabled.includes(value)
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
    return (
      <ResponsiveContainer height={330} width="100%">
        <LineChart
          data={this.state.chartData.GameStates}
          margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
        >
          {_.toPairs<string>(this.state.legendColours)
            .filter((pair) => !_.includes(this.state.disabled, pair[0]))
            .map((pair) => (
              <Line
                name={pair[0]}
                type="monotone"
                dataKey={this.state.datapaths[pair[0]]}
                stroke={pair[1]}
              />
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
          {this.getSeasonEnds().map((seasonEnd) => (
            <ReferenceLine
              x={seasonEnd}
              label="Season End"
              stroke="black"
              strokeDasharray="3 3"
            />
          ))}
          {_.toPairs<string>({ CriticalThreshold: '#e6eeff' })
            .filter((refArea) => !_.includes(this.state.disabled, refArea[0]))
            .map((refArea) => (
              <ReferenceArea
                y1={0}
                y2={this.state.chartData.Config.MinimumResourceThreshold}
                label={refArea[0]}
                stroke={refArea[1]}
                strokeOpacity={0.1}
              />
            ))}
          <Legend
            verticalAlign="top"
            align="center"
            height={20}
            wrapperStyle={{ top: 0, left: 25, right: 0, width: 'auto' }}
            payload={_.toPairs<string>(this.state.legendColours).map(
              (pair) => ({
                value: pair[0],
                color: pair[1],
                type: 'circle',
                id: `${pair[0]}${pair[1]}`,
              })
            )}
            content={this.renderCustomizedLegend}
          />
          <Brush dataKey="Turn" height={25} stroke="#2fa1c6" />
        </LineChart>
      </ResponsiveContainer>
    )
  }
}

export default LineRechartComponent

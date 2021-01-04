import React from "react";
import { LegendProps, Brush, ResponsiveContainer, LineChart, Line, YAxis, XAxis, CartesianGrid, Tooltip, Legend, Surface, Symbols, ReferenceLine } from 'recharts';
import outputJSON from '../../../output/output.json';
//import { OutputJSONType } from "../../../consts/types";
import _ from "lodash";

class LineRechartComponent extends React.Component<{}, any> {

    constructor(props) {
        super(props);
        this.state = {
            disabled: [],
            lineColours: {
                team1: "#0095FF",
                team2: "#FF0000",
                team3: "#802FF0",
                team4: "#00C49F",
                team5: "#FFBB28",
                team6: "#FF8042",
            },
            datapaths: {
                team1: "ClientInfos.Team1.Resources",
                team2: "ClientInfos.Team2.Resources",
                team3: "ClientInfos.Team3.Resources",
                team4: "ClientInfos.Team4.Resources",
                team5: "ClientInfos.Team5.Resources",
                team6: "ClientInfos.Team6.Resources",
            },
            chartData: outputJSON
        };
    }

    handleClick = dataKey => {
        this.setState({
            disabled: this.state.disabled.includes(dataKey)
              ? this.state.disabled.filter((obj) => obj !== dataKey)
              : this.state.disabled.concat(dataKey),
          });
      
    };

    renderCustomizedLegend = ({ payload }: LegendProps) => {
        return (
            <div className="customized-legend">
                {payload?.map(entry => {
                    const { value, color } = entry;
                    const { disabled } = this.state;
                    const active = disabled.includes(value);
                    const style = {
                        marginRight: 10,
                        colour: active ? "#AAA" : "#000"
                    };

                    return (
                        <span
                            className="legend-item"
                            onClick={() => this.handleClick(value)}
                            style={style}
                        >
                            <Surface width={10} height={10} viewBox={{x:0, y:0, width:10, height:10}} >
                                <Symbols cx={5} cy={5} type="circle" size={50} fill={color} />
                                {active && (
                                    <Symbols
                                        cx={5}
                                        cy={5}
                                        type="circle"
                                        size={25}
                                        fill={"#FFF"}
                                    />
                                )}
                            </Surface>
                            <span>{value}</span>
                        </span>
                    )
                })}
            </div>
        );
    };

    getSeasonEnds() {
        var seasonEnds: any[] = [];
        let i;
        for (i in this.state.chartData.GameStates) {
            if (this.state.chartData.GameStates[i].Environment.LastDisasterReport.Magnitude !== 0) {
                seasonEnds.push(this.state.chartData.GameStates[i].Turn - 1)
            }
        }
        return seasonEnds;
    };

    render() {
        return (
            <ResponsiveContainer height={330} width="95%">
                <LineChart data={this.state.chartData.GameStates}
                    margin={{ top: 5, right: 30, left: 20, bottom: 5 }}>
                    {_.toPairs<string>(this.state.lineColours)
                        .filter(pair => !_.includes(this.state.disabled, pair[0]))
                        .map(pair => (
                            <Line
                                name={pair[0]}
                                type="monotone"
                                dataKey={this.state.datapaths[pair[0]]}
                                stroke={pair[1]}
                            />
                        ))}
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="Turn" height={50} label={{ value: 'Turn', angle: 0, position: 'insideBottom', dy: -8 }} />
                    <YAxis label={{ value: 'Resources', angle: -90, position: 'insideLeft' }} />
                    <Tooltip />
                    {this.getSeasonEnds().map((seasonEnd) =>
                    (
                        <ReferenceLine
                            x={seasonEnd}
                            label="Season End"
                            stroke="red"
                            strokeDasharray="3 3"
                        />
                    ))}
                    <Legend
                        verticalAlign="top"
                        align="center"
                        height={20}
                        wrapperStyle={{ top: 0, left: 25, right: 0, width: 'auto' }}
                        payload={_.toPairs<string>(this.state.lineColours).map(pair => ({
                            value: pair[0],
                            color: pair[1], 
                            type: "circle",
                            id: `${pair[0]}${pair[1]}`
                        }))}

                        content={this.renderCustomizedLegend}
                    />
                    <Brush dataKey='Turn' height={25} stroke="#2fa1c6" />
                </LineChart>
            </ResponsiveContainer >
        )
    };
}

export default LineRechartComponent;

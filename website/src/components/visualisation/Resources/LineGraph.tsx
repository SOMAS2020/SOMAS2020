import React from "react";
import { Brush, ResponsiveContainer, LineChart, Line, YAxis, XAxis, CartesianGrid, Tooltip, Legend, Surface, Symbols, ReferenceLine } from 'recharts';
import outputJSON from '../../../output/output.json';
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
        const { disabled } = this.state;
        if (disabled.includes(dataKey)) {
            this.setState({
                disabled: this.state.disabled.filter(obj => obj !== dataKey)
            });
        }
        else {
            this.setState({
                disabled: this.state.disabled.concat(dataKey)
            });
        }
    };

    renderCustomizedLegend = ({ payload }) => {
        return (
            <div className="customized-legend">
                {payload.map(entry => {
                    const { dataKey, colour } = entry;
                    const { disabled } = this.state;
                    const active = disabled.includes(dataKey);
                    const style = {
                        marginRight: 10,
                        colour: active ? "#AAA" : "#000"
                    };

                    return (
                        <span
                            className="legend-item"
                            onClick={() => this.handleClick(dataKey)}
                            style={style}
                        >
                            <Surface width={10} height={10} viewBox="0 0 10 10" >
                                <Symbols cx={5} cy={5} type="circle" size={50} fill={colour} />
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
                            <span>{dataKey}</span>
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
                    {_.toPairs(this.state.lineColours)
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
                        payload={_.toPairs(this.state.lineColours).map(pair => ({
                            dataKey: pair[0],
                            colour: pair[1],
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
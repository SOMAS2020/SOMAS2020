import React from "react";
import styles from "./Roles.module.css";
import {
    BarChart,
    Bar,
    XAxis,
    YAxis,
    Tooltip,
    Legend,
    ResponsiveContainer,
} from "recharts";
import { ProcessedRoleData } from "./Util/RoleTypes";
import { processRoleData } from "./Util/ProcessedRoleData";
import { OutputJSONType } from "../../../consts/types";

// type CustomTooltipProps = {
//     active: boolean;
//     payload: [{ name: string; value: number; unit: string }];
//     label: string;
//     data: ProcessedRoleData;
// };

// const CustomTooltip = ({
//     active,
//     label,
//     data,
// }: CustomTooltipProps) => {
//     const getTurnsInRoles = (name: string): TurnsInRoles => (
//         data.filter((val) => (val.name === name)).reduce((acc, curr) => {
//             curr.turnsInRoles.map((tInR) => acc.incrementRoles(tInR.roles, tInR.turns));
//             return acc;
//         }, new TurnsInRoles())
//     );

//     if (active && data.length > 0) {
//         const turnsInRoles = getTurnsInRoles(label);
//         const turns = data[0].turnsInRoles.reduce((acc, curr) => acc += curr.turns, 0);

//         return (
//             <div className={styles.customTooltip}>
//                 <p className={styles.label}>{label}</p>
//                 {turnsInRoles.toPairs().map(([role, turnsInRole]) => (
//                     <p className={styles.content} key={role}>
//                         Turns as {role}: {turnsInRole} (
//                         {((turnsInRole * 100) / turns).toFixed(1)}%)
//                     </p>
//                 ))}
//             </div>
//         );
//     }

//     return null;
// };

const getRandomColor = () => {
    var letters = "0123456789ABCDEF";
    var color = "#";
    for (var i = 0; i < 6; i++) {
        color += letters[Math.floor(Math.random() * 16)];
    }
    return color;
}

const getNewColors = (teams: string[]): Map<string, string> => {
    let colorMap = new Map<string, string>();

    // TODO: initialise map with team colors (only use random color for unexpected team)

    teams.map((team) => {
        if (colorMap.has(team)) {
            return team;
        } else {
            colorMap.set(team, getRandomColor());
            return team;
        }
    });

    return colorMap;
}

const Roles = (props: { output: OutputJSONType }) => {

    const data = processRoleData(props.output);
    const teams = [ "Team1", "Team2", "Team3", "Team4", "Team5", "Team6" ];
    const colors = getNewColors(teams);

    console.log(data);

    return (
        <div className={styles.root}>
            <p className={styles.text}>Role Visualisation</p>
            <ResponsiveContainer height={460} width="100%">
                <BarChart data={data} layout="vertical">
                    <YAxis type="category" dataKey="role" />
                    <XAxis type="number" domain={[0, "dataMax"]} />
                    {/* <Tooltip
                        content={(props: CustomTooltipProps) =>
                            CustomTooltip({ ...props, data })
                        }
                    /> */}
                    <Legend
                        verticalAlign="top"
                        payload={teams.map((team, i) => ({
                            value: team,
                            type: "square",
                            id: `${team}${i}`,
                            color: colors.get(team)
                        }))}
                    />
                    {data[0].occupied.map((a, i) => [
                        teams.map((team) => (
                            <Bar
                                dataKey={`occupied[${i}].${team}`}
                                stackId="a"
                                fill={colors.get(team)}
                                key={`${i}${team}`}
                            />
                        ))
                    ])}
                </BarChart>
            </ResponsiveContainer>
            <p className={styles.graphLabel}>Turns</p>
        </div>
    );
};

export default Roles;

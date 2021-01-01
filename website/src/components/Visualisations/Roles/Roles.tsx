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
import { TurnsInRoles, ProcessedRoleData } from "./Util/RoleTypes";
import { getProcessedRoleData } from "./Util/ProcessedRoleData";

const presidentColor = "#00bbf9";
const judgeColor = "#fee440";
const speakerColor = "#f15bb5";
const noneColor = "#b2bec3";

type CustomTooltipProps = {
    active: boolean;
    payload: [{ name: string; value: number; unit: string }];
    label: string;
    data: ProcessedRoleData;
};

const getTurnsInRoles = (
    data: ProcessedRoleData,
    name: string
): TurnsInRoles => {
    const roles = data.find((d) => d.name === name)?.roles;
    if (roles) {
        const turnsAsPresident = roles.reduce((acc, a) => acc + a.president, 0);
        const turnsAsJudge = roles.reduce((acc, a) => acc + a.judge, 0);
        const turnsAsSpeaker = roles.reduce((acc, a) => acc + a.speaker, 0);
        const turnsAsNone = roles.reduce((acc, a) => acc + a.none, 0);
        return new TurnsInRoles(
            turnsAsPresident,
            turnsAsJudge,
            turnsAsSpeaker,
            turnsAsNone
        );
    } else {
        console.log(`[VIS ROLE] Could not find ${name} in data...`);
        return new TurnsInRoles();
    }
};

const CustomTooltip = (
    { active, payload, label, data }: CustomTooltipProps
) => {
    if (active) {
        const turnsInRoles = getTurnsInRoles(data, label);
        const turns =
            turnsInRoles.president +
            turnsInRoles.judge +
            turnsInRoles.speaker +
            turnsInRoles.none;

        return (
            <div className={styles.customTooltip}>
                <p className={styles.label}>{label}</p>
                <p className={styles.content}>
                    Turns as President: {turnsInRoles.president} (
                    {((turnsInRoles.president * 100) / turns).toFixed(1)}%)
                </p>
                <p className={styles.content}>
                    Turns as Judge: {turnsInRoles.judge} (
                    {((turnsInRoles.judge * 100) / turns).toFixed(1)}%)
                </p>
                <p className={styles.content}>
                    Turns as Speaker: {turnsInRoles.speaker} (
                    {((turnsInRoles.speaker * 100) / turns).toFixed(1)}%)
                </p>
                <p className={styles.content}>
                    Turns without power: {turnsInRoles.none} (
                    {((turnsInRoles.none * 100) / turns).toFixed(1)}%)
                </p>
            </div>
        );
    }

    return null;
};

const Roles = () => {
    const data = getProcessedRoleData();

    return (
        <div className={styles.root}>
            <ResponsiveContainer height={460} width="100%">
                <BarChart data={data} layout="vertical" margin={{ bottom: 30 }}>
                    <YAxis type="category" dataKey="name" />
                    <XAxis type="number" />
                    <Tooltip
                        content={(props: CustomTooltipProps) =>
                            CustomTooltip({...props, data})
                        }
                    />
                    <Legend
                        payload={[
                            {
                                value: "President",
                                type: "square",
                                id: "ID01",
                                color: presidentColor,
                            },
                            {
                                value: "Judge",
                                type: "square",
                                id: "ID02",
                                color: judgeColor,
                            },
                            {
                                value: "Speaker",
                                type: "square",
                                id: "ID03",
                                color: speakerColor,
                            },
                            {
                                value: "None",
                                type: "square",
                                id: "ID04",
                                color: noneColor,
                            },
                        ]}
                    />
                    {data[0].roles.map((_, i) => [
                        <Bar
                            dataKey={`roles[${i}].president`}
                            stackId="a"
                            fill={presidentColor}
                        />,
                        <Bar
                            dataKey={`roles[${i}].judge`}
                            stackId="a"
                            fill={judgeColor}
                        />,
                        <Bar
                            dataKey={`roles[${i}].speaker`}
                            stackId="a"
                            fill={speakerColor}
                        />,
                        <Bar
                            dataKey={`roles[${i}].none`}
                            stackId="a"
                            fill={noneColor}
                        />,
                    ])}
                </BarChart>
            </ResponsiveContainer>
            <p className={styles.text}>Role Visualisation</p>
        </div>
    );
};

export default Roles;

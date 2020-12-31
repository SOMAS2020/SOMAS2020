import React from "react";
import styles from "./Roles.module.css";
import { BarChart, Bar, XAxis, YAxis, Tooltip, Legend } from "recharts";
import { TurnsInRoles, ProcessedRoleData } from "./Util/RoleTypes";
import { getProcessedRoleData } from "./Util/ProcessedRoleData";

const presidentColor = "#00bbf9";
const judgeColor = "#fee440";
const speakerColor = "#f15bb5";
const noneColor = "#b2bec3";

const data = getProcessedRoleData();

// const data: ProcessedRoleData = [
//     {
//         name: "Team1",
//         roles: [
//             { president: 5, judge: 0, speaker: 0, none: 0 },
//             { president: 0, judge: 0, speaker: 0, none: 2 },
//             { president: 0, judge: 1, speaker: 0, none: 0 },
//             { president: 0, judge: 0, speaker: 0, none: 0 },
//         ],
//     },
//     {
//         name: "Team2",
//         roles: [
//             { president: 0, judge: 3, speaker: 0, none: 0 },
//             { president: 0, judge: 0, speaker: 0, none: 3 },
//             { president: 0, judge: 0, speaker: 1, none: 0 },
//             { president: 0, judge: 0, speaker: 0, none: 1 },
//         ],
//     },
//     {
//         name: "Team3",
//         roles: [
//             { president: 0, judge: 0, speaker: 5, none: 0 },
//             { president: 0, judge: 0, speaker: 0, none: 3 },
//             { president: 0, judge: 0, speaker: 0, none: 0 },
//             { president: 0, judge: 0, speaker: 0, none: 0 },
//         ],
//     },
//     {
//         name: "Team4",
//         roles: [
//             { president: 0, judge: 0, speaker: 0, none: 3 },
//             { president: 0, judge: 2, speaker: 0, none: 0 },
//             { president: 0, judge: 0, speaker: 0, none: 3 },
//             { president: 0, judge: 0, speaker: 0, none: 0 },
//         ],
//     },
//     {
//         name: "Team5",
//         roles: [
//             { president: 0, judge: 0, speaker: 0, none: 5 },
//             { president: 1, judge: 0, speaker: 0, none: 0 },
//             { president: 0, judge: 1, speaker: 0, none: 0 },
//             { president: 0, judge: 0, speaker: 1, none: 0 },
//         ],
//     },
//     {
//         name: "Team6",
//         roles: [
//             { president: 0, judge: 0, speaker: 0, none: 5 },
//             { president: 0, judge: 0, speaker: 1, none: 0 },
//             { president: 0, judge: 0, speaker: 0, none: 2 },
//             { president: 0, judge: 0, speaker: 0, none: 0 },
//         ],
//     },
// ];

type CustomTooltipProps = {
    active: boolean;
    payload: [{ name: string; value: number; unit: string }];
    label: string;
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
        return {
            president: turnsAsPresident,
            judge: turnsAsJudge,
            speaker: turnsAsSpeaker,
            none: turnsAsNone,
        };
    } else {
        console.log(`[VisRole] Could not find ${name} in data...`);
        return { president: 0, judge: 0, speaker: 0, none: 0 };
    }
};

const CustomTooltip = ({ active, payload, label }: CustomTooltipProps) => {
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
                    {(turnsInRoles.president * 100) / turns}%)
                </p>
                <p className={styles.content}>
                    Turns as Judge: {turnsInRoles.judge} (
                    {(turnsInRoles.judge * 100) / turns}%)
                </p>
                <p className={styles.content}>
                    Turns as Speaker: {turnsInRoles.speaker} (
                    {(turnsInRoles.speaker * 100) / turns}%)
                </p>
                <p className={styles.content}>
                    Turns without power: {turnsInRoles.none} (
                    {(turnsInRoles.none * 100) / turns}%)
                </p>
            </div>
        );
    }

    return null;
};

const Roles = () => {
    return (
        <BarChart
            width={500}
            height={300}
            data={data}
            margin={{
                top: 20,
                right: 30,
                left: 20,
                bottom: 5,
            }}
            layout="vertical"
        >
            <YAxis type="category" dataKey="name" />
            <XAxis type="number" />
            <Tooltip content={CustomTooltip} />
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
    );
};

export default Roles;

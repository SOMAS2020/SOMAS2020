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
import { processRoleData } from "./Util/ProcessedRoleData";
import { OutputJSONType } from "../../../consts/types";
import _ from "lodash";

type CustomTooltipProps = {
    active: boolean;
    payload: [{ name: string; value: number; unit: string }];
    label: string;
    data: ProcessedRoleData;
};

const CustomTooltip = ({
    active,
    payload,
    label,
    data,
}: CustomTooltipProps) => {
    const getTurnsInRoles = (name: string): TurnsInRoles => {
        const roles = data.find((d) => d.name === name)?.roles;
        if (roles) {
            const turnsAsPresident = roles.reduce(
                (acc, a) => acc + a.president,
                0
            );
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

    if (active) {
        const turnsInRoles = getTurnsInRoles(label);
        const turns =
            turnsInRoles.president +
            turnsInRoles.judge +
            turnsInRoles.speaker +
            turnsInRoles.none;

        return (
            <div className={styles.customTooltip}>
                <p className={styles.label}>{label}</p>
                {turnsInRoles.toPairs().map(([role, turnsInRole]) => (
                    <p className={styles.content} key={role}>
                        Turns as {role}: {turnsInRole} (
                        {((turnsInRole * 100) / turns).toFixed(1)}%)
                    </p>
                ))}
            </div>
        );
    }

    return null;
};

const Roles = (props: { output: OutputJSONType }) => {
    const colors = {
        president: "#00bbf9",
        judge: "#fee440",
        speaker: "#f15bb5",
        none: "#b2bec3",
    };

    const data = processRoleData(props.output);

    return (
        <div className={styles.root}>
            <p className={styles.text}>Role Visualisation</p>
            <ResponsiveContainer height={460} width="100%">
                <BarChart data={data} layout="vertical">
                    <YAxis type="category" dataKey="name" />
                    <XAxis type="number" domain={[0, "dataMax"]} />
                    <Tooltip
                        content={(props: CustomTooltipProps) =>
                            CustomTooltip({ ...props, data })
                        }
                    />
                    <Legend
                        verticalAlign="top"
                        payload={[
                            {
                                value: "President",
                                type: "square",
                                id: "ID01",
                                color: colors.president,
                            },
                            {
                                value: "Judge",
                                type: "square",
                                id: "ID02",
                                color: colors.judge,
                            },
                            {
                                value: "Speaker",
                                type: "square",
                                id: "ID03",
                                color: colors.speaker,
                            },
                            {
                                value: "None",
                                type: "square",
                                id: "ID04",
                                color: colors.none,
                            },
                        ]}
                    />
                    {data[0].roles.map((a, i) => [
                        _.toPairs(colors).map(([role, color]) => (
                            <Bar
                                dataKey={`roles[${i}].${role}`}
                                stackId="a"
                                fill={color}
                                key={`${i}${role}`}
                            />
                        )),
                    ])}
                </BarChart>
            </ResponsiveContainer>
            <p className={styles.graphLabel}>Turns</p>
        </div>
    );
};

export default Roles;

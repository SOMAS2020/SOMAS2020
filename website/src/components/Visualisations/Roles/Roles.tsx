import React, { useEffect, useState } from "react";
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
import { ProcessedRoleData, TeamAndTurns, RoleName } from "./Util/RoleTypes";
import { processRoleData } from "./Util/ProcessedRoleData";
import { OutputJSONType } from "../../../consts/types";

type CustomTooltipProps = {
  active: boolean;
  payload: [{ name: string; value: number; unit: string }];
  label: string;
  data: ProcessedRoleData;
  colors: Map<string, string>;
};

const CustomTooltip = ({ active, label, data, colors }: CustomTooltipProps) => {
  const getTurnsAsTeams = (role: RoleName): TeamAndTurns =>
    data
      .find((elem) => elem.role === role)
      ?.occupied?.reduce((acc, tAndT) => acc.add(tAndT), new TeamAndTurns()) ??
    new TeamAndTurns();

  if (active && data.length > 0) {
    const turnsAsTeams = getTurnsAsTeams(label as RoleName);
    const totalTurns = data[0].occupied.reduce(
      (acc, elem) => acc + elem.turns(),
      0
    );

    const newLabel = label === "Pres" ? "President" : label;

    return (
      <div className={styles.customTooltip}>
        <p className={styles.label}>{newLabel}</p>
        {turnsAsTeams.map((team, turns) => (
          <p
            className={styles.content}
            key={team}
            style={{ color: colors.get(team) }}
          >
            Turns as {team}: {turns} ({((turns * 100) / totalTurns).toFixed(1)}
            %)
          </p>
        ))}
      </div>
    );
  }

  return null;
};

const Roles = (props: { output: OutputJSONType }) => {
  const [data, setData] = useState(processRoleData(props.output));

  useEffect(() => {
    setData(processRoleData(props.output));
  }, [props.output]);

  const teams = ["Team1", "Team2", "Team3", "Team4", "Team5", "Team6"];
  const colors = new Map([
    ["Team1", "#0095FF"],
    ["Team2", "#FF0000"],
    ["Team3", "#802FF0"],
    ["Team4", "#00C49F"],
    ["Team5", "#FFBB28"],
    ["Team6", "#FF8042"],
  ]);

  return (
    <div className={styles.root}>
      <p className={styles.text}>Role Visualisation</p>
      <ResponsiveContainer height={460} width="100%">
        <BarChart data={data} layout="vertical">
          <YAxis type="category" dataKey="role" />
          <XAxis
            type="number"
            domain={[0, "dataMax"]}
            tickCount={20}
            allowDecimals={false}
          />
          <Tooltip
            content={(props: CustomTooltipProps) =>
              CustomTooltip({ ...props, data, colors })
            }
          />
          <Legend
            verticalAlign="top"
            payload={teams.map((team, i) => ({
              value: team,
              type: "square",
              id: `${team}${i}`,
              color: colors.get(team),
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
            )),
          ])}
        </BarChart>
      </ResponsiveContainer>
      <p className={styles.graphLabel}>Turns</p>
    </div>
  );
};

export default Roles;

import outputJSON from "../../../../output/output.json";
import { ProcessedRoleData, TurnsInRoles } from "./RoleTypes";

export const getProcessedRoleData = (): ProcessedRoleData => {
    if (outputJSON.GameStates.length === 0) return [];

    let allRoles: ProcessedRoleData = [];

    for (var id in outputJSON.GameStates[0].ClientInfos)
        allRoles.push({
            name: id,
            roles: [new TurnsInRoles()],
        });

    if (allRoles.length === 0) return [];

    /*
    METHOD:
        - For each gamestate
        - For each team in allRoles
        - Check if in power
        - If so increment power
        - If not increment none
        - standardise
        - return
    */

    return allRoles;
};

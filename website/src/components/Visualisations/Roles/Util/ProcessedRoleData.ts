import outputJSON from "../../../../output/output.json";
import {
    ProcessedRoleData,
    ProcessedRoleElement,
    TurnsInRoles,
} from "./RoleTypes";

const standardise = (allRoles: ProcessedRoleData): ProcessedRoleData => {
    const maxLength = allRoles.reduce(
        (acc, allRoleElem) =>
            allRoleElem.roles.length > acc ? allRoleElem.roles.length : acc,
        0
    );

    return allRoles.map((allRolesElem) => {
        for (let i = 0; i <= maxLength - allRolesElem.roles.length; i++) {
            allRolesElem.roles.push(new TurnsInRoles());
        }
        return allRolesElem;
    });
};

export const processRoleData = (data: typeof outputJSON) => {
    if (data.GameStates.length === 0) return [];

    let allRoles: ProcessedRoleData = [];

    for (var id in data.GameStates[0].ClientInfos)
        allRoles.push(new ProcessedRoleElement(id));

    if (allRoles.length === 0) return [];

    allRoles = data.GameStates.reduce(
        (allRolesNew, gameState) =>
            allRolesNew.map((allRolesElem) => {
                if (allRolesElem.name === gameState.PresidentID) {
                    allRolesElem.increment("president");
                } else if (allRolesElem.name === gameState.JudgeID) {
                    allRolesElem.increment("judge");
                } else if (allRolesElem.name === gameState.SpeakerID) {
                    allRolesElem.increment("speaker");
                } else {
                    allRolesElem.increment("none");
                }
                return allRolesElem;
            }),
        allRoles
    );

    return standardise(allRoles);
};

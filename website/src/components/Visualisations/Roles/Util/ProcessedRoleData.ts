import outputJSON from "../../../../output/output.json";
import { TurnsInRoles, ProcessedRoleData } from "./RoleTypes";

const hasBeenTouched = (tir: TurnsInRoles): boolean =>
    tir.president !== 0 ||
    tir.judge !== 0 ||
    tir.speaker !== 0 ||
    tir.none !== 0;

const idToIndex = (id: string): number => {
    switch (id) {
        case "Team1":
            return 0;
        case "Team2":
            return 1;
        case "Team3":
            return 2;
        case "Team4":
            return 3;
        case "Team5":
            return 4;
        case "Team6":
            return 5;
        default:
            return -1;
    }
};

export const getProcessedRoleData = (): ProcessedRoleData => {
    if (outputJSON.GameStates.length === 0) return [];

    let emptyTurnsInRoles: TurnsInRoles = {
        president: 0,
        judge: 0,
        speaker: 0,
        none: 0,
    };

    let allRoles: ProcessedRoleData = [];

    for (var id in outputJSON.GameStates[0].ClientInfos)
        allRoles.push({
            name: id,
            roles: [emptyTurnsInRoles],
        });

    if (allRoles.length === 0) return [];

    for (let gameState of outputJSON.GameStates) {
        let presidentID = gameState.PresidentID;
        let judgeID = gameState.JudgeID;
        let speakerID = gameState.SpeakerID;
    }

    return allRoles;
};

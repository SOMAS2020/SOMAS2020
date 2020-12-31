export type TurnsInRoles = {
    president: number;
    judge: number;
    speaker: number;
    none: number;
};

export type ProcessedRoleData = {
    name: string;
    roles: TurnsInRoles[];
}[];

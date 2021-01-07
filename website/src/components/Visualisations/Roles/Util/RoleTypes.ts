export type RoleName = "President" | "Judge" | "Speaker";

export class TeamAndTurns {
    Team1: number;
    Team2: number;
    Team3: number;
    Team4: number;
    Team5: number;
    Team6: number;

    constructor(team1: number = 0, team2: number = 0, team3: number = 0, team4: number = 0, team5: number = 0, team6: number = 0) {
        this.Team1 = team1;
        this.Team2 = team2;
        this.Team3 = team3;
        this.Team4 = team4;
        this.Team5 = team5;
        this.Team6 = team6;
    }

    set(key: string, val: number) {
        switch(key) {
            case "Team1": {
                this.Team1 = val;
                break;
            }
            case "Team2": {
                this.Team2 = val;
                break;
            }
            case "Team3": {
                this.Team3 = val;
                break;
            }
            case "Team4": {
                this.Team4 = val;
                break;
            }
            case "Team5": {
                this.Team5 = val;
                break;
            }
            case "Team6": {
                this.Team6 = val;
                break;
            }
        }
    }

    get(key: string): number {
        switch(key) {
            case "Team1": return this.Team1;
            case "Team2": return this.Team2;
            case "Team3": return this.Team3;
            case "Team4": return this.Team4;
            case "Team5": return this.Team5;
            case "Team6": return this.Team6;
            default: return 0;
        }
    }

    has(key: string): boolean {
        switch(key) {
            case "Team1": return this.Team1 !== 0;
            case "Team2": return this.Team2 !== 0;
            case "Team3": return this.Team3 !== 0;
            case "Team4": return this.Team4 !== 0;
            case "Team5": return this.Team5 !== 0;
            case "Team6": return this.Team6 !== 0;
            default: return false;
        }
    }

    increment(key: string, val: number = 1) {
        switch(key) {
            case "Team1": {
                this.Team1 += val;
                break;
            }
            case "Team2": {
                this.Team2 += val;
                break;
            }
            case "Team3": {
                this.Team3 += val;
                break;
            }
            case "Team4": {
                this.Team4 += val;
                break;
            }
            case "Team5": {
                this.Team5 += val;
                break;
            }
            case "Team6": {
                this.Team6 += val;
                break;
            }
        }
    }

    touched(): boolean {
        return this.Team1 !== 0
            || this.Team2 !== 0
            || this.Team3 !== 0
            || this.Team4 !== 0
            || this.Team5 !== 0
            || this.Team6 !== 0;
    }
};

export type ProcessedRoleElem = {
    role: RoleName,
    occupied: TeamAndTurns[];
};

export type ProcessedRoleData = ProcessedRoleElem[];

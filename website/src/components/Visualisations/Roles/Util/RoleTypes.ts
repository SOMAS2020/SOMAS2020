export type RoleName = 'Pres' | 'Judge' | 'Speaker'

export class TeamAndTurns {
    allTeams: Record<string, number>

    totalAgents: number

    constructor(totalAgents: number, preMadeDict: Record<string, number> = {}) {
        this.allTeams = preMadeDict
        this.totalAgents = totalAgents

        if (Object.keys(this.allTeams).length === 0) {
            for (let i = 0; i < totalAgents; i++) {
                this.allTeams[`Team${i + 1}`] = 0
            }
            this.allTeams.NotRun = 0
        }
    }

    set(key: string, val: number) {
        this.allTeams[key] = val
    }

    get(key: string): number {
        return this.allTeams[key]
    }

    has(key: string): boolean {
        return key in this.allTeams && this.allTeams[key] !== 0
    }

    increment(key: string, val: number = 1) {
        this.allTeams[key] += val
    }

    touched(): boolean {
        const teams = Object.values(this.allTeams)

        for (let i = 0; i < teams.length; i++) {
            if (teams[i] !== 0) {
                return true
            }
        }
        return false
    }

    turns(): number {
        return Object.values(this.allTeams).reduce((a, b) => a + b)
    }

    add(teamAndTurns: TeamAndTurns): TeamAndTurns {
        const adder: Record<string, number> = {}
        for (let i = 0; i < this.totalAgents; i++) {
            adder[`Team${i + 1}`] =
                this.allTeams[`Team${i + 1}`] +
                teamAndTurns.allTeams[`Team${i + 1}`]
        }
        adder.NotRun = this.allTeams.NotRun + teamAndTurns.allTeams.NotRun

        return new TeamAndTurns(this.totalAgents, adder)
    }

    map<T>(func: (team: string, turns: number) => T): T[] {
        const mapper: T[] = []
        for (let i = 0; i < this.totalAgents; i++) {
            mapper.push(func(`Team${i + 1}`, this.allTeams[`Team${i + 1}`]))
        }
        mapper.push(func('NotRun', this.allTeams.NotRun))

        return mapper
    }
}

export type ProcessedRoleElem = {
    role: RoleName
    occupied: TeamAndTurns[]
}

export type ProcessedRoleData = ProcessedRoleElem[]

export type IIGOInfo = {
    turn: number
    status: string
}

export type IIGOInfos = IIGOInfo[]

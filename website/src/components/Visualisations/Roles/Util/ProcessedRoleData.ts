import { OutputJSONType } from '../../../../consts/types'
import { ProcessedRoleData, TeamAndTurns } from './RoleTypes'

const standardise = (data: ProcessedRoleData): ProcessedRoleData => {
    const maxLength = data.reduce(
        (max, curr) =>
            curr.occupied.length > max ? curr.occupied.length : max,
        0
    )

    return data.map((elem) => {
        const { length } = elem.occupied
        for (let i = 0; i < maxLength - length; i++) {
            elem.occupied.push(new TeamAndTurns())
        }
        return elem
    })
}

const increment = (occupied: TeamAndTurns[], team: string): TeamAndTurns[] => {
    if (occupied.length > 0 && occupied[occupied.length - 1].has(team)) {
        occupied[occupied.length - 1].increment(team)
    } else {
        const teamAndTurns = new TeamAndTurns()
        teamAndTurns.set(team, 1)
        occupied.push(teamAndTurns)
    }
    return occupied
}

export const processRoleData = (data: OutputJSONType): ProcessedRoleData => {
    if (data.GameStates.length === 0) return []

    const retData: ProcessedRoleData = [
        {
            role: 'Pres',
            occupied: [],
        },
        {
            role: 'Judge',
            occupied: [],
        },
        {
            role: 'Speaker',
            occupied: [],
        },
    ]

    return standardise(
        retData.map((elem) => {
            elem.occupied = data.GameStates.reduce((acc, gameState) => {
                switch (elem.role) {
                    case 'Pres': {
                        elem.occupied = increment(
                            elem.occupied,
                            gameState.PresidentID
                        )
                        break
                    }
                    case 'Judge': {
                        elem.occupied = increment(
                            elem.occupied,
                            gameState.JudgeID
                        )
                        break
                    }
                    case 'Speaker': {
                        elem.occupied = increment(
                            elem.occupied,
                            gameState.SpeakerID
                        )
                        break
                    }
                    default:
                        break
                }
                return acc
            }, elem.occupied)
            return elem
        })
    )
}
export default processRoleData

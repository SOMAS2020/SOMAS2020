export type ForagingTurn = {
    turn: number
    deerInputResources: number
    deerNumParticipants: number
    deerNumCaught: number
    deerTotalUtility: number
    fishInputResources: number
    fishNumParticipants: number
    fishNumCaught: number
    fishTotalUtility: number
}

export type ForagingHistory = ForagingTurn[]

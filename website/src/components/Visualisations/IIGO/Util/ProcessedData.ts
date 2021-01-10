import { OutputJSONType } from '../../../../consts/types'

export type RuleType = {
    ruleName: string
    mutable: boolean
    linked: boolean
    variables: string[]
    history: { season: number; turn: number }[]
}

export const processRulesData = (data: OutputJSONType): RuleType[] => {
    if (data.GameStates.length === 0) return []

    // return CurrentRulesInPlay keys in term of seasons
    const rulesInSeasons = data.GameStates.map((episode) => {
        return {
            season: episode.Season,
            turn: episode.Turn,
            rules: Object.keys(episode.CurrentRulesInPlay),
        }
    })

    // return a list of rules of RuleType
    const rulesDict: RuleType[] = []
    // each season, do...
    data.GameStates.forEach((episode) => {
        // each rules in season, do...
        Object.keys(episode.CurrentRulesInPlay).forEach((rules) => {
            // add history
            const history: any = []
            rulesInSeasons.forEach((item) => {
                if (item.rules.includes(rules)) {
                    history.push({ season: item.season, turn: item.turn })
                }
            })
            rulesDict[rules] = {
                ruleName: episode.CurrentRulesInPlay[rules].RuleName,
                mutable: episode.CurrentRulesInPlay[rules].Mutable,
                linked: episode.CurrentRulesInPlay[rules].Linked,
                variables: episode.CurrentRulesInPlay[rules].RequiredVariables,
                history,
            }
        })
    })

    return Object.values(rulesDict)
}

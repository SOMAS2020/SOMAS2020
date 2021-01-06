import { OutputJSONType } from "../../../../consts/types";


export type RuleType = {
    ruleName: string;
    mutable: boolean;
    linked: boolean;
    variables: string[];
    history: {season: number; turn: number}[];
};

export const processRulesData = (data: OutputJSONType) :RuleType[]=> {
    if (data.GameStates.length === 0) return []

    // return CurrentRulesInPlay keys in term of seasons
    let rulesInSeasons = data.GameStates.map((episode) => {
        return {
            season: episode.Season,
            turn: episode.Turn,
            rules: Object.keys(episode.CurrentRulesInPlay)
        }
    })

    // return a list of rules of RuleType
    let rulesDict:any = {}
    // each season, do...
    data.GameStates.forEach((episode) => {
        // each rules in season, do...
        Object.keys(episode.CurrentRulesInPlay).forEach((rules) => {
            // add history
            let history_:any = []
            rulesInSeasons.forEach((item) => {
                if (item.rules.includes(rules)){
                    history_.push({season:item.season, turn:item.turn})
                }
            })
            rulesDict[rules] = {
                ruleName: episode.CurrentRulesInPlay[rules].RuleName,
                mutable: episode.CurrentRulesInPlay[rules].Mutable,
                linked: episode.CurrentRulesInPlay[rules].Linked,
                variables: episode.CurrentRulesInPlay[rules].RequiredVariables,
                history: history_,
            }
        })
    })

    return Object.values(rulesDict);
}


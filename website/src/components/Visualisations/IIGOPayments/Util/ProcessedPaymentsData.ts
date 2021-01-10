import { OutputJSONType, TeamName } from '../../../../consts/types'
import { ProcessedTaxData, TaxData } from './IIGOPaymentsTypes'

export const processPaymentsData = (data: OutputJSONType): ProcessedTaxData => {
    if (data.GameStates.length === 0) return []

    const retData: ProcessedTaxData = []
    const IIGOHistory = Object.entries(
        data.GameStates[data.GameStates.length - 1].IIGOHistory
    )

    const empty = (): TaxData => {
        return {
            expectedAlloc: 0,
            actualAlloc: 0,
            expectedTax: 0,
            actualTax: 0,
            expectedSanction: 0,
            actualSanction: 0,
        }
    }

    const TaxInfo: Map<TeamName, TaxData> = new Map<TeamName, TaxData>([
        [1, empty()],
        [2, empty()],
        [3, empty()],
        [4, empty()],
        [5, empty()],
        [6, empty()],
    ])

    IIGOHistory.forEach(([turnNumber, exchanges]) => {
        if (exchanges) {
            exchanges.forEach((teamAction) => {
                const type = teamAction.Pairs[0].VariableName
                switch (type) {
                    case 'IslandAllocation':
                        {
                            const team =
                                TeamName[
                                    teamAction.ClientID as keyof typeof TeamName
                                ]
                            const current: TaxData | undefined = TaxInfo.get(
                                team
                            )
                            if (current) {
                                current.actualAlloc +=
                                    teamAction.Pairs[0].Values[0]
                                current.expectedAlloc +=
                                    teamAction.Pairs[1].Values[0]
                                TaxInfo.set(team, current)
                            }
                        }
                        break
                    case 'IslandTaxContribution':
                        {
                            const team =
                                TeamName[
                                    teamAction.ClientID as keyof typeof TeamName
                                ]
                            const current: TaxData | undefined = TaxInfo.get(
                                team
                            )
                            if (current) {
                                current.actualTax +=
                                    teamAction.Pairs[0].Values[0]
                                current.expectedTax +=
                                    teamAction.Pairs[1].Values[0]
                                TaxInfo.set(team, current)
                            }
                        }
                        break
                    case 'SanctionPaid':
                        {
                            const team =
                                TeamName[
                                    teamAction.ClientID as keyof typeof TeamName
                                ]
                            const current: TaxData | undefined = TaxInfo.get(
                                team
                            )
                            if (teamAction.Pairs[0].Values[0] !== 0) {
                                console.log(turnNumber, teamAction)
                            }
                            if (current) {
                                current.actualSanction +=
                                    teamAction.Pairs[0].Values[0]
                                current.expectedSanction +=
                                    teamAction.Pairs[1].Values[0]
                                TaxInfo.set(team, current)
                            }
                        }
                        break
                    default:
                        break
                }
            })
        }
    })
    console.log(TaxInfo)

    TaxInfo.forEach((value, key) => {
        retData.push({
            name: key,
            actualTax: value.actualTax,
            expectedTax: value.expectedTax,
            actualAlloc: value.actualAlloc,
            expectedAlloc: value.expectedAlloc,
            actualSanction: value.actualSanction,
            expectedSanction: value.expectedSanction,
        })
    })

    return retData
}
export default processPaymentsData

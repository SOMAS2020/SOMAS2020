import { OutputJSONType } from '../../../../consts/types'
import { ProcessedTaxData, TaxData } from './IIGOPaymentsTypes'
import { TeamNameGen, numAgents } from '../../utils'

export const processPaymentsData = (data: OutputJSONType): ProcessedTaxData => {
    if (data.GameStates.length === 0) return []

    const retData: ProcessedTaxData = []
    const IIGOHistory = Object.entries(
        data.GameStates[data.GameStates.length - 1].IIGOHistory
    )

    const totalAgents = numAgents(data)
    const TeamName = TeamNameGen(totalAgents)

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

    const TaxInfo: Record<number, TaxData> = {}

    for (let i = 0; i < totalAgents; i++) {
        TaxInfo[i + 1] = empty()
    }

    IIGOHistory.forEach(([_, exchanges]) => {
        if (exchanges) {
            exchanges.forEach((teamAction) => {
                const team = TeamName[teamAction.ClientID]
                const current: TaxData | undefined = TaxInfo[team]
                const type = teamAction.Pairs[0].VariableName
                if (current) {
                    switch (type) {
                        case 'IslandAllocation':
                            current.actualAlloc += teamAction.Pairs[0].Values[0]
                            current.expectedAlloc +=
                                teamAction.Pairs[1].Values[0]
                            TaxInfo[team] = current
                            break
                        case 'IslandTaxContribution':
                            current.actualTax += teamAction.Pairs[0].Values[0]
                            current.expectedTax += teamAction.Pairs[1].Values[0]
                            TaxInfo[team] = current
                            break
                        case 'SanctionPaid':
                            current.actualSanction +=
                                teamAction.Pairs[0].Values[0]
                            current.expectedSanction +=
                                teamAction.Pairs[1].Values[0]
                            TaxInfo[team] = current
                            break
                        default:
                            break
                    }
                }
            })
        }
    })

    Object.entries(TaxInfo).forEach((entry) => {
        const [key, value] = entry
        retData.push({
            name: Number(key),
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

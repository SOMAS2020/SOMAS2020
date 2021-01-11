import { TeamName } from '../../../../consts/types'

export type RoleName = 'Pres' | 'Judge' | 'Speaker'

export type TaxData = {
    expectedTax: number
    actualTax: number
    actualAlloc: number
    expectedAlloc: number
    actualSanction: number
    expectedSanction: number
}

export type ProcessedTaxDataElem = {
    name: TeamName
    expectedTax: number
    actualTax: number
    actualAlloc: number
    expectedAlloc: number
    actualSanction: number
    expectedSanction: number
}

export type ProcessedTaxData = ProcessedTaxDataElem[]

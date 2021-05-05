import * as localForage from 'localforage'
import { VIS_OUTPUT } from '../../consts/localForage'
import { OutputJSONType } from '../../consts/types'

export const loadLocalVisOutput = async () => {
    const got: OutputJSONType | null | undefined = await localForage.getItem(
        VIS_OUTPUT
    )
    if (got) return got
    return undefined
}

export const storeLocalVisOutput = async (o: OutputJSONType) => {
    await localForage.setItem(VIS_OUTPUT, o)
}

export const clearLocalVisOutput = async () => {
    await localForage.removeItem(VIS_OUTPUT)
}

// roles atm to fix
// export const teamColors = new Map([
//     ['Team0', '#101D42'],
//     ['Team1', '#0095FF'],
//     ['Team2', '#FF0000'],
//     ['Team3', '#802FF0'],
//     ['Team4', '#00C49F'],
//     ['Team5', '#FFBB28'],
//     ['Team6', '#FF8042'],
//     ['Team7', '#FF8042'],
//     ['Team8', '#FF8042'],
// ])

export const numAgents = (data: OutputJSONType): number => {
    const cis = data.GameStates[0].ClientInfos
    return Object.keys(cis).length
}

export const TeamNameGen = (agents: number) => {
    const out = { CommonPool: 0 }
    for (let i = 0; i < agents; i++) {
        out[`Team${i + 1}`] = i + 1
    }
    return out
}

const hslToHex = (h: number, s: number, l: number) => {
    l /= 100
    const a = (s * Math.min(l, 1 - l)) / 100
    const f = (n) => {
        const k = (n + h / 30) % 12
        const color = l - a * Math.max(Math.min(k - 3, 9 - k, 1), -1)
        return Math.round(255 * color)
            .toString(16)
            .padStart(2, '0')
    }
    return `#${f(0)}${f(8)}${f(4)}`
}

export const generateColours = (allTeams: number) => {
    const legendColours = {
        CommonPool: '#ACE600',
        TotalResources: '#FF69B4',
        CriticalThreshold: '#B7B4B0',
    }

    for (let i = 0; i < allTeams; i++) {
        const hue = Math.floor(360 / 12) * i
        const saturation = 90 + Math.random() * 10
        const lightness = 50 + Math.random() * 10

        legendColours[`team${i + 1}`] = hslToHex(hue, saturation, lightness)
    }

    return legendColours
}

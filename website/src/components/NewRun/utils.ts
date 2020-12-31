import { notEmpty } from "../../utils";
import { Flag } from "../../wasmAPI";

export type GoIota = {
    i: number
    name: string
}

/**
 * Tries to get go iotas for the flag. This relies on the Usage string printing
 * separate lines in the format
 * /^([0-9]+): ([a-zA-Z]+)$/
 * for each iota, in each line.
 * 
 * @param flag The flag in question
 */
export const tryGetGoIota = (flag: Flag): GoIota[] | undefined => {
    const t = flag.Type
    if (
        t !== `int` &&
        t !== `int64` &&
        t !== `uint` &&
        t !== `uint64`
    ) {
        return undefined
    }

    const iotaRegex = /^([0-9]+): ([a-zA-Z]+)$/

    // try to parse usage for iota entries
    const iotas: GoIota[] = flag.Usage
        .split(`\n`)
        .map(line => {
            const m = line.match(iotaRegex)

            if (!m || m.length !== 3) return undefined

            return {
                i: parseInt(m[1]),
                name: m[2],
            }
        })
        .filter(notEmpty)

    if (iotas.length > 0) {
        return iotas
    }
    return undefined
}

/**
 * Get a string invalid reason for an int flag if present.
 * @param flag The flag in question
 */
const validateIntegerFlag = async (flag: Flag) => {
    const t = flag.Type
    if (
        t !== `int` &&
        t !== `int64` &&
        t !== `uint` &&
        t !== `uint64`
    ) {
        return undefined
    }

    const int32Max = BigInt(`2147483647`)
    const int32Min = -int32Max-BigInt(`1`)
    const uint32Max = BigInt(`4294967295`)
    const int64Max = BigInt(`9223372036854775807`)
    const int64Min = -int64Max-BigInt(`1`)
    const uint64Max = BigInt(`18446744073709551615`)

    if (!flag.Value.match(/^[0-9]+$/)) {
        return `"${flag.Value}" is not an integer`
    }
    const val = BigInt(flag.Value)
    switch (flag.Type) {
        case `int`:
            if (val < int32Min || val > int32Max) {
                return `int must be in signed 32-bit range [${String(int32Min)}, ${String(int32Max)}]`
            }
            break
        case `uint`:
            if (val < 0 || val > uint32Max) {
                return `uint must be in unsigned 32-bit range [0, ${String(uint32Max)}]`
            }
            break
        case `int64`:
            if (val < int64Min || val > int64Max) {
                return `int64 must be in signed 64-bit range [${String(int64Min)}, ${String(int64Max)}]`
            }
            break
        case `uint64`:
            if (val < 0 || val > uint64Max) {
                return `uint64 must be in unsigned 64-bit range [0, ${String(uint64Max)}]`
            }
            break
        default:
            return `Should not happen! Report a bug.`
    }
    return undefined
}

export const setFlagWithValidation = async (flag: Flag, newValue: string): Promise<Flag> => {
    const newValuedFlag: Flag = {
        ...flag,
        Value: newValue,
    }
    switch (flag.Type) {
        case `string`:
            return { 
                ...newValuedFlag,
                InvalidReason: undefined,
            }
        case `int`:
        case `int64`:
        case `uint`:
        case `uint64`:
            const intInvalid = await validateIntegerFlag(newValuedFlag)
            if (intInvalid) {
                return {
                    ...newValuedFlag,
                    InvalidReason: intInvalid,
                }
            }
            const goIotas = tryGetGoIota(newValuedFlag)
            if (goIotas) {
                const val = parseInt(newValue)
                if (!goIotas.map(g => g.i).includes(val)) {
                    return {
                        ...newValuedFlag,
                        InvalidReason: `${flag.Value} is not a valid iota value for this flag`
                    }
                }
            }
            return {
                ...newValuedFlag,
                InvalidReason: undefined,
            }
        case `float64`:
            // TODO
            return {
                ...newValuedFlag,
                InvalidReason: undefined,
            }
            
        case `bool`:
            // @ts-ignore this is weird.
            if (newValue !== `true` || newValue !== `false`) {
                return {
                    ...newValuedFlag,
                    InvalidReason: `Bool flags can only take "true" or "false"`,
                }
            }
            return {
                ...newValuedFlag,
                InvalidReason: undefined,
            }
        default: 
            return {
                ...newValuedFlag,
                InvalidReason: `Unknown type "${flag.Type}"`
            }
    }
}
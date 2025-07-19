import * as dotenv from "dotenv"
dotenv.config()

export function getEnvString(key :string, fallback: string) {
    const value = process.env[key];
    if (!value)
        return fallback
    return value
}

export function getEnvInt(key :string, fallback: number) {
    const value = process.env[key];
    if (!value)
        return fallback
    try {
        return parseInt(value)
    } catch {
        return fallback
    }
}

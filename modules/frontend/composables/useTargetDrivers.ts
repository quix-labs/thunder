export type TargetDriver = {
    config: {
        ID: string,
        name: string
        image?: string
        notes?: string[]
    },
    fields: Array<{
        name: string
        label: string
        type: string
        required: boolean
        help?: string
        default?: string
        min?: string
    }>
}
export type TargetDrivers = {
    [key: string]: TargetDriver
}

export default function useTargetDrivers() {
    return useGoFetch<TargetDrivers>('/target-drivers')
}
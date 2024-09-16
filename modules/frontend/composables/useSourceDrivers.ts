export type Driver = {
    config: {
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
    }>
}
export type SourceDrivers = {
    [key: string]: Driver
}


export default () => useGoLazyAsyncData<SourceDrivers>('source-drivers', '/source-drivers')

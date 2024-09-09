export type ProcessFormState = {
    activeTabs: "source" | "mapping" | "output"
    activeMappingPath?: string
    preventScroll: boolean
    shortMode: boolean
    liveReload: boolean
    showSamplePanel: boolean
}
export default function useProcessFormState() {
    return useState<ProcessFormState>('process-form', () => toRef({
        activeTabs: 'source',
        preventScroll: false,
        shortMode: false,
        liveReload: false,
        showSamplePanel: true
    }))
}
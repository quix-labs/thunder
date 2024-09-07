export type ProcessFormState = {
    activeTabs: "source" | "mapping" | "output"
    activeMappingPath?: string
    preventScroll: boolean
}
export default function useProcessFormState() {
    return useState<ProcessFormState>('process-form', () => toRef({
        activeTabs: 'source',
        preventScroll: false
    }))
}

export default function useTargets() {
    return useGoFetch<any>('/targets')
}
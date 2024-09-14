export default function useProcessorsStatus() {
    return useGoFetch<any>('/processors/status')
}
import useGoFetch from "./useGoFetch";

export default function useSources() {
    return useGoFetch<any>('/sources')
}
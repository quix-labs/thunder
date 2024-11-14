import {joinURL} from "ufo";

export const getGoApiUrl = (path?: string) => {
    return joinURL('/go-api', path || "")
}

const useGoFetch: typeof useFetch = (path, options = {}) => {
    return useFetch(getGoApiUrl(path.toString()), {server: false, ...options, lazy: true})
}
export default useGoFetch